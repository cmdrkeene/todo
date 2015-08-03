package todo

import "testing"

func TestListCreate(t *testing.T) {
	list := newList()
	command := newCreateList("test")
	events := list.Handle(command)

	if len(events) != 1 {
		t.Error("want 1 event")
	}

	event, ok := events[0].(listCreated)
	if !ok {
		t.Error("want", listCreated{})
		t.Error("got ", event)
	}

	if event.name != "test" {
		t.Error("want", "test")
		t.Error("got ", event.name)
	}
}
