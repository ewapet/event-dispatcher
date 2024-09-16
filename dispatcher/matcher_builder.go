package dispatcher

import (
	"github.com/ewapet/event-dispatcher/listener"
	"maps"
)

type MatcherBuilder[ID comparable, E any] struct {
	eventTypeListenerMap map[ID][]listener.Listener[ID, E]
	globalListeners      []listener.Listener[ID, E]
}

func NewMatcherBuilder[ID comparable, E any]() *MatcherBuilder[ID, E] {
	return &MatcherBuilder[ID, E]{
		eventTypeListenerMap: make(map[ID][]listener.Listener[ID, E]),
		globalListeners:      make([]listener.Listener[ID, E], 0),
	}
}

func (l *MatcherBuilder[ID, E]) AddGlobalListener(listenerToAdd listener.Listener[ID, E]) {
	l.globalListeners = append(l.globalListeners, listenerToAdd)
}

func (l *MatcherBuilder[ID, E]) AddListener(eventIDs []ID, listenerToAdd listener.Listener[ID, E]) {
	if len(eventIDs) == 0 {
		panic("invalid list of event types provided - is empty")
	}
	for _, eventType := range eventIDs {
		if l.eventTypeListenerMap[eventType] == nil {
			l.eventTypeListenerMap[eventType] = make([]listener.Listener[ID, E], 0)
		}
		l.eventTypeListenerMap[eventType] = append(l.eventTypeListenerMap[eventType], listenerToAdd)
	}
}

func (l *MatcherBuilder[ID, E]) Build() Matcher[ID, E] {
	listenerMapCopy := make(map[ID][]listener.Listener[ID, E])
	maps.Copy(listenerMapCopy, l.eventTypeListenerMap)

	globalListenersCopy := make([]listener.Listener[ID, E], len(l.globalListeners))
	copy(globalListenersCopy, l.globalListeners)

	return mapMatcher[ID, E]{
		eventTypeListenerMap: listenerMapCopy,
		globalListeners:      globalListenersCopy,
	}
}
