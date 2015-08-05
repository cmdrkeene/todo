package todo

import "testing"

func TestListUncheck(t *testing.T) {
	// given
	list := newList(
		newListCreated("a", "test"),
		newItemAdded("a", "b", "Buy Milk"),
		newItemChecked("a", "b"),
		newListCompleted("a"),
	)

	// when
	events := list.Handle(
		newUncheckItem("a", "b"),
	)

	// then
	matchEvents(
		t,
		events,
		newItemUnchecked("a", "b"),
		newListUncompleted("a"),
	)
}

func TestListCheck(t *testing.T) {
	// given
	list := newList(
		newListCreated("a", "test"),
		newItemAdded("a", "b", "Buy Milk"),
	)

	// when
	events := list.Handle(
		newCheckItem("a", "b"),
	)

	// then
	matchEvents(
		t,
		events,
		newItemChecked("a", "b"),
		newListCompleted("a"),
	)
}

func TestListAdd(t *testing.T) {
	// given
	list := newList(
		newListCreated("a", "test"),
	)

	// when
	events := list.Handle(
		newAddItem("a", "b", "Buy Milk"),
	)

	// then
	matchEvents(
		t,
		events,
		newItemAdded("a", "b", "Buy Milk"),
	)
}

func TestListCreate(t *testing.T) {
	// given
	list := newList()

	// when
	events := list.Handle(
		newCreateList("a", "test"),
	)

	// then
	matchEvents(
		t,
		events,
		newListCreated("a", "test"),
	)
}

func matchEvents(t *testing.T, got []event, want ...event) {
	if !eventsMatch(got, want...) {
		gotWant(t, got, want)
	}
}

func eventsMatch(got []event, want ...event) bool {
	if len(got) != len(want) {
		return false
	}

	for i, e := range want {
		if got[i] != e {
			return false
		}
	}

	return true
}

func gotWant(t *testing.T, got, want interface{}) {
	t.Errorf("got  %#v", got)
	t.Errorf("want %#v", want)
}
