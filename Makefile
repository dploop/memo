test:
	go test -v -cover -coverprofile=cover.out ./...

coverout:
	go tool cover -html=cover.out

bench:
	go test -v -bench=. -run=^$$ -benchtime=10s -cpuprofile=cpu.out -benchmem -memprofile=mem.out .

bench_clock:
	go test -v -bench=. -run=^$$ -benchtime=10s -cpuprofile=cpu.out -benchmem -memprofile=mem.out ./clock

cpuout:
	go tool pprof -http=: cpu.out

memout:
	go tool pprof -http=: mem.out

lint:
	golangci-config-generator
	golangci-lint run

install-gcg:
	go install github.com/dploop/golangci-config-generator/cmd/golangci-config-generator@latest

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.40.1
