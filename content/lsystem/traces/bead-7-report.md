# Bead 7: render

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 359s (5m)  
**Final exit criterion:** `grep -q 'func TestRender' render_test.go && go test -v -run TestRender ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** render  
**Output files:** render.go, render_test.go  
**Exit criteria:** `grep -q 'func TestRender' render_test.go && go test -v -run TestRender ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement `BBox(segs []Segment) (Box, bool)` and `RenderSVG(segs []Segment) string` in `render.go`.

1. `BBox`: Find min/max of all X and Y coordinates. Return `(Box{}, false)` if `segs` is empty.
2. `RenderSVG`: 
   - Generate a 800x800 SVG with a white background rect.
   - If `segs` is empty, return just the background.
   - Otherwise, calculate a uniform scale `s` and offsets `ox, oy` to fit the `BBox` within the canvas with a `Margin` (20px) border, centering the drawing.
   - Map turtle coordinates to SVG: `mapX(tx) = ox + (tx - b.MinX) * s`; `mapY(ty) = oy + (b.MaxY - ty) * s` (Y is flipped).
   - Generate a `<path>` with `d` attribute containing `M %.3f %.3f L %.3f %.3f` for each segment.

Pin — `render` bead, `BBox` + `RenderSVG`: `BBox([{0,0,2,0},{2,0,2,1}])` → `Box{0,0,2,1}, true`; the fit transform for that box (`CanvasSize 800`, `Margin 20`) is `s=380, ox=20, oy=210`, mapping `(0,0)→(20,590)`, `(2,1)→(780,210)`. `BBox(nil)` → `Box{}, false`. `RenderSVG([{0,0,1,0}])` → `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="800" viewBox="0 0 800 800"><rect width="800" height="800" fill="white"/><path d="M 20.000 400.000 L 780.000 400.000" fill="none" stroke="black" stroke-width="1"/></svg>`. `RenderSVG(nil)` → `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="800" viewBox="0 0 800 800"><rect width="800" height="800" fill="white"/></svg>`.

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `render` bead, `BBox` + `RenderSVG`:** `BBox([{0,0,2,0},{2,0,2,1}])` →
  `Box{0,0,2,1}, true`; the fit transform for that box (`CanvasSize 800`, `Margin 20`) is
  `s=380, ox=20, oy=210`, mapping `(0,0)→(20,590)`, `(2,1)→(780,210)`. `BBox(nil)` →
  `Box{}, false`. `RenderSVG([{0,0,1,0}])` →
  `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="800" viewBox="0 0 800 800"><rect width="800" height="800" fill="white"/><path d="M 20.000 400.000 L 780.000 400.000" fill="none" stroke="black" stroke-width="1"/></svg>`.
  `RenderSVG(nil)` →
  `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="800" viewBox="0 0 800 800"><rect width="800" height="800" fill="white"/></svg>`.

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 7 | success | 359s | no fire | 1/1 | PASS |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success):
grep -q 'func TestRender' render_test.go && go test -v -run TestRender ./... → exit 0
--- PASS: TestRender (0.00s)
PASS
ok	lsystem	0.420s


## Final Output Files

*State of output_files on disk at report time.*

### render.go

```go
package main

import (
	"fmt"
	"math"
)

const CanvasSize = 800
const Margin = 20

type Box struct{ MinX, MinY, MaxX, MaxY float64 }

func BBox(segs []Segment) (Box, bool) {
	if len(segs) == 0 {
		return Box{}, false
	}
	minX := segs[0].X1
	maxX := segs[0].X1
	minY := segs[0].Y1
	maxY := segs[0].Y1
	for _, s := range segs {
		if s.X1 < minX {
			minX = s.X1
		}
		if s.X2 < minX {
			minX = s.X2
		}
		if s.X1 > maxX {
			maxX = s.X1
		}
		if s.X2 > maxX {
			maxX = s.X2
		}
		if s.Y1 < minY {
			minY = s.Y1
		}
		if s.Y2 < minY {
			minY = s.Y2
		}
		if s.Y1 > maxY {
			maxY = s.Y1
		}
		if s.Y2 > maxY {
			maxY = s.Y2
		}
	}
	return Box{MinX: minX, MinY: minY, MaxX: maxX, MaxY: maxY}, true
}

func RenderSVG(segs []Segment) string {
	base := `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="800" viewBox="0 0 800 800"><rect width="800" height="800" fill="white"/>`
	if len(segs) == 0 {
		return base + `</svg>`
	}
	box, _ := BBox(segs)
	width := box.MaxX - box.MinX
	height := box.MaxY - box.MinY
	inner := float64(CanvasSize - 2*Margin)
	s := math.Inf(1)
	if width > 0 {
		s = math.Min(s, inner/width)
	}
	if height > 0 {
		s = math.Min(s, inner/height)
	}
	if math.IsInf(s, 1) {
		s = 1
	}
	ox := float64(Margin) + (inner-width*s)/2
	oy := float64(Margin) + (inner-height*s)/2

	d := ""
	for i, seg := range segs {
		x1 := ox + (seg.X1-box.MinX)*s
		y1 := oy + (box.MaxY-seg.Y1)*s
		x2 := ox + (seg.X2-box.MinX)*s
		y2 := oy + (box.MaxY-seg.Y2)*s
		part := fmt.Sprintf("M %.3f %.3f L %.3f %.3f", x1, y1, x2, y2)
		if i > 0 {
			d += " " + part
		} else {
			d = part
		}
	}
	return fmt.Sprintf(`%s<path d="%s" fill="none" stroke="black" stroke-width="1"/></svg>`, base, d)
}

```

