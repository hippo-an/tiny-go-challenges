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
	migrate -path db/migration -database "postgresql://root:secret@localhost:15432/booking?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:15432/booking?sslmode=disable" -verbose down

run:
	go run ./cmd/web

coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...

coverage-report:
	@echo "Generating coverage report..."
	go tool cover -html=coverage.out
	@echo "Coverage report generated: coverage.html"

.PHONY: rundb stopdb createdb dropdb migrateup migratedown run coverage coverage-report
