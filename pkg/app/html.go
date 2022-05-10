package app

import (
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/errors"
)

type htmlElement[T any] struct {
	tag           string
	isSelfClosing bool
	body          []UI
	attributes    map[string]string
	events        map[string]eventHandler
}

func (e htmlElement[T]) Text(v any) T {
	switch e.tag {
	case "textarea":
		return e.Attr("value", v)

	default:
		return e.Body(Text(v))
	}
}

func (e htmlElement[T]) Body(v ...UI) T {
	if e.isSelfClosing {
		panic(errors.New("setting html element body failed").
			Tag("reason", "self closing element can't have children").
			Tag("tag", e.tag),
		)
	}

	e.body = FilterUIElems(v...)
	return e.toHTMLInterface()
}

func (e htmlElement[T]) Attr(k string, v any) T {
	if e.attributes == nil {
		e.attributes = make(map[string]string)
	}

	switch k {
	case "style", "allow":
		var b strings.Builder
		b.WriteString(e.attributes[k])
		b.WriteString(toAttributeValue(v))
		b.WriteByte(';')
		e.attributes[k] = b.String()

	case "class":
		var b strings.Builder
		b.WriteString(e.attributes[k])
		if b.Len() != 0 {
			b.WriteByte(' ')
		}
		b.WriteString(toAttributeValue(v))
		e.attributes[k] = b.String()

	default:
		e.attributes[k] = toAttributeValue(v)
	}

	return e.toHTMLInterface()
}

func (e htmlElement[T]) On(event string, h EventHandler, scope ...any) T {
	if e.events == nil {
		e.events = make(map[string]eventHandler)
	}

	e.events[event] = eventHandler{
		event: event,
		scope: toPath(scope...),
		value: h,
	}

	return e.toHTMLInterface()
}

func (e htmlElement[T]) toHTMLInterface() T {
	var i any = e
	return i.(T)
}
