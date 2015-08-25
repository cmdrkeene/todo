package todo

import "os"

// fileEventStore stores and loads events from a file
// TODO lock stream while syncing?
type fileEventStore struct {
	path    string
	stream  eventStream
	records []eventRecord
}

func newFileEventStore(path string) *fileEventStore {
	return &fileEventStore{
		path:   path,
		stream: newMemoryStream(),
	}
}

func (s *fileEventStore) Append(id uuid, events ...event) {
	s.stream.Append(id, events...)
	s.appendRecords(id, events...)
}

func (s *fileEventStore) appendRecords(id uuid, events ...event) {
	for _, e := range events {
		s.records = append(s.records, newEventRecord(id, e))
	}
}

func (s *fileEventStore) Find(id uuid) []event {
	return s.stream.Find(id)
}

// Scan sends loaded records on new channel
func (s *fileEventStore) Scan() chan eventRecord {
	c := make(chan eventRecord, 32)
	for _, r := range s.records {
		c <- r
	}
	close(c)
	return c
}

// Load reads event records from path
func (s *fileEventStore) Load() error {
	file, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	// TODO serialization? json? gob?
}

// Sync writes event records to path
func (s *fileEventStore) Sync() error {
	file, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer file.Close()
}
