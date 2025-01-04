run:
	go run cmd/app/main.go

register:
	go run cmd/register/main.go

deregister:
	rm db/wpp.db

sqlc:
	sqlc generate -f infrastructure/database/sqlc.yml

migration:
	goose -dir infrastructure/database/migrations sqlite3 infrastructure/database/app.db create $(NAME) sql

up:
	goose -dir infrastructure/database/migrations sqlite3 infrastructure/database/app.db up

reset:
	goose -dir infrastructure/database/migrations sqlite3 infrastructure/database/app.db reset

status:
	goose -dir infrastructure/database/migrations sqlite3 infrastructure/database/app.db status

build:
	go build -o bin/app cmd/app/main.go && go build -o bin/register cmd/register/main.go
