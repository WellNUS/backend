include .env

postgres:
	docker run --name postgres14.3 -p $(POSTGRES_PORT):$(POSTGRES_PORT) -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) -d postgres:14.3-alpine

startdb:
	docker start postgres14.3

stopdb:
	docker stop postgres14.3

createdb:
	docker exec -it postgres14.3 createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) wellnus

dropdb:
	docker exec -it postgres14.3 dropdb wellnus

migrateup:
	migrate -path db/migration -database "$(POSTGRES_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(POSTGRES_URL)" -verbose down

unittest:
	go test $(shell go list ./unit_test/...| grep -v test_helper) -p 1

.PHONY: postgres startdb stopdb createdb dropdb migrateup migratedown unittest

