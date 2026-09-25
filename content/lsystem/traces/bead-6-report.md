# Bead 6: turtle

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 239s (3m)  
**Final exit criterion:** `grep -q 'func TestInterpret' turtle_test.go && go test -v -run TestInterpret ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** turtle  
**Output files:** turtle.go, turtle_test.go  
**Exit criteria:** `grep -q 'func TestInterpret' turtle_test.go && go test -v -run TestInterpret ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement `Interpret(mods []Module, angleDefault, stepDefault, headingStart float64) ([]Segment, error)` in `turtle.go`.

1. State: `x, y` (position), `h` (heading in degrees), and a stack of `(x, y, h)` triples. Start at `(0, 0, headingStart)`.
2. Process modules:
   - `F`: Move distance `d` (from `Args[0]` or `stepDefault`). Append `Segment` and update `x, y` using `math.Cos` and `math.Sin` (convert `h` to radians). 
   - `f`: Same as `F` but do not append a segment.
   - `+`: Turn by angle `a` (from `Args[0]` or `angleDefault`). `h += a`.
   - `-`: Turn by angle `a` (from `Args[0]` or `angleDefault`). `h -= a`.
   - `[`: Push current `(x, y, h)` to stack.
   - `]`: Pop and restore `x, y, h`. Error if stack is empty.
   - Others: Ignore.
3. Error if `F, f, +, -` have 2 or more arguments.
4. On error, return `(nil, err)`. If stack is non-empty at end, return segments as-is.

Pin — `turtle` bead, `Interpret`: start `(0,0)` heading `0`, `angleDefault 90`, `stepDefault 1`. `F(1) +(90) F(1)` → segments `{0,0,1,0}`, `{1,0,1,1}` (both exact). `F(1) [ +(90) F(1) ] F(1)` → **three** segments `{0,0,1,0}`, `{1,0,1,1}`, `{1,0,2,0}` (the `]` restores position to `(1,0)` and heading to `0`). `]` on an empty stack → `(nil, error)`. `F(1,2)` → `(nil, error)`.

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `turtle` bead, `Interpret`:** start `(0,0)` heading `0`, `angleDefault 90`,
  `stepDefault 1`. `F(1) +(90) F(1)` → segments `{0,0,1,0}`, `{1,0,1,1}` (both exact).
  `F(1) [ +(90) F(1) ] F(1)` → **three** segments `{0,0,1,0}`, `{1,0,1,1}`, `{1,0,2,0}`
  (the `]` restores position to `(1,0)` and heading to `0`). `]` on an empty stack →
  `(nil, error)`. `F(1,2)` → `(nil, error)`.

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 6 | success | 239s | no fire | 0/0 | not run |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success):
No commands were run.
Post-execution test run (go test ./...): all tests passed


## Final Output Files

*State of output_files on disk at report time.*

### turtle.go

```go
package main

import (
	"fmt"
	"math"
)

type Segment struct{ X1, Y1, X2, Y2 float64 }

func Interpret(mods []Module, angleDefault, stepDefault, headingStart float64) ([]Segment, error) {
	x, y := 0.0, 0.0
	h := headingStart
	type state struct{ x, y, h float64 }
	stack := []state{}
	segments := []Segment{}

	for _, m := range mods {
		sym := m.Sym
		args := m.Args
		if sym == 'F' || sym == 'f' || sym == '+' || sym == '-' {
			if len(args) >= 2 {
				return nil, fmt.Errorf("too many args")
			}
		}
		switch sym {
		case 'F':
			d := stepDefault
			if len(args) == 1 {
				d = args[0]
			}
			x1, y1 := x, y
			rad := h * math.Pi / 180
			x += d * math.Cos(rad)
			y += d * math.Sin(rad)
			segments = append(segments, Segment{X1: x1, Y1: y1, X2: x, Y2: y})
		case 'f':
			d := stepDefault
			if len(args) == 1 {
				d = args[0]
			}
			rad := h * math.Pi / 180
			x += d * math.Cos(rad)
			y += d * math.Sin(rad)
		case '+':
			a := angleDefault
			if len(args) == 1 {
				a = args[0]
			}
			h += a
		case '-':
			a := angleDefault
			if len(args) == 1 {
				a = args[0]
			}
			h -= a
		case '[':
			stack = append(stack, state{x: x, y: y, h: h})
		case ']':
			if len(stack) == 0 {
				return nil, fmt.Errorf("stack empty")
			}
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			x, y, h = top.x, top.y, top.h
		default:
			// ignore
		}
	}
	return segments, nil
}

```

