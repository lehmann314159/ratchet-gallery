# Bead 1: expr

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 373s (6m)  
**Final exit criterion:** `grep -q 'func TestExpr' expr_test.go && go test -v -run TestExpr ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** expr  
**Output files:** expr.go, expr_test.go  
**Exit criteria:** `grep -q 'func TestExpr' expr_test.go && go test -v -run TestExpr ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement the expression sub-language in `expr.go`. 

1. Implement `Lex(src string) ([]Token, error)`: a single left-to-right scan. Whitespace is skipped. `+ - * / ( ) ,` are single tokens. Digit/dot runs are `TokNum`. Letter runs are `TokIdent`. Any other character is an error.
2. Implement `ParseExpr(toks []Token, pos int) (Expr, int, error)`: a recursive descent parser for the grammar:
   - expr := term (('+' | '-') term)*
   - term := factor (('*' | '/') factor)*
   - factor := NUMBER | IDENT | '(' expr ')' | '-' factor
   All operators are left-associative. Unary minus is a factor prefix. It must return the index of the first token NOT consumed.
3. Implement `EvalExpr(e Expr, env map[string]float64) (float64, error)`: recursive evaluation. `Ref` looks up names in `env` (error if absent). Division by exactly 0.0 is an error.

Pin — `expr` bead, `ParseExpr` + `EvalExpr`: with the stated environment, the parsed expression evaluates to: `2+3*4` → `14`; `-2*3` → `-6`; `1/3` → `0.3333333333333333`; `10-3-2` → `5` (left-associative, NOT `9`); `2-(3-1)` → `0`; `s*0.6+a` with `{s:2, a:10}` → `11.2`; `-s` with `{s:5}` → `-5`; `2*-3` → `-6`. `s+1` under the empty environment → error; `1/0` → error.

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `expr` bead, `ParseExpr` + `EvalExpr`:** with the stated environment, the parsed
  expression evaluates to: `2+3*4` → `14`; `-2*3` → `-6`; `1/3` → `0.3333333333333333`;
  `10-3-2` → `5` (left-associative, NOT `9`); `2-(3-1)` → `0`; `s*0.6+a` with
  `{s:2, a:10}` → `11.2`; `-s` with `{s:5}` → `-5`; `2*-3` → `-6`. `s+1` under the empty
  environment → error; `1/0` → error.

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 1 | success | 373s | no fire | 0/0 | not run |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success):
No commands were run.


## Final Output Files

*State of output_files on disk at report time.*

### expr.go

