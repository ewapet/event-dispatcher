// Package dispatcher contains functionality for uni-directional event dispatching
package dispatcher

import (
	"github.com/ewapet/event-dispatcher/listener"
)

type dispatcher[ID comparable, E any] struct {
	listenerMatcher Matcher[ID, E]
}

// New creates an EventDispatcher using the provided matcher.
//
// If nil is provided, an EventDispatcher having no listeners is created
func New[ID comparable, E any](listenerMatcher Matcher[ID, E]) EventDispatcher[ID, E] {
	if listenerMatcher == nil {
		return NewFrom[ID, E](nil, nil)
	}
	return &dispatcher[ID, E]{
		listenerMatcher: listenerMatcher,
	}
}

// NewFrom creates an EventDispatcher using the provided global and event-specific listeners.
//
// Global listeners are always executed no matter what event ID is passed. Listeners triggered by a specific event ID
// can be specified as a map in the second parameter
//
// Nil can be passed for either parameter to specify that no listeners are present
func NewFrom[ID comparable, E any](
	globalListeners []listener.Listener[ID, E],
	listenerMap listener.Map[ID, E],
) EventDispatcher[ID, E] {

	matcherBuilder := NewMatcherBuilder[ID, E]()
	for _, currentGlobalListener := range globalListeners {
		matcherBuilder.AddGlobalListener(currentGlobalListener)
	}
	for eventType, currentListener := range listenerMap {
		matcherBuilder.AddListener([]ID{eventType}, currentListener)
	}
	return New[ID, E](matcherBuilder.Build())
}

// NewFromFunc creates an EventDispatcher using one or more of the provided closures.
//
// Each closure is converted into a global listener -- all will be invoked for each event that's dispatched
func NewFromFunc[ID comparable, E any](closure func(ID, E), additional ...func(ID, E)) EventDispatcher[ID, E] {
	closuresToBuild := []func(ID, E){closure}
	closuresToBuild = append(closuresToBuild, additional...)
	globalListeners := make([]listener.Listener[ID, E], 0, len(closuresToBuild))

	for _, currentClosure := range closuresToBuild {
		if currentClosure == nil {
			panic("invalid closure provided - is nil")
		}
		globalListeners = append(globalListeners, listener.NewFromFunc[ID, E](currentClosure))
	}

	return NewFrom[ID, E](globalListeners, nil)
}

// NewFromListener creates an EventDispatcher using one or more provided listeners
//
// Each listener is converted into a global listener -- all will be invoked for each event that's dispatched.
// Note: Nil listeners are not permitted and the function will panic if one is encountered.
func NewFromListener[ID comparable, E any](targetListener listener.Listener[ID, E], additional ...listener.Listener[ID, E]) EventDispatcher[ID, E] {
	listenersToAdd := []listener.Listener[ID, E]{targetListener}
	listenersToAdd = append(listenersToAdd, additional...)
	globalListeners := make([]listener.Listener[ID, E], 0, len(listenersToAdd))

	for _, currentListener := range listenersToAdd {
		if currentListener == nil {
			panic("invalid listener provided - is nil")
		}
		globalListeners = append(globalListeners, currentListener)
	}

	return NewFrom[ID, E](globalListeners, nil)
}

// NewFromMap creates an EventDispatcher where the specified listeners are only triggered for the corresponding event ID they're mapped to.
//
// Each listener is converted into a global listener -- all will be invoked for each event that's dispatched.
// Note: Nil maps or listeners are not permitted and the function will panic if one is encountered.
func NewFromMap[ID comparable, E any](listenerMap map[ID]listener.Listener[ID, E]) EventDispatcher[ID, E] {
	if listenerMap == nil {
		panic("invalid map provided - is nil")
	}
	return NewFrom[ID, E](nil, listenerMap)
}

// Dispatch invokes all relevant listeners for the provided event ID.
//
// If no events are specified alongside the ID, then the zero value of the event type is dispatched instead.
// If multiple event objects are specified, each listener is invoked with a single event sequentially.
func (d dispatcher[ID, E]) Dispatch(eventID ID, eventPayloads ...E) {
	listenerTargets := d.listenerMatcher.Match(eventID)
	if len(eventPayloads) == 0 {
		eventPayloads = append(eventPayloads, *new(E))
	}
	for _, currentEvent := range eventPayloads {
		for _, currentListenerTarget := range listenerTargets {
			currentListenerTarget.Receive(eventID, currentEvent)
		}
	}
}
