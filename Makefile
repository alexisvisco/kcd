lint:
	golangci-lint run

cov:
	go test -v -coverpkg=./...  -covermode=count -coverprofile=coverage.out ./...
