.PHONY: all clean fmt lint generate test

all: clean fmt lint test
publish: lint clean

clean:
	go clean -i .
	rm coverage.* || true
	echo ''>go.sum
	go mod tidy

fmt:
	gofmt -w .
	goimports -w .

lint: clean
	go generate ./...
	#go vet ./...
	#golint ./...
	golangci-lint run -c golangci.yml
	gocyclo .

generate:
	go generate ./...

test: clean generate
	go test -gcflags=all=-l -coverpkg=./... -coverprofile=coverage.data ./... -run=^TestUnit.*$
	go tool cover -html=coverage.data -o coverage.html