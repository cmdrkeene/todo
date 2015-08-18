package todo

type memoryBus struct{}

func newMemoryBus() *memoryBus {
	return &memoryBus{}
}

func (b *memoryBus) Publish(r eventRecord) {
}

func (b *memoryBus) Subscribe() chan eventRecord {
	return make(chan eventRecord)
}
