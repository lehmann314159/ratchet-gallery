# Bead 4: save

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 187s (3m)  
**Final exit criterion:** `grep -q 'func TestSaveImage' save_test.go && go test -v -run TestSaveImage ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** save  
**Output files:** save.go, save_test.go  
**Exit criteria:** `grep -q 'func TestSaveImage' save_test.go && go test -v -run TestSaveImage ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement image saving in `save.go`.

1. Implement `SaveImage(p Params, dir string, now time.Time) (string, error)`:
   - Lock `saveMu` for the duration.
   - Generate filename: `name := fmt.Sprintf("%019d-%s.png", now.UnixNano(), TypeName(p.Type))`.
   - Render image using `Render(p, ImageWidth, ImageHeight)`.
   - Write PNG to `filepath.Join(dir, name)` using `png.Encode`.
   - Return `name` on success, or `("", err)` on failure.
   - Pin: SaveImage(DefaultParams(Mandelbrot), tmpDir, time.Unix(0, 1700000000123456789)) -> ("1700000000123456789-mandelbrot.png", nil). The file must be a 600x600 PNG.
2. Implement `ListSaved(dir string) ([]SavedImage, error)`:
   - Read directory. If it doesn't exist, return `(nil, nil)`.
   - Collect all files ending in `.png` (case-sensitive).
   - Sort filenames in descending order (newest first).
   - Return as `SavedImage` slice where Title == Filename.
   - Pin: ListSaved returns .png entries sorted by filename descending (newest first) and returns (nil, nil) for a missing directory.

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `save` bead, `SaveImage`/`ListSaved`:**
  `SaveImage(DefaultParams(Mandelbrot), tmpDir, time.Unix(0, 1700000000123456789))` →
  `("1700000000123456789-mandelbrot.png", nil)` (name = `fmt.Sprintf("%019d-%s.png", 1700000000123456789, "mandelbrot")`,
  and 1700000000123456789 is already 19 digits so no zero-padding is added); the file
  decodes as a 600×600 PNG. `ListSaved` returns `.png` entries sorted by filename
  **descending** (newest first) and returns `(nil, nil)` for a missing directory.

### Revision 2 — created by REVISE_PENDING

**Title:** save  
**Output files:** save.go, save_test.go  
**Exit criteria:** `grep -q 'func TestSaveImage' save_test.go && go test -v -run TestSaveImage ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement image saving in `save.go`.

1. Implement `SaveImage(p Params, dir string, now time.Time) (string, error)`:
   - Lock `saveMu` for the duration.
   - Generate filename: `name := fmt.Sprintf("%019d-%s.png", now.UnixNano(), TypeName(p.Type))`.
   - Render image using `Render(p, 600, 600)`.
   - Write PNG to `filepath.Join(dir, name)` using `png.Encode`.
   - Return `name` on success, or `("", err)` on failure.
   - Pin: SaveImage(DefaultParams(Mandelbrot), tmpDir, time.Unix(0, 1700000000123456789)) -> ("1700000000123456789-mandelbrot.png", nil). The file must be a 600x600 PNG.
2. Implement `ListSaved(dir string) ([]SavedImage, error)`:
   - Read directory. If it doesn't exist, return `(nil, nil)`.
   - Collect all files ending in `.png` (case-sensitive).
   - Sort filenames in descending order (newest first).
   - Return as `SavedImage` slice where Title == Filename.
   - Pin: ListSaved returns .png entries sorted by filename descending (newest first) and returns `(nil, nil)` for a missing directory.

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `save` bead, `SaveImage`/`ListSaved`**:
  `SaveImage(DefaultParams(Mandelbrot), tmpDir, time.Unix(0, 1700000000123456789))` →
  `("1700000000123456789-mandelbrot.png", nil)` (name = `fmt.Sprintf("%019d-%s.png", 1700000000123456789, "mandelbrot")`,
  and 1700000000123456789 is already 19 digits so no zero-padding is added); the file
  decodes as a 600×600 PNG. `ListSaved` returns `.png` entries sorted by filename
  **descending** (newest first) and returns `(nil, nil)` for a missing directory.

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 4 | success | 187s | no fire | 0/0 | not run |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success):
  No commands run. All tests passed in post-execution test run. TestSaveImage not explicitly run during execution.

## Final Output Files

*State of output_files on disk at report time.*

### save.go

```go
package main

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const SavedDir = "saved"

var saveMu sync.Mutex

type SavedImage struct {
	Filename string
	Title    string
}

func SaveImage(p Params, dir string, now time.Time) (string, error) {
	saveMu.Lock()
	defer saveMu.Unlock()

	name := fmt.Sprintf("%019d-%s.png", now.UnixNano(), TypeName(p.Type))
	img := Render(p, 600, 600)

	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
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

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".png") {
			files = append(files, name)
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})

	result := make([]SavedImage, len(files))
	for i, name := range files {
		result[i] = SavedImage{Filename: name, Title: name}
	}
	return result, nil
}

