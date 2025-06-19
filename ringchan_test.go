package ringchan

import (
	"testing"
	"time"
)

func TestRingChanBasic(t *testing.T) {
	rc := New[int](3)

	go func() {
		for i := 1; i <= 5; i++ {
			rc.In <- i
		}
		rc.Close()
	}()

	time.Sleep(50 * time.Millisecond)

	l := rc.Len()
	if l != 3 {
		t.Fatalf("expected Len()=%v, got %v", 3, l)
	}

	var got []int
	for v := range rc.Out {
		got = append(got, v)
	}

	// Only last 3 values should be kept due to overwrite
	want := []int{3, 4, 5}
	if len(got) != len(want) {
		t.Fatalf("expected %v values, got %v", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("expected %v at index %d, got %v", want[i], i, got[i])
		}
	}
}

func TestRingChanBlockingReceive(t *testing.T) {
	rc := New[int](1)

	go func() {
		time.Sleep(100 * time.Millisecond)
		rc.In <- 42
		rc.Close()
	}()

	val := <-rc.Out
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}
}

func TestRingChanRangeAfterClose(t *testing.T) {
	rc := New[string](2)

	rc.In <- "foo"
	rc.In <- "bar"
	rc.Close()

	var results []string
	for v := range rc.Out {
		results = append(results, v)
	}

	if len(results) != 2 || results[0] != "foo" || results[1] != "bar" {
		t.Errorf("unexpected results: %v", results)
	}
}
