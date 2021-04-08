format:
	go fmt ./...

lint:
	golangci-lint run ./...

test:
	go clean -testcache && go test -race $$(go list ./... | grep -v tests)

test-v:
	go clean -testcache && go test -race $$(go list ./... | grep -v tests) -v

all:
	$(MAKE) format
	$(MAKE) lint
	$(MAKE) test