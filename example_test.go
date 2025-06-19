package ringchan_test

import (
	"fmt"
	"time"

	"github.com/floatdrop/ringchan"
)

func ExampleNew() {
	rc := ringchan.New[string](3)

	go func() {
		inputs := []string{"A", "B", "C", "D", "E"}
		for _, v := range inputs {
			rc.In <- v
		}
		rc.Close()
	}()

	time.Sleep(50 * time.Millisecond)

	for v := range rc.Out {
		fmt.Println("Got:", v)
	}

	// Output:
	// Got: C
	// Got: D
	// Got: E
}
