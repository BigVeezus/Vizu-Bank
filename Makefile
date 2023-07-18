build:
	@go build -o bin/vizubank

run: build
	@./bin/vizubank

test:
	@go test -v ./...