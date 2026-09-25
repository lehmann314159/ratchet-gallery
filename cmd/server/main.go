package main

import (
	"log"
	"net/http"
	"os"

	"ratchet-gallery/internal/content"
	"ratchet-gallery/internal/handlers"
)

func main() {
	contentDir := os.Getenv("CONTENT_DIR")
	if contentDir == "" {
		contentDir = "content"
	}
	store := content.NewStore(contentDir)
	h := handlers.New(store, "templates")

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /{$}", h.Index)
	mux.HandleFunc("GET /p/{slug}", h.Project)
	mux.HandleFunc("GET /p/{slug}/traces/{file}", h.Trace)
	mux.HandleFunc("/live/{slug}/", h.Proxy)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8070"
	}
	log.Printf("ratchet-gallery listening on :%s (content dir: %s)", port, contentDir)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
