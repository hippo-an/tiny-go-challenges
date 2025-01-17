APP_NAME = 'ghosth'

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: go-install-air
go-install-air:  ## Install air using 'go install'
	go install github.com/air-verse/air@latest

.PHONY: go-install-templ
go-install-templ:
	go install github.com/a-h/templ/cmd/templ@latest

.PHONY: get-install-tailwindcss
get-install-tailwindcss:
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
	chmod +x tailwindcss-macos-arm64
	mv tailwindcss-macos-arm64 twc

.PHONY: tailwind-watch
tailwind-watch:
	./twc -i ./static/css/custom.css -o ./static/css/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	./twc -i ./static/css/custom.css -o ./static/css/style.css

.PHONY: build
build: tailwind-build templ-generate
	go build -o ./tmp/$(APP_NAME) ./cmd/main.go

.PHONY: watch
watch:
	air -c .air.toml

.PHONY: templ-generate
templ-generate:
	templ generate ./templates

.PHONY: templ-watch
templ-watch:
	@templ generate -watch -proxy=http://localhost:9000
