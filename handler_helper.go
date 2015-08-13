package todo

import "net/http"

const get = "GET"
const post = "POST"

func methods(fn http.HandlerFunc, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, m := range methods {
			if r.Method != m {
				methodNotAllowed(w)
			}
		}
		fn(w, r)
	}
}

func ok(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func badRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func methodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
