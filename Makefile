lint:
	golangci-lint run

test:
	go test -v -cover -coverprofile=cover.out ./...

coverout:
	go tool cover -html=cover.out

bench:
	go test -v -bench=. -run=^$$ -benchtime=10s -cpuprofile=cpu.out

bench_clock:
	go test -v -bench=./clock -run=^$$ -benchtime=10s -cpuprofile=cpu.out

cpuout:
	go tool pprof -http=: cpu.out

