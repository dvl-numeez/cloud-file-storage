build:
	@go build -o bin/cloudfilestorage

run:build
	@./bin/cloudfilestorage

test:
	@go test -v ./...
