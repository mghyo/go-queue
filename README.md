# Go Queue

A thread-safe, generic FIFO queue implementation for Go with configurable capacity limits.

[![Go Reference](https://pkg.go.dev/badge/github.com/yourusername/go-queue.svg)](https://pkg.go.dev/github.com/yourusername/go-queue)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/go-queue)](https://goreportcard.com/report/github.com/yourusername/go-queue)

## Features

- **Generic**: Works with any type using Go's type parameters
- **Thread-safe**: Safe for concurrent use across multiple goroutines
- **FIFO semantics**: First-In-First-Out behavior
- **Configurable capacity**: Set maximum size or use unlimited capacity
- **Zero dependencies**: Uses only Go standard library

## Installation

```bash
go get github.com/yourusername/go-queue
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/yourusername/go-queue"
)

func main() {
    // Create a new queue
    q := queue.New[int]()
    
    // Add items (FIFO)
    q.Enqueue(1)
    q.Enqueue(2)
    q.Enqueue(3)
    
    // Remove items in order
    for q.Size() > 0 {
        val, _ := q.Dequeue()
        fmt.Println(val) // Prints: 1, 2, 3
    }
}
```

## Usage

### Basic Operations

```go
q := queue.New[string]()

// Add items to back of queue
q.Enqueue("first")
q.Enqueue("second")

// Check front item without removing
front, err := q.Peek() // Returns "first"

// Remove items from front
val, err := q.Dequeue() // Returns "first"
val, err = q.Dequeue()  // Returns "second"

// Check size
fmt.Println(q.Size()) // 0
```

### Capacity-Limited Queue

```go
// Create queue with max capacity of 3
q := queue.New[int](queue.WithCapacity[int](3))

q.Enqueue(1) // OK
q.Enqueue(2) // OK  
q.Enqueue(3) // OK
err := q.Enqueue(4) // Returns queue.ErrOverflow
```

### Error Handling

```go
q := queue.New[int]()

// Empty queue operations return ErrUnderflow
val, err := q.Dequeue()
if errors.Is(err, queue.ErrUnderflow) {
    fmt.Println("Queue is empty")
}

val, err = q.Peek()
if errors.Is(err, queue.ErrUnderflow) {
    fmt.Println("Queue is empty")
}
```

## API Reference

### Types

```go
type Queue[T any] interface {
    Enqueue(val T) error   // Add item to back
    Dequeue() (T, error)   // Remove item from front  
    Size() int             // Current number of items
    Peek() (T, error)      // View front item without removing
}
```

### Functions

```go
// Create new queue
func New[T any](opts ...Option[T]) Queue[T]

// Set maximum capacity (-1 for unlimited)
func WithCapacity[T any](cap int) Option[T]
```

### Constants & Errors

```go
const UnlimitedCapacity = -1

var ErrOverflow = errors.New("queue overflow")   // Queue is full
var ErrUnderflow = errors.New("queue underflow") // Queue is empty
```

## Performance

- **Enqueue**: O(1) amortized
- **Dequeue**: O(n) - shifts all elements
- **Peek**: O(1)
- **Size**: O(1)

> **Note**: This implementation prioritizes simplicity over performance. For high-throughput applications requiring O(1) dequeue, consider using a circular buffer implementation.

## Thread Safety

All operations are thread-safe and can be used concurrently:

```go
q := queue.New[int]()
var wg sync.WaitGroup

// Multiple producers
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(val int) {
        defer wg.Done()
        q.Enqueue(val)
    }(i)
}

// Multiple consumers  
for i := 0; i < 5; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        q.Dequeue()
    }()
}

wg.Wait()
```

## Testing

```bash
go test              # Run tests
go test -race        # Run with race detection
go test -bench=.     # Run benchmarks
go test -cover       # Check coverage
```

## Requirements

- Go 1.18+ (for generics support)

## License

MIT License - see [LICENSE](LICENSE) file for details.