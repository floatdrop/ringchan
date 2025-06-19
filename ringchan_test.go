package ringchan

import (
	"testing"
	"time"
)

func TestRingChanBasic(t *testing.T) {
	input := make(chan int, 5)
	rc := New(input, 3)

	go func() {
		for i := 1; i <= 5; i++ {
			input <- i
		}
		close(input)
	}()

	time.Sleep(50 * time.Millisecond)

	l := len(rc.C)
	if l != 3 {
		t.Fatalf("expected Len()=%v, got %v", 3, l)
	}

	var got []int
	for v := range rc.C {
		got = append(got, v)
	}

	// Only last 3 values should be kept due to overwrite
	want := []int{3, 4, 5}
	if len(got) != len(want) {
		t.Fatalf("expected %v values, got %v", len(want), len(got))
	}
	if rc.Dropped != 2 {
		t.Fatalf("expected %d values to be dropped, got %d", 2, rc.Dropped)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("expected %v at index %d, got %v", want[i], i, got[i])
		}
	}
}

func TestRingChanBlockingReceive(t *testing.T) {
	input := make(chan int, 1)
	rc := New(input, 1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		input <- 42
		close(input)
	}()

	val := <-rc.C
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestRingChanRangeAfterClose(t *testing.T) {
	input := make(chan string, 2)
	rc := New(input, 2)

	input <- "foo"
	input <- "bar"
	close(input)

	var results []string
	for v := range rc.C {
		results = append(results, v)
	}

	if len(results) != 2 || results[0] != "foo" || results[1] != "bar" {
		t.Errorf("unexpected results: %v", results)
	}
}
