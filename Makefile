# Env var controls: NUM_PARSERS, PARSE_CHUNK_SIZE_MB, PROFILE
.PHONY: run
run: build
	@echo "NUM_PARSERS: $(if $(NUM_PARSERS),$(NUM_PARSERS),<unset>)"
	@echo "PARSE_CHUNK_SIZE_MB: $(if $(PARSE_CHUNK_SIZE_MB),$(PARSE_CHUNK_SIZE_MB),<unset>)"
	@echo "PROFILE: $(if $(PROFILE),$(PROFILE),<unset>)"
	@bash -c 'diff measurements.out <(time ./bin/1brc-go)'

# Approach taken from 1brc evaluation.sh. Requires hyperfine to be installed.
.PHONY: evaluate
evaluate: build
	@echo "NUM_PARSERS: $(if $(NUM_PARSERS),$(NUM_PARSERS),<unset>)"
	@echo "PARSE_CHUNK_SIZE_MB: $(if $(PARSE_CHUNK_SIZE_MB),$(PARSE_CHUNK_SIZE_MB),<unset>)"
	@echo "PROFILE: $(if $(PROFILE),$(PROFILE),<unset>)"
	@hyperfine --warmup 0 --runs 5 --export-json timing.json "./bin/1brc-go 2>&1"

build:
	@echo "build..."
	@go build -pgo=auto -o bin/1brc-go main.go

.PHONY: pprof
pprof:
	@go tool pprof -http :8080 $(f)