### render_test.go

```go
package main

import "testing"

func TestRender(t *testing.T) {
	segs := []Segment{{X1: 0, Y1: 0, X2: 2, Y2: 0}, {X1: 2, Y1: 0, X2: 2, Y2: 1}}
	box, ok := BBox(segs)
	if !ok {
		t.Fatalf("BBox expected ok=true")
	}
	expBox := Box{MinX: 0, MinY: 0, MaxX: 2, MaxY: 1}
	if box != expBox {
		t.Fatalf("BBox got %+v want %+v", box, expBox)
	}

	boxNil, okNil := BBox(nil)
	if okNil {
		t.Fatalf("BBox(nil) expected ok=false")
	}
	if boxNil != (Box{}) {
		t.Fatalf("BBox(nil) expected zero Box, got %+v", boxNil)
	}

	boxEmpty, okEmpty := BBox([]Segment{})
	if okEmpty {
		t.Fatalf("BBox([]) expected ok=false")
	}
	if boxEmpty != (Box{}) {
		t.Fatalf("BBox([]) expected zero Box, got %+v", boxEmpty)
	}

	svg := RenderSVG([]Segment{{X1: 0, Y1: 0, X2: 1, Y2: 0}})
	expectedSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="800" viewBox="0 0 800 800"><rect width="800" height="800" fill="white"/><path d="M 20.000 400.000 L 780.000 400.000" fill="none" stroke="black" stroke-width="1"/></svg>`
	if svg != expectedSVG {
		t.Fatalf("RenderSVG single segment mismatch\ngot  %s\nwant %s", svg, expectedSVG)
	}

	svgNil := RenderSVG(nil)
	expectedNil := `<svg xmlns="http://www.w3.org/2000/svg" width="800" height="800" viewBox="0 0 800 800"><rect width="800" height="800" fill="white"/></svg>`
	if svgNil != expectedNil {
		t.Fatalf("RenderSVG(nil) mismatch\ngot  %s\nwant %s", svgNil, expectedNil)
	}

	svgEmpty := RenderSVG([]Segment{})
	if svgEmpty != expectedNil {
		t.Fatalf("RenderSVG([]) mismatch\ngot  %s\nwant %s", svgEmpty, expectedNil)
	}
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-7-attempt-1.log`*

```
		x2 := ox + (seg.X2-box.MinX)*s
		y2 := oy + (box.MaxY-seg.Y2)*s
		part := fmt.Sprintf("M %.3f %.3f L %.3f %.3f", x1, y1, x2, y2)
		if i > 0 {
			d += " " + part
		} else {
			d = part
		}
	}
	return fmt.Sprintf(`%s<path d="%s" fill="none" stroke="black" stroke-width="1"/></svg>`, base, d)
}
 path:render.go]]
[result]
ok: wrote 1698 bytes to render.go
[TURN 2]
Now run exit criteria.

Exit criteria: grep -q 'func TestRender' render_test.go && go test -v -run TestRender ./...

We need to run command.

[tool: run_command map[command:grep -q 'func TestRender' render_test.go && go test -v -run TestRender ./...]]
[result]
stdout:
=== RUN   TestRender
--- PASS: TestRender (0.00s)
PASS
ok  	lsystem	0.420s

exit: 0
[TURN 3]
Test passes.

Now confirm file exists.[tool: run_command map[command:ls render.go]]
[result]
stdout:
render.go

exit: 0
[TURN 4]
All good.

We should output final message.

All exit criteria passed.

`render.go` now implements `BBox` and `RenderSVG` as specified. The `TestRender` suite passes:

```
=== RUN   TestRender
--- PASS: TestRender (0.00s)
PASS
ok  	lsystem	0.420s
```

Output file exists:
- `render.go`

[done — exit criteria already satisfied on disk; no write needed]

```
