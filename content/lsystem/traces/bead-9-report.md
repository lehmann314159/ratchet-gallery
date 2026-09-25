# Bead 9: save

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 126s (2m)  
**Final exit criterion:** `grep -q 'func TestSave' save_test.go && go test -v -run TestSave ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** save  
**Output files:** save.go, save_test.go  
**Exit criteria:** `grep -q 'func TestSave' save_test.go && go test -v -run TestSave ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement `SaveSVG` and `ListSaved` in `save.go`.

1. `SaveSVG(svg, dir, now)`: Hold `saveMu`. Generate filename `fmt.Sprintf("%019d.svg", now.UnixNano())`. Write SVG bytes to `filepath.Join(dir, name)`. Return filename or error.
2. `ListSaved(dir)`: Read directory. Collect all `*.svg` files as `SavedImage{Filename: name, Title: name}`. Sort them by filename in descending order (newest first). Return `(nil, nil)` if directory is missing.

Pin — `save` bead, `SaveSVG` / `ListSaved`: `SaveSVG("<svg/>", tmpDir, time.Unix(0, 1700000000123456789))` → `("1700000000123456789.svg", nil)` and that file exists in `tmpDir` with body `<svg/>`. After writing `0000000000000000001.svg` and `0000000000000000002.svg` into a temp dir, `ListSaved` returns them in the order `["0000000000000000002.svg", "0000000000000000001.svg"]` and ignores any non-`.svg` file. `ListSaved` on a missing path → `(nil, nil)`.

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `save` bead, `SaveSVG` / `ListSaved`:**
  `SaveSVG("<svg/>", tmpDir, time.Unix(0, 1700000000123456789))` →
  `("1700000000123456789.svg", nil)` and that file exists in `tmpDir` with body
  `<svg/>`. After writing `0000000000000000001.svg` and `0000000000000000002.svg` into a
  temp dir, `ListSaved` returns them in the order
  `["0000000000000000002.svg", "0000000000000000001.svg"]` and ignores any non-`.svg`
  file. `ListSaved` on a missing path → `(nil, nil)`.

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 9 | success | 126s | no fire | 0/0 | not run |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success):
No commands were run.
All tests passed.


## Final Output Files

*State of output_files on disk at report time.*

### save.go

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

const SavedDir = "saved"

var saveMu sync.Mutex

type SavedImage struct {
	Filename string
	Title    string
}

func SaveSVG(svg string, dir string, now time.Time) (string, error) {
	saveMu.Lock()
	defer saveMu.Unlock()
	name := fmt.Sprintf("%019d.svg", now.UnixNano())
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(svg), 0644); err != nil {
		return "", err
	}
	return name, nil
}

func ListSaved(dir string) ([]SavedImage, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var imgs []SavedImage
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if len(name) >= 4 && name[len(name)-4:] == ".svg" {
			imgs = append(imgs, SavedImage{Filename: name, Title: name})
		}
	}
	sort.Slice(imgs, func(i, j int) bool {
		return imgs[i].Filename > imgs[j].Filename
	})
	return imgs, nil
}

```

### save_test.go

```go
package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSave(t *testing.T) {
	t.Run("SaveSVG writes file with correct name and body", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "save_test")
		if err != nil {
			t.Fatalf("MkdirTemp: %v", err)
		}
		defer os.RemoveAll(dir)

		svg := "<svg/>"
		now := time.Unix(0, 1700000000123456789)
		name, err := SaveSVG(svg, dir, now)
		if err != nil {
			t.Fatalf("SaveSVG error: %v", err)
		}
		if name != "1700000000123456789.svg" {
			t.Fatalf("expected filename 1700000000123456789.svg, got %s", name)
		}
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("ReadFile error: %v", err)
		}
		if string(data) != svg {
			t.Fatalf("file body mismatch: got %q want %q", string(data), svg)
		}
	})

	t.Run("ListSaved returns descending order and ignores non-svg", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "save_test")
		if err != nil {
			t.Fatalf("MkdirTemp: %v", err)
		}
		defer os.RemoveAll(dir)

		f1 := "0000000000000000001.svg"
		f2 := "0000000000000000002.svg"
		if err := os.WriteFile(filepath.Join(dir, f1), []byte{}, 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, f2), []byte{}, 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "readme.txt"), []byte{}, 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}

		imgs, err := ListSaved(dir)
		if err != nil {
			t.Fatalf("ListSaved error: %v", err)
		}
		if len(imgs) != 2 {
			t.Fatalf("expected 2 images, got %d", len(imgs))
		}
		if imgs[0].Filename != f2 || imgs[0].Title != f2 {
			t.Fatalf("first image expected %s, got %+v", f2, imgs[0])
		}
		if imgs[1].Filename != f1 || imgs[1].Title != f1 {
			t.Fatalf("second image expected %s, got %+v", f1, imgs[1])
		}
	})

	t.Run("ListSaved missing dir returns nil nil", func(t *testing.T) {
		missing := "/tmp/save_missing_dir_1234567890"
		os.RemoveAll(missing)
		imgs, err := ListSaved(missing)
		if err != nil {
			t.Fatalf("expected nil error for missing dir, got %v", err)
		}
		if imgs != nil {
			t.Fatalf("expected nil slice for missing dir, got %v", imgs)
		}
	})
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-9-attempt-1.log`*

```
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if len(name) >= 4 && name[len(name)-4:] == ".svg" {
			imgs = append(imgs, SavedImage{Filename: name, Title: name})
		}
	}
	sort.Slice(imgs, func(i, j int) bool {
		return imgs[i].Filename > imgs[j].Filename
	})
	return imgs, nil
}
 path:save.go]]
[result]
ok: wrote 996 bytes to save.go
[TURN 4]
Now verify exit criteria: grep -q 'func TestSave' save_test.go && go test -v -run TestSave ./...

Run command.[tool: run_command map[command:grep -q 'func TestSave' save_test.go && go test -v -run TestSave ./...]]
[result]
stdout:
=== RUN   TestSave
=== RUN   TestSave/SaveSVG_writes_file_with_correct_name_and_body
=== RUN   TestSave/ListSaved_returns_descending_order_and_ignores_non-svg
=== RUN   TestSave/ListSaved_missing_dir_returns_nil_nil
--- PASS: TestSave (0.00s)
    --- PASS: TestSave/SaveSVG_writes_file_with_correct_name_and_body (0.00s)
    --- PASS: TestSave/ListSaved_returns_descending_order_and_ignores_non-svg (0.00s)
    --- PASS: TestSave/ListSaved_missing_dir_returns_nil_nil (0.00s)
PASS
ok  	lsystem	0.448s

exit: 0
[TURN 5]
Tests pass. Need confirm file exists.[tool: run_command map[command:ls save.go]]
[result]
stdout:
save.go

exit: 0
[TURN 6]
All good.All exit criteria pass.

- `save.go` exists and implements `SaveSVG` with mutex protection, zero-padded 19-digit filename, and file write.
- `ListSaved` reads the directory, filters `*.svg`, returns `SavedImage` with `Filename`/`Title`, sorts descending, and returns `(nil, nil)` for missing paths.

Tests:
```
=== RUN   TestSave
=== RUN   TestSave/SaveSVG_writes_file_with_correct_name_and_body
=== RUN   TestSave/ListSaved_returns_descending_order_and_ignores_non-svg
=== RUN   TestSave/ListSaved_missing_dir_returns_nil_nil
--- PASS: TestSave
PASS
```
Output file verified: `save.go` present.
[done — exit criteria already satisfied on disk; no write needed]

```
