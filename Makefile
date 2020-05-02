test:
	go test -v -cover -coverprofile=cover.out ./...

coverout:
	go tool cover -html=cover.out

bench:
	go test -v -bench=. -run=^$$ -benchtime=10s -cpuprofile=cpu.out .

bench_clock:
	go test -v -bench=. -run=^$$ -benchtime=10s -cpuprofile=cpu.out ./clock

cpuout:
	go tool pprof -http=: cpu.out

lint:
	golangci-config-generator
	golangci-lint run

install-gcg:
	go install github.com/dploop/golangci-config-generator/cmd/golangci-config-generator@latest

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0
