rundb:
	docker image pull postgres:16.2-alpine3.19
	docker container run \
         --name postgres16 \
         -p 5432:5432 \
         -e POSTGRES_USER=root \
         -e POSTGRES_PASSWORD=secret \
         -d --rm \
        postgres:16.2-alpine3.19

stopdb:
	docker container stop postgres16

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root mom_bank

dropdb:
	docker exec -it postgres16 dropdb mom_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/mom_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/mom_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

.PHONY: rundb stopdb createdb dropdb migrateup migratedown sqlc