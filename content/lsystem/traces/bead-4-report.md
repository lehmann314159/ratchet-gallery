# Bead 4: grammar-system

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 222s (3m)  
**Final exit criterion:** `grep -q 'func TestParseSystem' grammar_system_test.go && go test -v -run TestParseSystem ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** grammar-system  
**Output files:** grammar_system.go, grammar_system_test.go  
**Exit criteria:** `grep -q 'func TestParseSystem' grammar_system_test.go && go test -v -run TestParseSystem ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement `ParseSystem(src string) (System, error)` in `grammar_system.go`.

1. Split `src` by line. Ignore blank lines and lines starting with `#`.
2. If a line contains `->`, call `parseRule(line, n)` and add to `System.Rules`.
3. Otherwise, the line must be `key: value`. Keys are `angle`, `step`, `heading`, `iterations`, `axiom`.
   - `angle`, `step`, `heading`: `strconv.ParseFloat`.
   - `iterations`: `strconv.Atoi`.
   - `axiom`: `parseModuleSequence(value)`. An error here is a parse error for line `n`.
4. Error on unknown keys or lines that are neither comments, rules, nor key-value pairs. All errors must include the 1-based line number.
5. Ensure `axiom` is present; otherwise, return an error.
6. Clamp `Iterations` to `[0, 12]` before returning.
7. Set defaults: `Angle 90`, `Step 1`, `Heading 0`, `Iterations 0`.

### Revision 2 — created by REVISE_PENDING

**Title:** grammar-system  
**Output files:** grammar_system.go, grammar_system_test.go  
**Exit criteria:** `grep -q 'func TestParseSystem' grammar_system_test.go && go test -v -run TestParseSystem ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement `ParseSystem(src string) (System, error)` in `grammar_system.go`.

1. Split `src` by line. Ignore blank lines and lines starting with `#`.
2. If a line contains `->`, call the already-implemented `parseRule(line, n)` from `grammar_rules.go` and add to `System.Rules`.
3. Otherwise, the line must be `key: value`. Keys are `angle`, `step`, `heading`, `iterations`, `axiom`.
   - `angle`, `step`, `heading`: `strconv.ParseFloat`.
   - `iterations`: `strconv.Atoi`.
   - `axiom`: `parseModuleSequence(value)`. An error here is a parse error for line `n`.
4. Error on unknown keys or lines that are neither comments, rules, nor key-value pairs. All errors must include the 1-based line number.
5. Ensure `axiom` is present; otherwise, return an error.
6. Clamp `Iterations` to `[0, 12]` before returning.
7. Set defaults: `Angle 90`, `Step 1`, `Heading 0`, `Iterations 0`.

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 4 | success | 222s | no fire | 0/0 | not run |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success):


## Final Output Files

*State of output_files on disk at report time.*

### grammar_system.go

```go
package main

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseSystem(src string) (System, error) {
	sys := System{
		Angle:   90,
		Step:    1,
		Heading: 0,
		Iterations: 0,
	}
	axiomSet := false
	lines := strings.Split(src, "\n")
	for i, raw := range lines {
		n := i + 1
		line := strings.TrimRight(raw, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.Contains(line, "->") {
			rule, err := parseRule(line, n)
			if err != nil {
				return System{}, err
			}
			sys.Rules = append(sys.Rules, rule)
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return System{}, fmt.Errorf("%d: invalid line", n)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "angle":
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return System{}, fmt.Errorf("%d: %v", n, err)
			}
			sys.Angle = f
		case "step":
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return System{}, fmt.Errorf("%d: %v", n, err)
			}
			sys.Step = f
		case "heading":
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return System{}, fmt.Errorf("%d: %v", n, err)
			}
			sys.Heading = f
		case "iterations":
			v, err := strconv.Atoi(val)
			if err != nil {
				return System{}, fmt.Errorf("%d: %v", n, err)
			}
			sys.Iterations = v
		case "axiom":
			mods, err := parseModuleSequence(val)
			if err != nil {
				return System{}, fmt.Errorf("%d: %v", n, err)
			}
			sys.Axiom = mods
			axiomSet = true
		default:
			return System{}, fmt.Errorf("%d: unknown key %s", n, key)
		}
	}
	if !axiomSet {
		return System{}, fmt.Errorf("axiom missing")
	}
	if sys.Iterations < 0 {
		sys.Iterations = 0
	} else if sys.Iterations > 12 {
		sys.Iterations = 12
	}
	return sys, nil
}

```

### grammar_system_test.go

