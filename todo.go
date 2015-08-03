package todo

import "time"

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

func (s *listService) get(list uuid) *list {
	return nil
}

func (s *listService) ChangeName(list uuid, name string) error {
	return nil
}

func (s *listService) AddItem(list uuid, title string) (uuid, error) {
	return "", nil
}

func (s *listService) RemoveItem(list, item uuid) error {
	return nil
}

func (s *listService) CheckItem(list, item uuid) error {
	return nil
}

func (s *listService) UncheckItem(list, item uuid) error {
	return nil
}

type list struct {
	id    uuid
	state struct {
		name  string
		items []item
	}
	// created -> incomplete -> complete
}

func newList() *list {
	return &list{}
}

func (l *list) Handle(c command) []event {
	switch command := c.(type) {
	case createList:
		return []event{newListCreated(command)}
	default:
		panic("unknown command")
	}
}
func (l *list) applyHistory([]event) {}
func (l *list) apply(event) []event  { return nil }

func (l *list) complete() bool {
	for _, i := range l.state.items {
		if !i.checked {
			return false
		}
	}
	return true
}

type item struct {
	id      uuid
	title   string
	checked bool
}

type createList struct {
	id   uuid
	name string
}

func (c createList) AggregateID() uuid {
	return c.id
}

func newCreateList(name string) createList {
	if name == "" {
		panic("name can't be blank")
	}
	return createList{
		id:   newID(),
		name: name,
	}
}

type listCreated struct {
	id   uuid
	name string
	when time.Time
}

func newListCreated(command createList) listCreated {
	return listCreated{
		id:   command.id,
		name: command.name,
		when: time.Now(),
	}
}

type renameList struct {
	id   string
	name string
}

type listRenamed struct{}

type deleteList struct{}
type listDeleted struct{}

type addItem struct{}
type itemAdded struct{}

type checkItem struct{}
type itemChecked struct{}

type uncheckItem struct{}
type itemUnchecked struct{}

type removeItem struct{}
type itemRemoved struct{}

type listCompleted struct{}
