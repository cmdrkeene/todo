package todo

import "time"

type uuid string

func newID() uuid {
	return "new"
}

type aggregate interface {
	Handle(command) []event
}

type command interface {
	AggregateID() uuid
}

type event interface{}

type eventStream interface {
	Append(uuid, ...event)
	Find(uuid) []event
}

type eventBus interface {
	Publish(event)
	Subscribe(event) chan event
}

type eventRecorder interface {
	FindRecords(uuid) []eventRecord
	FindRecordsByEventType(string) []eventRecord
	FindRecordsByEventID(uuid) []eventRecord
}

type eventRecord struct {
	aggregateID uuid
	data        event
	eventID     uuid
	eventType   string
	occurred    time.Time
}

type eventRecordScanner interface {
	Scan() chan eventRecord
}
