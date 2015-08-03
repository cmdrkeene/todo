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
	name  string
	items []item
	//state created -> incomplete -> complete
}

func newList(history ...event) *list {
	l := &list{}
	l.applyHistory(history)
	return l
}

func (l *list) Handle(c command) []event {
	switch command := c.(type) {
	case createList:
		return l.apply(command.Event())
	case addItem:
		return l.apply(command.Event())
	default:
		panic("unknown command")
	}
}

func (l *list) apply(e event) []event {
	switch event := e.(type) {
	case listCreated:
		return l.applyCreated(event)
	case itemAdded:
		return l.applyItemAdded(event)
	default:
		panic("unknown event")
	}
}

func (l *list) applyHistory(history []event) {
	for _, e := range history {
		l.apply(e)
	}
}

func (l *list) applyCreated(e listCreated) []event {
	l.id = e.id
	l.name = e.name
	return []event{e}
}

func (l *list) applyItemAdded(e itemAdded) []event {
	l.items = append(
		l.items,
		newItem(e.item, e.title),
	)
	return []event{e}
}

func (l *list) complete() bool {
	for _, i := range l.items {
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

func newItem(id uuid, title string) item {
	return item{
		id:      id,
		title:   title,
		checked: false,
	}
}

type createList struct {
	id   uuid
	name string
}

func (c createList) AggregateID() uuid {
	return c.id
}

func (c createList) Event() event {
	return newListCreated(
		c.id,
		c.name,
	)
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
}

func newListCreated(id uuid, name string) listCreated {
	return listCreated{
		id:   id,
		name: name,
	}
}

type renameList struct {
	id   string
	name string
}

type listRenamed struct{}

type deleteList struct{}
type listDeleted struct{}

type addItem struct {
	list  uuid
	item  uuid
	title string
}

func (command addItem) AggregateID() uuid {
	return command.list
}

func (command addItem) Event() event {
	return newItemAdded(
		command.list,
		command.item,
		command.title,
	)
}

func newAddItem(list uuid, title string) addItem {
	return addItem{
		list:  list,
		item:  newID(),
		title: title,
	}
}

type itemAdded struct {
	list  uuid
	item  uuid
	title string
}

func newItemAdded(list, item uuid, title string) itemAdded {
	return itemAdded{
		list:  list,
		item:  item,
		title: title,
	}
}

type checkItem struct{}
type itemChecked struct{}

type uncheckItem struct{}
type itemUnchecked struct{}

type removeItem struct{}
type itemRemoved struct{}

type listCompleted struct{}
