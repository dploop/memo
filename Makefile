test:
	go test -v -cover -coverprofile=coverage.out ./...

cover:
	go tool cover -html=coverage.out

bench:
	go test -v  -bench . -run ^$$ ./...

lint:
	golangci-lint run
