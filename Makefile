postgres:
	docker run --name postgres14.3 -p 5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14.3-alpine

startdb:
	docker start postgres14.3

stopdb:
	docker stop postgres14.3

createdb:
	docker exec -it postgres14.3 createdb --username=root --owner=root wellnus

dropdb:
	docker exec -it postgres14.3 dropdb wellnus

migrateup:
	migrate -path db/migration -database "postgresql://root:password@0.0.0.0:49730/wellnus?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:password@0.0.0.0:49730/wellnus?sslmode=disable" -verbose down

.PHONY: postgres startdb stopdb createdb dropdb migrateup migratedown

