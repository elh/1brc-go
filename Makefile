.PHONY: run
run: build
	@export PROFILE=true; bash -c 'diff measurements.out <(time ./bin/1brc-go)'

# Approach taken from 1brc evaluation.sh. Requires hyperfine to be installed.
.PHONY: evaluate
evaluate: build
	@hyperfine --warmup 0 --runs 5 --export-json timing.json "./bin/1brc-go 2>&1"

build:
	@go build -o bin/1brc-go main.go

.PHONY: pprof
pprof:
	@go tool pprof -http :8080 $(f)
