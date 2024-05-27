rundb:
	docker container run \
         --name mom-bank \
         -p 5432:5432 \
         -e POSTGRES_USER=root \
         -e POSTGRES_PASSWORD=secret \
         -d --rm \
        postgres:16.2-alpine3.19

stopdb:
	docker container stop mom-bank

createdb:
	docker exec -it mom-bank createdb --username=root --owner=root mom_bank

dropdb:
	docker exec -it mom-bank dropdb mom_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/mom_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/mom_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test -v -cover ./...

.PHONY: rundb stopdb createdb dropdb migrateup migratedown sqlc test server