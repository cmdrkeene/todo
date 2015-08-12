package todo

// listService issues commands against an event stream
type listService struct {
	stream eventStream
}

func newListService(s eventStream) *listService {
	return &listService{stream: s}
}

func (s *listService) Create(name string) uuid {
	id := newID()
	s.issue(newCreateList(id, name))
	return id
}

func (s *listService) Rename(list uuid, newName string) {
	s.issue(newRenameList(list, newName))
}

func (s *listService) Add(list uuid, title string) uuid {
	id := newID()
	s.issue(newAddItem(list, id, title))
	return id
}

func (s *listService) Remove(list, item uuid) {
	s.issue(newRemoveItem(list, item))
}

func (s *listService) Check(list, item uuid) {
	s.issue(newCheckItem(list, item))
}

func (s *listService) Uncheck(list, item uuid) {
	s.issue(newUncheckItem(list, item))
}

func (s *listService) issue(cmd command) {
	s.save(s.dispatch(cmd))
}

func (s *listService) dispatch(cmd command) (uuid, []event) {
	switch cmd.(type) {
	case createList:
		return s.handle(newList(), cmd)
	default:
		return s.handle(s.existingList(cmd.AggregateID()), cmd)
	}
}

func (s *listService) handle(list *list, cmd command) (uuid, []event) {
	return list.id, list.Handle(cmd)
}

func (s *listService) existingList(id uuid) *list {
	return newList(s.stream.Find(id)...)
}

func (s *listService) save(list uuid, changes []event) {
	s.stream.Append(list, changes...)
}
