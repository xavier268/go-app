package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObserver2(t *testing.T) {
	foo := &foo{}
	c := NewClientTester(foo)
	defer c.Close()

	o := newObserver2(foo)
	require.Equal(t, foo, o.source)

	t.Run("set while", func(t *testing.T) {
		o.While(func() bool { return true })
		require.Len(t, o.conditions, 1)
	})

	t.Run("set on change", func(t *testing.T) {
		o.OnChange(func() {})
		require.Len(t, o.onChanges, 1)
	})

	t.Run("set status", func(t *testing.T) {
		var code int
		o.Status(&code)
		require.Equal(t, &code, o.status)
	})

	t.Run("set nil status panics", func(t *testing.T) {
		require.Panics(t, func() {
			o.Status(nil)
		})
	})

	t.Run("set error", func(t *testing.T) {
		var err error
		o.Error(&err)
		require.Equal(t, &err, o.err)
	})

	t.Run("set nil error panics", func(t *testing.T) {
		require.Panics(t, func() {
			o.Error(nil)
		})
	})

	t.Run("set value", func(t *testing.T) {
		var str string
		o.Value(&str)
		require.Equal(t, &str, o.value)
	})

	t.Run("set nil value panics", func(t *testing.T) {
		require.Panics(t, func() {
			o.Value(nil)
		})
	})

	t.Run("set non pointer value panics", func(t *testing.T) {
		var str string
		require.Panics(t, func() {
			o.Value(str)
		})
	})
}
