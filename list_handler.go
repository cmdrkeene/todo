package todo

import "net/http"

type listHandler struct {
	router  *http.ServeMux
	service *listService
}

func newListHandler(s *listService) *listHandler {
	return &listHandler{service: s}
}

func (h *listHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/List/Add":
		h.add(w, r)
	case "/List/Check":
		h.check(w, r)
	case "/List/Create":
		h.create(w, r)
	case "/List/Remove":
		h.remove(w, r)
	case "/List/Uncheck":
		h.uncheck(w, r)
	default:
		notFound(w)
		return
	}
}

func (h *listHandler) add(w http.ResponseWriter, r *http.Request) {
	h.sendID(
		w,
		h.service.Add(
			h.listID(r),
			h.title(r),
		),
	)
}

func (h *listHandler) create(w http.ResponseWriter, r *http.Request) {
	h.sendID(
		w,
		h.service.Create(
			h.name(r),
		),
	)
}

func (h *listHandler) sendID(w http.ResponseWriter, id uuid) {
	w.Write([]byte(id))
}

func (h *listHandler) remove(w http.ResponseWriter, r *http.Request) {
	list, item, valid := h.removeParams(r)
	if valid {
		h.service.Remove(list, item)
		ok(w)
	} else {
		badRequest(w)
	}
}

func (h *listHandler) removeParams(r *http.Request) (list, item uuid, valid bool) {
	return h.listID(r),
		h.itemID(r),
		h.allPresent(h.listID(r), h.itemID(r))
}

func (h *listHandler) allPresent(all ...uuid) bool {
	for _, u := range all {
		if u == "" {
			return false
		}
	}
	return true
}

func (h *listHandler) rename(w http.ResponseWriter, r *http.Request) {
	h.service.Rename(h.listID(r), h.name(r))
	ok(w)
}

func (h *listHandler) check(w http.ResponseWriter, r *http.Request) {
	h.service.Check(h.listID(r), h.itemID(r))
	ok(w)
}

func (h *listHandler) uncheck(w http.ResponseWriter, r *http.Request) {
	h.service.Uncheck(h.listID(r), h.itemID(r))
	ok(w)
}

func (h *listHandler) name(r *http.Request) string {
	const nameKey = "name"
	return r.FormValue(nameKey)
}

func (h *listHandler) title(r *http.Request) string {
	const titleKey = "title"
	return r.FormValue(titleKey)
}

func (h *listHandler) listID(r *http.Request) uuid {
	const listKey = "list"
	return uuid(r.FormValue(listKey))
}

func (h *listHandler) itemID(r *http.Request) uuid {
	const itemKey = "item"
	return uuid(r.FormValue(itemKey))
}
