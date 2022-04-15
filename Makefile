EXECUTABLE_NAME := metrics

.PHONY: build
build: tidy
	@CGO_ENABLED=0 go build -o bin/$(EXECUTABLE_NAME) ./cmd/server

.PHONY: run
run: build
	@bin/$(EXECUTABLE_NAME)

.PHONY: tidy
tidy:
	@go fmt ./...
	@go mod tidy -v

.PHONY: test
test:
	@go test -v -count=1 -coverprofile=coverage.out ./...

.PHONY: coverage-func
coverage-func: test
	@go tool cover -func=coverage.out

.PHONY: coverage-html
coverage-html: test
	@go tool cover -html=coverage.out

.PHONY: clean
clean:
	@rm -rf coverage.out bin/
	@go clean