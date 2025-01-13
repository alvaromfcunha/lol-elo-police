run:
	go run cmd/app/main.go

register:
	go run cmd/register/main.go

deregister:
	rm db/wpp.db

sqlc:
	rm -rf internal/generated/database
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

prod:
	mkdir out
	mkdir out/bin
	mkdir out/infrastructure
	mkdir out/infrastructure/database
	mkdir out/infrastructure/config
	mkdir out/infrastructure/template
	goose -dir infrastructure/database/migrations sqlite3 out/infrastructure/database/app.db up
	cp infrastructure/database/whatsapp.db out/infrastructure/database/whatsapp.db
	cp infrastructure/config/.env.example out/infrastructure/config/.env
	cp infrastructure/template/messages.txt out/infrastructure/template/messages.txt
	go build -o out/bin/app cmd/app/main.go
