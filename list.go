package todo

import "fmt"

type list struct {
	id           uuid
	name         string
	items        []item
	stateMachine *listStateMachine
}

func newList(history ...event) *list {
	l := &list{}
	l.stateMachine = newListStateMachine()
	l.applyHistory(history)
	return l
}

func (l *list) Handle(command command) (events []event) {
	return l.applyHistory(l.dispatch(command))
}

func (l *list) dispatch(c command) []event {
	switch command := c.(type) {
	case addItem:
		return l.handleAddItem(command)
	case checkItem:
		return l.handleCheckItem(command)
	case createList:
		return l.handleCreateList(command)
	case uncheckItem:
		return l.handleUncheckItem(command)
	default:
		panic(fmt.Sprintf("unknown command %#v", command))
	}
}

func (l *list) handleUncheckItem(command uncheckItem) []event {
	events := []event{
		newItemUnchecked(command.list, command.item),
	}
	if l.complete() {
		events = append(events, newListUncompleted(command.list))
	}
	return events
}

func (l *list) handleCreateList(command createList) []event {
	return []event{
		newListCreated(command.id, command.name),
	}
}

func (l *list) handleAddItem(command addItem) []event {
	return []event{
		newItemAdded(command.list, command.item, command.title),
	}
}

func (l *list) handleCheckItem(command checkItem) []event {
	events := []event{
		newItemChecked(command.list, command.item),
	}
	if l.lastUnchecked() {
		events = append(events, newListCompleted(command.list))
	}
	return events
}

func (l *list) apply(e event) {
	switch event := e.(type) {
	case listCreated:
		l.applyCreated(event)
	case itemAdded:
		l.applyItemAdded(event)
	case itemChecked:
		l.applyItemChecked(event)
	case itemUnchecked:
		l.applyItemUnchecked(event)
	case listCompleted:
		l.applyCompleted(event)
	case listUncompleted:
		l.applyUncompleted(event)
	default:
		panic(fmt.Sprintf("unknown event %#v", event))
	}
}

func (l *list) applyHistory(history []event) []event {
	for _, e := range history {
		l.apply(e)
	}
	return history
}

func (l *list) applyCreated(e listCreated) {
	l.id = e.id
	l.name = e.name
}

func (l *list) applyItemAdded(e itemAdded) {
	l.stateMachine.mustTransition(lsIncomplete)
	l.items = append(
		l.items,
		newItem(e.item, e.title),
	)
}

func (l *list) applyItemChecked(e itemChecked) {
	l.checkItem(e.item)
}

func (l *list) applyItemUnchecked(e itemUnchecked) {
	l.uncheckItem(e.item)
}

func (l *list) applyCompleted(e listCompleted) {
	l.stateMachine.mustTransition(lsCompleted)
}

func (l *list) applyUncompleted(e listUncompleted) {
	l.stateMachine.mustTransition(lsIncomplete)
}

func (l *list) checkItem(item uuid) {
	l.setItem(item, true)
}

func (l *list) uncheckItem(item uuid) {
	l.setItem(item, false)
}

func (l *list) setItem(item uuid, value bool) {
	for i, current := range l.items {
		if current.id == item {
			if current.checked == value {
				panic(fmt.Sprintf("item already %v", value))
			}
		}
		l.items[i].checked = value
	}
}

func (l *list) complete() bool {
	return l.numUnchecked() == 0
}

func (l *list) lastUnchecked() bool {
	return l.numUnchecked() == 1
}

func (l *list) numUnchecked() int {
	var sum int
	for _, current := range l.items {
		if !current.checked {
			sum++
		}
	}
	return sum
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

func newCreateList(id uuid, name string) createList {
	if name == "" {
		panic("name can't be blank")
	}

	return createList{
		id:   id,
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
	list uuid
	name string
}

type listRenamed struct {
	list uuid
	name string
}

type deleteList struct {
	list uuid
}

type listDeleted struct {
	list uuid
}

type addItem struct {
	list  uuid
	item  uuid
	title string
}

func (command addItem) AggregateID() uuid {
	return command.list
}

func newAddItem(list, item uuid, title string) addItem {
	return addItem{
		list:  list,
		item:  item,
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

func newUncheckItem(list, item uuid) uncheckItem {
	return uncheckItem{
		list: list,
		item: item,
	}
}

func (command uncheckItem) AggregateID() uuid {
	return command.list
}

type itemUnchecked struct {
	list uuid
	item uuid
}

func newItemUnchecked(list, item uuid) itemUnchecked {
	return itemUnchecked{
		list: list,
		item: item,
	}
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

type listUncompleted struct {
	list uuid
}

func newListUncompleted(list uuid) listUncompleted {
	return listUncompleted{list: list}
}

type listState int

const (
	lsInitial listState = iota
	lsIncomplete
	lsCompleted
)

type listStateMachine struct {
	state listState
}

func newListStateMachine() *listStateMachine {
	return &listStateMachine{
		state: lsInitial,
	}
}

func (m *listStateMachine) transition(to listState) bool {
	if m.canTransition(to) {
		m.state = to
		return true
	} else {
		return false
	}
}

func (m *listStateMachine) canTransition(to listState) bool {
	switch to {
	case lsCompleted:
		return m.state == lsIncomplete
	case lsIncomplete:
		return m.state == lsInitial || m.state == lsCompleted
	case lsInitial:
		return false
	default:
		panic(fmt.Sprintln("unknown state", to))
	}
}

func (l *listStateMachine) mustTransition(to listState) {
	if !l.transition(to) {
		panic(fmt.Sprintf("failed transition from %#v to %#v", l.state, to))
	}
}
