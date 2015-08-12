package todo

type memoryStream struct {
	stream map[uuid][]event
}

func newMemoryStream() *memoryStream {
	return &memoryStream{
		stream: make(map[uuid][]event),
	}
}

func (s *memoryStream) Append(id uuid, changes ...event) {
	s.stream[id] = append(s.stream[id], changes...)
}

func (s *memoryStream) Find(id uuid) []event {
	return s.stream[id]
}
