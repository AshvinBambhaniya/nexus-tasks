// Package routinewrapper provides a way to wrap goroutines with a global handler.
package routinewrapper

import (
	"sync"
)

var handle func()
var _once sync.Once

// Init initializes the global handle function.
func Init(fn func()) {
	_once.Do(func() {
		// this sets the global handle function
		handle = fn
	})
}

// RoutineGenerator wraps a function call with the global handle function.
func RoutineGenerator(fn func()) {
	defer handle()
	fn()
}