```go
package main

import (
	"fmt"
	"strconv"
)

type TokenKind int

const (
	TokNum TokenKind = iota
	TokIdent
	TokPlus
	TokMinus
	TokStar
	TokSlash
	TokLParen
	TokRParen
	TokComma
)

type Token struct {
	Kind TokenKind
	Text string
}

func Lex(src string) ([]Token, error) {
	var toks []Token
	i := 0
	n := len(src)
	for i < n {
		c := src[i]
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			i++
			continue
		}
		switch c {
		case '+':
			toks = append(toks, Token{Kind: TokPlus, Text: "+"})
			i++
			continue
		case '-':
			toks = append(toks, Token{Kind: TokMinus, Text: "-"})
			i++
			continue
		case '*':
			toks = append(toks, Token{Kind: TokStar, Text: "*"})
			i++
			continue
		case '/':
			toks = append(toks, Token{Kind: TokSlash, Text: "/"})
			i++
			continue
		case '(':
			toks = append(toks, Token{Kind: TokLParen, Text: "("})
			i++
			continue
		case ')':
			toks = append(toks, Token{Kind: TokRParen, Text: ")"})
			i++
			continue
		case ',':
			toks = append(toks, Token{Kind: TokComma, Text: ","})
			i++
			continue
		}
		if (c >= '0' && c <= '9') || c == '.' {
			start := i
			for i < n && ((src[i] >= '0' && src[i] <= '9') || src[i] == '.') {
				i++
			}
			text := src[start:i]
			toks = append(toks, Token{Kind: TokNum, Text: text})
			continue
		}
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			start := i
			for i < n && ((src[i] >= 'a' && src[i] <= 'z') || (src[i] >= 'A' && src[i] <= 'Z')) {
				i++
			}
			text := src[start:i]
			toks = append(toks, Token{Kind: TokIdent, Text: text})
			continue
		}
		return nil, fmt.Errorf("invalid character: %c", c)
	}
	return toks, nil
}

type Expr interface{ isExpr() }

type NumLit struct{ Val float64 }

func (NumLit) isExpr() {}

type Ref struct{ Name string }

func (Ref) isExpr() {}

type Neg struct{ X Expr }

func (Neg) isExpr() {}

type BinOp struct {
	Op   byte
	L, R Expr
}

func (BinOp) isExpr() {}

func ParseExpr(toks []Token, pos int) (Expr, int, error) {
	var parseExpr func(int) (Expr, int, error)
	var parseTerm func(int) (Expr, int, error)
	var parseFactor func(int) (Expr, int, error)

	parseFactor = func(p int) (Expr, int, error) {
		if p >= len(toks) {
			return nil, p, fmt.Errorf("unexpected end in factor")
		}
		tok := toks[p]
		switch tok.Kind {
		case TokNum:
			v, err := strconv.ParseFloat(tok.Text, 64)
			if err != nil {
				return nil, p, fmt.Errorf("invalid number %s", tok.Text)
			}
			return NumLit{Val: v}, p + 1, nil
		case TokIdent:
			return Ref{Name: tok.Text}, p + 1, nil
		case TokLParen:
			e, next, err := parseExpr(p + 1)
			if err != nil {
				return nil, p, err
			}
			if next >= len(toks) || toks[next].Kind != TokRParen {
				return nil, next, fmt.Errorf("expected ')'")
			}
			return e, next + 1, nil
		case TokMinus:
			x, next, err := parseFactor(p + 1)
			if err != nil {
				return nil, p, err
			}
			return Neg{X: x}, next, nil
		default:
			return nil, p, fmt.Errorf("unexpected token %v", tok.Kind)
		}
	}

	parseTerm = func(p int) (Expr, int, error) {
		left, p, err := parseFactor(p)
		if err != nil {
			return nil, p, err
		}
		for p < len(toks) {
			tok := toks[p]
			if tok.Kind == TokStar || tok.Kind == TokSlash {
				op := byte(tok.Text[0])
				p++
				right, next, err := parseFactor(p)
				if err != nil {
					return nil, p, err
				}
				left = BinOp{Op: op, L: left, R: right}
				p = next
			} else {
				break
			}
		}
		return left, p, nil
	}

	parseExpr = func(p int) (Expr, int, error) {
		left, p, err := parseTerm(p)
		if err != nil {
			return nil, p, err
		}
		for p < len(toks) {
			tok := toks[p]
			if tok.Kind == TokPlus || tok.Kind == TokMinus {
				op := byte(tok.Text[0])
				p++
				right, next, err := parseTerm(p)
				if err != nil {
					return nil, p, err
				}
				left = BinOp{Op: op, L: left, R: right}
				p = next
			} else {
				break
			}
		}
		return left, p, nil
	}

	return parseExpr(pos)
}

func EvalExpr(e Expr, env map[string]float64) (float64, error) {
	switch x := e.(type) {
	case NumLit:
		return x.Val, nil
	case Ref:
		v, ok := env[x.Name]
		if !ok {
			return 0, fmt.Errorf("unbound variable %s", x.Name)
		}
		return v, nil
	case Neg:
		v, err := EvalExpr(x.X, env)
		if err != nil {
			return 0, err
		}
		return -v, nil
	case BinOp:
		l, err := EvalExpr(x.L, env)
		if err != nil {
			return 0, err
		}
		r, err := EvalExpr(x.R, env)
		if err != nil {
			return 0, err
		}
		switch x.Op {
		case '+':
			return l + r, nil
		case '-':
			return l - r, nil
		case '*':
			return l * r, nil
		case '/':
			if r == 0.0 {
				return 0, fmt.Errorf("division by zero")
			}
			return l / r, nil
		default:
			return 0, fmt.Errorf("unknown operator %c", x.Op)
		}
	default:
		return 0, fmt.Errorf("unknown expr")
	}
}

```

### expr_test.go

