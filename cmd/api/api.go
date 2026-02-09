package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
	"ontopsolutions.net/gasperlf/social/docs"
	"ontopsolutions.net/gasperlf/social/internal/auth"
	"ontopsolutions.net/gasperlf/social/internal/mailer"
	"ontopsolutions.net/gasperlf/social/internal/ratelimiter"
	"ontopsolutions.net/gasperlf/social/internal/store"
	"ontopsolutions.net/gasperlf/social/internal/store/cache"
)

type application struct {
	config        config
	store         store.Storage
	cacheStore    cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redisCfg    redisConfig
	rateLimiter ratelimiter.Config
}

type redisConfig struct {
	addr     string
	password string
	db       int
	enabled  bool
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	user string
	pass string
}

type mailConfig struct {
	fromEmail string
	exp       time.Duration
	sendGrid  sendGridConfig
	mailTrap  mailTrapConfig
}

type sendGridConfig struct {
	apiKey string
}

type mailTrapConfig struct {
	apiKey string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func mount(app *application) http.Handler {

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(app.RateLimiterMiddleware)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.With(app.BasicMidleware()).
			Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.DeletePostHandler)
				r.Patch("/", app.CheckPostOwnership("moderator", app.UpdatePostHandler))
				r.Post("/comments", app.CheckPostOwnership("admin", app.createCommentPostHandler))
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})

		})
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Host = app.config.apiURL

	srvr := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		s := <-c

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("caught shutdown signal", "signal", s.String())

		shutdown <- srvr.Shutdown(ctx)
	}()

	app.logger.Infow("server has started at", "addr", app.config.addr, "env", app.config.env, "version", version)

	err := srvr.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown

	if err != nil {
		return err
	}

	return nil
}
