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

Validation Library::
github.com/go-playground/validator/v10

to run config for swagger

https://github.com/swaggo/http-swagger?tab=readme-ov-file
check env variables and export

```bash
swag -version
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

```

Logging Library::

```bash
https://github.com/uber-go/zap
```

CORS:

```bash
go get github.com/go-chi/cors
```

performance test

https://www.npmjs.com/package/autocannon

```bash
npx autocannon -c 100 -d 10 http://localhost:4000/api/v1/posts
```
