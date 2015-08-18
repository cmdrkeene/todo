package todo

import (
	"fmt"
	"time"
)

type memoryScanner struct {
	records  []eventRecord
	recorded time.Time
}

func newMemoryScanner() *memoryScanner {
	return &memoryScanner{
		recorded: time.Now(),
	}
}

func (s *memoryScanner) addEvents(events ...event) {
	for _, e := range events {
		s.records = append(s.records, eventRecord{
			occurred:  s.recorded,
			eventID:   newID(),
			eventType: fmt.Sprintf("%t", e),
			event:     e,
		})
	}
}

func (s *memoryScanner) Scan() chan eventRecord {
	c := make(chan eventRecord, 32)
	for _, r := range s.records {
		c <- r
	}
	close(c)
	return c
}
