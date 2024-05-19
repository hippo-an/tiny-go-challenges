run: build
	@./bin/app

build:
	@go build -o ./bin/app ./cmd/app/*.go

templgenerate:
	@templ generate

devrun:
	@air

templproxy:
	@templ generate -watch -proxy=http://localhost:8080

twc:
	@sh < ./install_tailwindcss.sh

css:
	@./twc -i ./views/css/app.css -o ./public/styles.css --watch