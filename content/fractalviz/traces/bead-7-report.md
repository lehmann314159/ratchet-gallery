# Bead 7: main

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 207s (3m)  
**Final exit criterion:** `go build ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** main  
**Output files:** main.go  
**Exit criteria:** `go build ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement the entry point in `main.go`.

1. Initialize `templates = InitTemplates()`.
2. Create the `SavedDir` directory using `os.MkdirAll(SavedDir, 0o755)`. Log fatal on error.
3. Configure a `http.ServeMux` with routes:
   - GET / -> HandleIndex
   - POST /render -> HandleRender
   - POST /select -> HandleSelect
   - GET /fractal.png -> HandleImage
   - POST /save -> HandleSave
   - GET /saved/ -> http.StripPrefix("/saved/", http.FileServer(http.Dir(SavedDir)))
4. Start the server with `http.ListenAndServe(":8080", mux)`.

### Revision 2 — created by REVISE_PENDING

**Title:** main  
**Output files:** main.go  
**Exit criteria:** `go build ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement the entry point in `main.go`.

1. Initialize `templates = InitTemplates()`.
2. Create the `SavedDir` directory using `os.MkdirAll(SavedDir, 0o755)`. Log fatal on error.
3. Configure a `http.ServeMux` with routes:
   - GET / -> `HandleIndex` (defined in handlers.go)
   - POST /render -> `HandleRender` (defined in handlers.go)
   - POST /select -> `HandleSelect` (defined in handlers.go)
   - GET /fractal.png -> `HandleImage` (defined in handlers.go)
   - POST /save -> `HandleSave` (defined in handlers.go)
   - GET /saved/ -> `http.StripPrefix("/saved/", http.FileServer(http.Dir(SavedDir)))`
4. Start the server with `http.ListenAndServe(":8080", mux)`. Ensure all handler functions from handlers.go are used directly.

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 7 | success | 207s | no fire | 1/1 | not run |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success):
No commands run. All tests passed. Main.go present. Undeclared files: do_not_use_this_test.go, go.mod, survey.md. Workspace note: discarded 1 file(s).

## Final Output Files

*State of output_files on disk at report time.*

### main.go

```go
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var templates *template.Template

func main() {
	templates = InitTemplates()
	if err := os.MkdirAll(SavedDir, 0o755); err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleIndex)
	mux.HandleFunc("/render", HandleRender)
	mux.HandleFunc("/select", HandleSelect)
	mux.HandleFunc("/fractal.png", HandleImage)
	mux.HandleFunc("/save", HandleSave)
	mux.Handle("/saved/", http.StripPrefix("/saved/", http.FileServer(http.Dir(SavedDir))))
	log.Fatal(http.ListenAndServe(":8080", mux))
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-7-attempt-1.log`*

```

Let's write again with var.

[tool: write_file map[content:package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var templates *template.Template

func main() {
	templates = InitTemplates()
	if err := os.MkdirAll(SavedDir, 0o755); err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleIndex)
	mux.HandleFunc("/render", HandleRender)
	mux.HandleFunc("/select", HandleSelect)
	mux.HandleFunc("/fractal.png", HandleImage)
	mux.HandleFunc("/save", HandleSave)
	mux.Handle("/saved/", http.StripPrefix("/saved/", http.FileServer(http.Dir(SavedDir))))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
 path:main.go]]
[result]
ok: wrote 581 bytes to main.go
[TURN 5]
Now build again.[tool: run_command map[command:go build ./...]]
[result]
exit: 0
[TURN 6]
Build passes.

Now check file exists.[tool: run_command map[command:ls main.go]]
[result]
stdout:
main.go

exit: 0
[TURN 7]
All good.

Exit criteria satisfied.

We should send final message.All exit criteria pass. The entry point is implemented in `main.go`:

- `templates` is initialized via `InitTemplates()`
- `SavedDir` is created with `os.MkdirAll(SavedDir, 0o755)` and fatal on error
- `http.ServeMux` routes are configured for `/`, `/render`, `/select`, `/fractal.png`, `/save`, and `/saved/`
- Server starts on `:8080`

`go build ./...` succeeds and `main.go` exists on disk.
[done — exit criteria already satisfied on disk; no write needed]
[workspace] discarded 1 file(s) written/modified outside output_files (not copied to project): fractalviz

```
