package queue

import "sync"

type Queue[T any] interface {
	Enqueue(val T) error
	Dequeue() (T, error)
	Size() int
	Peek() (T, error)
}

type queue[T any] struct {
	mu       sync.RWMutex
	capacity int
	items    []T
}

func New[T any](opts ...Option[T]) Queue[T] {
	return newQueue(opts...)
}

func newQueue[T any](opts ...Option[T]) *queue[T] {
	s := &queue[T]{
		capacity: UnlimitedCapacity,
	}
	for _, opt := range opts {
		opt(s)
	}

	s.items = make([]T, 0)

	return s
}

func (q *queue[T]) Enqueue(val T) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.capacity >= 0 && len(q.items)+1 > q.capacity {
		return ErrOverflow
	}

	q.items = append(q.items, val)

	return nil
}

func (q *queue[T]) Dequeue() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		var zero T
		return zero, ErrUnderflow
	}

	result := q.items[0]
	q.items = q.items[1:]

	return result, nil
}

func (q *queue[T]) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.items)
}

func (q *queue[T]) Peek() (T, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	sz := len(q.items)
	if sz == 0 {
		var zero T
		return zero, ErrUnderflow
	}

	return q.items[0], nil
}