```go
package main

import (
	"strings"
	"testing"
)

func TestParseSystem(t *testing.T) {
	src := "angle: 90\nstep: 2\nheading: 45\niterations: 5\naxiom: F\nF -> FF\n"
	sys, err := ParseSystem(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sys.Angle != 90 {
		t.Errorf("Angle = %v, want 90", sys.Angle)
	}
	if sys.Step != 2 {
		t.Errorf("Step = %v, want 2", sys.Step)
	}
	if sys.Heading != 45 {
		t.Errorf("Heading = %v, want 45", sys.Heading)
	}
	if sys.Iterations != 5 {
		t.Errorf("Iterations = %v, want 5", sys.Iterations)
	}
	if len(sys.Axiom) != 1 || sys.Axiom[0].Sym != 'F' {
		t.Errorf("Axiom mismatch: %+v", sys.Axiom)
	}
	if len(sys.Rules) != 1 || sys.Rules[0].Head != 'F' {
		t.Errorf("Rules mismatch: %+v", sys.Rules)
	}

	src2 := "axiom: A\n"
	sys2, err := ParseSystem(src2)
	if err != nil {
		t.Fatalf("unexpected error for defaults: %v", err)
	}
	if sys2.Angle != 90 || sys2.Step != 1 || sys2.Heading != 0 || sys2.Iterations != 0 {
		t.Errorf("defaults wrong: %+v", sys2)
	}

	srcLow := "axiom: A\niterations: -5\n"
	sysLow, err := ParseSystem(srcLow)
	if err != nil {
		t.Fatalf("unexpected error for clamp low: %v", err)
	}
	if sysLow.Iterations != 0 {
		t.Errorf("clamp low failed, got %d", sysLow.Iterations)
	}

	srcHigh := "axiom: A\niterations: 20\n"
	sysHigh, err := ParseSystem(srcHigh)
	if err != nil {
		t.Fatalf("unexpected error for clamp high: %v", err)
	}
	if sysHigh.Iterations != 12 {
		t.Errorf("clamp high failed, got %d", sysHigh.Iterations)
	}

	_, err = ParseSystem("angle: 90\n")
	if err == nil {
		t.Errorf("expected error for missing axiom")
	}

	_, err = ParseSystem("foo: bar\naxiom: A\n")
	if err == nil || !strings.Contains(err.Error(), "1:") {
		t.Errorf("expected unknown key error with line number, got %v", err)
	}

	_, err = ParseSystem("garbage\naxiom: A\n")
	if err == nil || !strings.Contains(err.Error(), "1:") {
		t.Errorf("expected invalid line error with line number, got %v", err)
	}

	srcBlank := "# comment\n\nangle: 30\n\naxiom: B\n"
	sysBlank, err := ParseSystem(srcBlank)
	if err != nil {
		t.Fatalf("unexpected error for blanks/comments: %v", err)
	}
	if sysBlank.Angle != 30 {
		t.Errorf("angle not parsed with blanks/comments, got %v", sysBlank.Angle)
	}

	_, err = ParseSystem("axiom: F(1,2)\n")
	if err != nil {
		t.Errorf("expected no parse error for axiom with multiple args, got %v", err)
	}
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-4-attempt-1.log`*

```
		case "iterations":
			v, err := strconv.Atoi(val)
			if err != nil {
				return System{}, fmt.Errorf("%d: %v", n, err)
			}
			sys.Iterations = v
		case "axiom":
			mods, err := parseModuleSequence(val)
			if err != nil {
				return System{}, fmt.Errorf("%d: %v", n, err)
			}
			sys.Axiom = mods
			axiomSet = true
		default:
			return System{}, fmt.Errorf("%d: unknown key %s", n, key)
		}
	}
	if !axiomSet {
		return System{}, fmt.Errorf("axiom missing")
	}
	if sys.Iterations < 0 {
		sys.Iterations = 0
	} else if sys.Iterations > 12 {
		sys.Iterations = 12
	}
	return sys, nil
}
 path:grammar_system.go]]
[result]
ok: wrote 1862 bytes to grammar_system.go
[TURN 3]
Now run exit criteria.[tool: run_command map[command:grep -q 'func TestParseSystem' grammar_system_test.go && go test -v -run TestParseSystem ./...]]
[result]
stdout:
=== RUN   TestParseSystem
--- PASS: TestParseSystem (0.00s)
PASS
ok  	lsystem	0.438s

exit: 0
[TURN 4]
Test passes. Need to ensure file exists.[tool: run_command map[command:ls grammar_system.go]]
[result]
stdout:
grammar_system.go

exit: 0
[TURN 5]
All good.All exit criteria pass. The implementation is in place.

```bash
grep -q 'func TestParseSystem' grammar_system_test.go && go test -v -run TestParseSystem ./...
```
PASS

Output file exists:
- `grammar_system.go`

[done — exit criteria already satisfied on disk; no write needed]

```
