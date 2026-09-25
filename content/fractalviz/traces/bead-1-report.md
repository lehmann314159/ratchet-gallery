# Bead 1: fractal-core

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 295s (4m)  
**Final exit criterion:** `grep -q 'func TestEscape' fractal_test.go && go test -v -run TestEscape ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** fractal-core  
**Output files:** fractal.go, fractal_test.go  
**Exit criteria:** `grep -q 'func TestEscape' fractal_test.go && go test -v -run TestEscape ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement the escape-time fractal core in `fractal.go`. 

1. Implement `TypeName(t)` to return lowercase identifiers: Mandelbrot->"mandelbrot", Julia->"julia", BurningShip->"burningship", Multibrot->"multibrot". Panic on out-of-range types.
2. Implement `ParseType(s)` as the exact inverse of `TypeName` (case-sensitive). Return `(Mandelbrot, error)` for unknown strings.
3. Implement `DefaultParams(t)` returning the specific values for each type: 
   - Mandelbrot: CenterRe=-0.5, CenterIm=0, Zoom=1, MaxIter=100, EscapeR=2, JuliaRe=-0.123, JuliaIm=0.745, Exponent=3
   - Julia: CenterRe=0, CenterIm=0, Zoom=1, MaxIter=100, EscapeR=2, JuliaRe=-0.123, JuliaIm=0.745, Exponent=3
   - Burning Ship: CenterRe=-0.5, CenterIm=-0.5, Zoom=1, MaxIter=100, EscapeR=2, JuliaRe=-0.123, JuliaIm=0.745, Exponent=3
   - Multibrot: CenterRe=0, CenterIm=0, Zoom=1, MaxIter=100, EscapeR=2, JuliaRe=-0.123, JuliaIm=0.745, Exponent=3
4. Implement `Escape(p Params, point complex128)` using the following logic:
   - If Type == Julia: z = point, c = complex(p.JuliaRe, p.JuliaIm)
   - Else: z = 0, c = point
   - Loop n from 0 to p.MaxIter-1:
     - If real(z)*real(z) + imag(z)*imag(z) > p.EscapeR*p.EscapeR: return n
     - Update z based on type:
       - Mandelbrot/Julia: z = z*z + c
       - Burning Ship: z = complex(math.Abs(real(z)), math.Abs(imag(z))); z = z*z + c
       - Multibrot: compute z^p.Exponent by repeated multiplication (z * z * ...), then z = zp + c
   - Return p.MaxIter if loop completes.

Verification Pins:
- Mandelbrot (MaxIter=100, EscapeR=2): point 0+0i -> 100; -1+0i -> 100; 2+2i -> 1; 1+0i -> 3 (NOT 2, NOT 4 — the escape test is real²+imag² > 4 strict, evaluated before each update; orbit 0 -> 1 -> 2 -> 5).
- Julia (MaxIter=100, EscapeR=2, JuliaRe=-0.123, JuliaIm=0.745): point 0+0i -> 100; point 10+10i -> 0 (z0 = point already outside, escapes at n=0); point 2+0i -> 1. Julia uses z0 = point, c = (JuliaRe, JuliaIm).
- Burning Ship (MaxIter=100, EscapeR=2): point 0+0i -> 100; point -0.5-0.5i -> 100; point 2+2i -> 1; point 1+1i -> 2. Update abs's Re z and Im z before squaring.
- Multibrot (MaxIter=100, EscapeR=2, Exponent=3): point 0+0i -> 100; point 1+0i -> 3; point 1.5+0i -> 2. Power is z*z*z by repeated multiplication, not cmplx.Pow.

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `fractal-core` bead, `Escape` (Multibrot, Exponent = 3):** with `MaxIter = 100`,
  `EscapeR = 2`, `Exponent = 3`: `point = 0+0i` → `100`; `point = 1+0i` → `3`;
  `point = 1.5+0i` → `2`. Power is `z*z*z` by repeated multiplication, not `cmplx.Pow`.

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 1 | success | 295s | no fire | 0/0 | not run |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success):
No commands were run. All tests passed in post-execution test run.

