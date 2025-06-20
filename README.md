# RingChan

[![CI](https://github.com/floatdrop/ringchan/actions/workflows/ci.yaml/badge.svg)](https://github.com/floatdrop/ringchan/actions/workflows/ci.yaml)
![Coverage](https://img.shields.io/badge/Coverage-90.9%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/github.com/floatdrop/ringchan)](https://goreportcard.com/report/github.com/floatdrop/ringchan)
[![Go Reference](https://pkg.go.dev/badge/github.com/floatdrop/ringchan.svg)](https://pkg.go.dev/github.com/floatdrop/ringchan)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**RingChan** is a thread-safe, fixed-capacity ring buffer implemented as a channel in Go. It mimics Go's channel behavior while providing ring-buffer semantics — meaning **new items overwrite the oldest when full**.

## Features

- Fixed-size buffer with overwrite behavior
- Range-friendly: can be iterated using `for ... range`
- Safe for concurrent producers and consumers

## Installation

```bash
go get github.com/floatdrop/ringchan
```

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/floatdrop/ringchan"
)

func main() {
	input := make(chan string, 5)
	ring := ringchan.New(input, 3)

	go func() {
		inputs := []string{"A", "B", "C", "D", "E"}
		for _, v := range inputs {
			input <- v
		}
		close(input)
	}()

	time.Sleep(50 * time.Millisecond)

	for v := range ring.C {
		fmt.Println("Got:", v)
	}

	// Output:
	// Got: C
	// Got: D
	// Got: E
}
```

## Benchmarks

```bash
go test -bench=. -benchmem
```

```
goos: darwin
goarch: arm64
pkg: github.com/floatdrop/ringchan
cpu: Apple M1 Pro
BenchmarkSingleSender-10       	 7097070	       167.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkParallelSenders-10    	 4145682	       295.0 ns/op	       0 B/op	       0 allocs/op
PASS
coverage: 90.9% of statements
ok  	github.com/floatdrop/ringchan	3.050s
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
