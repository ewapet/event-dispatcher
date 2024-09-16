package dispatcher

import (
	"github.com/ewapet/event-dispatcher/listener"
)

type mapMatcher[ID comparable, E any] struct {
	eventTypeListenerMap map[ID][]listener.Listener[ID, E]
	globalListeners      []listener.Listener[ID, E]
}

func (m mapMatcher[ID, E]) Match(t ID) []listener.Listener[ID, E] {
	listenersToReturn := make([]listener.Listener[ID, E], 0, len(m.globalListeners))
	listenersToReturn = append(listenersToReturn, m.globalListeners...)
	if m.eventTypeListenerMap[t] != nil {
		listenersToReturn = append(listenersToReturn, m.eventTypeListenerMap[t]...)
	}
	return listenersToReturn
}
