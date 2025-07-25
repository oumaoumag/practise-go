package main

import "errors"

// MyQueue is a generic queue implementation using a slice
type MyQueue[T any] struct {
	slice []T
}

// Enqueue add an element to the end of the queue.
func (q *MyQueue[T]) Enqueue(value T) {
	q.slice = append(q.slice, value)
}

// Dequeue removes and returns the element feom the front of the queue
func (q *MyQueue[T]) Dequeue() (T, error) {
	if len(q.slice) == 0 {
		var zero T
		return zero, errors.New("queue is empty")
	}

	value := q.slice[0]
	q.slice = q.slice[1:]
	return value, nil

}