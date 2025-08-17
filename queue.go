// Package queue provides a thread-safe, generic queue implementation with configurable capacity limits.
//
// The queue follows FIFO (First-In-First-Out) semantics and supports any type through Go generics.
// All operations are safe for concurrent use across multiple goroutines.
//
// Example usage:
//
//	q := queue.New[int]()
//	q.Enqueue(1)
//	q.Enqueue(2)
//	val, err := q.Dequeue() // returns 1, nil
package queue

import "sync"

// Queue defines the interface for a generic queue data structure.
// All operations are thread-safe and support any type T.
type Queue[T any] interface {
	// Enqueue adds an item to the back of the queue.
	// Returns ErrOverflow if the queue is at capacity.
	Enqueue(val T) error

	// Dequeue removes and returns the front item from the queue.
	// Returns ErrUnderflow if the queue is empty.
	Dequeue() (T, error)

	// Size returns the current number of items in the queue.
	Size() int

	// Peek returns the front item without removing it from the queue.
	// Returns ErrUnderflow if the queue is empty.
	Peek() (T, error)
}

// New creates a new queue with the specified options.
// If no options are provided, creates an unlimited capacity queue.
//
// Example:
//
//	q := queue.New[int]()                           // Unlimited capacity
//	q := queue.New[int](queue.WithCapacity[int](10)) // Capacity of 10
func New[T any](opts ...Option[T]) Queue[T] {
	return newQueue(opts...)
}

type queue[T any] struct {
	mu       sync.RWMutex
	capacity int
	items    []T
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
