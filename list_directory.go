package todo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// listDirectory is a Data Transfer Object that enumerates lists
type listDirectory struct {
	entries []listDirectoryEntry
}

func newListDirectory(entries ...listDirectoryEntry) *listDirectory {
	return &listDirectory{
		entries: entries,
	}
}

type listDirectoryEntry struct {
	Created time.Time
	ID      uuid
	Name    string
}

func newListDirectoryEntry(created time.Time, id uuid, name string) listDirectoryEntry {
	return listDirectoryEntry{
		Created: created,
		ID:      id,
		Name:    name,
	}
}

func (d *listDirectory) Add(entries ...listDirectoryEntry) {
	d.entries = append(d.entries, entries...)
}

func (d *listDirectory) Remove(list uuid) {
	d.removeAtIndex(d.index(list))
}

func (d *listDirectory) removeAtIndex(index int) {
	d.entries = append(
		d.entries[:index],
		d.entries[index+1:]...,
	)
}

func (d *listDirectory) index(list uuid) int {
	for i, l := range d.entries {
		if l.ID == list {
			return i
		}
	}
	panic(
		fmt.Sprintf("list %s not in directory %#v", list, d.entries),
	)
}

func (d *listDirectory) Rename(list uuid, name string) {
	d.entries[d.index(list)].Name = name
}

type listDirectoryWriter struct {
	bus       eventBus
	directory *listDirectory
	scanner   eventScanner
}

func newListDirectoryWriter(d *listDirectory, b eventBus, s eventScanner) *listDirectoryWriter {
	return &listDirectoryWriter{
		bus:       b,
		directory: d,
		scanner:   s,
	}
}

func (l *listDirectoryWriter) Restore() {
	for record := range l.scanner.Scan() {
		l.handle(record)
	}
}

func (l *listDirectoryWriter) Subscribe() {
	for record := range l.bus.Subscribe() {
		l.handle(record)
	}
}

func (l *listDirectoryWriter) handle(record eventRecord) {
	switch event := record.event.(type) {
	case listCreated:
		l.directory.Add(l.newEntry(record, event))
	case listRenamed:
		l.directory.Rename(event.list, event.name)
	case listDeleted:
		l.directory.Remove(event.list)
	default:
		// ignore
	}
}

func (l *listDirectoryWriter) newEntry(record eventRecord, event listCreated) listDirectoryEntry {
	return newListDirectoryEntry(
		record.occurred,
		event.id,
		event.name,
	)
}

type listDirectoryHandler struct {
	directory *listDirectory
}

func (l *listDirectoryHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	l.encode(w, l.directory)
}

func (l *listDirectoryHandler) encode(w http.ResponseWriter, d *listDirectory) {
	panicOnError(
		json.NewEncoder(w).Encode(d),
	)
}
