package todo

import "testing"

func TestListAddItem(t *testing.T) {
	list := newList(
		newListCreated("1", "test"),
	)
	command := newAddItem(list.id, "Buy Milk")

	events := list.Handle(command)

	if !eventsMatch(events, command.Event()) {
		wantGot(t, command.Event(), events)
	}
}

func TestListCreate(t *testing.T) {
	command := newCreateList("test")

	events := newList().Handle(command)

	if !eventsMatch(events, command.Event()) {
		wantGot(t, command.Event(), events)
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

func wantGot(t *testing.T, want, got interface{}) {
	t.Error("want", want)
	t.Error("got ", got)
}
