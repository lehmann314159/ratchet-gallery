package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"ratchet-gallery/internal/content"
)

type Handler struct {
	store     *content.Store
	templates map[string]*template.Template
}

func New(store *content.Store, templateDir string) *Handler {
	layoutFiles, _ := filepath.Glob(filepath.Join(templateDir, "layouts", "*.html"))
	pageFiles, _ := filepath.Glob(filepath.Join(templateDir, "*.html"))

	templates := make(map[string]*template.Template)
	for _, page := range pageFiles {
		name := filepath.Base(page)
		files := append([]string{page}, layoutFiles...)
		templates[name] = template.Must(template.New(name).ParseFiles(files...))
	}

	return &Handler{store: store, templates: templates}
}

func (h *Handler) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, ok := h.templates[name]
	if !ok {
		http.Error(w, "template not found: "+name, http.StatusInternalServerError)
		return
	}
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
