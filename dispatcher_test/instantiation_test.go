package dispatcher_test

import (
	"github.com/ewapet/event-dispatcher/dispatcher"
	"github.com/ewapet/event-dispatcher/listener"
	"testing"
)

func Test_GivenCreatingADispatcherViaNew_WhenNilIsPassed_ThenAnEmptyListenerIsCreated(t *testing.T) {
	subject := dispatcher.New[any, any](nil)
	subject.Dispatch("testing.testing")
}

func Test_GivenCreatingADispatcherViaNewFromFunc_WhenANilClosureIsPassed_ThenTheFunctionPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected a panic but none occurred")
		}
	}()
	dispatcher.NewFromFunc[string, string](nil)
}

func Test_GivenCreatingADispatcherViaNewFromFunc_WhenANilClosureIsPassedAfterOthers_ThenTheFunctionPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected a panic but none occurred")
		}
	}()
	dispatcher.NewFromFunc[string, string](
		func(s string, s2 string) {},
		nil,
	)
}

func Test_GivenCreatingADispatcherViaNewFromFunc_WhenAClosureIsPassed_ThenAGlobalListenerIsCreated(t *testing.T) {
	expectedCallCount := 2
	actualCallCount := 0

	subject := dispatcher.NewFromFunc(func(ID string, event string) {
		actualCallCount += 1
	})

	subject.Dispatch("testing.event1", "event")
	subject.Dispatch("testing.event2", "event")

	if expectedCallCount != actualCallCount {
		t.Errorf("expected to be called %#v, but was called %#v", expectedCallCount, actualCallCount)
	}
}

func Test_GivenCreatingADispatcherViaNewFromFunc_WhenTwoClosuresArePassed_ThenTwoGlobalListenersAreCreated(t *testing.T) {
	expectedCallCount := 4
	actualCallCount := 0

	subject := dispatcher.NewFromFunc(
		func(ID string, event string) {
			actualCallCount += 1
		},
		func(ID string, event string) {
			actualCallCount += 1
		},
	)

	subject.Dispatch("testing.event1", "event")
	subject.Dispatch("testing.event2", "event")

	if expectedCallCount != actualCallCount {
		t.Errorf("expected to be called %#v, but was called %#v", expectedCallCount, actualCallCount)
	}
}

func Test_GivenCreatingADispatcherViaNewFromListener_WhenANilClosureIsPassed_ThenTheFunctionPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected a panic but none occurred")
		}
	}()
	dispatcher.NewFromListener[string, string](nil)
}

func Test_GivenCreatingADispatcherViaNewFromListener_WhenANilClosureIsPassedAfterOthers_ThenTheFunctionPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected a panic but none occurred")
		}
	}()
	dispatcher.NewFromListener(listener.Zero[string, string](), nil)
}

func Test_GivenCreatingADispatcherViaNewFromListener_WhenAListenerIsPassed_ThenAGlobalListenerIsCreated(t *testing.T) {
	expectedCallCount := 2
	actualCallCount := 0

	subject := dispatcher.NewFromListener(listener.NewFromFunc(func(ID string, event string) {
		actualCallCount += 1
	}))

	subject.Dispatch("testing.event1", "event")
	subject.Dispatch("testing.event2", "event")

	if expectedCallCount != actualCallCount {
		t.Errorf("expected to be called %#v, but was called %#v", expectedCallCount, actualCallCount)
	}
}

func Test_GivenCreatingADispatcherViaNewFromListener_WhenTwoClosuresArePassed_ThenTwoGlobalListenersAreCreated(t *testing.T) {
	expectedCallCount := 4
	actualCallCount := 0

	subject := dispatcher.NewFromListener(
		listener.NewFromFunc(func(ID string, event string) {
			actualCallCount += 1
		}),
		listener.NewFromFunc(func(ID string, event string) {
			actualCallCount += 1
		}),
	)

	subject.Dispatch("testing.event1", "event")
	subject.Dispatch("testing.event2", "event")

	if expectedCallCount != actualCallCount {
		t.Errorf("expected to be called %#v, but was called %#v", expectedCallCount, actualCallCount)
	}
}

