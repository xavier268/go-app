package app

import (
	"reflect"
	"sync"
)

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
}

func newObserver2(source UI) *observer2 {
	return &observer2{
		source: source,
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
}

type Notifier interface {
	Option(o StateOption) Notifier

	Broadcast() Notifier

	Persist() Notifier

	Encrypt() Notifier

	Status(code int) Notifier

	Error(err error)

	Value(v interface{})
}

// StateOption represents an option applied when a state is set.
type StateOption func(*State)

type State2 struct {
	mutex sync.RWMutex

	key        string
	statusCode int
	err        error
	value      interface{}
	observers  map[*observer2]struct{}
}

func (s *State2) AddObserver(o *observer2) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.observers[o] = struct{}{}
}

func (s *State2) removeObserver(o *observer2) {
	delete(s.observers, o)
}
