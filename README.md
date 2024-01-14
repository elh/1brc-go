# 1brc-go

An excuse to play around with Go Profiling. Runs in ~7s on my 2023 Apple M2 Pro Macbook with 10 cores and 16GB RAM.

```bash
make run
# builds the go executable and runs the following underlying test

diff data/measurements.out <(time MEASUREMENTS_PATH=data/measurements.txt ./bin/main)
# MEASUREMENTS_PATH=data/measurements.txt ./bin/main  38.93s user 4.71s system 628% cpu 6.949 total

# MEASUREMENTS_PATH is a required file path
# PROFILE enables profiling if set to "true"
go run main.go
```

Add your measurement files and expected results to `data/`.

## Log
Measured locally on my machine.

* 140s - Naive implementation. Reading lines with `Scanner.Text`.
* 83s - Use `Scanner.ReadSlice` to reduce a massive string allocation.
* 54s - Use `file.Read` directly to iterate bytes in a single scan. Avoid `bytes` helper functions.
* 45s - Use a larger buffer. Remove a redundant map lookup. Use `unsafe.String` to avoid string allocations.
* 35s - Implement a custom, simplified replacement for `strconv.ParseFloat`.
    * Most time just goes into scanning byte slices. Biggest non-scanning bottleneck is `map[string]*Stats` lookups. Decently pleased with single threaded optimizations. Now, to make concurrent.
* 7s - Concurrently read and parse the file. See `numParsers` and `parseChunkSize` parameters.

### Ideas
* Optimize concurrency more scientifically
* Custom map/set implementation
* Memory arena (or other optimizations of locality)

## Reading
* https://github.com/DataDog/go-profiler-notes/blob/main/guide/README.md#cpu-profiler ⭐️
* https://pkg.go.dev/bufio
* https://pkg.go.dev/unsafe#String
* https://go.dev/src/arena/arena.go
