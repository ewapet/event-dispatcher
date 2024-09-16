package dispatcher

import (
	"github.com/ewapet/event-dispatcher/listener"
)

// An EventDispatcher invokes listeners corresponding to the provided event ID and event objects.
type EventDispatcher[ID comparable, E any] interface {
	Dispatch(ID, ...E)
}

// A Matcher identifies relevant listeners for the given event ID, returning them as a slice.
type Matcher[ID comparable, E any] interface {
	Match(ID) []listener.Listener[ID, E]
}
