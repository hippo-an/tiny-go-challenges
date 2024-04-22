rundb:
	docker image pull postgres:16.2-alpine3.19
	docker container run \
         --name booking-back \
         -p 15432:5432 \
         -e POSTGRES_USER=root \
         -e POSTGRES_PASSWORD=secret \
         -d --rm \
        postgres:16.2-alpine3.19

stopdb:
	docker container stop booking-back

createdb:
	docker exec -it booking-back createdb --username=root --owner=root booking

dropdb:
	docker exec -it booking-back dropdb booking

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:15432/booking?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:15432/booking?sslmode=disable" -verbose down 1


.PHONY: rundb stopdb createdb dropdb migrateup migratedown
