# 1brc-go

Parse and aggregate 1B rows of text very quickly using Go for the [1BRC challenge](https://github.com/gunnarmorling/1brc).

Mainly done as an excuse to play around with `pprof`.

<sub>_NOTE: Originally developed at [elh/1brc-go](https://github.com/elh/1brc-go) and copied into gunnarmorling/1brc._</sub>

<br>

## Performance

On Jan 14, as measured on my 2023 Apple M2 Pro Macbook with 1brc's `evaluate.sh` script, this performs very competitively against the top Java and Go submissions. At the very least this means, I have put in some effort optimizing this for the machine I have in front of me. There is a pretty funny lesson to be learned here that everyone's solution is the fastest on their own machine :)

Props to AlexanderYastrebov who is pushing on enabling language agnostic solutions in the leaderboard. Very curious now to see where this stacks up if I can run it on the test environment. I haven't looked at other Go solution code yet, but the discussion his started would motivate me to write my own.

| # | Result (m:s.ms) | Implementation     | JDK | Submitter     | Notes     |
|---|-----------------|--------------------|-----|---------------|-----------|
| 1 | 00:06.876 | [link](https://github.com/elh/1brc-go)| 🔷 Go 1.21.5 | [Eugene Huang](https://github.com/elh) | Mine! 👈 |
| 2 | 00:07.602 | [link](https://gist.github.com/corlinp/176a97c58099bca36bcd5679e68f9708)| 🔷 Go 1.21.5 | [Corlin Palmer](https://github.com/corlinp) | Go implementation |
| 3 | 00:13.765 | [link](https://github.com/gunnarmorling/1brc/blob/main/src/main/java/dev/morling/onebrc/CalculateAverage_royvanrijn.java)| 21.0.1-graal | [Roy van Rijn](https://github.com/royvanrijn) | GraalVM native binary |
|   | 00:13.989 | [link](https://github.com/gunnarmorling/1brc/blob/main/src/main/java/dev/morling/onebrc/CalculateAverage_artsiomkorzun.java)| 21.0.1-graal | [Artsiom Korzun](https://github.com/artsiomkorzun) |  |
|   | 00:14.044 | [link](https://github.com/gunnarmorling/1brc/blob/main/src/main/java/dev/morling/onebrc/CalculateAverage_thomaswue.java)| 21.0.1-graal | [Thomas Wuerthinger](https://github.com/thomaswue) | GraalVM native binary |
|   | 00:14.464 | [link](https://github.com/AlexanderYastrebov/1brc/tree/go-implementation/src/main/go)| 🔷 Go 1.21.5 | [Alexander Yastrebov](https://github.com/AlexanderYastrebov) | Go implementation |
|   | 00:14.839 | [link](https://github.com/gunnarmorling/1brc/blob/main/src/main/java/dev/morling/onebrc/CalculateAverage_mtopolnik.java)| 21.0.1-graal | [Marko Topolnik](https://github.com/mtopolnik) |  |
|   | 00:14.949 | [link](https://github.com/gunnarmorling/1brc/blob/main/src/main/java/dev/morling/onebrc/CalculateAverage_merykittyunsafe.java)| 21.0.1-open | [merykittyunsafe](https://github.com/merykittyunsafe) |  |
|   | 00:16.075 | [link](https://github.com/jkroepke/1brc-go/tree/main/go)| 🔷 Go 1.21.5 | [Jan-Otto Kröpke](https://github.com/jkroepke) | Go implementation |

```bash
make evaluate
# Benchmark 1: ./bin/1brc-go 2>&1
#   Time (mean ± σ):      6.626 s ±  0.170 s    [User: 38.109 s, System: 3.072 s]
#   Range (min … max):    6.331 s …  6.749 s    5 runs
```

<br>

## Usage

```bash
# builds the go executable and runs the following test
make run

# Result:
# ./bin/1brc-go  37.90s user 3.98s system 590% cpu 7.089 total
diff measurements.out <(time ./bin/1brc-go)

# takes an optional single command line arg for the measurements file path. If not provided, defaults to `measurements.txt`
# PROFILE is optional and enables profiling if "true"
go run main.go
```

Add a gitignore-d `measurements.txt` file in the project root or use command line arg to override the path.

<br>

## Log
Times measured on 2023 Apple M2 Pro Macbook with 10 cores and 16GB RAM.

* 140s - Naive implementation. Reading lines with `Scanner.Text`.
* 83s - Using `Scanner.ReadSlice` to reduce a massive string allocation.
* 54s - Using `file.Read` directly and avoided some `bytes` helper functions to iterate bytes in a single scan.
* 45s - Using a far larger buffer to reduce IO latency. Removed a redundant map lookup. Using `unsafe.String` to avoid additional string allocs.
* 35s - Implemented a custom, simplified replacement for `strconv.ParseFloat`.
    * Most of the time just goes into scanning byte slices. Biggest other bottleneck is `map[string]*Stats` lookups. Decently pleased with single threaded optimizations so now making it concurrent.
* 7s - Concurrently read and parse the file as byte index addressed chunks, then merge results in a single final goroutine. See tunable `numParsers` and `parseChunkSize` parameters.

### Ideas
* Optimize concurrency.
    * Pipeline more work.
    * Tune concurrency and read buffer size parameters. I picked current values for my machine very quickly and unscientifically.
* Custom map/set implementation for string->float
* Memory arena (or other optimizations of locality)?
* Faster way to iterate bytes?

<br>

## Reading
* https://github.com/DataDog/go-profiler-notes/blob/main/guide/README.md#cpu-profiler ⭐️
* https://pkg.go.dev/bufio
* https://pkg.go.dev/unsafe#String
* https://go.dev/src/arena/arena.go

<br>

---
Copyright © 2024 Eugene Huang
