package listener

// A Listener receives events and executes corresponding logic
type Listener[ID comparable, E any] interface {
	Receive(ID, E)
}

// Map is a shorthand type for a map of event IDs to their corresponding listeners
type Map[ID comparable, E any] map[ID]Listener[ID, E]

// NewFromFunc creates a new listener from the provided function
func NewFromFunc[ID comparable, E any](listenerFunc func(ID, E)) Listener[ID, E] {
	return newClosureAdapter[ID, E](listenerFunc)
}

// Zero creates a new listener that doesn't perform any actions
func Zero[ID comparable, E any]() Listener[ID, E] {
	return NewFromFunc(func(id ID, e E) {})
}
