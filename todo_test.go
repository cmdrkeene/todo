package todo

import "testing"

func TestListCreate(t *testing.T) {
	list := newList()
	command := newCreateList("test")
	events := list.Handle(command)

	if !eventsMatch(events, newListCreated(command)) {
		t.Error("events not equal")
	}
}

func eventsMatch(subject []event, expected ...event) bool {
	if len(subject) != len(expected) {
		return false
	}

	for i, e := range expected {
		if subject[i] != e {
			return false
		}
	}

	return true
}
