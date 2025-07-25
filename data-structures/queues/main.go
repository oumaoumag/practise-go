package main

import (
	
)

// MyQueue is a generic queue implementation using a slice
type Queue[T any] struct {
	slice []T
}