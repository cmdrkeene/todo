package todo

import "testing"

func TestListCheckItem(t *testing.T) {
	list := newList(
		newListCreated("a", "test"),
		newItemAdded("a", "b", "Buy Milk"),
	)
	command := newCheckItem(list.id, list.items[0].id)

	events := list.Handle(command)

	if !eventsMatch(events, command.Event()) {
		gotWant(t, events, command.Event())
	}
}

func TestListUncheckItem(t *testing.T) {}

func TestListAddItem(t *testing.T) {
	list := newList(
		newListCreated("a", "test"),
	)
	command := newAddItem(list.id, "Buy Milk")

	events := list.Handle(command)

	if !eventsMatch(events, command.Event()) {
		gotWant(t, events, command.Event())
	}
}

func TestListCreate(t *testing.T) {
	command := newCreateList("test")

	events := newList().Handle(command)

	if !eventsMatch(events, command.Event()) {
		gotWant(t, events, command.Event())
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

func gotWant(t *testing.T, got, want interface{}) {
	t.Error("got ", got)
	t.Error("want", want)
}