```

### save_test.go

```go
package main

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveImage(t *testing.T) {
	t.Run("SaveImage_Pin", func(t *testing.T) {
		tmpDir := t.TempDir()
		name, err := SaveImage(DefaultParams(Mandelbrot), tmpDir, time.Unix(0, 1700000000123456789))
		if err != nil {
			t.Fatalf("SaveImage returned error: %v", err)
		}
		if name != "1700000000123456789-mandelbrot.png" {
			t.Fatalf("unexpected name %q", name)
		}
		if _, err := os.Stat(filepath.Join(tmpDir, name)); err != nil {
			t.Fatalf("file not created: %v", err)
		}
	})

	t.Run("SaveImage_FileIs600x600PNG", func(t *testing.T) {
		tmpDir := t.TempDir()
		name, err := SaveImage(DefaultParams(Mandelbrot), tmpDir, time.Unix(0, 1700000000123456789))
		if err != nil {
			t.Fatalf("SaveImage error: %v", err)
		}
		f, err := os.Open(filepath.Join(tmpDir, name))
		if err != nil {
			t.Fatalf("open error: %v", err)
		}
		defer f.Close()
		img, err := png.Decode(f)
		if err != nil {
			t.Fatalf("png decode error: %v", err)
		}
		b := img.Bounds()
		if b.Dx() != 600 || b.Dy() != 600 {
			t.Fatalf("unexpected bounds %v", b)
		}
		if !b.Eq(image.Rect(0, 0, 600, 600)) {
			t.Fatalf("bounds not 0,0,600,600: %v", b)
		}
	})

	t.Run("ListSaved_SortedDescending", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, "0000000000000000001-julia.png"), []byte{}, 0644)
		os.WriteFile(filepath.Join(tmpDir, "0000000000000000002-julia.png"), []byte{}, 0644)
		os.WriteFile(filepath.Join(tmpDir, "ignore.txt"), []byte{}, 0644)
		list, err := ListSaved(tmpDir)
		if err != nil {
			t.Fatalf("ListSaved error: %v", err)
		}
		if len(list) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(list))
		}
		if list[0].Filename != "0000000000000000002-julia.png" || list[0].Title != "0000000000000000002-julia.png" {
			t.Fatalf("first entry wrong: %+v", list[0])
		}
		if list[1].Filename != "0000000000000000001-julia.png" || list[1].Title != "0000000000000000001-julia.png" {
			t.Fatalf("second entry wrong: %+v", list[1])
		}
	})

	t.Run("ListSaved_MissingDir", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "doesnotexist")
		os.RemoveAll(missing)
		list, err := ListSaved(missing)
		if err != nil {
			t.Fatalf("expected nil error for missing dir, got %v", err)
		}
		if list != nil {
			t.Fatalf("expected nil slice for missing dir, got %v", list)
		}
	})

	t.Run("ListSaved_IgnoresNonPng", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, "a.png"), []byte{}, 0644)
		os.WriteFile(filepath.Join(tmpDir, "b.jpg"), []byte{}, 0644)
		os.WriteFile(filepath.Join(tmpDir, "c.PNG"), []byte{}, 0644)
		list, err := ListSaved(tmpDir)
		if err != nil {
			t.Fatalf("ListSaved error: %v", err)
		}
		if len(list) != 1 {
			t.Fatalf("expected 1 png entry, got %d", len(list))
		}
		if list[0].Filename != "a.png" {
			t.Fatalf("unexpected filename %q", list[0].Filename)
		}
	})
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-4-attempt-1.log`*

```
 path:save.go]]
[result]
ok: wrote 1202 bytes to save.go
[TURN 3]
Now verify exit criteria: grep -q 'func TestSaveImage' save_test.go && go test -v -run TestSaveImage ./...

Run command.[tool: run_command map[command:grep -q 'func TestSaveImage' save_test.go && go test -v -run TestSaveImage ./...]]
[result]
stdout:
=== RUN   TestSaveImage
=== RUN   TestSaveImage/SaveImage_Pin
=== RUN   TestSaveImage/SaveImage_FileIs600x600PNG
=== RUN   TestSaveImage/ListSaved_SortedDescending
=== RUN   TestSaveImage/ListSaved_MissingDir
=== RUN   TestSaveImage/ListSaved_IgnoresNonPng
--- PASS: TestSaveImage (0.08s)
    --- PASS: TestSaveImage/SaveImage_Pin (0.05s)
    --- PASS: TestSaveImage/SaveImage_FileIs600x600PNG (0.04s)
    --- PASS: TestSaveImage/ListSaved_SortedDescending (0.00s)
    --- PASS: TestSaveImage/ListSaved_MissingDir (0.00s)
    --- PASS: TestSaveImage/ListSaved_IgnoresNonPng (0.00s)
PASS
ok  	fractalviz	0.490s

exit: 0
[TURN 4]
Tests pass.

Now confirm every file listed in Output Files exists on disk. Output Files: save.go. Exists.

Run ls to check.[tool: run_command map[command:ls -l save.go]]
[result]
stdout:
-rw-r--r--@ 1 mike  staff  1202 Sep  6 12:39 save.go

exit: 0
[TURN 5]
All good.

Final message.All exit criteria passed.

- `save.go` implements `SaveImage` with `saveMu` locking, filename `"%019d-%s.png"`, rendering at 600×600 and PNG encoding.
- `ListSaved` returns `nil,nil` for missing directories, collects case-sensitive `.png` files, sorts filenames descending, and returns `SavedImage` with `Title == Filename`.

Tests:
```
=== RUN   TestSaveImage
=== RUN   TestSaveImage/SaveImage_Pin
=== RUN   TestSaveImage/SaveImage_FileIs600x600PNG
=== RUN   TestSaveImage/ListSaved_SortedDescending
=== RUN   TestSaveImage/ListSaved_MissingDir
=== RUN   TestSaveImage/ListSaved_IgnoresNonPng
--- PASS: TestSaveImage
PASS
```

Output file exists:
`save.go` is present on disk.
[done — exit criteria already satisfied on disk; no write needed]

```
