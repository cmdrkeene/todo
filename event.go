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

type eventBus interface {
	Publish(event)
	Subscribe(event) chan event
}
