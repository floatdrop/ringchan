package ringchan

type Ring[T any] struct {
	// C is the channel to receive values from.
	C <-chan T
}

// New creates a ring-buffered channel with fixed capacity from incoming channel.
// When full, new inserts will drop the oldest items to make space.
func New[T any](in <-chan T, size int) *Ring[T] {
	out := make(chan T, size)

	rc := &Ring[T]{
		C: out,
	}

	go func() {
		defer close(out)

		for v := range in {
			select {
			case out <- v:
			default:
				// Do non-blocking receive to drop the oldest item
				// if the buffer is full. This avoids blocking in case of an empty buffer.
				select {
				case <-out:
					out <- v
				default:
					out <- v
				}
			}
		}
	}()

	return rc
}
