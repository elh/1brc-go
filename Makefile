.PHONY: run
run: build
	@export MEASUREMENTS_PATH=data/measurements.txt; export PROFILE=true; bash -c 'diff data/measurements.out <(time ./bin/main)'

.PHONY: run-small
run-small: build
	@export MEASUREMENTS_PATH=data/measurements_1m.txt; export PROFILE=true; bash -c 'diff data/measurements_1m.out <(time ./bin/main)'

build:
	@go build -o bin main.go

.PHONY: pprof
pprof:
	@go tool pprof -http :8080 $(f)
