package todo

import (
	"reflect"
	"testing"
	"time"
)

func TestListDirectory(t *testing.T) {
	// given
	directory := newListDirectory()
	now := time.Now()

	// when
	directory.Add(newListDirectoryEntry(now, "a", "test a"))
	directory.Add(newListDirectoryEntry(now, "b", "test b"))
	directory.Rename("b", "test b+")
	directory.Remove("a")

	// then
	mustDeepEqual(
		t,
		directory.entries,
		[]listDirectoryEntry{newListDirectoryEntry(now, "b", "test b+")},
	)
}

func TestListDirectoryWriterRestore(t *testing.T) {
	// given
	scanner := newMemoryScanner()
	scanner.addEvents(
		newListCreated("a", "test a"),
		newListCreated("b", "test b"),
		newListRenamed("b", "test b+"),
		newListDeleted("b"),
	)
	writer := newListDirectoryWriter(
		newListDirectory(),
		newMemoryBus(),
		scanner,
	)

	// when
	writer.Restore()

	// then
	want := newListDirectory()
	want.Add(newListDirectoryEntry(scanner.recorded, "a", "test a"))
	mustDeepEqual(
		t,
		writer.directory,
		want,
	)
}

func TestListDirectoryHandler(t *testing.T) {}
func TestListMemoryRepository(t *testing.T) {}

func mustDeepEqual(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		gotWant(t, got, want)
	}
}
