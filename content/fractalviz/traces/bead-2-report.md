# Bead 2: render

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 257s (4m)  
**Final exit criterion:** `grep -q 'func TestRender' render_test.go && go test -v -run TestRender ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** render  
**Output files:** render.go, render_test.go  
**Exit criteria:** `grep -q 'func TestRender' render_test.go && go test -v -run TestRender ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement the rendering logic in `render.go`.

1. Implement `PixelToComplex(p Params, px, py, width, height int)`:
   - viewH = 4.0 / p.Zoom
   - viewW = viewH * width / height
   - re = p.CenterRe + (px + 0.5 - width/2) / width * viewW
   - im = p.CenterIm - (py + 0.5 - height/2) / height * viewH
   - Return complex(re, im).
   - Pin: with Params{CenterRe: -0.5, CenterIm: 0, Zoom: 1}, width=height=600: pixel (0,0) -> (-2.4966666666666666, +1.9966666666666666); pixel (300,300) -> (-0.49666666666666665, -0.0033333333333333335); pixel (599,599) -> (+1.4966666666666666, -1.9966666666666666). py = 0 is the top = largest imaginary; the + 0.5 pixel-centre offset is required.
2. Implement `Color(n, maxIter int)`:
   - If n >= maxIter: return opaque black {0, 0, 0, 255}.
   - Else: v = uint8(255 * n / maxIter) (integer division). Return {v, v, 255, 255}.
   - Pin: with maxIter=100: Color(0, 100) -> {0, 0, 255, 255} (NOT black); Color(50, 100) -> {127, 127, 255, 255}; Color(99, 100) -> {252, 252, 255, 255}; Color(100, 100) -> {0, 0, 0, 255}.
3. Implement `Render(p Params, width, height int)`:
   - Create a new RGBA image of size width x height.
   - For each pixel (px, py), set color to `Color(Escape(p, PixelToComplex(p, px, py, width, height)), p.MaxIter)`.
   - Pin: Render(DefaultParams(Mandelbrot), 600, 600) has bounds image.Rect(0, 0, 600, 600); pixel (300, 300) is {0, 0, 0, 255}; pixel (0, 0) is {2, 2, 255, 255} (escapes at n=1).

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `render` bead, `Render`:** `Render(DefaultParams(Mandelbrot), 600, 600)` has
  bounds `image.Rect(0, 0, 600, 600)`; pixel `(300, 300)` is `color.RGBA{0, 0, 0, 255}`;
  pixel `(0, 0)` is `color.RGBA{2, 2, 255, 255}` (escapes at n=1).

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 2 | success | 257s | no fire | 0/0 | not run |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success): no failures to record.

## Final Output Files

*State of output_files on disk at report time.*

### render.go

```go
package main

import (
	"image"
	"image/color"
)

const ImageWidth = 600
const ImageHeight = 600

func PixelToComplex(p Params, px, py, width, height int) complex128 {
	viewH := 4.0 / p.Zoom
	viewW := viewH * float64(width) / float64(height)
	re := p.CenterRe + (float64(px)+0.5-float64(width)/2.0)/float64(width)*viewW
	im := p.CenterIm - (float64(py)+0.5-float64(height)/2.0)/float64(height)*viewH
	return complex(re, im)
}

func Color(n, maxIter int) color.RGBA {
	if n >= maxIter {
		return color.RGBA{0, 0, 0, 255}
	}
	v := uint8(255 * n / maxIter)
	return color.RGBA{v, v, 255, 255}
}

func Render(p Params, width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			c := PixelToComplex(p, px, py, width, height)
			n := Escape(p, c)
			col := Color(n, p.MaxIter)
			img.Set(px, py, col)
		}
	}
	return img
}

```

### render_test.go

