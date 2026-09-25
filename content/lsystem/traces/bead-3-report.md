# Bead 3: grammar-rules

**Status:** succeeded  
**Attempts:** 1  
**Wall time:** 248s (4m)  
**Final exit criterion:** `grep -q 'func TestParseRule' grammar_rules_test.go && go test -v -run TestParseRule ./...`  

---

## Spec History

### Revision 1 — created by DECOMPOSE_SPEC

**Title:** grammar-rules  
**Output files:** grammar_rules.go, grammar_rules_test.go  
**Exit criteria:** `grep -q 'func TestParseRule' grammar_rules_test.go && go test -v -run TestParseRule ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement `parseRule(line string, lineNum int) (Rule, error)` in `grammar_rules.go`.

1. Split the trimmed line on the first `->`. 
2. The head part: the substring before the first `(` (or all of it if no `(`) must be exactly one ASCII letter. If a `(` follows, parse the comma-separated formal parameter names. Error on empty head, multi-letter head, unclosed `(`, or empty parameter names.
3. The body part: everything after the first `->` is parsed using `parseBodyModuleSequence`.
4. Every error returned must be prefixed with the 1-based `lineNum`.

Pin — `grammar-rules` bead, parseRule malformed rule head: `A -> B` parses; `AB -> C` is an error (multi-letter head); `-> B` is an error, **not** a panic (empty head). Every error carries a 1-based line number.

Pin — `grammar-rules` bead, rule head vs. formal parameters: in a rule line `head -> body`, the rule head is the substring of `head` **before the first `(`** (or all of it when there is no `(`), and that substring must be exactly one ASCII letter. Any `(name, name, …)` that follows the letter — still before the `->` — is the formal-parameter list. Parameters are **never** parsed from the body. Worked: `A(s) -> F(s)` → `Rule.Head == 'A'`, `Rule.Params == ["s"]`, arity 1, body `F(s)`; `A -> B` → `Head 'A'`, `Params` empty; `A(x, y) -> …` → `Params ["x", "y"]`. An unclosed `(` in the head is an error; an empty parameter name (`A() -> …`, `A(,) -> …`) is an error.

Design doc Decomposition Notes pin for this bead (verbatim, must be followed exactly, not re-derived or paraphrased):
- **Pin — `grammar-rules` bead, `parseRule` malformed rule head:** `A -> B` parses;
  `AB -> C` is an error (multi-letter head); `-> B` is an error, **not** a panic (empty
  head). Every error carries a 1-based line number.

- **Pin — `grammar-rules` bead, rule head vs. formal parameters:** in a rule line
  `head -> body`, the rule head is the substring of `head` **before the first `(`** (or
  all of it when there is no `(`), and that substring must be exactly one ASCII letter.
  Any `(name, name, …)` that follows the letter — still before the `->` — is the
  formal-parameter list. Parameters are **never** parsed from the body. Worked:
  `A(s) -> F(s)` → `Rule.Head == 'A'`, `Rule.Params == ["s"]`, arity 1, body `F(s)`;
  `A -> B` → `Head 'A'`, `Params` empty; `A(x, y) -> …` → `Params ["x", "y"]`. An
  unclosed `(` in the head is an error; an empty parameter name (`A() -> …`,
  `A(,) -> …`) is an error.

### Revision 2 — created by REVISE_PENDING

**Title:** grammar-rules  
**Output files:** grammar_rules.go, grammar_rules_test.go  
**Exit criteria:** `grep -q 'func TestParseRule' grammar_rules_test.go && go test -v -run TestParseRule ./...`  
**Execution budget:** 900s  
**Monitor override:** honor  

Implement `parseRule(line string, lineNum int) (Rule, error)` in `grammar_rules.go`.

1. Split the trimmed line on the first `->`. 
2. The head part: the substring before the first `(` (or all of it if no `(`) must be exactly one ASCII letter. If a `(` follows, parse the comma-separated formal parameter names. Error on empty head, multi-letter head, unclosed `(`, or empty parameter names.
3. The body part: everything after the first `->` is parsed using `parseBodyModuleSequence` from `grammar_modules.go`.
4. Every error returned must be prefixed with the 1-based `lineNum`.

Pin — `grammar-rules` bead, parseRule malformed rule head: `A -> B` parses; `AB -> C` is an error (multi-letter head); `-> B` is an error, **not** a panic (empty head). Every error carries a 1-based line number.

Pin — `grammar-rules` bead, rule head vs. formal parameters: in a rule line `head -> body`, the rule head is the substring of `head` **before the first `(`** (or all of it when there is no `(`), and that substring must be exactly one ASCII letter. Any `(name, name, …)` that follows the letter — still before the `->` — is the formal-parameter list. Parameters are **never** parsed from the body. Worked: `A(s) -> F(s)` → `Rule.Head == 'A'`, `Rule.Params == [

## Attempt History

| # | Execution ID | Termination | Duration | Monitor | write_file ok/total | Last test result |
|---|---|---|---|---|---|---|
| 1 | 3 | success | 248s | no fire | 0/0 | not run |

## ADJUDICATE Decisions

### After attempt 1 → declare_success

## Compressed History

Attempt 1 (success):
No commands were run. All tests passed in post-execution test run.


## Final Output Files

*State of output_files on disk at report time.*

### grammar_rules.go

```go
package main

import (
	"fmt"
	"strings"
)

func parseRule(line string, lineNum int) (Rule, error) {
	lineTrim := strings.TrimSpace(line)
	if lineTrim == "" {
		return Rule{}, fmt.Errorf("%d: empty line", lineNum)
	}
	parts := strings.SplitN(lineTrim, "->", 2)
	if len(parts) != 2 {
		return Rule{}, fmt.Errorf("%d: missing ->", lineNum)
	}
	headRaw := strings.TrimSpace(parts[0])
	bodyRaw := strings.TrimSpace(parts[1])

	if headRaw == "" {
		return Rule{}, fmt.Errorf("%d: empty head", lineNum)
	}

	idx := strings.Index(headRaw, "(")
	var headChar byte
	var params []string

	if idx == -1 {
		// no params
		headStr := strings.TrimSpace(headRaw)
		if len(headStr) != 1 || !isLetter(headStr[0]) {
			return Rule{}, fmt.Errorf("%d: invalid head", lineNum)
		}
		headChar = headStr[0]
		params = []string{}
	} else {
		headPart := strings.TrimSpace(headRaw[:idx])
		paramPart := strings.TrimSpace(headRaw[idx:])
		if len(headPart) != 1 || !isLetter(headPart[0]) {
			return Rule{}, fmt.Errorf("%d: invalid head", lineNum)
		}
		headChar = headPart[0]
		if !strings.HasPrefix(paramPart, "(") || !strings.HasSuffix(paramPart, ")") {
			return Rule{}, fmt.Errorf("%d: unclosed paren in head", lineNum)
		}
		inner := strings.TrimSpace(paramPart[1 : len(paramPart)-1])
		if inner == "" {
			return Rule{}, fmt.Errorf("%d: empty parameter list", lineNum)
		}
		rawParams := strings.Split(inner, ",")
		params = make([]string, 0, len(rawParams))
		for _, p := range rawParams {
			pTrim := strings.TrimSpace(p)
			if pTrim == "" {
				return Rule{}, fmt.Errorf("%d: empty parameter name", lineNum)
			}
			params = append(params, pTrim)
		}
	}

	bodyMods, err := parseBodyModuleSequence(bodyRaw)
	if err != nil {
		return Rule{}, fmt.Errorf("%d: %v", lineNum, err)
	}

	return Rule{
		Head:   headChar,
		Params: params,
		Body:   bodyMods,
	}, nil
}

```

### grammar_rules_test.go

```go
package main

import "testing"

func TestParseRule(t *testing.T) {
	// valid simple rule
	r, err := parseRule("A -> B", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Head != 'A' {
		t.Fatalf("head mismatch")
	}
	if len(r.Params) != 0 {
		t.Fatalf("params not empty")
	}
	if len(r.Body) != 1 || r.Body[0].Sym != 'B' {
		t.Fatalf("body mismatch")
	}

	// valid with single param
	r, err = parseRule("A(s) -> F(s)", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Head != 'A' {
		t.Fatalf("head mismatch")
	}
	if len(r.Params) != 1 || r.Params[0] != "s" {
		t.Fatalf("params mismatch")
	}
	if len(r.Body) != 1 || r.Body[0].Sym != 'F' {
		t.Fatalf("body sym mismatch")
	}
	if len(r.Body[0].Args) != 1 {
		t.Fatalf("body args count")
	}

	// valid with multiple params
	r, err = parseRule("A(x, y) -> B", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Head != 'A' {
		t.Fatalf("head mismatch")
	}
	if len(r.Params) != 2 || r.Params[0] != "x" || r.Params[1] != "y" {
		t.Fatalf("params mismatch")
	}

	// multi-letter head error
	_, err = parseRule("AB -> C", 4)
	if err == nil {
		t.Fatalf("expected error for multi-letter head")
	}

	// empty head error
	_, err = parseRule("-> B", 5)
	if err == nil {
		t.Fatalf("expected error for empty head")
	}

	// unclosed paren in head
	_, err = parseRule("A(s -> B", 6)
	if err == nil {
		t.Fatalf("expected error for unclosed paren")
	}

	// empty param name
	_, err = parseRule("A() -> B", 7)
	if err == nil {
		t.Fatalf("expected error for empty param")
	}

	// empty param name with comma
	_, err = parseRule("A(,) -> B", 8)
	if err == nil {
		t.Fatalf("expected error for empty param comma")
	}

	// error prefix contains line number
	_, err = parseRule("AB -> C", 9)
	if err == nil {
		t.Fatalf("expected error")
	}
	if len(err.Error()) < 2 || err.Error()[0:2] != "9:" {
		t.Fatalf("error not prefixed with line number")
	}
}

```

## Last Trace Excerpt

*Final 60 lines of `bead-3-attempt-1.log`*

```
exit: 0
[TURN 5]
Test passes. Need to ensure output file exists.[tool: run_command map[command:ls grammar_rules.go]]
[result]
stdout:
grammar_rules.go

exit: 0
[TURN 6]
All good. Possibly need to handle empty head error message prefix? Test only checks prefix contains line number. Our error for empty head is fmt.Errorf("%d: empty head", lineNum). Good.

Check multi-letter head error: we return invalid head. That's fine.

Unclosed paren: we return unclosed paren in head.

Empty param name: we return empty parameter name. Test expects error for A() -> B and A(,) -> B. For A(), inner == "" -> we return empty parameter list. Test expects error for empty param. Might be okay. For A(,), inner = ","? Actually inner = ""? Wait paramPart = "()", inner = "" -> empty parameter list error. Good.

For A(,), inner = ","? Wait paramPart = "(,)", inner = ",". rawParams split -> ["", ""]. pTrim empty -> error empty parameter name. Good.

Potential edge: head part with spaces like "A (s)". Our idx finds '(' at position after space. headPart = strings.TrimSpace(headRaw[:idx]) => "A". Good. paramPart = strings.TrimSpace(headRaw[idx:]) => "(s)". Works.

Potential edge: head part with trailing spaces before '('? Works.

Potential edge: head part with multiple '('? Not needed.

Potential edge: head letter must be exactly one ASCII letter. Our check len(headStr)!=1 or !isLetter. Good.

Potential edge: head part may contain extra characters after head letter before '('? Eg "Afoo(s)". headPart = "Afoo". len !=1 -> error. Good.

Potential edge: head part may be "A(s)extra"? headRaw[:idx] = "A", paramPart = "(s)extra". HasPrefix '(' true, HasSuffix ')'? No, suffix is "a"? Actually paramPart ends with "extra". HasSuffix ')' false -> unclosed paren error. Might be okay.

Potential edge: body parsing error prefix. We wrap with fmt.Errorf("%d: %v", lineNum, err). Good.

Potential edge: empty line error? Not required.

Potential edge: missing -> error. Might be okay.

Test passes.

Now final message.All checks pass. The implementation is in place.

```bash
grep -q 'func TestParseRule' grammar_rules_test.go && go test -v -run TestParseRule ./...
```

Output:
```
=== RUN   TestParseRule
--- PASS: TestParseRule (0.00s)
PASS
ok  	lsystem	0.436s
```

Output file exists:
```
grammar_rules.go
```

[done — exit criteria already satisfied on disk; no write needed]

```
