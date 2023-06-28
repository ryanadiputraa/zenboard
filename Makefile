postgres:
	docker run --name postgres15.3 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -p 5432:5432 -d postgres:15.3-alpine

createdb:
	docker exec -it postgres15.3 createdb --username=root --owner=root zenboard

dropdb:
	docker exec -it postgres15.3 dropdb zenboard

migrateup:
	migrate -path pkg/db/migration -database "postgresql://root:root@localhost:5432/zenboard?sslmode=disable" -verbose up

migratedown:
	migrate -path pkg/db/migration -database "postgresql://root:root@localhost:5432/zenboard?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb server