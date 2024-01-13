# 1brc-go

An excuse to play around with Go Profiling.

```bash
make run
# builds the go executable and runs the following underlying test

diff data/measurements_1b.out <(time MEASUREMENTS_PATH=data/measurements_1b.txt ./bin/main)
# MEASUREMENTS_PATH=data/measurements_1b.txt PROFILE=true ./bin/main  30.25s user 1.62s system 91% cpu 35.002 total
```

Add your measurement files and expected results to `data/`.

## Log
Measured locally on my 2023 Apple M2 Pro Macbook.

* 140s - Naive implementation. Reading lines with `Scanner.Text`.
* 83s - Use `Scanner.ReadSlice` to reduce a massive string allocation.
* 54s - Use `file.Read` directly to iterate bytes in a single scan. Avoid `bytes` helper functions.
* 45s - Use a larger buffer. Remove a redundant map lookup.
* 35s - Implement a custom, simplified replacement for `strconv.ParseFloat`.

Biggest bottleneck now is the `map[string]*Stats` lookups. Decently pleased with single threaded optimizations. Now, to make concurrent.

### Ideas
* Concurrency
* Custom map/set implementation
* Memory arena (or other optimizations of locality)

## Reading
* https://github.com/DataDog/go-profiler-notes/blob/main/guide/README.md#cpu-profiler ⭐️
* https://pkg.go.dev/bufio
* https://pkg.go.dev/unsafe#String
* https://go.dev/src/arena/arena.go
