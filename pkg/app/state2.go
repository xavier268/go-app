package app

import (
	"sync"
)

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
