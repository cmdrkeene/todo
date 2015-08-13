package todo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkListHandlerCreate(b *testing.B) {
	b.ReportAllocs()
	handler := newListHandler(newListService(newMemoryStream()))
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/List/Create?name=bench", nil)
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(w, r)
	}
}
