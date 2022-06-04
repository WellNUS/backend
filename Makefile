postgres:
	docker run --name postgres14.3 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14.3-alpine

createdb:
	docker exec -it postgres14.3 createdb --username=root --owner=root wellnus

dropdb:
	docker exec -it postgres14.3 dropdb wellnus

migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/wellnus?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5432/wellnus?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown

