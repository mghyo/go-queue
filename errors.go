package queue

import "errors"

var (
	// ErrOverflow is returned when attempting to enqueue an item to a queue
	// that has reached its maximum capacity.
	//
	// This error occurs when:
	//   - The queue was created with WithCapacity option
	//   - The current size equals the specified capacity
	//   - Enqueue() is called on the full queue
	//
	// Example:
	//
	//	q := queue.New[int](queue.WithCapacity[int](2))
	//	q.Enqueue(1) // OK
	//	q.Enqueue(2) // OK
	//	err := q.Enqueue(3) // Returns ErrOverflow
	//	if errors.Is(err, queue.ErrOverflow) {
	//		fmt.Println("Queue is full")
	//	}
	ErrOverflow = errors.New("queue overflow")

	// ErrUnderflow is returned when attempting to dequeue or peek at an empty queue.
	//
	// This error occurs when:
	//   - Dequeue() is called on an empty queue
	//   - Peek() is called on an empty queue
	//
	// When this error is returned, the operation returns the zero value for type T.
	//
	// Example:
	//
	//	q := queue.New[int]()
	//	val, err := q.Dequeue() // Returns 0, ErrUnderflow
	//	if errors.Is(err, queue.ErrUnderflow) {
	//		fmt.Println("Queue is empty")
	//	}
	ErrUnderflow = errors.New("queue underflow")
)
