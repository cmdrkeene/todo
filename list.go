package todo

import "fmt"

type list struct {
	id           uuid
	name         string
	items        *itemCollection
	stateMachine *listStateMachine
}

func newList(history ...event) *list {
	l := &list{}
	l.items = newItemCollection()
	l.stateMachine = newListStateMachine()
	l.applyHistory(history)
	return l
}

func (l *list) Handle(command command) (events []event) {
	return l.applyHistory(l.dispatch(command))
}

func (l *list) applyHistory(history []event) []event {
	for _, e := range history {
		l.apply(e)
	}
	return history
}

func (l *list) apply(e event) {
	switch event := e.(type) {
	case listCreated:
		l.applyCreated(event)
	case itemAdded:
		l.applyItemAdded(event)
	case itemChecked:
		l.applyItemChecked(event)
	case itemRemoved:
		l.applyItemRemoved(event)
	case itemUnchecked:
		l.applyItemUnchecked(event)
	case listCompleted:
		l.applyCompleted(event)
	case listEmptied:
		l.applyEmptied(event)
	case listUncompleted:
		l.applyUncompleted(event)
	case listRenamed:
		l.applyRenamed(event)
	default:
		panic(fmt.Sprintf("unknown event %#v", event))
	}
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
	case removeItem:
		return l.handleRemoveItem(command)
	case renameList:
		return l.handleRenameList(command)
	default:
		panic(fmt.Sprintf("unknown command %#v", command))
	}
}

func (l *list) applyRenamed(e listRenamed) {
	l.name = e.name
}

func (l *list) applyCreated(e listCreated) {
	l.id = e.id
	l.name = e.name
}

func (l *list) applyEmptied(e listEmptied) {
	if !l.items.Empty() {
		panic("items not empty")
	}

	l.stateMachine.mustTransition(lsEmpty)
}

func (l *list) applyItemAdded(e itemAdded) {
	l.items.Add(
		newItem(e.item, e.title),
	)

	if l.items.Size() == 1 {
		l.stateMachine.mustTransition(lsIncomplete)
	}
}

func (l *list) applyItemChecked(e itemChecked) {
	l.items.Set(e.item, true)
}

func (l *list) applyItemUnchecked(e itemUnchecked) {
	l.items.Set(e.item, false)
}

func (l *list) applyCompleted(e listCompleted) {
	l.stateMachine.mustTransition(lsCompleted)
}

func (l *list) applyUncompleted(e listUncompleted) {
	l.stateMachine.mustTransition(lsIncomplete)
}

func (l *list) applyItemRemoved(e itemRemoved) {
	l.items.Remove(e.item)
}

func (l *list) handleRenameList(command renameList) []event {
	return []event{
		newListRenamed(command.list, command.name),
	}
}

func (l *list) handleUncheckItem(command uncheckItem) []event {
	events := []event{
		newItemUnchecked(command.list, command.item),
	}

	if l.items.Size() == 1 {
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

	if l.items.Unchecked() == 1 {
		events = append(events, newListCompleted(command.list))
	}

	return events
}

func (l *list) handleRemoveItem(command removeItem) []event {
	events := []event{
		newItemRemoved(command.list, command.item),
	}
	if l.items.Size() == 1 {
		events = append(events, newListEmptied(command.list))
	}
	return events
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

type itemCollection struct {
	items []item
}

func newItemCollection() *itemCollection {
	return &itemCollection{}
}

func (c *itemCollection) Add(item item) {
	c.items = append(c.items, item)
}

func (c *itemCollection) Set(item uuid, value bool) {
	index := c.indexOf(item)
	if c.items[index].checked == value {
		panic(fmt.Sprintf("item already %v", value))
	}
	c.items[index].checked = value
}

func (c *itemCollection) Remove(item uuid) {
	index := c.indexOf(item)
	c.items = append(c.items[:index], c.items[index+1:]...)
}

func (c *itemCollection) Empty() bool {
	return c.Size() == 0
}

func (c *itemCollection) Size() int {
	return len(c.items)
}

func (c *itemCollection) Unchecked() int {
	var sum int
	for _, current := range c.items {
		if !current.checked {
			sum++
		}
	}
	return sum
}

func (c *itemCollection) indexOf(item uuid) int {
	for i, current := range c.items {
		if current.id == item {
			return i
		}
	}
	panic("item not in list")
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

func newRenameList(list uuid, name string) renameList {
	return renameList{
		list: list,
		name: name,
	}
}

func (command renameList) AggregateID() uuid {
	return command.list
}

type listRenamed struct {
	list uuid
	name string
}

func newListRenamed(list uuid, name string) listRenamed {
	return listRenamed{
		list: list,
		name: name,
	}
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

func newRemoveItem(list, item uuid) removeItem {
	return removeItem{
		list: list,
		item: item,
	}
}

func (command removeItem) AggregateID() uuid {
	return command.list
}

type itemRemoved struct {
	list uuid
	item uuid
}

func newItemRemoved(list, item uuid) itemRemoved {
	return itemRemoved{
		list: list,
		item: item,
	}
}

type listEmptied struct {
	list uuid
}

func newListEmptied(list uuid) listEmptied {
	return listEmptied{
		list: list,
	}
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
	lsEmpty listState = iota
	lsIncomplete
	lsCompleted
)

type listStateMachine struct {
	state listState
}

func newListStateMachine() *listStateMachine {
	return &listStateMachine{
		state: lsEmpty,
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
		return m.state == lsEmpty || m.state == lsCompleted
	case lsEmpty:
		return m.state == lsCompleted || m.state == lsIncomplete
	default:
		panic(fmt.Sprintln("unknown state", to))
	}
}

func (l *listStateMachine) mustTransition(to listState) {
	if !l.transition(to) {
		panic(fmt.Sprintf("failed transition from %#v to %#v", l.state, to))
	}
}