```go
package main

import (
	"image"
	"image/color"
	"testing"
)

func TestRender(t *testing.T) {
	t.Run("PixelToComplex", func(t *testing.T) {
		p := Params{CenterRe: -0.5, CenterIm: 0, Zoom: 1}
		c00 := PixelToComplex(p, 0, 0, 600, 600)
		re00 := real(c00)
		im00 := imag(c00)
		if re00 != -2.4966666666666666 || im00 != 1.9966666666666666 {
			t.Fatalf("pixel (0,0) = (%v,%v), want (-2.4966666666666666,1.9966666666666666)", re00, im00)
		}
		c33 := PixelToComplex(p, 300, 300, 600, 600)
		re33 := real(c33)
		im33 := imag(c33)
		if re33 != -0.49666666666666665 || im33 != -0.0033333333333333335 {
			t.Fatalf("pixel (300,300) = (%v,%v), want (-0.49666666666666665,-0.0033333333333333335)", re33, im33)
		}
		c99 := PixelToComplex(p, 599, 599, 600, 600)
		re99 := real(c99)
		im99 := imag(c99)
		if re99 != 1.4966666666666666 || im99 != -1.9966666666666666 {
			t.Fatalf("pixel (599,599) = (%v,%v), want (1.4966666666666666,-1.9966666666666666)", re99, im99)
		}
	})
	t.Run("Color", func(t *testing.T) {
		c0 := Color(0, 100)
		if c0 != (color.RGBA{R: 0, G: 0, B: 255, A: 255}) {
			t.Fatalf("Color(0,100) = %v, want {0,0,255,255}", c0)
		}
		c50 := Color(50, 100)
		if c50 != (color.RGBA{R: 127, G: 127, B: 255, A: 255}) {
			t.Fatalf("Color(50,100) = %v, want {127,127,255,255}", c50)
		}
		c99 := Color(99, 100)
		if c99 != (color.RGBA{R: 252, G: 252, B: 255, A: 255}) {
			t.Fatalf("Color(99,100) = %v, want {252,252,255,255}", c99)
		}
		c100 := Color(100, 100)
		if c100 != (color.RGBA{R: 0, G: 0, B: 0, A: 255}) {
			t.Fatalf("Color(100,100) = %v, want {0,0,0,255}", c100)
		}
	})
	t.Run("Render", func(t *testing.T) {
		img := Render(DefaultParams(Mandelbrot), 600, 600)
		if img == nil {
			t.Fatalf("Render returned nil")
		}
		b := img.Bounds()
		if b != image.Rect(0, 0, 600, 600) {
			t.Fatalf("bounds = %v, want image.Rect(0,0,600,600)", b)
		}
		p300 := img.RGBAAt(300, 300)
		if p300 != (color.RGBA{R: 0, G: 0, B: 0, A: 255}) {
			t.Fatalf("pixel (300,300) = %v, want {0,0,0,255}", p300)
		}
		p00 := img.RGBAAt(0, 0)
		if p00 != (color.RGBA{R: 2, G: 2, B: 255, A: 255}) {
			t.Fatalf("pixel (0,0) = %v, want {2,2,255,255}", p00)
		}
	})
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-2-attempt-1.log`*

```
			col := Color(n, p.MaxIter)
			img.Set(px, py, col)
		}
	}
	return img
}
 path:render.go]]
[result]
ok: wrote 915 bytes to render.go
[TURN 3]
Now run exit criteria: grep -q 'func TestRender' render_test.go && go test -v -run TestRender ./...

Let's run.[tool: run_command map[command:grep -q 'func TestRender' render_test.go && go test -v -run TestRender ./...]]
[result]
stdout:
=== RUN   TestRender
=== RUN   TestRender/PixelToComplex
=== RUN   TestRender/Color
=== RUN   TestRender/Render
--- PASS: TestRender (0.02s)
    --- PASS: TestRender/PixelToComplex (0.00s)
    --- PASS: TestRender/Color (0.00s)
    --- PASS: TestRender/Render (0.02s)
PASS
ok  	fractalviz	0.414s

exit: 0
[TURN 4]
Tests pass.

Now confirm file exists.[tool: run_command map[command:ls render.go]]
[result]
stdout:
render.go

exit: 0
[TURN 5]
All good.

We should ensure no other files created. Output files list only render.go. Good.

Final message.All exit criteria passed.

- `render.go` has been implemented with `PixelToComplex`, `Color`, and `Render` matching the specified pins.
- Tests pass:
```
=== RUN   TestRender
=== RUN   TestRender/PixelToComplex
=== RUN   TestRender/Color
=== RUN   TestRender/Render
--- PASS: TestRender
```

Output file exists:
```
render.go
```

[done — exit criteria already satisfied on disk; no write needed]

```
