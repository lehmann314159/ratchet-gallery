package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Proxy forwards everything under /live/{slug}/ to that project's container
// unmodified — the backend app owns its own prefix-aware routes/HTML (its
// routePrefix must equal "/live/<slug>", baked in at onboarding time), since
// neither Caddy nor this proxy can rewrite paths inside an HTML response.
func (h *Handler) Proxy(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	m, err := h.store.Load(slug)
	if err != nil || m.Backend == "" {
		http.NotFound(w, r)
		return
	}
	target := &url.URL{Scheme: "http", Host: m.Backend}
	httputil.NewSingleHostReverseProxy(target).ServeHTTP(w, r)
}
