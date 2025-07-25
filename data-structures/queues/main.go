package main

import (
	
)

// MyQueue is a generic queue implementation using a slice
type Queue[T any] struct {
	slice []T
}

// Enqueue add an element to the end of the queue.
func (q *Queue[T]) Enqueue(value T) {
	q.slice = append(q.slice, value)
}

