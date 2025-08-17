package queue

type Option[T any] func(*queue[T])

const (
	UnlimitedCapacity = -1
)

func WithCapacity[T any](cap int) Option[T] {
	return func(s *queue[T]) {
		if cap < UnlimitedCapacity {
			panic("cannot specify arbitrary negative capacity")
		}
		s.capacity = cap
	}
}
