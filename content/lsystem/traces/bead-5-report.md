# Bead 5: rewrite

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 379s (6m)  
**Final exit criterion:** `grep -q 'func TestDerive' rewrite_test.go && go test -v -run TestDerive ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** rewrite  
**Output files:** rewrite.go, rewrite_test.go  
**Exit criteria:** `grep -q 'func TestDerive' rewrite_test.go && go test -v -run TestDerive ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement `Derive(sys System, iterations int) ([]Module, error)` in `rewrite.go`.

1. Perform `iterations` parallel rewrite steps. In each step, iterate over the current string and build a new string.
2. For each module `m`:
   - If `m.Sym` is `[` or `]`, copy it verbatim.
   - Find the first rule where `rule.Head == m.Sym` AND `len(rule.Params) == len(m.Args)`. 
   - If no match, copy `m` verbatim (identity).
   - If match, bind `Params[i] -> Args[i]` and evaluate each `BodyModule` in the rule's body using `EvalExpr`. 
3. If the resulting string exceeds `MaxModules` (2,000,000), return error `"system too large (>2000000 modules)"`.
4. The result must not alias `sys.Axiom`. Even if `iterations == 0`, return a fresh slice with copied modules.

Pin — `rewrite` bead, `Derive` (non-parametric): `axiom: A`, `A -> AB`, `B -> A`. Symbol strings: step 0 `A`; step 1 `AB`; step 2 `ABA`; step 3 `ABAAB`; step 4 `ABAABABA`.

Pin — `rewrite` bead, `Derive` (parametric): `axiom: A(1)`, `A(s) -> F(s) A(s*0.5)`. Step 1 modules `F(1) A(0.5)`; step 2 `F(1) F(0.5) A(0.25)`; step 3 `F(1) F(0.5) F(0.25) A(0.125)`.

Pin — `rewrite` bead, `Derive` (branch + arity): `axiom: A(1)`, `A(s) -> F(s)[+(25)A(s*0.6)][-(25)A(s*0.6)]` → step 1 module sequence `F(1) [ +(25) A(0.6) ] [ -(25) A(0.6) ]`. Separately: `axiom: A A(2)`, `A(s) -> F(s)` → step 1 `A F(2)` (bare `A` unmatched, copied).

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `rewrite` bead, `Derive` (non-parametric):** `axiom: A`, `A -> AB`, `B -> A`.
  Symbol strings: step 0 `A`; step 1 `AB`; step 2 `ABA`; step 3 `ABAAB`; step 4
  `ABAABABA`.

- **Pin — `rewrite` bead, `Derive` (parametric):** `axiom: A(1)`, `A(s) -> F(s) A(s*0.5)`.
  Step 1 modules `F(1) A(0.5)`; step 2 `F(1) F(0.5) A(0.25)`; step 3
  `F(1) F(0.5) F(0.25) A(0.125)`.

- **Pin — `rewrite` bead, `Derive` (branch + arity):** `axiom: A(1)`,
  `A(s) -> F(s)[+(25)A(s*0.6)][-(25)A(s*0.6)]` → step 1 module sequence
  `F(1) [ +(25) A(0.6) ] [ -(25) A(0.6) ]`. Separately: `axiom: A A(2)`, `A(s) -> F(s)`
  → step 1 `A F(2)` (bare `A` unmatched, copied).

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 5 | success | 379s | no fire | 0/0 | not run |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success): no test run, all tests passed in post-execution test run

## Final Output Files

*State of output_files on disk at report time.*

### rewrite.go

```go
package main

import "fmt"

const MaxModules = 2000000

func Derive(sys System, iterations int) ([]Module, error) {
    // Check initial size
    if len(sys.Axiom) > MaxModules {
        return nil, fmt.Errorf("system too large (>2000000 modules)")
    }
    // Create fresh copy of axiom
    current := make([]Module, len(sys.Axiom))
    for i, m := range sys.Axiom {
        argsCopy := make([]float64, len(m.Args))
        copy(argsCopy, m.Args)
        current[i] = Module{Sym: m.Sym, Args: argsCopy}
    }
    if iterations == 0 {
        return current, nil
    }
    for iter := 0; iter < iterations; iter++ {
        next := make([]Module, 0, len(current)*2)
        for _, m := range current {
            if m.Sym == '[' || m.Sym == ']' {
                argsCopy := make([]float64, len(m.Args))
                copy(argsCopy, m.Args)
                next = append(next, Module{Sym: m.Sym, Args: argsCopy})
                continue
            }
            matched := false
            for _, rule := range sys.Rules {
                if rule.Head == m.Sym && len(rule.Params) == len(m.Args) {
                    matched = true
                    env := make(map[string]float64, len(rule.Params))
                    for i, p := range rule.Params {
                        env[p] = m.Args[i]
                    }
                    for _, bm := range rule.Body {
                        argsVals := make([]float64, len(bm.Args))
                        for j, expr := range bm.Args {
                            v, err := EvalExpr(expr, env)
                            if err != nil {
                                return nil, err
                            }
                            argsVals[j] = v
                        }
                        next = append(next, Module{Sym: bm.Sym, Args: argsVals})
                    }
                    break
                }
            }
            if !matched {
                argsCopy := make([]float64, len(m.Args))
                copy(argsCopy, m.Args)
                next = append(next, Module{Sym: m.Sym, Args: argsCopy})
            }
        }
        if len(next) > MaxModules {
            return nil, fmt.Errorf("system too large (>2000000 modules)")
        }
        current = next
    }
    return current, nil
}

```

