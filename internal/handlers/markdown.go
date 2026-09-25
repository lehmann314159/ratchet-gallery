package handlers

import (
	"bytes"
	"html/template"
	"os"

	"github.com/yuin/goldmark"
)

// renderMarkdown reads and converts a design-doc/survey file. Missing or
// unreadable files render as empty rather than failing the whole page —
// not every onboarded project will have both docs.
func (h *Handler) renderMarkdown(path string) template.HTML {
	if path == "" {
		return ""
	}
	src, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	if err := goldmark.Convert(src, &buf); err != nil {
		return ""
	}
	return template.HTML(buf.String())
}
