package riot

import (
	"testing"
)

type LoadEvent struct {
	id string
}

type SaveEvent struct {
	id string
}

func TestObservable(t *testing.T) {
	sink := NewSink()

	// exercises:

	first_handler := 0

	listener_1 := sink.On(func(loaded LoadEvent) {
		t.Logf("LoadEvent handler #1: %s\n", loaded.id)
		first_handler += 1
	})

	second_handler := 0

	listener_2 := sink.On(func(loaded LoadEvent) {
		t.Logf("LoadEvent handler #2: %s\n", loaded.id)
		second_handler += 1
	})

	first_once_handler := 0

	sink.Once(func(saved SaveEvent) {
		t.Logf("SaveEvent once-handler #1: %s\n", saved.id)
		first_once_handler += 1
	})

	third_handler := 0

	sink.On(func(saved SaveEvent) {
		t.Logf("SaveEvent handler: %s\n", saved.id)
		third_handler += 1
	})

	sink.Send(LoadEvent{id: "LoadEvent#1"})
	sink.Send(SaveEvent{id: "SaveEvent#1"})
	sink.Send(SaveEvent{id: "SaveEvent#2"})

	second_once_handler := 0

	sink.Once(func(saved SaveEvent) {
		t.Logf("SaveEvent once-handler #2: %s\n", saved.id)
		second_once_handler += 1
	})

	sink.Send(SaveEvent{id: "SaveEvent#3"})

	sink.Off(listener_1)

	sink.Send(LoadEvent{id: "LoadEvent#2"}) // should go to LoadEvent handler #2 only

	sink.Off(listener_2)

	sink.Send(LoadEvent{id: "LoadEvent#3"}) // nobody is listening

	// assertions:

	if first_handler != 1 {
		t.Errorf("expected 1 call to LoadEvent handler #1, got %d", first_handler)
	}

	if second_handler != 2 {
		t.Errorf("expected 2 calls to LoadEvent handler #2, got %d", second_handler)
	}

	if first_once_handler != 1 {
		t.Errorf("expected 1 call to SaveEvent once-handler #1, got %d", first_once_handler)
	}

	if third_handler != 3 {
		t.Errorf("expected 3 calls to SaveEvent handler, got %d", third_handler)
	}

	if second_once_handler != 1 {
		t.Errorf("expected 1 call to SaveEvent once-handler #2, got %d", second_once_handler)
	}
}
