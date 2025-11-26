.PHONY: build build-server build-client docker-build docker-build-server docker-build-client \
        deploy undeploy logs-server logs-client clean test

# Variables
IMAGE_REGISTRY ?=
SERVER_IMAGE ?= smoker-server
CLIENT_IMAGE ?= smoker-client
TAG ?= latest

# Local build
build: build-server build-client

build-server:
	go build -o bin/server ./cmd/server

build-client:
	go build -o bin/client ./cmd/client

# Docker build
docker-build: docker-build-server docker-build-client

docker-build-server:
	docker build --target server -t $(SERVER_IMAGE):$(TAG) .

docker-build-client:
	docker build --target client -t $(CLIENT_IMAGE):$(TAG) .

# Push images (optional, for remote registry)
docker-push: docker-push-server docker-push-client

docker-push-server:
	docker push $(IMAGE_REGISTRY)/$(SERVER_IMAGE):$(TAG)

docker-push-client:
	docker push $(IMAGE_REGISTRY)/$(CLIENT_IMAGE):$(TAG)

# Kubernetes deployment
deploy:
	kubectl apply -f deployments/

undeploy:
	kubectl delete -f deployments/ --ignore-not-found

# Logs
logs-server:
	kubectl logs -l app=smoker-server -f

logs-client:
	kubectl logs -l app=smoker-client -f --max-log-requests=10

# Testing
test:
	go test -v ./...

# Local run
run-server:
	go run ./cmd/server

run-client:
	SERVER_HOST=localhost go run ./cmd/client

# Clean
clean:
	rm -rf bin/
	docker rmi $(SERVER_IMAGE):$(TAG) $(CLIENT_IMAGE):$(TAG) 2>/dev/null || true

# Help
help:
	@echo "Available targets:"
	@echo "  build            - Build server and client binaries locally"
	@echo "  docker-build     - Build Docker images for server and client"
	@echo "  deploy           - Deploy to Kubernetes"
	@echo "  undeploy         - Remove from Kubernetes"
	@echo "  logs-server      - Stream server logs"
	@echo "  logs-client      - Stream client logs"
	@echo "  test             - Run tests"
	@echo "  run-server       - Run server locally"
	@echo "  run-client       - Run client locally (connects to localhost)"
	@echo "  clean            - Remove binaries and Docker images"
