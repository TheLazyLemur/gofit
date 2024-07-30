tidy:
	@go mod tidy

gen:
	@sqlc generate
	@templ generate ./...

build: gen tidy
	@go build -v -o ./bin/gofit ./src/cmd/gofit/...

run: gen tidy
	@go run -v ./src/cmd/gofit/...
