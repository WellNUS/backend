include .env

start:
	docker compose -p $(DOCKER_STACK_NAME) up --build

close:
	docker compose -p $(DOCKER_STACK_NAME) down

postgres:
	docker run --name $(DOCKER_STACK_NAME)-postgres-1 -p $(DB_PORT):$(DB_PORT) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:14.3-alpine

startdb:
	docker start $(DOCKER_STACK_NAME)-postgres-1

stopdb:
	docker stop $(DOCKER_STACK_NAME)-postgres-1

createdb:
	docker exec -it $(DOCKER_STACK_NAME)-postgres-1 createdb --username=$(DB_USER) --owner=$(DB_USER) wellnus

dropdb:
	docker exec -it $(DOCKER_STACK_NAME)-postgres-1 dropdb wellnus

migrateup:
	migrate -path db/migration -database "$(DB_ADDRESS)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_ADDRESS)" -verbose down

unittest:
	go test $(shell go list ./unit_test/...| grep -v test_helper) -p 1

.PHONY: server close postgres startdb stopdb createdb dropdb migrateup migratedown unittest

