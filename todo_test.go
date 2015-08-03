package todo

import "testing"

func TestListCheck(t *testing.T) {
	list := newList(
		newListCreated("a", "test"),
		newItemAdded("a", "b", "Buy Milk"),
	)
	command := newCheckItem(list.id, list.items[0].id)

	events := list.Handle(command)

	checked := command.Event()
	completed := newListCompleted(list.id)
	if !eventsMatch(events, checked, completed) {
		gotWant(t, events, command.Event())
	}

	if !list.complete() {
		t.Error("list should be complete")
	}
}

func TestListUncheck(t *testing.T) {}

func TestListAdd(t *testing.T) {
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
