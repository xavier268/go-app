package app

import "reflect"

type Observer2 interface {
	// Defines a condition that reports whether the observer keeps observing the
	// associated state. Multiple conditions can be defined by successively
	// calling While().
	While(condition func() bool) Observer2

	// Executes the given function on the UI goroutine when the observed value
	// changes. Multiple functions can be executed by successively calling
	// OnChange().
	OnChange(fn func()) Observer2

	// Stores the status code associated with the observed state into the given
	// code. Panics when the given code is nil.
	Status(code *int) Observer2

	// Stores the error associated with the observed state into the given error.
	// Panics when the given error is nil.
	Error(err *error) Observer2

	// Stores the value associated with the observed state into the given
	// receiver. Panics when the receiver is not a pointer or nil.
	//
	// The receiver is updated each time the associated state changes. It is
	// unchanged when its pointed value has a different type than the associated
	// state value.
	Value(recv interface{})
}

type observer2 struct {
	source     UI
	conditions []func() bool
	onChanges  []func()
	status     *int
	err        *error
	value      interface{}
	subscribe  func(*observer2)
}

func newObserver2(source UI, subscribe func(*observer2)) *observer2 {
	return &observer2{
		source:    source,
		subscribe: subscribe,
	}
}

func (o *observer2) While(condition func() bool) Observer2 {
	o.conditions = append(o.conditions, condition)
	return o
}

func (o *observer2) OnChange(fn func()) Observer2 {
	o.onChanges = append(o.onChanges, fn)
	return o
}

func (o *observer2) Status(code *int) Observer2 {
	if code == nil {
		panic("observer status is nil")
	}
	o.status = code
	return o
}

func (o *observer2) Error(err *error) Observer2 {
	if err == nil {
		panic("observer error is nil")
	}
	o.err = err
	return o
}

func (o *observer2) Value(v interface{}) {
	if v == nil {
		panic("observer value is nil")
	}

	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		panic("observer value is not a pointer")
	}

	o.value = v
	o.subscribe(o)
}
