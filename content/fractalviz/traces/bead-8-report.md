# Bead 8: integration

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 146s (2m)  
**Final exit criterion:** `grep -q 'func TestIntegration' integration_test.go && go test -v -run TestIntegration ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** integration  
**Output files:** integration_test.go  
**Exit criteria:** `grep -q 'func TestIntegration' integration_test.go && go test -v -run TestIntegration ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement a final integration test in `integration_test.go`.

1. Use `httptest.NewServer` with the real mux configured in `main.go`.
2. Perform a `GET /fractal.png?type=mandelbrot` request.
3. Assert:
   - HTTP status 200.
   - Content-Type is `image/png`.
   - Body `png.Decode`s successfully to an image with bounds `image.Rect(0, 0, 600, 600)`.
   - Pin: The pixel at (300, 300) is opaque black (R=0, G=0, B=0, A=255).

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `integration` bead:** `GET /fractal.png?type=mandelbrot` → HTTP 200,
  `Content-Type: image/png`, body `png.Decode`s to bounds `image.Rect(0, 0, 600, 600)`,
  pixel `(300, 300)` opaque black.

### Revision 2 — created by REVISE_PENDING

**Title:** integration  
**Output files:** integration_test.go  
**Exit criteria:** `grep -q 'func TestIntegration' integration_test.go && go test -v -run TestIntegration ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement a final integration test in `integration_test.go`.

1. Use `httptest.NewServer` with the real mux configured in `main.go`.
2. Perform a `GET /fractal.png?type=mandelbrot` request.
3. Assert:
   - HTTP status 200.
   - Content-Type is `image/png`.
   - Body `png.Decode`s successfully to an image with bounds `image.Rect(0, 0, 600, 600)`.
   - Pin: The pixel at (300, 300) is opaque black (R=0, G=0, B=0, A=255).

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `integration` bead:** `GET /fractal.png?type=mandelbrot` → HTTP 200,
  `Content-Type: image/png`, body `png.Decode`s to bounds `image.Rect(0, 0, 600, 600)`,
  pixel `(300, 300)` opaque black.

### Revision 3 — created by REVISE_PENDING

**Title:** integration  
**Output files:** integration_test.go  
**Exit criteria:** `grep -q 'func TestIntegration' integration_test.go && go test -v -run TestIntegration ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement a final integration test in `integration_test.go`.

1. Use `httptest.NewServer` with the real mux configured in `main.go`.
2. Perform a `GET /fractal.png?type=mandelbrot` request.
3. Assert:
   - HTTP status 200.
   - Content-Type is `image/png`.
   - Body `png.Decode`s successfully to an image with bounds `image.Rect(0, 0, 600, 600)`.
   - Pixel at (300, 300) is opaque black (R=0, G=0, B=0, A=255).

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `integration` bead:** `GET /fractal.png?type=mandelbrot` → HTTP 200,
  `Content-Type: image/png`, body `png.Decode`s to bounds `image.Rect(0, 0, 600, 600)`,
  pixel `(300, 300)` opaque black.

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 8 | success | 146s | no fire | 0/0 | not run |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success):
No commands were run.
Post-execution test run (go test ./...): all tests passed


## Final Output Files

*State of output_files on disk at report time.*

### integration_test.go

```go
package main

import (
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegration(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleIndex)
	mux.HandleFunc("/render", HandleRender)
	mux.HandleFunc("/select", HandleSelect)
	mux.HandleFunc("/fractal.png", HandleImage)
	mux.HandleFunc("/save", HandleSave)
	mux.Handle("/saved/", http.StripPrefix("/saved/", http.FileServer(http.Dir(SavedDir))))
	InitTemplates()
	server := httptest.NewServer(mux)
	defer server.Close()

	resp, err := http.Get(server.URL + "/fractal.png?type=mandelbrot")
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}
	img, err := png.Decode(resp.Body)
	if err != nil {
		t.Fatalf("png.Decode failed: %v", err)
	}
	if img.Bounds() != image.Rect(0, 0, 600, 600) {
		t.Fatalf("unexpected bounds %v", img.Bounds())
	}
	c := color.RGBAModel.Convert(img.At(300, 300)).(color.RGBA)
	if c.R != 0 || c.G != 0 || c.B != 0 || c.A != 255 {
		t.Fatalf("pixel (300,300) = %+v, want opaque black", c)
	}
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-8-attempt-1.log`*

```
We could just leave as is. But maybe they want us to ensure the test uses httptest.NewServer with real mux configured in main.go. The current test creates mux manually, but it's same as main. Could be acceptable.

Potential nuance: main.go registers handlers in order: "/", "/render", "/select", "/fractal.png", "/save", "/saved/". The test does same order. Good.

One nuance: In main.go, mux.HandleFunc("/", HandleIndex). In Go's ServeMux, pattern "/" matches all. But they also register other patterns. That's fine.

Potential nuance: The test uses InitTemplates() before server. Good.

Maybe they want to ensure the test uses the same mux as main.go, which also creates SavedDir. Not needed.

I think it's fine.

But to be safe, maybe we should rewrite file to ensure it matches spec exactly? The spec says "Use httptest.NewServer with the real mux configured in main.go." Could be interpreted as we should create mux exactly as main.go does, which we did.

Maybe they want to avoid manual handler registration and instead use a function that returns mux? But main.go doesn't expose mux.

Could be okay.

Check if test passes with -run TestIntegration. It passed.

Now confirm file exists.[tool: run_command map[command:ls integration_test.go]]
[result]
stdout:
integration_test.go

exit: 0
[TURN 4]
All good.

Potential improvement: The spec says "Perform a GET /fractal.png?type=mandelbrot request." The test does http.Get(server.URL + "/fractal.png?type=mandelbrot"). Good.

Assert HTTP status 200, Content-Type image/png, body png.Decode to image with bounds image.Rect(0,0,600,600), pixel at (300,300) opaque black.

Done.

Maybe we should ensure the test uses the exact same imports? It's fine.

Potential edge: color.RGBAModel.Convert returns color.Color, type assertion to color.RGBA may panic if not RGBA. But png.Decode returns image.Image which could be *image.RGBA or *image.Paletted? Render returns *image.RGBA, png.Encode encodes it, png.Decode decodes to *image.RGBA? Actually png.Decode returns image.Image, could be *image.RGBA. Convert should work.

Test passed.

Thus final message.

We should output final message confirming.

All exit criteria pass. The integration test is present and passing.

```bash
grep -q 'func TestIntegration' integration_test.go && go test -v -run TestIntegration ./...
```
```
=== RUN   TestIntegration
--- PASS: TestIntegration (0.06s)
PASS
ok   fractalviz 0.521s
```

`integration_test.go` exists and implements the required `GET /fractal.png?type=mandelbrot` integration check with HTTP 200, `Content-Type: image/png`, bounds `image.Rect(0,0,600,600)`, and opaque black pixel at (300,300).
[done — exit criteria already satisfied on disk; no write needed]

```
