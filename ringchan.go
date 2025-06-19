package ringchan

type Ring[T any] struct {
	in  chan T
	out chan T

	// In is the channel to send values to.
	In chan<- T
	// Out is the channel to receive values from.
	Out <-chan T
}

// New creates a ring-buffered channel with fixed capacity.
// When full, new inserts will drop the oldest items to make space.
func New[T any](capacity int) *Ring[T] {
	in := make(chan T, capacity)
	out := make(chan T, capacity)

	rc := &Ring[T]{
		in:  in,
		out: out,
		In:  in,
		Out: out,
	}

	go rc.run()
	return rc
}

func (rc *Ring[T]) run() {
	defer close(rc.out)

	for v := range rc.in {
		select {
		case rc.out <- v:
		default:
			// Do non-blocking receive to drop the oldest item
			// if the buffer is full. This avoids blocking in case of an empty buffer.
			select {
			case <-rc.out:
				rc.out <- v
			default:
				rc.out <- v
			}
		}
	}
}

// Len returns the number of items currently in the ring buffer.
func (rc *Ring[T]) Len() int {
	return len(rc.out)
}

// Close closes the input channel, ending the ring.
func (rc *Ring[T]) Close() {
	close(rc.in)
}
