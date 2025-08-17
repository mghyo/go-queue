package queue

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNew(t *testing.T) {
	q := New[int]()
	if q == nil {
		t.Fatal("New() returned nil")
	}

	if size := q.Size(); size != 0 {
		t.Errorf("New queue size = %d, want 0", size)
	}
}

func TestNewWithCapacity(t *testing.T) {
	q := New[int](WithCapacity[int](5))
	if q == nil {
		t.Fatal("New() with capacity returned nil")
	}

	if size := q.Size(); size != 0 {
		t.Errorf("New queue with capacity size = %d, want 0", size)
	}
}

func TestEnqueueDequeue(t *testing.T) {
	q := New[int]()

	// Test enqueue
	err := q.Enqueue(42)
	if err != nil {
		t.Errorf("Enqueue(42) error = %v, want nil", err)
	}

	if size := q.Size(); size != 1 {
		t.Errorf("Size after enqueue = %d, want 1", size)
	}

	// Test dequeue
	val, err := q.Dequeue()
	if err != nil {
		t.Errorf("Dequeue() error = %v, want nil", err)
	}

	if val != 42 {
		t.Errorf("Dequeue() = %d, want 42", val)
	}

	if size := q.Size(); size != 0 {
		t.Errorf("Size after dequeue = %d, want 0", size)
	}
}

func TestFIFOBehavior(t *testing.T) {
	q := New[int]()
	values := []int{1, 2, 3, 4, 5}

	// Enqueue all values
	for _, v := range values {
		err := q.Enqueue(v)
		if err != nil {
			t.Errorf("Enqueue(%d) error = %v, want nil", v, err)
		}
	}

	if size := q.Size(); size != len(values) {
		t.Errorf("Size after enqueues = %d, want %d", size, len(values))
	}

	// Dequeue all values (should be in same order - FIFO)
	for i, expected := range values {
		val, err := q.Dequeue()
		if err != nil {
			t.Errorf("Dequeue() error = %v, want nil", err)
		}
		if val != expected {
			t.Errorf("Dequeue() at position %d = %d, want %d", i, val, expected)
		}
	}

	if size := q.Size(); size != 0 {
		t.Errorf("Size after all dequeues = %d, want 0", size)
	}
}

func TestPeek(t *testing.T) {
	q := New[string]()

	// Test peek on empty queue
	_, err := q.Peek()
	if !errors.Is(err, ErrUnderflow) {
		t.Errorf("Peek() on empty queue error = %v, want ErrUnderflow", err)
	}

	// Enqueue and peek
	err = q.Enqueue("first")
	if err != nil {
		t.Errorf("Enqueue() error = %v, want nil", err)
	}

	val, err := q.Peek()
	if err != nil {
		t.Errorf("Peek() error = %v, want nil", err)
	}
	if val != "first" {
		t.Errorf("Peek() = %q, want %q", val, "first")
	}

	// Size should remain the same after peek
	if size := q.Size(); size != 1 {
		t.Errorf("Size after peek = %d, want 1", size)
	}

	// Enqueue another value
	err = q.Enqueue("second")
	if err != nil {
		t.Errorf("Enqueue() error = %v, want nil", err)
	}

	// Peek should still return the first value (front of queue)
	val, err = q.Peek()
	if err != nil {
		t.Errorf("Peek() error = %v, want nil", err)
	}
	if val != "first" {
		t.Errorf("Peek() after second enqueue = %q, want %q", val, "first")
	}

	// Dequeue first item
	dequeued, err := q.Dequeue()
	if err != nil {
		t.Errorf("Dequeue() error = %v, want nil", err)
	}
	if dequeued != "first" {
		t.Errorf("Dequeue() = %q, want %q", dequeued, "first")
	}

	// Now peek should return "second"
	val, err = q.Peek()
	if err != nil {
		t.Errorf("Peek() after dequeue error = %v, want nil", err)
	}
	if val != "second" {
		t.Errorf("Peek() after dequeue = %q, want %q", val, "second")
	}
}

func TestDequeueUnderflow(t *testing.T) {
	q := New[int]()

	val, err := q.Dequeue()
	if !errors.Is(err, ErrUnderflow) {
		t.Errorf("Dequeue() on empty queue error = %v, want ErrUnderflow", err)
	}

	// Check that zero value is returned
	if val != 0 {
		t.Errorf("Dequeue() on empty queue value = %d, want 0 (zero value)", val)
	}
}

