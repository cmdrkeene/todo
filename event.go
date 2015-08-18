package todo

import "time"

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

// eventRecord is a storage mechanism for high level events
type eventRecord struct {
	aggregateID uuid
	event       event
	eventID     uuid
	eventType   string
	occurred    time.Time
}

type eventBus interface {
	Publish(eventRecord)
	Subscribe() chan eventRecord
}

type eventRecorder interface {
	// FindRecords(uuid) []eventRecord
	FindRecordsByEventType(string) []eventRecord
	FindRecordsByEventID(uuid) []eventRecord
}

type eventScanner interface {
	Scan() chan eventRecord
}
