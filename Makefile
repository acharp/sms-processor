ensure:
	dep ensure

build:
	@go build -o bin/api cmd/api/*.go


run: build
	@./bin/api $(ARGS)
		

TESTING_OPTS = -failfast
test:
	go test $(TESTING_OPTS) ./...