func TestCapacityOverflow(t *testing.T) {
	q := New[int](WithCapacity[int](2))

	// Enqueue up to capacity
	err := q.Enqueue(1)
	if err != nil {
		t.Errorf("Enqueue(1) error = %v, want nil", err)
	}

	err = q.Enqueue(2)
	if err != nil {
		t.Errorf("Enqueue(2) error = %v, want nil", err)
	}

	// Try to exceed capacity
	err = q.Enqueue(3)
	if !errors.Is(err, ErrOverflow) {
		t.Errorf("Enqueue(3) exceeding capacity error = %v, want ErrOverflow", err)
	}

	// Size should still be 2
	if size := q.Size(); size != 2 {
		t.Errorf("Size after overflow attempt = %d, want 2", size)
	}

	// Dequeue one item, then should be able to enqueue again
	_, err = q.Dequeue()
	if err != nil {
		t.Errorf("Dequeue() after overflow error = %v, want nil", err)
	}

	err = q.Enqueue(3)
	if err != nil {
		t.Errorf("Enqueue(3) after dequeue error = %v, want nil", err)
	}
}

func TestUnlimitedCapacity(t *testing.T) {
	q := New[int]() // Default is unlimited

	// Enqueue many items
	for i := 0; i < 1000; i++ {
		err := q.Enqueue(i)
		if err != nil {
			t.Errorf("Enqueue(%d) with unlimited capacity error = %v, want nil", i, err)
		}
	}

	if size := q.Size(); size != 1000 {
		t.Errorf("Size with unlimited capacity = %d, want 1000", size)
	}
}

func TestDifferentTypes(t *testing.T) {
	t.Run("string queue", func(t *testing.T) {
		q := New[string]()
		_ = q.Enqueue("first")
		_ = q.Enqueue("second")

		val, _ := q.Dequeue()
		if val != "first" {
			t.Errorf("Dequeue() = %q, want %q", val, "first")
		}

		val, _ = q.Dequeue()
		if val != "second" {
			t.Errorf("Dequeue() = %q, want %q", val, "second")
		}
	})

	t.Run("struct queue", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		q := New[Person]()
		person1 := Person{Name: "Alice", Age: 30}
		person2 := Person{Name: "Bob", Age: 25}

		_ = q.Enqueue(person1)
		_ = q.Enqueue(person2)

		val, err := q.Dequeue()
		if err != nil {
			t.Errorf("Dequeue() error = %v, want nil", err)
		}
		if val != person1 {
			t.Errorf("Dequeue() = %+v, want %+v", val, person1)
		}

		val, err = q.Dequeue()
		if err != nil {
			t.Errorf("Dequeue() error = %v, want nil", err)
		}
		if val != person2 {
			t.Errorf("Dequeue() = %+v, want %+v", val, person2)
		}
	})
}

func TestConcurrency(t *testing.T) {
	q := New[int]()
	const numGoroutines = 100
	const numOperations = 100

	var wg sync.WaitGroup

	// Start multiple goroutines enqueueing values
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_ = q.Enqueue(start*numOperations + j)
			}
		}(i)
	}

	wg.Wait()

	// Check that all items were enqueued
	expectedSize := numGoroutines * numOperations
	if size := q.Size(); size != expectedSize {
		t.Errorf("Size after concurrent enqueues = %d, want %d", size, expectedSize)
	}

	// Start multiple goroutines dequeuing values
	results := make(chan int, expectedSize)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				if val, err := q.Dequeue(); err == nil {
					results <- val
				}
			}
		}()
	}

	wg.Wait()
	close(results)

	// Count results
	count := 0
	seen := make(map[int]bool)
	for val := range results {
		count++
		if seen[val] {
			t.Errorf("Duplicate value dequeued: %d", val)
		}
		seen[val] = true
	}

	if count != expectedSize {
		t.Errorf("Dequeued %d items, want %d", count, expectedSize)
	}

	// Queue should be empty
	if size := q.Size(); size != 0 {
		t.Errorf("Size after concurrent dequeues = %d, want 0", size)
	}
}

