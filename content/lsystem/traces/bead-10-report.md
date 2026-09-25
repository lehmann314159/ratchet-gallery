# Bead 10: templates

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 620s (10m)  
**Final exit criterion:** `grep -q 'func TestTemplates' templates_test.go && go test -v -run TestTemplates ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** templates  
**Output files:** templates.go, templates_test.go  
**Exit criteria:** `grep -q 'func TestTemplates' templates_test.go && go test -v -run TestTemplates ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement `InitTemplates`, `RenderPage`, and `RenderResult` in `templates.go`.

1. `InitTemplates()`: Parse two inline Go string literals into a `*template.Template` named `"page"` and `"result"`. Panic on parse error.
2. `RenderPage(w, v)`: Execute the `"page"` template with `PageView`.
3. `RenderResult(w, v)`: Execute the `"result"` template with `PageView`.

Template `"result"` must be wrapped in `<div id="app">` and contain the form (with `<select>`, `<textarea>`, `<input>`, and buttons), the error message, the SVG visualization, and the gallery. Template `"page"` must be a full HTML document that includes `"result"` via `{{template "result" .}}` and loads HTMX from `https://unpkg.com/htmx.org@1.9.12`.

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 10 | success | 620s | no fire | 1/1 | not run |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success):
No commands run. All tests passed in post-execution run. TestTemplates not executed during this attempt.

## Final Output Files

*State of output_files on disk at report time.*

### templates.go

```go
package main

import (
	"html/template"
	"net/http"
)

func InitTemplates() *template.Template {
	templates = template.New("page")
	pageTmpl := `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>L-System Studio</title>
<script src="https://unpkg.com/htmx.org@1.9.12"></script>
</head>
<body>
{{template "result" .}}
</body>
</html>`
	resultTmpl := `<div id="app">
<form method="post" action="/render">
<select name="preset">
{{range .Presets}}<option value="{{.Name}}">{{.Label}}</option>{{end}}
</select>
<textarea name="source">{{.Source}}</textarea>
<input type="number" name="iterations" value="{{.Iterations}}">
<button type="submit">Render</button>
<button type="button">Save</button>
</form>
{{if .Error}}<div class="error">{{.Error}}</div>{{end}}
<div class="viz">{{.SVG}}</div>
<div class="gallery">
{{range .Saved}}<a href="/saved/{{.Filename}}">{{.Title}}</a>{{end}}
</div>
</div>`
	if _, err := templates.Parse(pageTmpl); err != nil {
		panic(err)
	}
	if _, err := templates.New("result").Parse(resultTmpl); err != nil {
		panic(err)
	}
	return templates
}

func RenderPage(w http.ResponseWriter, v PageView) {
	templates.ExecuteTemplate(w, "page", v)
}

func RenderResult(w http.ResponseWriter, v PageView) {
	templates.ExecuteTemplate(w, "result", v)
}

```

### templates_test.go

```go
package main

import (
	"html/template"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTemplates(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("InitTemplates panicked: %v", r)
		}
	}()
	tmpl := InitTemplates()
	if tmpl == nil {
		t.Fatalf("InitTemplates returned nil")
	}
	if tmpl.Lookup("page") == nil {
		t.Fatalf("page template missing")
	}
	if tmpl.Lookup("result") == nil {
		t.Fatalf("result template missing")
	}

	rec := httptest.NewRecorder()
	pv := PageView{
		Source:     "axiom: F",
		SVG:        template.HTML("<svg></svg>"),
		Iterations: 0,
		Presets:    []Preset{},
		Saved:      []SavedImage{},
		Error:      "",
	}
	RenderPage(rec, pv)
	body := rec.Body.String()
	if !strings.Contains(body, "https://unpkg.com/htmx.org@1.9.12") {
		t.Fatalf("page missing htmx load")
	}
	if !strings.Contains(body, "<html") {
		t.Fatalf("page not full html document")
	}
	if !strings.Contains(body, `id="app"`) {
		t.Fatalf("page missing app div from result inclusion")
	}

	rec2 := httptest.NewRecorder()
	pv2 := PageView{
		Source:     "test",
		SVG:        template.HTML("<svg>raw</svg>"),
		Iterations: 1,
		Presets: []Preset{
			{Name: "koch", Label: "Koch curve", Source: "axiom: F\n"},
		},
		Saved: []SavedImage{
			{Filename: "a.svg", Title: "a"},
		},
		Error: "boom",
	}
	RenderResult(rec2, pv2)
	out := rec2.Body.String()
	if !strings.Contains(out, `id="app"`) {
		t.Fatalf("result missing div id app wrapper")
	}
	if !strings.Contains(out, "<form") {
		t.Fatalf("result missing form")
	}
	if !strings.Contains(out, "<select") {
		t.Fatalf("result missing select")
	}
	if !strings.Contains(out, "<textarea") {
		t.Fatalf("result missing textarea")
	}
	if !strings.Contains(out, "<input") {
		t.Fatalf("result missing input")
	}
	if !strings.Contains(out, "<button") {
		t.Fatalf("result missing buttons")
	}
	if !strings.Contains(out, "boom") {
		t.Fatalf("error message not rendered")
	}
	if !strings.Contains(out, "<svg>raw</svg>") {
		t.Fatalf("SVG not rendered raw, likely escaped")
	}
	if !strings.Contains(out, `class="viz"`) {
		t.Fatalf("result missing viz div")
	}
	if !strings.Contains(out, "a.svg") {
		t.Fatalf("gallery not rendered")
	}
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-10-attempt-1.log`*

```
	return templates
}

func RenderPage(w http.ResponseWriter, v PageView) {
	templates.ExecuteTemplate(w, "page", v)
}

func RenderResult(w http.ResponseWriter, v PageView) {
	templates.ExecuteTemplate(w, "result", v)
}
 path:templates.go]]
[result]
ok: wrote 1277 bytes to templates.go
[TURN 5]
Now run test again.[tool: run_command map[command:grep -q 'func TestTemplates' templates_test.go && go test -v -run TestTemplates ./...]]
[result]
stdout:
=== RUN   TestTemplates
--- PASS: TestTemplates (0.00s)
PASS
ok  	lsystem	0.461s

exit: 0
[TURN 6]
Test passes.

Now check exit criteria: each check passes AND every Output File exists on disk. Output file templates.go exists. Let's verify ls.[tool: run_command map[command:ls templates.go]]
[result]
stdout:
templates.go

exit: 0
[TURN 7]
All good.

We should ensure no other files created. Output files list only templates.go. Good.

Now final message.

All exit criteria passed.

The implementation parses the two inline templates, panics on parse errors, and renders both page and result views with the required structure.

```bash
grep -q 'func TestTemplates' templates_test.go && go test -v -run TestTemplates ./...
```
```
=== RUN   TestTemplates
--- PASS: TestTemplates (0.00s)
PASS
ok  	lsystem	0.461s
```

Output file exists:
```
templates.go
```

[done — exit criteria already satisfied on disk; no write needed]

```
