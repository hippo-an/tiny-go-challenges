rundb:
	@echo "Run database container...."
	@docker container run \
         --name mom-bank-db \
		 --network mombank-network \
         -p 5432:5432 \
         -e POSTGRES_USER=root \
         -e POSTGRES_PASSWORD=secret \
         -d --rm \
        postgres:16.2-alpine3.19
	@echo "Started Postgres database container!"

stopdb:
	@echo "Stopping database container..."
	@docker container stop mom-bank-db

createdb:
	@docker exec -it mom-bank-db createdb --username=root --owner=root mom_bank
	@echo "database mom_bank successfully started!"

dropdb:
	@docker exec -it mom-bank-db dropdb mom_bank

migrateup:
	@migrate -path db/migration -database "postgresql://root:secret@localhost:5432/mom_bank?sslmode=disable" -verbose up

muone:
	@migrate -path db/migration -database "postgresql://root:secret@localhost:5432/mom_bank?sslmode=disable" -verbose up 1

migratedown:
	@migrate -path db/migration -database "postgresql://root:secret@localhost:5432/mom_bank?sslmode=disable" -verbose down

mdone:
	@migrate -path db/migration -database "postgresql://root:secret@localhost:5432/mom_bank?sslmode=disable" -verbose down 1

sqlc:
	@sqlc generate

server:
	@go run main.go

test:
	@go test -v -cover -race ./...

test-cover:
	@go test -v -coverprofile=coverage.out -race ./...
	@go tool cover -html=coverage.out

mock:
	@mockgen -package mockdb -destination db/mock/store.go github.com/hippo-an/tiny-go-challenges/mombank/db/sqlc Store

container:
	@docker container run --name mombank -p 8080:9999 -e GIN_MODE=debug --rm --network mombank-network -e DB_SOURCE="postgresql://root:secret@mom-bank-db:5432/mom_bank?sslmode=disable" mombank:latest

.PHONY: rundb stopdb createdb dropdb migrateup migratedown sqlc test server mock test-cover mdone muone