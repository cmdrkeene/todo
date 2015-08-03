package todo

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

type eventStore interface {
	History(uuid) []event
	Add(...event)
}

// eventStore interface is necessarily simple, but
// we must record timing and version information
// we could use a vector clock approach and panic when recorded version
// is equal to or greater than our newly incremented version
// it's simpler to keep timing off the raw event structs
// internal store should use a container

type eventBus interface {
	Publish(event)
	Subscribe(event) chan event
}
