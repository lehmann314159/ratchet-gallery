package handlers

import "net/http"

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	projects, err := h.store.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, "index.html", map[string]any{"Projects": projects})
}

func (h *Handler) Project(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	m, err := h.store.Load(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	h.render(w, "project.html", map[string]any{
		"Project":   m,
		"DesignDoc": h.renderMarkdown(h.store.DocPath(slug, m.DesignDoc)),
		"Survey":    h.renderMarkdown(h.store.DocPath(slug, m.Survey)),
	})
}
