APP_NAME = 'ghosth'

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: go-install-air
go-install-air:  ## Install air using 'go install'
	go install github.com/air-verse/air@latest

.PHONY: go-install-templ
go-install-templ: ## Install templ using 'go install'
	go install github.com/a-h/templ/cmd/templ@latest

.PHONY: install-tailwindcss
install-tailwindcss: ## Install tailwindcss binary
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
	chmod +x tailwindcss-macos-arm64
	mv tailwindcss-macos-arm64 twc

.PHONY: tailwind-build
tailwind-build: ## Build tailwind css with minify option
	./twc -i ./static/css/custom.css -o ./static/css/style.css --minify

.PHONY: build
build: tailwind-build templ-generate  ## Build tailwind, templ, and server
	go build -o ./tmp/$(APP_NAME) ./cmd/main.go


.PHONY: templ-generate
templ-generate:  ## Generate go code based on templ file
	templ generate -path ./templates

##############################################################################

live/tailwind:
	./twc -i ./static/css/custom.css -o ./static/css/style.css --minify --watch

live/server:
	air -c .air.toml

live/templ:
	templ generate -path ./templates -watch --proxy=http://localhost:9000 -v

live/static:
	air --build.cmd "templ generate --notify-proxy" \
		--build.bin "true" \
		--build.delay "100" \
		--build.exclude_dir "" \
		--build.include_dir "static" \
		--build.include_ext "js,css"

.PHONY: live
live:  ## Live reloading with server, tailwind, templ, static files
	@make -j4 live/server live/tailwind  live/templ live/static