### turtle_test.go

```go
package main

import "testing"

func TestInterpret(t *testing.T) {
	t.Run("BasicFPlusF", func(t *testing.T) {
		mods := []Module{{Sym: 'F', Args: []float64{1}}, {Sym: '+', Args: []float64{90}}, {Sym: 'F', Args: []float64{1}}}
		segs, err := Interpret(mods, 90, 1, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(segs) != 2 {
			t.Fatalf("expected 2 segments, got %d", len(segs))
		}
		if segs[0] != (Segment{X1: 0, Y1: 0, X2: 1, Y2: 0}) {
			t.Fatalf("seg0 mismatch: %+v", segs[0])
		}
		if segs[1] != (Segment{X1: 1, Y1: 0, X2: 1, Y2: 1}) {
			t.Fatalf("seg1 mismatch: %+v", segs[1])
		}
	})
	t.Run("StackRestore", func(t *testing.T) {
		mods := []Module{{Sym: 'F', Args: []float64{1}}, {Sym: '[', Args: []float64{}}, {Sym: '+', Args: []float64{90}}, {Sym: 'F', Args: []float64{1}}, {Sym: ']', Args: []float64{}}, {Sym: 'F', Args: []float64{1}}}
		segs, err := Interpret(mods, 90, 1, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(segs) != 3 {
			t.Fatalf("expected 3 segments, got %d", len(segs))
		}
		if segs[0] != (Segment{X1: 0, Y1: 0, X2: 1, Y2: 0}) {
			t.Fatalf("seg0 mismatch: %+v", segs[0])
		}
		if segs[1] != (Segment{X1: 1, Y1: 0, X2: 1, Y2: 1}) {
			t.Fatalf("seg1 mismatch: %+v", segs[1])
		}
		if segs[2] != (Segment{X1: 1, Y1: 0, X2: 2, Y2: 0}) {
			t.Fatalf("seg2 mismatch: %+v", segs[2])
		}
	})
	t.Run("EmptyStackError", func(t *testing.T) {
		mods := []Module{{Sym: ']', Args: []float64{}}}
		segs, err := Interpret(mods, 90, 1, 0)
		if err == nil {
			t.Fatalf("expected error for empty stack pop")
		}
		if segs != nil {
			t.Fatalf("expected nil segments on error")
		}
	})
	t.Run("TooManyArgsError", func(t *testing.T) {
		mods := []Module{{Sym: 'F', Args: []float64{1, 2}}}
		segs, err := Interpret(mods, 90, 1, 0)
		if err == nil {
			t.Fatalf("expected error for too many args")
		}
		if segs != nil {
			t.Fatalf("expected nil segments on error")
		}
	})
	t.Run("BareSymbols", func(t *testing.T) {
		mods := []Module{{Sym: 'F', Args: []float64{}}, {Sym: '+', Args: []float64{}}, {Sym: 'F', Args: []float64{}}}
		segs, err := Interpret(mods, 90, 1, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(segs) != 2 {
			t.Fatalf("expected 2 segments, got %d", len(segs))
		}
		if segs[0] != (Segment{X1: 0, Y1: 0, X2: 1, Y2: 0}) {
			t.Fatalf("seg0 mismatch")
		}
		if segs[1] != (Segment{X1: 1, Y1: 0, X2: 1, Y2: 1}) {
			t.Fatalf("seg1 mismatch")
		}
	})
	t.Run("IgnoreOther", func(t *testing.T) {
		mods := []Module{{Sym: 'X', Args: []float64{}}, {Sym: 'F', Args: []float64{1}}}
		segs, err := Interpret(mods, 90, 1, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(segs) != 1 {
			t.Fatalf("expected 1 segment, got %d", len(segs))
		}
		if segs[0] != (Segment{X1: 0, Y1: 0, X2: 1, Y2: 0}) {
			t.Fatalf("seg mismatch: %+v", segs[0])
		}
	})
	t.Run("LowercaseFNoSegment", func(t *testing.T) {
		mods := []Module{{Sym: 'f', Args: []float64{1}}, {Sym: 'F', Args: []float64{1}}}
		segs, err := Interpret(mods, 90, 1, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(segs) != 1 {
			t.Fatalf("expected 1 segment, got %d", len(segs))
		}
		if segs[0] != (Segment{X1: 1, Y1: 0, X2: 2, Y2: 0}) {
			t.Fatalf("seg mismatch: %+v", segs[0])
		}
	})
	t.Run("LeftoverBracket", func(t *testing.T) {
		mods := []Module{{Sym: 'F', Args: []float64{1}}, {Sym: '[', Args: []float64{}}}
		segs, err := Interpret(mods, 90, 1, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(segs) != 1 {
			t.Fatalf("expected 1 segment, got %d", len(segs))
		}
		if segs[0] != (Segment{X1: 0, Y1: 0, X2: 1, Y2: 0}) {
			t.Fatalf("seg mismatch")
		}
	})
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-6-attempt-1.log`*

