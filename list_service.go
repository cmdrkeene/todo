package todo

// receives Commands and routes it to the appropriate aggregate
type listService struct {
	store eventStore
}

func (s *listService) Create(name string) (uuid, error) {
	list := &list{}
	s.store.Add(
		list.Handle(
			&createList{
				name: name,
				id:   newID(),
			},
		),
	)
	return list.id, nil
}

func (s *listService) Rename(list uuid, name string) error {
	return nil
}

func (s *listService) Add(list uuid, title string) (uuid, error) {
	return "", nil
}

func (s *listService) Remove(list, item uuid) error {
	return nil
}

func (s *listService) Check(list, item uuid) error {
	return nil
}

func (s *listService) Uncheck(list, item uuid) error {
	return nil
}
