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
	traces, err := h.store.Traces(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	headline, err := h.store.Headline(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.render(w, "project.html", map[string]any{
		"Project":   m,
		"DesignDoc": h.renderMarkdown(h.store.DocPath(slug, m.DesignDoc)),
		"Survey":    h.renderMarkdown(h.store.DocPath(slug, m.Survey)),
		"Traces":    traces,
		"Headline":  headline,
	})
}

// Trace renders one build-forensics file for a project — always Markdown
// today (Ratchet's bead reports), rendered the same way as the design doc.
func (h *Handler) Trace(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	file := r.PathValue("file")
	if _, err := h.store.Load(slug); err != nil {
		http.NotFound(w, r)
		return
	}
	path := h.store.TracePath(slug, file)
	content := h.renderMarkdown(path)
	if content == "" {
		http.NotFound(w, r)
		return
	}
	h.render(w, "trace.html", map[string]any{
		"Slug":    slug,
		"File":    file,
		"Content": content,
	})
}
