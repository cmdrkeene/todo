package todo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListCommandHandlerCreate(t *testing.T) {
	// given
	stream := newMemoryStream()
	handler := newListCommandHandler(newListService(stream))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/List/Create?name=test", nil)

	// when
	handler.ServeHTTP(w, r)

	// then
	if w.Code != 200 {
		gotWant(t, w.Code, 200)
	}
	list := uuid(w.Body.String())
	matchEvents(
		t,
		stream.Find(list),
		newListCreated(list, "test"),
	)
}
