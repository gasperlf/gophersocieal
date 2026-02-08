package main

import (
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"ontopsolutions.net/gasperlf/social/internal/auth"
	"ontopsolutions.net/gasperlf/social/internal/db"
	"ontopsolutions.net/gasperlf/social/internal/env"
	"ontopsolutions.net/gasperlf/social/internal/mailer"
	"ontopsolutions.net/gasperlf/social/internal/store"
	"ontopsolutions.net/gasperlf/social/internal/store/cache"
)

const version = "0.0.1"

//	@title			GopherSocial API
//	@version		0.0.1
//	@description	API for GopherSocial, a social network for gophers.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {

	cfg := config{
		addr:        env.GetString("API_ADDR", ":8081"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8081"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:pass@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       env.GetInt("REDIS_DB", 0),
			enabled:  env.GetBool("REDIS_ENABLED", true),
		},
		env: env.GetString("APP_ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAIL_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: "username",
				pass: "password",
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "!·.00000999iiillllfr13434343"),
				exp:    time.Hour * 24 * 3, //3 days
				iss:    "gophersocial",
			},
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	// Redis client initialization would go here if needed
	var rdb *redis.Client
	if cfg.redisCfg.enabled {
		rdb = cache.NewRedisClient(
			cfg.redisCfg.addr,
			cfg.redisCfg.password,
			cfg.redisCfg.db,
		)

		defer func() {
			if err := rdb.Close(); err != nil {
				logger.Error("failed to close redis client", "error", err)
			}
		}()
		logger.Info("redis cache initialized")
	}

	logger.Info("redis client initialized")

	store := store.NewStorage(db)
	cacheStore := cache.NewRedisStorage(rdb)

	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
		cacheStore:    cacheStore,
	}

	mux := mount(app)
	logger.Fatal(app.run(mux))

}