```
Now run exit criteria.[tool: run_command map[command:grep -q 'func TestInterpret' turtle_test.go && go test -v -run TestInterpret ./...]]
[result]
stdout:
=== RUN   TestInterpret
=== RUN   TestInterpret/BasicFPlusF
=== RUN   TestInterpret/StackRestore
=== RUN   TestInterpret/EmptyStackError
=== RUN   TestInterpret/TooManyArgsError
=== RUN   TestInterpret/BareSymbols
=== RUN   TestInterpret/IgnoreOther
=== RUN   TestInterpret/LowercaseFNoSegment
=== RUN   TestInterpret/LeftoverBracket
--- PASS: TestInterpret (0.00s)
    --- PASS: TestInterpret/BasicFPlusF (0.00s)
    --- PASS: TestInterpret/StackRestore (0.00s)
    --- PASS: TestInterpret/EmptyStackError (0.00s)
    --- PASS: TestInterpret/TooManyArgsError (0.00s)
    --- PASS: TestInterpret/BareSymbols (0.00s)
    --- PASS: TestInterpret/IgnoreOther (0.00s)
    --- PASS: TestInterpret/LowercaseFNoSegment (0.00s)
    --- PASS: TestInterpret/LeftoverBracket (0.00s)
PASS
ok  	lsystem	0.433s

exit: 0
[TURN 6]
All tests pass. Need to confirm file exists.[tool: run_command map[command:ls turtle.go]]
[result]
stdout:
turtle.go

exit: 0
[TURN 7]
Done.All exit criteria passed. The implementation in `turtle.go` now correctly handles movement, turning, stack push/pop, defaults, argument validation, and error cases.

```bash
grep -q 'func TestInterpret' turtle_test.go && go test -v -run TestInterpret ./...
```
```
=== RUN   TestInterpret
=== RUN   TestInterpret/BasicFPlusF
=== RUN   TestInterpret/StackRestore
=== RUN   TestInterpret/EmptyStackError
=== RUN   TestInterpret/TooManyArgsError
=== RUN   TestInterpret/BareSymbols
=== RUN   TestInterpret/IgnoreOther
=== RUN   TestInterpret/LowercaseFNoSegment
=== RUN   TestInterpret/LeftoverBracket
--- PASS: TestInterpret (0.00s)
PASS
ok  	lsystem	0.433s
```

Output file exists:
```
turtle.go
```

[done — exit criteria already satisfied on disk; no write needed]

```
