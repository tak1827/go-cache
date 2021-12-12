lint:
	go vet ./...

fmt:
	gofmt -w -l .

test:
	go test ./... -race
