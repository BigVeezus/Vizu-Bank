build:
	@go build -o bin/vizubank

nodemon:
	nodemon --watch './**/*.go' --signal SIGTERM --exec go build -o bin/vizubank

run: build
	@./bin/vizubank

test:
	@go test -v ./...