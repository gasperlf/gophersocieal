package main

import (
	"go.uber.org/zap"
	"ontopsolutions.net/gasperlf/social/internal/db"
	"ontopsolutions.net/gasperlf/social/internal/env"
	"ontopsolutions.net/gasperlf/social/internal/store"
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
		addr:   env.GetString("API_ADDR", ":8081"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8081"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:pass@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("APP_ENV", "development"),
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

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}

	mux := mount(app)
	logger.Fatal(app.run(mux))

}
