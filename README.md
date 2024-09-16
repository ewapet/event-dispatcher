# Event Dispatcher

## Overview

A simple, flexible, and generic event-dispatcher implementation. 

**Example:** Create a global listener that accepts all incoming events and outputs their values.
```go

// Creates a new dispatcher with a single global listener accepting all events
myDispatcher := dispatcher.NewFromFunc(func(eventID uint, eventData string) {
	
    fmt.Sprintf("You dispatched event ID %s with data %v", eventID, eventData)

})

myDispatcher.dispatch("event.message", "Hello World")
```

### Feature Preview

1) **Dispatch Anything** - the library does not dictate what your events should look like. Only that they're of a consistent type.
1) **Extensible** - Inject custom dispatching logic to suit your own needs, whether deterministic or non-deterministic.

## Usage

### Dispatching Events

Event dispatching is comprised of two components: the **EventID** and it's associated **EventData**. Each listener is called
once per **EventData** entry provided. 

Users are free to dispatch an event ID alone (of type `comparable`) without any associated data. In that case, the event
data is the zero value of the corresponding generic type.

**Example:** Dispatching only an EventID
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
    // listener is called twice with an eventID of "event.greeting"
    // eventData is "hello" on first call, and "bonjour" on second call
    ...
})
myDispatcher.dispatch("event.greeting", "hello", "bonjour")
```

### Defining Listeners

```go
EventDispatcher[ID comparable, E any]
```

**Example**: Creating a dispatcher with one global closure listener 
```go
broadcaster := dispatcher.NewFromFunc(func(eventID string, eventData string) {
    fmt.Println("Hello!")
})

broadcaster.dispatch()
```

**Example**: Creating a dispatcher with multiple closure listeners

```go
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

broadcaster.dispatch("event.greeting")
```