```go
package main

import (
	"fmt"
	"testing"
)

func TestExpr(t *testing.T) {
	// Lex basic tokenization and whitespace skipping
	toks, err := Lex("2 + 3 * 4")
	if err != nil {
		t.Fatalf("Lex unexpected error: %v", err)
	}
	if len(toks) != 5 {
		t.Fatalf("expected 5 tokens, got %d", len(toks))
	}
	if toks[0].Kind != TokNum || toks[0].Text != "2" {
		t.Fatalf("token 0 mismatch")
	}
	if toks[1].Kind != TokPlus {
		t.Fatalf("token 1 not Plus")
	}
	if toks[2].Kind != TokNum || toks[2].Text != "3" {
		t.Fatalf("token 2 mismatch")
	}
	if toks[3].Kind != TokStar {
		t.Fatalf("token 3 not Star")
	}
	if toks[4].Kind != TokNum || toks[4].Text != "4" {
		t.Fatalf("token 4 mismatch")
	}

	// Lex error on unknown character
	_, err = Lex("2 $ 3")
	if err == nil {
		t.Fatalf("expected Lex error for invalid character")
	}

	// Helper to parse and evaluate
	eval := func(src string, env map[string]float64) (float64, error) {
		toks, err := Lex(src)
		if err != nil {
			return 0, err
		}
		e, next, err := ParseExpr(toks, 0)
		if err != nil {
			return 0, err
		}
		if next != len(toks) {
			return 0, fmt.Errorf("parse did not consume all tokens, next=%d len=%d", next, len(toks))
		}
		return EvalExpr(e, env)
	}

	// 2+3*4 → 14
	v, err := eval("2+3*4", map[string]float64{})
	if err != nil {
		t.Fatalf("eval 2+3*4 error: %v", err)
	}
	if v != 14 {
		t.Fatalf("2+3*4 = %v, want 14", v)
	}

	// -2*3 → -6
	v, err = eval("-2*3", map[string]float64{})
	if err != nil {
		t.Fatalf("eval -2*3 error: %v", err)
	}
	if v != -6 {
		t.Fatalf("-2*3 = %v, want -6", v)
	}

	// 1/3 → 0.3333333333333333
	v, err = eval("1/3", map[string]float64{})
	if err != nil {
		t.Fatalf("eval 1/3 error: %v", err)
	}
	if v != 0.3333333333333333 {
		t.Fatalf("1/3 = %v, want 0.3333333333333333", v)
	}

	// 10-3-2 → 5 left-associative
	v, err = eval("10-3-2", map[string]float64{})
	if err != nil {
		t.Fatalf("eval 10-3-2 error: %v", err)
	}
	if v != 5 {
		t.Fatalf("10-3-2 = %v, want 5", v)
	}

	// 2-(3-1) → 0
	v, err = eval("2-(3-1)", map[string]float64{})
	if err != nil {
		t.Fatalf("eval 2-(3-1) error: %v", err)
	}
	if v != 0 {
		t.Fatalf("2-(3-1) = %v, want 0", v)
	}

	// s*0.6+a with env
	v, err = eval("s*0.6+a", map[string]float64{"s": 2, "a": 10})
	if err != nil {
		t.Fatalf("eval s*0.6+a error: %v", err)
	}
	if v != 11.2 {
		t.Fatalf("s*0.6+a = %v, want 11.2", v)
	}

	// -s with env
	v, err = eval("-s", map[string]float64{"s": 5})
	if err != nil {
		t.Fatalf("eval -s error: %v", err)
	}
	if v != -5 {
		t.Fatalf("-s = %v, want -5", v)
	}

	// 2*-3 → -6
	v, err = eval("2*-3", map[string]float64{})
	if err != nil {
		t.Fatalf("eval 2*-3 error: %v", err)
	}
	if v != -6 {
		t.Fatalf("2*-3 = %v, want -6", v)
	}

	// s+1 under empty env → error
	_, err = eval("s+1", map[string]float64{})
	if err == nil {
		t.Fatalf("expected error for unbound s")
	}

	// 1/0 → error
	_, err = eval("1/0", map[string]float64{})
	if err == nil {
		t.Fatalf("expected error for division by zero")
	}
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-1-attempt-1.log`*

```
		case '/':
			if r == 0.0 {
				return 0, fmt.Errorf("division by zero")
			}
			return l / r, nil
		default:
			return 0, fmt.Errorf("unknown operator %c", x.Op)
		}
	default:
		return 0, fmt.Errorf("unknown expr")
	}
}
 path:expr.go]]
[result]
ok: wrote 4725 bytes to expr.go
[TURN 2]
File written. Now run exit criteria.

Run command: grep -q 'func TestExpr' expr_test.go && go test -v -run TestExpr ./...

Let's run.[tool: run_command map[command:grep -q 'func TestExpr' expr_test.go && go test -v -run TestExpr ./...]]
[result]
stdout:
=== RUN   TestExpr
--- PASS: TestExpr (0.00s)
PASS
ok  	lsystem	0.440s

exit: 0
[TURN 3]
Test passes. Need to confirm every Output File exists on disk. Output Files: expr.go. Exists.

We should also run ls to check.[tool: run_command map[command:ls expr.go]]
[result]
stdout:
expr.go

exit: 0
[TURN 4]
All good.

We should send final message.All exit criteria passed.

- `expr.go` has been implemented with `Lex`, `ParseExpr`, and `EvalExpr`
- Tests pass:

```
=== RUN   TestExpr
--- PASS: TestExpr (0.00s)
PASS
ok  	lsystem	0.440s
```

Output file exists:
```
expr.go
```

[done — exit criteria already satisfied on disk; no write needed]

```