func TestMixedConcurrentOperations(t *testing.T) {
	q := New[int]()
	const numGoroutines = 50
	const numOperations = 200

	var wg sync.WaitGroup

	// Concurrent enqueues
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_ = q.Enqueue(start*numOperations + j)
			}
		}(i)
	}

	// Concurrent dequeues (fewer than enqueues)
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations/2; j++ {
				_, _ = q.Dequeue() // Ignore errors for this test
			}
		}()
	}

	// Concurrent peeks
	for i := 0; i < numGoroutines/4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_, _ = q.Peek() // Ignore errors
			}
		}()
	}

	// Concurrent size checks
	for i := 0; i < numGoroutines/4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				q.Size()
			}
		}()
	}

	wg.Wait()

	// Just verify no panics occurred and queue is in valid state
	size := q.Size()
	if size < 0 {
		t.Errorf("Invalid size after mixed concurrent operations: %d", size)
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("zero capacity", func(t *testing.T) {
		q := New[int](WithCapacity[int](0))
		err := q.Enqueue(1)
		if !errors.Is(err, ErrOverflow) {
			t.Errorf("Enqueue() with zero capacity error = %v, want ErrOverflow", err)
		}
	})

	t.Run("negative capacity (should panic)", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("WithCapacity(-5) should panic, but it didn't")
			}
		}()

		// This should panic
		New[int](WithCapacity[int](-5))

		// If we get here, the test should fail
		t.Error("Expected panic did not occur")
	})

	t.Run("explicit unlimited capacity", func(t *testing.T) {
		q := New[int](WithCapacity[int](UnlimitedCapacity))

		// Should work like unlimited capacity
		for i := 0; i < 10; i++ {
			err := q.Enqueue(i)
			if err != nil {
				t.Errorf("Enqueue(%d) with UnlimitedCapacity error = %v, want nil", i, err)
			}
		}

		if size := q.Size(); size != 10 {
			t.Errorf("Size with UnlimitedCapacity = %d, want 10", size)
		}
	})

	t.Run("multiple options", func(t *testing.T) {
		// Last option should win
		q := New[int](
			WithCapacity[int](5),
			WithCapacity[int](3),
		)

		// Should only be able to enqueue 3 items
		for i := 0; i < 3; i++ {
			err := q.Enqueue(i)
			if err != nil {
				t.Errorf("Enqueue(%d) error = %v, want nil", i, err)
			}
		}

		err := q.Enqueue(3)
		if !errors.Is(err, ErrOverflow) {
			t.Errorf("Enqueue(3) exceeding final capacity error = %v, want ErrOverflow", err)
		}
	})
}

func TestQueueWithCapacityStress(t *testing.T) {
	const capacity = 100
	q := New[int](WithCapacity[int](capacity))

	var enqueueCount, dequeueCount int64
	var wg sync.WaitGroup

	// Producer goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				err := q.Enqueue(id*1000 + j)
				if err == nil {
					atomic.AddInt64(&enqueueCount, 1)
				}
			}
		}(i)
	}

	// Consumer goroutines
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 500; j++ {
				_, err := q.Dequeue()
				if err == nil {
					atomic.AddInt64(&dequeueCount, 1)
				}
			}
		}()
	}

	wg.Wait()

	// Verify consistency
	finalSize := int64(q.Size())
	expectedSize := enqueueCount - dequeueCount

	if finalSize != expectedSize {
		t.Errorf("Size inconsistency: got %d, expected %d (enqueued: %d, dequeued: %d)",
			finalSize, expectedSize, enqueueCount, dequeueCount)
	}

	// Size should never exceed capacity
	if finalSize > capacity {
		t.Errorf("Size exceeded capacity: %d > %d", finalSize, capacity)
	}
}

// Benchmark tests
func BenchmarkEnqueue(b *testing.B) {
	q := New[int]()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = q.Enqueue(i)
	}
}

func BenchmarkDequeue(b *testing.B) {
	q := New[int]()
	// Pre-populate queue
	for i := 0; i < b.N; i++ {
		_ = q.Enqueue(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Dequeue()
	}
}

func BenchmarkPeek(b *testing.B) {
	q := New[int]()
	_ = q.Enqueue(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Peek()
	}
}

func BenchmarkSize(b *testing.B) {
	q := New[int]()
	for i := 0; i < 1000; i++ {
		_ = q.Enqueue(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Size()
	}
}

func BenchmarkMixedOperations(b *testing.B) {
	q := New[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		switch i % 4 {
		case 0:
			_ = q.Enqueue(i)
		case 1:
			_, _ = q.Dequeue()
		case 2:
			_, _ = q.Peek()
		case 3:
			q.Size()
		}
	}
}
