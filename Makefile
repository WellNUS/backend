include .env

postgres:
	docker run --name postgres14.3 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14.3-alpine

startdb:
	docker start postgres14.3

stopdb:
	docker stop postgres14.3

createdb:
	docker exec -it postgres14.3 createdb --username=root --owner=root wellnus

dropdb:
	docker exec -it postgres14.3 dropdb wellnus

migrateup:
	migrate -path db/migration -database "$(POSTGRES_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(POSTGRES_URL)" -verbose down

unittest:
	go test $(shell go list ./unit_test/...| grep -v test_helper) -p 1

.PHONY: postgres startdb stopdb createdb dropdb migrateup migratedown unittest