### rewrite_test.go

```go
package main

import "testing"

func TestDerive(t *testing.T) {
	src := "axiom: A\nA -> AB\nB -> A\n"
	sys, err := ParseSystem(src)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	mods, err := Derive(sys, 0)
	if err != nil {
		t.Fatalf("derive0 err: %v", err)
	}
	if len(mods) != 1 || mods[0].Sym != 'A' {
		t.Fatalf("step0 fail")
	}
	mods[0].Sym = 'Z'
	if sys.Axiom[0].Sym != 'A' {
		t.Fatalf("alias detected")
	}

	mods, err = Derive(sys, 1)
	if err != nil {
		t.Fatalf("derive1 err: %v", err)
	}
	if len(mods) != 2 || mods[0].Sym != 'A' || mods[1].Sym != 'B' {
		t.Fatalf("step1 fail")
	}

	mods, err = Derive(sys, 2)
	if err != nil {
		t.Fatalf("derive2 err: %v", err)
	}
	expected := []byte{'A', 'B', 'A'}
	if len(mods) != 3 {
		t.Fatalf("step2 len")
	}
	for i, ch := range expected {
		if mods[i].Sym != ch {
			t.Fatalf("step2 mismatch at %d", i)
		}
	}

	mods, err = Derive(sys, 3)
	if err != nil {
		t.Fatalf("derive3 err: %v", err)
	}
	expected = []byte{'A', 'B', 'A', 'A', 'B'}
	if len(mods) != 5 {
		t.Fatalf("step3 len")
	}
	for i, ch := range expected {
		if mods[i].Sym != ch {
			t.Fatalf("step3 mismatch at %d", i)
		}
	}

	mods, err = Derive(sys, 4)
	if err != nil {
		t.Fatalf("derive4 err: %v", err)
	}
	expected = []byte{'A', 'B', 'A', 'A', 'B', 'A', 'B', 'A'}
	if len(mods) != 8 {
		t.Fatalf("step4 len")
	}
	for i, ch := range expected {
		if mods[i].Sym != ch {
			t.Fatalf("step4 mismatch at %d", i)
		}
	}

	src2 := "axiom: A(1)\nA(s) -> F(s) A(s*0.5)\n"
	sys2, err := ParseSystem(src2)
	if err != nil {
		t.Fatalf("parse2 err: %v", err)
	}
	mods, err = Derive(sys2, 1)
	if err != nil {
		t.Fatalf("derive p1 err: %v", err)
	}
	if len(mods) != 2 || mods[0].Sym != 'F' || mods[0].Args[0] != 1 || mods[1].Sym != 'A' || mods[1].Args[0] != 0.5 {
		t.Fatalf("p1 fail")
	}
	mods, err = Derive(sys2, 2)
	if err != nil {
		t.Fatalf("derive p2 err: %v", err)
	}
	if len(mods) != 3 || mods[0].Sym != 'F' || mods[0].Args[0] != 1 || mods[1].Sym != 'F' || mods[1].Args[0] != 0.5 || mods[2].Sym != 'A' || mods[2].Args[0] != 0.25 {
		t.Fatalf("p2 fail")
	}
	mods, err = Derive(sys2, 3)
	if err != nil {
		t.Fatalf("derive p3 err: %v", err)
	}
	if len(mods) != 4 || mods[0].Sym != 'F' || mods[0].Args[0] != 1 || mods[1].Sym != 'F' || mods[1].Args[0] != 0.5 || mods[2].Sym != 'F' || mods[2].Args[0] != 0.25 || mods[3].Sym != 'A' || mods[3].Args[0] != 0.125 {
		t.Fatalf("p3 fail")
	}

	src3 := "axiom: A(1)\nA(s) -> F(s)[+(25)A(s*0.6)][-(25)A(s*0.6)]\n"
	sys3, err := ParseSystem(src3)
	if err != nil {
		t.Fatalf("parse3 err: %v", err)
	}
	mods, err = Derive(sys3, 1)
	if err != nil {
		t.Fatalf("derive branch err: %v", err)
	}
	if len(mods) != 9 {
		t.Fatalf("branch len %d", len(mods))
	}
	if mods[0].Sym != 'F' || mods[0].Args[0] != 1 {
		t.Fatalf("branch0")
	}
	if mods[1].Sym != '[' {
		t.Fatalf("branch1")
	}
	if mods[2].Sym != '+' || mods[2].Args[0] != 25 {
		t.Fatalf("branch2")
	}
	if mods[3].Sym != 'A' || mods[3].Args[0] != 0.6 {
		t.Fatalf("branch3")
	}
	if mods[4].Sym != ']' {
		t.Fatalf("branch4")
	}
	if mods[5].Sym != '[' {
		t.Fatalf("branch5")
	}
	if mods[6].Sym != '-' || mods[6].Args[0] != 25 {
		t.Fatalf("branch6")
	}
	if mods[7].Sym != 'A' || mods[7].Args[0] != 0.6 {
		t.Fatalf("branch7")
	}
	if mods[8].Sym != ']' {
		t.Fatalf("branch8")
	}

	src4 := "axiom: A A(2)\nA(s) -> F(s)\n"
	sys4, err := ParseSystem(src4)
	if err != nil {
		t.Fatalf("parse4 err: %v", err)
	}
	mods, err = Derive(sys4, 1)
	if err != nil {
		t.Fatalf("derive arity err: %v", err)
	}
	if len(mods) != 2 {
		t.Fatalf("arity len")
	}
	if mods[0].Sym != 'A' || len(mods[0].Args) != 0 {
		t.Fatalf("arity0")
	}
	if mods[1].Sym != 'F' || mods[1].Args[0] != 2 {
		t.Fatalf("arity1")
	}

	src5 := "axiom: F\n"
	sys5, err := ParseSystem(src5)
	if err != nil {
		t.Fatalf("parse5 err: %v", err)
	}
	mods, err = Derive(sys5, 5)
	if err != nil {
		t.Fatalf("derive identity err: %v", err)
	}
	if len(mods) != 1 || mods[0].Sym != 'F' {
		t.Fatalf("identity fail")
	}

	bigAxiom := make([]Module, 2000001)
	for i := range bigAxiom {
		bigAxiom[i] = Module{Sym: 'A', Args: []float64{}}
	}
	sysBig := System{Axiom: bigAxiom}
	_, err = Derive(sysBig, 0)
	if err == nil {
		t.Fatalf("expected error for too large")
	}
	if err.Error() != "system too large (>2000000 modules)" {
		t.Fatalf("wrong error: %v", err)
	}
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-5-attempt-1.log`*

