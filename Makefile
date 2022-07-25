include .env

# Commands using docker compose
composeup:
	docker compose -p $(DOCKER_COMPOSE_NAME) up --build

composedown:
	docker compose -p $(DOCKER_COMPOSE_NAME) down

# Commands not using docker compose
initpg:
	docker run --name $(DB_NAME) -p $(DB_PORT):$(DB_PORT) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres:14.3-alpine

deinitpg:
	docker rm $(DB_NAME)

startdb:
	docker start $(DB_NAME)

stopdb:
	docker stop $(DB_NAME)

createdb:
	docker exec -it $(DB_NAME) createdb --username=$(DB_USER) --owner=$(DB_USER) wellnus

dropdb:
	docker exec -it $(DB_NAME) dropdb wellnus

migrateup:
	migrate -path db/migration -database "$(DB_ADDRESS)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_ADDRESS)" -verbose down

unittest:
	go test $(shell go list ./unit_test/...| grep -v test_helper) -p 1

.PHONY: composeup composedown initpg deinitpg startdb stopdb createdb dropdb migrateup migratedown unittest

