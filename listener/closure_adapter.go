package listener

type closureAdapter[ID comparable, E any] struct {
	internalFunc func(ID, E)
}

func newClosureAdapter[ID comparable, E any](listenerClosure func(ID, E)) *closureAdapter[ID, E] {
	return &closureAdapter[ID, E]{
		internalFunc: listenerClosure,
	}
}

func (a closureAdapter[ID, E]) Receive(eventType ID, incomingEvent E) {
	a.internalFunc(eventType, incomingEvent)
}
