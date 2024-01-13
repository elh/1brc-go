package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sort"
	"strings"
	"time"
	"unsafe"
)

var (
	measurementsPath = os.Getenv("MEASUREMENTS_PATH")
	shouldProfile    = os.Getenv("PROFILE") == "true"

	profileTypes = []string{"goroutine", "allocs"} // "heap", "threadcreate", "block", "mutex"
)

type Stats struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

// rounding floats to 1 decimal place with 0.05 rounding up to 0.1
func round(x float64) float64 {
	return math.Floor((x+0.05)*10) / 10
}

// parseFloatFast is a high performance float parser using the assumption that
// the byte slice will always have a single decimal digit.
func parseFloatFast(bs []byte) float64 {
	var intStartIdx int // is negative?
	if bs[0] == '-' {
		intStartIdx = 1
	}

	v := float64(bs[len(bs)-1]-'0') / 10 // single decimal digit

	place := 1.0
	for i := len(bs) - 3; i >= intStartIdx; i-- { // integer part
		v += float64(bs[i]-'0') * place
		place *= 10
	}

	if intStartIdx == 1 {
		v *= -1
	}
	return v
}

func main() {
	if shouldProfile {
		nowUnix := time.Now().Unix()
		os.MkdirAll(fmt.Sprintf("profiles/%d", nowUnix), 0755)
		for _, profileType := range profileTypes {
			file, _ := os.Create(fmt.Sprintf("profiles/%d/%s.%s.pprof", nowUnix, filepath.Base(measurementsPath), profileType))
			defer file.Close()
			defer pprof.Lookup(profileType).WriteTo(file, 0)
		}

		file, _ := os.Create(fmt.Sprintf("profiles/%d/%s.cpu.pprof", nowUnix, filepath.Base(measurementsPath)))
		defer file.Close()
		pprof.StartCPUProfile(file)
		defer pprof.StopCPUProfile()
	}

	f, err := os.Open(measurementsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// result data
	names := make([]string, 0, 10000)
	stats := make(map[string]*Stats, 10000)

	bs := make([]byte, 1024*1024*1024) // file byte buffer. NOTE: sizing?
	remainderBs := make([]byte, 100)   // remainder unparsed from the last buffer
	lastName := make([]byte, 100)      // last name parsed
	var remainderLen, lastNameLen int

	isScanningName := true // currently scanning name or value?
	for {
		// load the buffer
		n, err := f.Read(bs)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// tick tock between parsing names and values; accummulating stats and
		// keeping track of unparsed remainders
		var idx, start int
		for {
			if isScanningName {
				for idx < n {
					if bs[idx] == ';' {
						nameBs := bs[start:idx]
						// TODO: handle remainder outside of this tight loop?
						if remainderLen > 0 {
							nameBs = append(remainderBs[:remainderLen], nameBs...)
							remainderLen = 0
						}
						lastNameLen = copy(lastName, nameBs)

						idx++
						start = idx
						isScanningName = false
						break
					}
					idx++
				}
			} else {
				for idx < n {
					if bs[idx] == '\n' {
						valueBs := bs[start:idx]
						// TODO: handle remainder outside of this tight loop?
						if remainderLen > 0 {
							valueBs = append(remainderBs[:remainderLen], valueBs...)
							remainderLen = 0
						}
						value := parseFloatFast(valueBs)

						nameUnsafe := unsafe.String(&lastName[0], lastNameLen)
						if s, ok := stats[nameUnsafe]; !ok {
							name := string(lastName[:lastNameLen]) // actually allocate string
							stats[name] = &Stats{Min: value, Max: value, Sum: value, Count: 1}
							names = append(names, name)
						} else {
							if value < s.Min {
								s.Min = value
							}
							if value > s.Max {
								s.Max = value
							}
							s.Sum += value
							s.Count++
						}

						idx++
						start = idx
						isScanningName = true
						break
					}
					idx++
				}
			}
			if idx >= n {
				break
			}
		}

		if start < n {
			remainderLen = copy(remainderBs, bs[start:])
		}
	}

	// sorted alphabetically for output
	sort.Strings(names)

	// build up result
	var builder strings.Builder
	for i, name := range names {
		s := stats[name]
		avg := round(s.Sum / float64(s.Count))
		builder.WriteString(fmt.Sprintf("%s=%.1f/%.1f/%.1f", name, s.Min, avg, s.Max))
		if i < len(names)-1 {
			builder.WriteString(", ")
		}
	}

	// print result to stdout
	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(writer, "{%s}\n", builder.String())
	writer.Flush()
}
