# Bead 2: grammar-modules

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 434s (7m)  
**Final exit criterion:** `grep -q 'func TestModuleSequence' grammar_modules_test.go && go test -v -run TestModuleSequence ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** grammar-modules  
**Output files:** grammar_modules.go, grammar_modules_test.go  
**Exit criteria:** `grep -q 'func TestModuleSequence' grammar_modules_test.go && go test -v -run TestModuleSequence ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement the shared L-system types and module-sequence parsing in `grammar_modules.go`.

1. Define the types `Module`, `BodyModule`, `Rule`, and `System` as specified in the survey.
2. Implement `parseBodyModuleSequence(src string) ([]BodyModule, error)`: a left-to-right scan. Whitespace is skipped. `[` and `]` are individual modules. A letter or `+`/`-` is a symbol; if followed by `(`, parse the argument list up to the matching `)` (handling nested parens). Split the argument text on top-level commas and parse each with `Lex` and `ParseExpr`. Store these as `[]Expr` (unevaluated). Error on malformed modules, unbalanced parens, or trailing/double commas.
3. Implement `parseModuleSequence(src string) ([]Module, error)`: calls `parseBodyModuleSequence`, then evaluates every argument expression using `EvalExpr` with an empty environment. Any reference to a parameter is an error.

Pin — `grammar-modules` bead, module-sequence length (`axiom` values and rule bodies): a module sequence is a **flat list**; `[` and `]` are each their own one-element module in it (`Sym` `'['` or `']'`, no args), never merged with a neighbour, with each other, or into an argument, and never nested at parse time. Worked: `F[+]` → **4** modules `F`, `[`, `+`, `]`; `F[+F]F` → **6** modules `F`, `[`, `+`, `F`, `]`, `F`; `F(1)[+(25)F(1)]` → **5** modules `F(1)`, `[`, `+(25)`, `F(1)`, `]`; `[]` → **2**.

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `grammar-modules` bead, module-sequence length (`axiom` values and rule bodies):**
  a module sequence is a **flat list**; `[` and `]` are each their own one-element
  module in it (`Sym` `'['` or `']'`, no args), never merged with a neighbour, with
  each other, or into an argument, and never nested at parse time. Worked:
  `F[+]` → **4** modules `F`, `[`, `+`, `]`; `F[+F]F` → **6** modules
  `F`, `[`, `+`, `F`, `]`, `F`; `F(1)[+(25)F(1)]` → **5** modules
  `F(1)`, `[`, `+(25)`, `F(1)`, `]`; `[]` → **2**. Matching `[` to `]` is the turtle's
  job (`Interpret`'s stack), never the parser's. (Consistent with the `rewrite` step-1
  worked example above, which lists `F(1)`, `[`, `+(25)`, `A(0.6)`, `]`, `[`, `-(25)`,
  `A(0.6)`, `]` as nine modules.)

### Revision 2 — created by REVISE_PENDING

**Title:** grammar-modules  
**Output files:** grammar_modules.go, grammar_modules_test.go  
**Exit criteria:** `grep -q 'func TestModuleSequence' grammar_modules_test.go && go test -v -run TestModuleSequence ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement the shared L-system types and module-sequence parsing in `grammar_modules.go`.

1. Define the types `Module`, `BodyModule`, `Rule`, and `System` as specified in the survey.
2. Implement `parseBodyModuleSequence(src string) ([]BodyModule, error)`: a left-to-right scan. Whitespace is skipped. `[` and `]` are individual modules. A letter or `+`/`-` is a symbol; if followed by `(`, parse the argument list up to the matching `)` (handling nested parens). Split the argument text on top-level commas and parse each with `Lex` and `ParseExpr` (already implemented in `expr.go`) as `[]Expr` (unevaluated). Error on malformed modules, unbalanced parens, or trailing/double commas.
3. Implement `parseModuleSequence(src string) ([]Module, error)`: calls `parseBodyModuleSequence`, then evaluates every argument expression using `EvalExpr` (already implemented in `expr.go`) with an empty environment. Any reference to a parameter is an error.

Pin — `grammar-modules` bead, module-sequence length (`axiom` values and rule bodies): a module sequence is a **flat list**; `[` and `]` are each their own one-element module in it (`Sym` `'['` or `']'`, no args), never merged with a neighbour, with each other, or into an argument, and never nested at parse time. Worked: `F[+]` → **4** modules `F`, `[`, `+`, `]`; `F[+F]F` → **6** modules `F`, `[`, `+`, `F`, `]`, `F`; `F(1)[+(25)F(1)]` → **5** modules `F(1)`, `[`, `+(25)`, `F(1)`, `]`; `[]` → **2**.

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `grammar-modules` bead, module-sequence length (`axiom` values and rule bodies):** a module sequence is a **flat list**; `[` and `]` are each their own one-element
  module in it (`Sym` `'['` or `']'`, no args), never merged with a neighbour, with
  each other, or into an argument, and never nested at parse time. Worked:
  `F[+]` → **4** modules `F`, `[`, `+`, `]`; `F[+F]F` → **6** modules
  `F`, `[`, `+`, `F`, `]`, `F`; `F(1)[+(25)F(1)]` → **5** modules
  `F(1)`, `[`, `+(25)`, `F(1)`, `]`; `[]` → **2**. Matching `[` to `]` is the turtle's
  job (`Interpret`'s stack), never the parser's. (Consistent with the `rewrite` step-1
  worked example above, which lists `F(1)`, `[`, `+(25)`, `A(0.6)`, `]`, `[`, `-(25)`,
  `A(0.6)`, `]` as nine modules.)

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 2 | success | 434s | no fire | 0/0 | PASS |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success):
grep -q 'func TestModuleSequence' grammar_modules_test.go && go test -v -run TestModuleSequence ./...
=== RUN   TestModuleSequence
--- PASS: TestModuleSequence (0.00s)
PASS
ok   lsystem    0.433s


## Final Output Files

*State of output_files on disk at report time.*

### grammar_modules.go

```go
package main

import (
	"fmt"
	"strings"
)

type Module struct {
	Sym  byte
	Args []float64
}

type BodyModule struct {
	Sym  byte
	Args []Expr
}

type Rule struct {
	Head   byte
	Params []string
	Body   []BodyModule
}

type System struct {
	Angle      float64
	Step       float64
	Heading    float64
	Iterations int
	Axiom      []Module
	Rules      []Rule
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func isLetter(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

func parseBodyModuleSequence(src string) ([]BodyModule, error) {
	n := len(src)
	i := 0
	mods := []BodyModule{}
	for i < n {
		c := src[i]
		if isSpace(c) {
			i++
			continue
		}
		if c == '[' || c == ']' {
			mods = append(mods, BodyModule{Sym: c, Args: nil})
			i++
			continue
		}
		if isLetter(c) || c == '+' || c == '-' {
			sym := c
			i++
			// skip whitespace to check for '('
			j := i
			for j < n && isSpace(src[j]) {
				j++
			}
			if j < n && src[j] == '(' {
				// find matching ')'
				depth := 1
				k := j + 1
				for k < n && depth > 0 {
					ch := src[k]
					if ch == '(' {
						depth++
					} else if ch == ')' {
						depth--
					}
					k++
				}
				if depth != 0 {
					return nil, fmt.Errorf("unbalanced parens")
				}
				argText := src[j+1 : k-1]
				args := []Expr{}
				if strings.TrimSpace(argText) != "" {
					parts := []string{}
					start := 0
					depth2 := 0
					for idx := 0; idx < len(argText); idx++ {
						ch := argText[idx]
						if ch == '(' {
							depth2++
						} else if ch == ')' {
							depth2--
						} else if ch == ',' && depth2 == 0 {
							part := argText[start:idx]
							if strings.TrimSpace(part) == "" {
								return nil, fmt.Errorf("trailing/double comma")
							}
							parts = append(parts, part)
							start = idx + 1
						}
					}
					lastPart := argText[start:]
					if strings.TrimSpace(lastPart) == "" {
						return nil, fmt.Errorf("trailing comma")
					}
					parts = append(parts, lastPart)

					for _, p := range parts {
						toks, err := Lex(p)
						if err != nil {
							return nil, err
						}
						expr, next, err := ParseExpr(toks, 0)
						if err != nil {
							return nil, err
						}
						if next != len(toks) {
							return nil, fmt.Errorf("extra tokens in expression")
						}
						args = append(args, expr)
					}
				}
				mods = append(mods, BodyModule{Sym: sym, Args: args})
				i = k
				continue
			}
			// no args
			mods = append(mods, BodyModule{Sym: sym, Args: []Expr{}})
			continue
		}
		return nil, fmt.Errorf("invalid character %c", c)
	}
	return mods, nil
}

func parseModuleSequence(src string) ([]Module, error) {
	bodyMods, err := parseBodyModuleSequence(src)
	if err != nil {
		return nil, err
	}
	env := map[string]float64{}
	mods := make([]Module, len(bodyMods))
	for idx, bm := range bodyMods {
		argsVals := make([]float64, len(bm.Args))
		for j, e := range bm.Args {
			v, err := EvalExpr(e, env)
			if err != nil {
				return nil, err
			}
			argsVals[j] = v
		}
		mods[idx] = Module{Sym: bm.Sym, Args: argsVals}
	}
	return mods, nil
}

```

### grammar_modules_test.go

```go
package main

import "testing"

func TestModuleSequence(t *testing.T) {
	mods, err := parseBodyModuleSequence("F[+]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mods) != 4 {
		t.Fatalf("expected 4 modules, got %d", len(mods))
	}
	if mods[0].Sym != 'F' || len(mods[0].Args) != 0 {
		t.Fatalf("mod0 mismatch")
	}
	if mods[1].Sym != '[' {
		t.Fatalf("mod1 mismatch")
	}
	if mods[2].Sym != '+' {
		t.Fatalf("mod2 mismatch")
	}
	if mods[3].Sym != ']' {
		t.Fatalf("mod3 mismatch")
	}

	mods, err = parseBodyModuleSequence("F[+F]F")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mods) != 6 {
		t.Fatalf("expected 6 modules, got %d", len(mods))
	}

	mods, err = parseBodyModuleSequence("F(1)[+(25)F(1)]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mods) != 5 {
		t.Fatalf("expected 5 modules, got %d", len(mods))
	}
	if mods[0].Sym != 'F' || len(mods[0].Args) != 1 {
		t.Fatalf("mod0 args")
	}
	if mods[0].Args[0] == nil {
		t.Fatalf("mod0 arg not parsed as Expr")
	}
	if mods[1].Sym != '[' {
		t.Fatalf("mod1")
	}
	if mods[2].Sym != '+' || len(mods[2].Args) != 1 {
		t.Fatalf("mod2")
	}
	if mods[2].Args[0] == nil {
		t.Fatalf("mod2 arg not parsed as Expr")
	}
	if mods[3].Sym != 'F' || len(mods[3].Args) != 1 {
		t.Fatalf("mod3")
	}
	if mods[4].Sym != ']' {
		t.Fatalf("mod4")
	}

	mods, err = parseBodyModuleSequence("[]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mods) != 2 {
		t.Fatalf("expected 2 modules, got %d", len(mods))
	}

	mods, err = parseBodyModuleSequence(" F [ + F ] ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mods) != 5 {
		t.Fatalf("expected 5 modules with whitespace, got %d", len(mods))
	}

	ms, err := parseModuleSequence("F(1)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ms) != 1 || ms[0].Sym != 'F' || len(ms[0].Args) != 1 || ms[0].Args[0] != 1 {
		t.Fatalf("module eval mismatch")
	}

	ms, err = parseModuleSequence("F(1+2*3)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ms[0].Args[0] != 7 {
		t.Fatalf("expected 7, got %v", ms[0].Args[0])
	}

	_, err = parseModuleSequence("F(s)")
	if err == nil {
		t.Fatalf("expected error for unbound param")
	}

	_, err = parseBodyModuleSequence("F(1")
	if err == nil {
		t.Fatalf("expected error for unbalanced parens")
	}

	_, err = parseBodyModuleSequence("F(1,)")
	if err == nil {
		t.Fatalf("expected error for trailing comma")
	}

	_, err = parseBodyModuleSequence("F(1,,2)")
	if err == nil {
		t.Fatalf("expected error for double comma")
	}
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-2-attempt-1.log`*

```
			v, err := EvalExpr(e, env)
			if err != nil {
				return nil, err
			}
			argsVals[j] = v
		}
		mods[idx] = Module{Sym: bm.Sym, Args: argsVals}
	}
	return mods, nil
}
 path:grammar_modules.go]]
[result]
ok: wrote 3103 bytes to grammar_modules.go
[TURN 2]
Now run exit criteria.

Check grep and go test.

[tool: run_command map[command:grep -q 'func TestModuleSequence' grammar_modules_test.go && go test -v -run TestModuleSequence ./...]]
[result]
stdout:
=== RUN   TestModuleSequence
--- PASS: TestModuleSequence (0.00s)
PASS
ok  	lsystem	0.433s

exit: 0
[TURN 3]
Test passes.

Now confirm file exists.

[tool: run_command map[command:ls grammar_modules.go]]
[result]
stdout:
grammar_modules.go

exit: 0
[TURN 4]
All good.

Now final message.All exit criteria pass.

- `grammar_modules.go` exists and implements `parseBodyModuleSequence` and `parseModuleSequence` with the required flat module-sequence semantics.
- `go test -v -run TestModuleSequence ./...` passes:

```
=== RUN   TestModuleSequence
--- PASS: TestModuleSequence (0.00s)
PASS
ok  	lsystem	0.433s
```

Output file verified:
```
grammar_modules.go
```

[done — exit criteria already satisfied on disk; no write needed]

```
