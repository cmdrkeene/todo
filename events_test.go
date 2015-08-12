package todo

import (
	"fmt"
	"testing"
)

// nerfEvent simulates a nonsensical and unhandleable event to aid coverage
type nerfEvent struct{}

// nerfCommand simulates a nonsensical and unhandleable command to aid coverage
type nerfCommand struct{}

func (nerfCommand) AggregateID() uuid {
	return "nerf"
}

func matchEvents(t *testing.T, got []event, want ...event) {
	if !eventsMatch(got, want...) {
		prettyGotWant(t, got, want)
	}
}

func prettyGotWant(t *testing.T, got, want interface{}) {
	gotWant(t, prettyEvents(got), prettyEvents(want))
}

func prettyEvents(events interface{}) string {
	pretty := "\n"
	for _, e := range toEventSlice(events) {
		pretty += fmt.Sprintf("%#v\n", e)
	}
	return pretty
}

func toEventSlice(v interface{}) []event {
	events, ok := v.([]event)
	if !ok {
		panic("prettyEvents passed something other then []event")
	}
	return events
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
	t.Errorf("got  %s", got)
	t.Errorf("want %s", want)
}
