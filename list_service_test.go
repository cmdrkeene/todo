package todo

import "testing"

func TestListService(t *testing.T) {
	// given
	stream := newMemoryStream()
	service := newListService(stream)

	// when
	groceries := service.Create("groceries")
	milk := service.Add(groceries, "Milk")
	bread := service.Add(groceries, "Bread")
	eggs := service.Add(groceries, "Eggs")
	service.Check(groceries, milk)
	service.Uncheck(groceries, milk)
	service.Remove(groceries, bread)
	service.Rename(groceries, "Shopping List")
	// service.Delete(groceries)

	// then
	matchEvents(
		t,
		stream.Find(groceries),
		newListCreated(groceries, "groceries"),
		newItemAdded(groceries, milk, "Milk"),
		newItemAdded(groceries, bread, "Bread"),
		newItemAdded(groceries, eggs, "Eggs"),
		newItemChecked(groceries, milk),
		newItemUnchecked(groceries, milk),
		newItemRemoved(groceries, bread),
		newListRenamed(groceries, "Shopping List"),
		// newListDeleted(groceries),
	)
}