```
                    }
                    break
                }
            }
            if !matched {
                argsCopy := make([]float64, len(m.Args))
                copy(argsCopy, m.Args)
                next = append(next, Module{Sym: m.Sym, Args: argsCopy})
            }
        }
        if len(next) > MaxModules {
            return nil, fmt.Errorf("system too large (>2000000 modules)")
        }
        current = next
    }
    return current, nil
}
 path:rewrite.go]]
[result]
ok: wrote 2303 bytes to rewrite.go
[TURN 4]
Now run exit criteria: grep -q 'func TestDerive' rewrite_test.go && go test -v -run TestDerive ./...

Let's run.[tool: run_command map[command:grep -q 'func TestDerive' rewrite_test.go && go test -v -run TestDerive ./...]]
[result]
stdout:
=== RUN   TestDerive
--- PASS: TestDerive (0.02s)
PASS
ok  	lsystem	0.455s

exit: 0
[TURN 5]
Test passes. Now confirm file exists.[tool: run_command map[command:ls rewrite.go]]
[result]
stdout:
rewrite.go

exit: 0
[TURN 6]
All good.All exit criteria pass. The implementation is in place.

```bash
grep -q 'func TestDerive' rewrite_test.go && go test -v -run TestDerive ./...
```
Output:
```
=== RUN   TestDerive
--- PASS: TestDerive (0.02s)
PASS
ok  	lsystem	0.455s
```

Output file exists:
```
rewrite.go
```

[done — exit criteria already satisfied on disk; no write needed]

```