func Test_GivenCreatingADispatcherViaNewFromMapWithOneListener_WhenAnEventIsDispatchedThatMatches_ThenTheListenerIsCalled(t *testing.T) {
	expectedCallCount := 1
	actualCallCount := 0

	subject := dispatcher.NewFromMap(listener.Map[string, string]{

		"testing.event1": listener.NewFromFunc(func(id string, e string) {
			actualCallCount += 1
		}),
	})

	subject.Dispatch("testing.event1", "event")

	if expectedCallCount != actualCallCount {
		t.Errorf("expected to be called %#v, but was called %#v", expectedCallCount, actualCallCount)
	}
}

func Test_GivenCreatingADispatcherViaNewFromMapWithOneListener_WhenTwoEventsAreDispatchedThatMatch_ThenTheListenerIsCalled(t *testing.T) {
	expectedCallCount := 2
	actualCallCount := 0

	subject := dispatcher.NewFromMap(listener.Map[string, string]{

		"testing.event1": listener.NewFromFunc(func(id string, e string) {
			actualCallCount += 1
		}),
	})

	subject.Dispatch("testing.event1", "event1")
	subject.Dispatch("testing.event1", "event2")

	if expectedCallCount != actualCallCount {
		t.Errorf("expected to be called %#v, but was called %#v", expectedCallCount, actualCallCount)
	}
}

func Test_GivenCreatingADispatcherViaNewFromMapWithOneListener_WhenOneEventIsDispatchedThatDoesNotMatch_ThenTheListenerIsNotCalled(t *testing.T) {
	expectedCallCount := 0
	actualCallCount := 0

	subject := dispatcher.NewFromMap(listener.Map[string, string]{
		"testing.event1": listener.NewFromFunc(func(id string, e string) {
			actualCallCount += 1
		}),
	})

	subject.Dispatch("testing.event2", "event1")

	if expectedCallCount != actualCallCount {
		t.Errorf("expected to be called %#v, but was called %#v", expectedCallCount, actualCallCount)
	}
}

func Test_GivenCreatingADispatcherViaNewFromMapWithTwoListener_WhenOneEventIsDispatchedThatDoesNotMatch_ThenTheListenerIsNotCalled(t *testing.T) {
	expectedCallCount := 0
	actualCallCount := 0

	subject := dispatcher.NewFromMap(listener.Map[string, string]{
		"testing.event1": listener.NewFromFunc(func(id string, e string) {
			actualCallCount += 1
		}),

		"testing.event2": listener.NewFromFunc(func(id string, e string) {
			actualCallCount += 1
		}),
	})

	subject.Dispatch("testing.event0")

	if expectedCallCount != actualCallCount {
		t.Errorf("expected to be called %#v, but was called %#v", expectedCallCount, actualCallCount)
	}
}

func Test_GivenCreatingADispatcherViaNewFromMapWithTwoListener_WhenOneEventIsDispatchedThatMatchesTheFirst_ThenTheFirstIsCalled(t *testing.T) {
	expectedCallCount := 1
	actualCallCount := 0

	subject := dispatcher.NewFromMap(listener.Map[string, string]{
		"testing.event1": listener.NewFromFunc(func(id string, e string) {
			actualCallCount += 1
		}),

		"testing.event2": listener.NewFromFunc(func(id string, e string) {
			actualCallCount += 1
		}),
	})

	subject.Dispatch("testing.event1")

	if expectedCallCount != actualCallCount {
		t.Errorf("expected to be called %#v, but was called %#v", expectedCallCount, actualCallCount)
	}
}

func Test_GivenCreatingADispatcherViaNewFromMapWithTwoListener_WhenOneEventIsDispatchedThatMatchesTheSecond_ThenTheSecondIsCalled(t *testing.T) {
	expectedCallCount := 1
	actualCallCount := 0

	subject := dispatcher.NewFromMap(listener.Map[string, string]{
		"testing.event1": listener.NewFromFunc(func(id string, e string) {
			actualCallCount += 1
		}),

		"testing.event2": listener.NewFromFunc(func(id string, e string) {
			actualCallCount += 1
		}),
	})

	subject.Dispatch("testing.event2")

	if expectedCallCount != actualCallCount {
		t.Errorf("expected to be called %#v, but was called %#v", expectedCallCount, actualCallCount)
	}
}
