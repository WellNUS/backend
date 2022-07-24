include .env

postgres:
	docker run --name postgres14.3 -p $(DB_PORT):$(DB_PORT) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:14.3-alpine

startdb:
	docker start postgres14.3

stopdb:
	docker stop postgres14.3

createdb:
	docker exec -it postgres14.3 createdb --username=$(DB_USER) --owner=$(DB_USER) wellnus

dropdb:
	docker exec -it postgres14.3 dropdb wellnus

migrateup:
	migrate -path db/migration -database "$(DB_ADDRESS)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_ADDRESS)" -verbose down

unittest:
	go test $(shell go list ./unit_test/...| grep -v test_helper) -p 1

build:
	go build -o ./bin/main main.go

main: build
	./bin/main

.PHONY: postgres startdb stopdb createdb dropdb migrateup migratedown unittest build main

