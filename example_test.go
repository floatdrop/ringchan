package ringchan_test

import (
	"fmt"
	"time"

	"github.com/floatdrop/ringchan"
)

func ExampleNew() {
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