## Final Output Files

*State of output_files on disk at report time.*

### fractal.go

```go
package main

import (
	"fmt"
	"math"
)

type FractalType int

const (
	Mandelbrot FractalType = iota
	Julia
	BurningShip
	Multibrot
)

type Params struct {
	Type     FractalType
	CenterRe float64
	CenterIm float64
	Zoom     float64
	MaxIter  int
	EscapeR  float64
	JuliaRe  float64
	JuliaIm  float64
	Exponent int
}

func TypeName(t FractalType) string {
	switch t {
	case Mandelbrot:
		return "mandelbrot"
	case Julia:
		return "julia"
	case BurningShip:
		return "burningship"
	case Multibrot:
		return "multibrot"
	default:
		panic("invalid fractal type")
	}
}

func ParseType(s string) (FractalType, error) {
	switch s {
	case "mandelbrot":
		return Mandelbrot, nil
	case "julia":
		return Julia, nil
	case "burningship":
		return BurningShip, nil
	case "multibrot":
		return Multibrot, nil
	default:
		return Mandelbrot, fmt.Errorf("unknown fractal type %s", s)
	}
}

func DefaultParams(t FractalType) Params {
	p := Params{
		Type:     t,
		Zoom:     1,
		MaxIter:  100,
		EscapeR:  2,
		JuliaRe:  -0.123,
		JuliaIm:  0.745,
		Exponent: 3,
	}
	switch t {
	case Mandelbrot:
		p.CenterRe = -0.5
		p.CenterIm = 0
	case Julia:
		p.CenterRe = 0
		p.CenterIm = 0
	case BurningShip:
		p.CenterRe = -0.5
		p.CenterIm = -0.5
	case Multibrot:
		p.CenterRe = 0
		p.CenterIm = 0
	default:
		// unknown type, return zero values
	}
	return p
}

func Escape(p Params, point complex128) int {
	var z complex128
	var c complex128
	if p.Type == Julia {
		z = point
		c = complex(p.JuliaRe, p.JuliaIm)
	} else {
		z = 0
		c = point
	}
	escapeR2 := p.EscapeR * p.EscapeR
	for n := 0; n < p.MaxIter; n++ {
		re := real(z)
		im := imag(z)
		if re*re+im*im > escapeR2 {
			return n
		}
		switch p.Type {
		case Mandelbrot, Julia:
			z = z*z + c
		case BurningShip:
			z = complex(math.Abs(real(z)), math.Abs(imag(z)))
			z = z*z + c
		case Multibrot:
			zp := z
			// repeated multiplication to compute z^Exponent
			for i := 1; i < p.Exponent; i++ {
				zp = zp * z
			}
			z = zp + c
		default:
			// unknown type, treat as no update
		}
	}
	return p.MaxIter
}

```

### fractal_test.go

