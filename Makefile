include .env

all: dev

dev: composedown startdb
	go run main.go

prod: composedown
	docker compose up -d

migrateup: startdb
	migrate -path db/migration -database "$(DB_ADDRESS)" -verbose up

migratedown: startdb
	migrate -path db/migration -database "$(DB_ADDRESS)" -verbose down

unittest: startdb
	go test $(shell go list ./unit_test/...| grep -v test_helper) -p 1

startdb:
	docker compose up -d db

composedown:
	docker compose down --rmi local

purgedb:
	sudo chmod -R 0777 ./.db_data/ && rm -rf ./.db_data/ || echo "No .db_data/ to purge"

.PHONY: all dev prod migrateup migratedown startdb composeDown purgeDB unittest

