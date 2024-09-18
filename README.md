# Event Dispatcher

## Overview

A simple, versatile, and generic event-dispatcher implementation. 

### Goals

1) **Dispatch Anything** - the library doesn't dictate what your events look like. Only that they're consistent. Or type 'any'
1) **Extensible** - Inject custom dispatching logic to suit your own needs, whether deterministic or non-deterministic.

```go
// Type: EventDispatcher[string, []rune]
myDispatcher := dispatcher.NewFromFunc(func(eventID string, eventData []rune) {

    fmt.Printf("You dispatched event ID %s with data %v", eventID, eventData)

})

myDispatcher.Dispatch("event.message", []rune("Hello World"))
```

## Usage

### A Word About Listeners

Though the dispatcher API provides helper functions for creating listeners -- like `NewFromFunc` displayed above -- you can
add your own custom listeners so long as they implement the following interface.

```go
type Listener[ID comparable, E any] interface {
    Receive(ID, E)
}
```

### Using the Default Dispatcher Implementation

The in-built dispatcher supports Global Listeners and ID Subscribed listeners (E.g. listeners that are only triggered when a specific EventID is fired).

#### Global Listeners

Global listeners are always triggered, no matter what event ID is dispatched. If you need simple broadcasting functionality
-- perhaps as a wrapper for multithreaded behavior -- you can define one or more closure-based listeners using `NewFromFunc`.

To add a custom listener, implement the interface in the previous section and use `dispatcher.NewFromListener(...)`

**Example**: Creating a dispatcher with one global closure listener.
```go
// Creates a new dispatcher of generic type EventDispatcher[string, string]
broadcaster := dispatcher.NewFromFunc(func(eventID string, eventData string) {
    fmt.Println("Hello!")
})

broadcaster.Dispatch("any.thing")
```

**Example**: Creating a dispatcher with multiple global closure listeners.
```go
// Creates a new dispatcher of generic type EventDispatcher[string, string]
broadcaster := dispatcher.NewFromFunc(
    func(eventID string, eventData string) {
        fmt.Println("Hello!")
    },
    func(eventID string, eventData string) {
        fmt.Println("Hola!")
    },
    func(eventID string, eventData string) {
        fmt.Println("Bonjour!")
    },
)

broadcaster.Dispatch("event.greeting")
```

**Example**: Creating a dispatcher with a custom listener.
```go
// Implement the listener interface
myListener := ...
dispatcher.NewFromListener(myListener)
```

#### ID Subscribed Listeners

The current implementation of the dispatcher supports simple map-based event-subscribing via `dispatcher.NewFromMap()`.

Where each listener is only triggered for a specific event ID.

**Example**: Creating a map based listener

```go
// Create a dispatcher with two listeners monitoring two distinct event IDs
myDispatcher := dispatcher.NewFromMap(listener.Map[string, any]{

    "event.greeting": listener.NewFromFunc(func(id string, e any) {
        fmt.Println("Hello")
    }),

    "event.farewell": listener.NewFromFunc(func(id string, e any) {
        fmt.Println("Goodbye")
    }),
})

myDispatcher.Dispatch("event.greeting") // Hello is output
```

### Mixing Listener Types

To specify both global and subscription based events, use `dispatcher.NewFrom()`

```go
// Creates a dispatcher with no global listeners, and one subscribed to "event.greeting"
myDispatcher := dispatcher.NewFrom(nil, listener.Map[string, any]{
    "event.greeting": listener.NewFromFunc(func(id string, e any) {
        fmt.Println("Hello")
    }),
})

myDispatcher.Dispatch("event.greeting")
```

### Dispatching Events

Event dispatching is comprised of two components: the **EventID** (of type comparable) and it's associated **EventData**. 
Each listener is called once per **EventData** entry provided. Event data is optional. If none is provided, listeners
are called once with the zero value of the generic type.

**Example:** Dispatching only an EventID without any data
```go
myDispatcher := dispatcher.NewFromFunc(func(eventID string, eventData []time.Time) {
    // eventData is nil, the zero value of a slice
    ...
})

myDispatcher.dispatch("event.example")
```

**Example:** Dispatching an event with one associated data entry

```go
myDispatcher := dispatcher.NewFromFunc(func(eventID string, eventData string) {
    // eventData is "hello", the zero value of string
    ...
})

myDispatcher.dispatch("event.greeting", "hello")
```

**Example:** Dispatching an event with multiple associated data entries

```go
myDispatcher := dispatcher.NewFromFunc(func(eventID string, eventData string) {
    // listener is called twice with an eventID of "event.greeting" both times
    // eventData is "hello" on first call, and "bonjour" on second call
    ...
})
myDispatcher.dispatch("event.greeting", "hello", "bonjour")
```

### Custom Event Dispatching

All in-built logic above can be replaced by injecting a custom `Matcher` interface implementation.

To do so, use `dispatcher.New()` and pass the matcher

```go
type Matcher[ID comparable, E any] interface {
	Match(ID) []listener.Listener[ID, E]
}
```

```go
myMatcher := ...
dispatcher.New(myMatcher)
```





