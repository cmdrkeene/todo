package todo

import "net/http"

type listCommandHandler struct {
	service *listService
}

func newListCommandHandler(s *listService) *listCommandHandler {
	return &listCommandHandler{service: s}
}

func (h *listCommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h *listCommandHandler) add(w http.ResponseWriter, r *http.Request) {
	h.sendID(
		w,
		h.service.Add(
			h.listID(r),
			h.title(r),
		),
	)
}

func (h *listCommandHandler) create(w http.ResponseWriter, r *http.Request) {
	h.sendID(
		w,
		h.service.Create(
			h.name(r),
		),
	)
}

func (h *listCommandHandler) sendID(w http.ResponseWriter, id uuid) {
	w.Write([]byte(id))
}

func (h *listCommandHandler) remove(w http.ResponseWriter, r *http.Request) {
	list, item, valid := h.removeParams(r)
	if valid {
		h.service.Remove(list, item)
		ok(w)
	} else {
		badRequest(w)
	}
}

func (h *listCommandHandler) removeParams(r *http.Request) (list, item uuid, valid bool) {
	return h.listID(r),
		h.itemID(r),
		h.allPresent(h.listID(r), h.itemID(r))
}

func (h *listCommandHandler) allPresent(all ...uuid) bool {
	for _, u := range all {
		if u == "" {
			return false
		}
	}
	return true
}

func (h *listCommandHandler) rename(w http.ResponseWriter, r *http.Request) {
	h.service.Rename(h.listID(r), h.name(r))
	ok(w)
}

func (h *listCommandHandler) check(w http.ResponseWriter, r *http.Request) {
	h.service.Check(h.listID(r), h.itemID(r))
	ok(w)
}

func (h *listCommandHandler) uncheck(w http.ResponseWriter, r *http.Request) {
	h.service.Uncheck(h.listID(r), h.itemID(r))
	ok(w)
}

func (h *listCommandHandler) name(r *http.Request) string {
	const nameKey = "name"
	return r.FormValue(nameKey)
}

func (h *listCommandHandler) title(r *http.Request) string {
	const titleKey = "title"
	return r.FormValue(titleKey)
}

func (h *listCommandHandler) listID(r *http.Request) uuid {
	const listKey = "list"
	return uuid(r.FormValue(listKey))
}

func (h *listCommandHandler) itemID(r *http.Request) uuid {
	const itemKey = "item"
	return uuid(r.FormValue(itemKey))
}
