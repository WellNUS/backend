include .env

all: dev

dev: startdb
	go run main.go

prod: 
	docker compose up

migrateup: startdb
	migrate -path db/migration -database "$(DB_ADDRESS)" -verbose up

migratedown: startdb
	migrate -path db/migration -database "$(DB_ADDRESS)" -verbose down

startdb:
	docker compose up -d db

composeDown:
	docker compose down --rmi local

purgeDB:
	sudo chmod -R 0777 ./.db_data/ && rm -rf ./.db_data/ || echo "No .db_data/ to purge"

unittest:
	go test $(shell go list ./unit_test/...| grep -v test_helper) -p 1

.PHONY: all dev prod migrateup migratedown startdb composeDown purgeDB unittest

