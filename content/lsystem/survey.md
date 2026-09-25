# Survey — lsystem

## Module
module: lsystem
package: main

## Files
- expr.go
- grammar_modules.go
- grammar_rules.go
- grammar_system.go
- rewrite.go
- turtle.go
- render.go
- studio.go
- save.go
- templates.go
- handlers.go
- main.go

## Types and Constants

```go
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

type Expr interface{ isExpr() }

type NumLit struct{ Val float64 }

type Ref struct{ Name string }

type Neg struct{ X Expr }

type BinOp struct {
	Op   byte
	L, R Expr
}

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

const MaxModules = 2000000

type Segment struct{ X1, Y1, X2, Y2 float64 }

const CanvasSize = 800

const Margin = 20

type Box struct{ MinX, MinY, MaxX, MaxY float64 }

const SavedDir = "saved"

type SavedImage struct {
	Filename string
	Title    string
}

const DefaultIterations = 4

type Preset struct {
	Name   string
	Label  string
	Source string
}

type PageView struct {
	Source     string
	SVG        template.HTML
	Iterations int
	Presets    []Preset
	Saved      []SavedImage
	Error      string
}
```

## Package-Level Variables

```go
// from save.go
var saveMu sync.Mutex

// from handlers.go
var Presets []Preset

// from main.go
var templates *template.Template
```

## do_not_use_this_test.go

```go
package main

var (
	_ = TokNum
	_ = TokIdent
	_ = TokPlus
	_ = TokMinus
	_ = TokStar
	_ = TokSlash
	_ = TokLParen
	_ = TokRParen
	_ = TokComma
	_ = Lex
	_ = ParseExpr
	_ = EvalExpr
	_ = ParseSystem
	_ = MaxModules
	_ = Derive
	_ = Interpret
	_ = CanvasSize
	_ = Margin
	_ = BBox
	_ = RenderSVG
	_ = Studio
	_ = SavedDir
	_ = SaveSVG
	_ = ListSaved
	_ = InitTemplates
	_ = RenderPage
	_ = RenderResult
	_ = DefaultIterations
	_ = Presets
	_ = HandleIndex
	_ = HandleRender
	_ = HandleSelect
	_ = HandleSave
)

```

## File Declarations

### expr.go

```go
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

func Lex(src string) ([]Token, error) { return nil, nil }

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

func ParseExpr(toks []Token, pos int) (Expr, int, error) { return nil, 0, nil }

func EvalExpr(e Expr, env map[string]float64) (float64, error) { return 0, nil }
```

### grammar_modules.go

```go
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

func parseBodyModuleSequence(src string) ([]BodyModule, error) { return nil, nil }

func parseModuleSequence(src string) ([]Module, error) { return nil, nil }
```

### grammar_rules.go

```go
func parseRule(line string, lineNum int) (Rule, error) { return Rule{}, nil }
```

### grammar_system.go

```go
func ParseSystem(src string) (System, error) { return System{}, nil }
```

### rewrite.go

```go
const MaxModules = 2000000

func Derive(sys System, iterations int) ([]Module, error) { return nil, nil }
```

### turtle.go

```go
type Segment struct{ X1, Y1, X2, Y2 float64 }

func Interpret(mods []Module, angleDefault, stepDefault, headingStart float64) ([]Segment, error) { return nil, nil }
```

### render.go

```go
const CanvasSize = 800
const Margin = 20

type Box struct{ MinX, MinY, MaxX, MaxY float64 }

func BBox(segs []Segment) (Box, bool) { return Box{}, false }

func RenderSVG(segs []Segment) string { return "" }
```

### studio.go

```go
func Studio(source string, iterations int) (string, error) { return "", nil }
```

### save.go

```go
const SavedDir = "saved"

var saveMu sync.Mutex

type SavedImage struct {
	Filename string
	Title    string
}

func SaveSVG(svg string, dir string, now time.Time) (string, error) { return "", nil }

func ListSaved(dir string) ([]SavedImage, error) { return nil, nil }
```

### templates.go

```go
func InitTemplates() *template.Template { return nil }

func RenderPage(w http.ResponseWriter, v PageView) {}

func RenderResult(w http.ResponseWriter, v PageView) {}
```

### handlers.go

```go
const DefaultIterations = 4

type Preset struct {
	Name   string
	Label  string
	Source string
}

var Presets []Preset

type PageView struct {
	Source     string
	SVG        template.HTML
	Iterations int
	Presets    []Preset
	Saved      []SavedImage
	Error      string
}

func assemble(source string, svg string, iterations int, errMsg string) PageView { return PageView{} }

func HandleIndex(w http.ResponseWriter, r *http.Request) {}

func HandleRender(w http.ResponseWriter, r *http.Request) {}

func HandleSelect(w http.ResponseWriter, r *http.Request) {}

func HandleSave(w http.ResponseWriter, r *http.Request) {}
```

### main.go

```go
var templates *template.Template

func main() {}
```
