# Setup
setup:
	sudo chmod +x ./bin/tests.sh
	go mod tidy
.PHONY: setup

# Run gofmt
format:
	go fmt ./...
.PHONY: format

# Run linter
lint:
	golangci-lint run ./...
.PHONY: lint

# Test uses race and coverage
test:
	go clean -testcache && go test -race $$(go list ./... | grep -v tests | grep -v examples | grep -v res | grep -v mocks) -coverprofile=coverage.out -covermode=atomic
.PHONY: test

# Test with -v
test-v:
	go clean -testcache && go test -race -v $$(go list ./... | grep -v tests | grep -v examples | grep -v res | grep -v mocks) -coverprofile=coverage.out -covermode=atomic
.PHONY: test-v

# Runs real world tests for a driver or all drivers.
# See ./bin/tests.sh for example usage.
test-driver:
	go clean -testcache && ./bin/tests.sh $(driver)
.PHONY: test-driver

# Run all the tests and opens the coverage report
cover: test
	go tool cover -html=coverage.out
.PHONY: cover

# Make mocks keeping directory tree
mocks:
	rm -rf internal/mocks && mockery --all --keeptree --output ./internal/mocks && mv ./internal/mocks/internal/* ./internal/mocks
.PHONY: mocks

doc:
	godoc -http localhost:8080
.PHONY: doc

# Make format, lint and test
all:
	$(MAKE) format
	$(MAKE) lint
	$(MAKE) test
