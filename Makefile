format:
	go fmt ./...

lint:
	golangci-lint run ./

test:
	go clean -testcache && go test -race $$(go list ./... | grep -v tests)

test-v:
	go clean -testcache && go test -race $$(go list ./... | grep -v tests) -v

test-cover:
	go clean -testcache && go test -v -cover -race $$(go list ./... | grep -v tests)

all:
	$(MAKE) format
	$(MAKE) lint
	$(MAKE) test

travis:
	$(MAKE) format
	$(MAKE) lint
	$(MAKE) test-cover