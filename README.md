## repositories

Middleware::
https://github.com/go-chi/chi

Hot Reloading::
https://github.com/air-verse/air

ENV Management::
https://direnv.net/

Database Migration::
https://github.com/golang-migrate/migrate
https://github.com/pressly/goose

Migration Commands::

```bash
migrate create -seq -ext sql -dir ./cmd/migrate/migrations create_users

migrate -path=./cmd/migrate/migrations -database="postgres://admin:pass@localhost/social?sslmode=disable" up
migrate -path=./cmd/migrate/migrations -database="postgres://admin:pass@localhost/social?sslmode=disable" down

make migration [posts_create]

```
