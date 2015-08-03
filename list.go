package todo

import "fmt"

type list struct {
	id    uuid
	name  string
	items []item
	state listState
}

type listState int

const (
	lsInitial listState = iota
	lsIncomplete
	lsCompleted
)

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
	case checkItem:
		return l.apply(command.Event())
	default:
		panic(fmt.Sprintf("unknown command %#v", command))
	}
}

func (l *list) apply(e event) []event {
	switch event := e.(type) {
	case listCreated:
		return l.applyCreated(event)
	case itemAdded:
		return l.applyItemAdded(event)
	case itemChecked:
		return l.applyItemChecked(event)
	default:
		panic(fmt.Sprintf("unknown event %#v", event))
	}
}

func (l *list) applyHistory(history []event) {
	for _, e := range history {
		l.apply(e)
	}
}

func (l *list) applyCreated(e listCreated) []event {
	if l.state != lsInitial {
		panic("list already created")
	}
	l.id = e.id
	l.name = e.name
	l.state = lsIncomplete
	return []event{e}
}

func (l *list) applyItemAdded(e itemAdded) []event {
	l.items = append(
		l.items,
		newItem(e.item, e.title),
	)
	return []event{e}
}

func (l *list) applyItemChecked(e itemChecked) []event {
	l.checkItem(e.item)
	return l.checkCompleted(e)
}

func (l *list) applyCompleted(e listCompleted) []event {
	l.state = lsCompleted
	return []event{e}
}

func (l *list) checkCompleted(e itemChecked) []event {
	if l.complete() {
		return []event{e, newListCompleted(l.id)}
	} else {
		return []event{e}
	}
}

func (l *list) checkItem(item uuid) {
	for i, current := range l.items {
		if current.id == item {
			if current.checked {
				panic("item already checked")
			}
			current.checked = true
			l.items[i] = current
		}
	}
}

func (l *list) complete() bool {
	for _, current := range l.items {
		if !current.checked {
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

type checkItem struct {
	list uuid
	item uuid
}

func (command checkItem) AggregateID() uuid {
	return command.list
}

func (command checkItem) Event() event {
	return newItemChecked(command.list, command.item)
}

func newCheckItem(list, item uuid) checkItem {
	return checkItem{
		item: item,
		list: list,
	}
}

type itemChecked struct {
	list uuid
	item uuid
}

func newItemChecked(list, item uuid) itemChecked {
	return itemChecked{
		item: item,
		list: list,
	}
}

type uncheckItem struct {
	list uuid
	item uuid
}

type itemUnchecked struct {
	list uuid
	item uuid
}

type removeItem struct {
	list uuid
	item uuid
}

type itemRemoved struct {
	list uuid
	item uuid
}

type listCompleted struct {
	list uuid
}

func newListCompleted(list uuid) listCompleted {
	return listCompleted{list: list}
}