```go
package main

import "testing"

func TestEscape(t *testing.T) {
	t.Run("Mandelbrot", func(t *testing.T) {
		p := DefaultParams(Mandelbrot)
		p1 := p
		if got := Escape(p1, 0+0i); got != 100 {
			t.Errorf("Mandelbrot 0+0i: got %d, want 100", got)
		}
		p2 := p
		if got := Escape(p2, -1+0i); got != 100 {
			t.Errorf("Mandelbrot -1+0i: got %d, want 100", got)
		}
		p3 := p
		if got := Escape(p3, 2+2i); got != 1 {
			t.Errorf("Mandelbrot 2+2i: got %d, want 1", got)
		}
		p4 := p
		if got := Escape(p4, 1+0i); got != 3 {
			t.Errorf("Mandelbrot 1+0i: got %d, want 3", got)
		}
	})
	t.Run("Julia", func(t *testing.T) {
		p := DefaultParams(Julia)
		p1 := p
		if got := Escape(p1, 0+0i); got != 100 {
			t.Errorf("Julia 0+0i: got %d, want 100", got)
		}
		p2 := p
		if got := Escape(p2, 10+10i); got != 0 {
			t.Errorf("Julia 10+10i: got %d, want 0", got)
		}
		p3 := p
		if got := Escape(p3, 2+0i); got != 1 {
			t.Errorf("Julia 2+0i: got %d, want 1", got)
		}
	})
	t.Run("BurningShip", func(t *testing.T) {
		p := DefaultParams(BurningShip)
		p1 := p
		if got := Escape(p1, 0+0i); got != 100 {
			t.Errorf("BurningShip 0+0i: got %d, want 100", got)
		}
		p2 := p
		if got := Escape(p2, -0.5-0.5i); got != 100 {
			t.Errorf("BurningShip -0.5-0.5i: got %d, want 100", got)
		}
		p3 := p
		if got := Escape(p3, 2+2i); got != 1 {
			t.Errorf("BurningShip 2+2i: got %d, want 1", got)
		}
		p4 := p
		if got := Escape(p4, 1+1i); got != 2 {
			t.Errorf("BurningShip 1+1i: got %d, want 2", got)
		}
	})
	t.Run("Multibrot", func(t *testing.T) {
		p := DefaultParams(Multibrot)
		p1 := p
		if got := Escape(p1, 0+0i); got != 100 {
			t.Errorf("Multibrot 0+0i: got %d, want 100", got)
		}
		p2 := p
		if got := Escape(p2, 1+0i); got != 3 {
			t.Errorf("Multibrot 1+0i: got %d, want 3", got)
		}
		p3 := p
		if got := Escape(p3, 1.5+0i); got != 2 {
			t.Errorf("Multibrot 1.5+0i: got %d, want 2", got)
		}
	})
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-1-attempt-1.log`*

```
    --- PASS: TestEscape/Multibrot (0.00s)
PASS
ok  	fractalviz	0.414s

exit: 0
[TURN 4]
Tests pass.

Now confirm file exists.[tool: run_command map[command:ls fractal.go]]
[result]
stdout:
fractal.go

exit: 0
[TURN 5]
All good.

We should ensure TypeName panic on out-of-range. Could test quickly but not needed.

Potential nuance: ParseType returns Mandelbrot, error for unknown. We return Mandelbrot, error. Good.

DefaultParams: For Mandelbrot, CenterRe=-0.5, CenterIm=0. Good.

One nuance: For Burning Ship default params, CenterIm=-0.5. Yes.

Potential nuance: Escape for Burning Ship: we applied abs to z before squaring, but spec says update abs's Re z and Im z before squaring. That's what we did. However, some implementations apply abs to z before checking escape? No.

Potential nuance: For Burning Ship, the abs should be applied to z each iteration before squaring, but also the z used for next iteration's escape check is after squaring+ c. That's fine.

Potential nuance: For Multibrot, repeated multiplication: we used zp = z; for i:=1; i < p.Exponent; i++ { zp = zp * z }. That's z^(i+1). For exponent 3, i=1,2 => zp = z*z, then z*z*z. Good.

Edge: exponent =0? Then loop doesn't run, zp = z, z = z + c. Might be weird but not tested.

All good.

Now final message.All exit criteria pass. The implementation is complete.

```bash
grep -q 'func TestEscape' fractal_test.go && go test -v -run TestEscape ./...
```
```
=== RUN   TestEscape
=== RUN   TestEscape/Mandelbrot
=== RUN   TestEscape/Julia
=== RUN   TestEscape/BurningShip
=== RUN   TestEscape/Multibrot
--- PASS: TestEscape (0.00s)
    --- PASS: TestEscape/Mandelbrot (0.00s)
    --- PASS: TestEscape/Julia (0.00s)
    --- PASS: TestEscape/BurningShip (0.00s)
    --- PASS: TestEscape/Multibrot (0.00s)
PASS
ok  	fractalviz	0.414s
```

Output file exists:
- `fractal.go`

[done — exit criteria already satisfied on disk; no write needed]

```
