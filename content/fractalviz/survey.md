# Survey — fractalviz

## Module
module: fractalviz
package: main

## Files
- fractal.go
- render.go
- params.go
- save.go
- templates.go
- handlers.go
- main.go

## Types and Constants

```go
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

const ImageWidth = 600

const ImageHeight = 600

const SavedDir = "saved"

type SavedImage struct {
	Filename string
	Title    string
}

type TypeOption struct {
	Value    string
	Label    string
	Selected bool
}

type PageView struct {
	Params       Params
	ImageURL     string
	Types        []TypeOption
	Saved        []SavedImage
	ShowJulia    bool
	ShowExponent bool
}
```

## Package-Level Variables

```go
// from save.go
var saveMu sync.Mutex

// from main.go
var templates *template.Template
```

## do_not_use_this_test.go

```go
package main

var (
	_ = Mandelbrot
	_ = Julia
	_ = BurningShip
	_ = Multibrot
	_ = TypeName
	_ = ParseType
	_ = DefaultParams
	_ = Escape
	_ = ImageWidth
	_ = ImageHeight
	_ = PixelToComplex
	_ = Color
	_ = Render
	_ = ParseParams
	_ = ImageQuery
	_ = SavedDir
	_ = SaveImage
	_ = ListSaved
	_ = InitTemplates
	_ = RenderPage
	_ = RenderResult
	_ = HandleIndex
	_ = HandleRender
	_ = HandleImage
	_ = HandleSave
	_ = HandleSelect
)

```

## File Declarations

### fractal.go

```go
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
	return ""
}

func ParseType(s string) (FractalType, error) {
	return Mandelbrot, nil
}

func DefaultParams(t FractalType) Params {
	return Params{}
}

func Escape(p Params, point complex128) int {
	return 0
}
```

### render.go

```go
const ImageWidth = 600
const ImageHeight = 600

func PixelToComplex(p Params, px, py, width, height int) complex128 {
	return 0
}

func Color(n, maxIter int) color.RGBA {
	return color.RGBA{}
}

func Render(p Params, width, height int) *image.RGBA {
	return nil
}
```

### params.go

```go
func ParseParams(v url.Values) (Params, error) {
	return Params{}, nil
}

func ImageQuery(p Params) string {
	return ""
}
```

### save.go

```go
const SavedDir = "saved"

var saveMu sync.Mutex

type SavedImage struct {
	Filename string
	Title    string
}

func SaveImage(p Params, dir string, now time.Time) (string, error) {
	return "", nil
}

func ListSaved(dir string) ([]SavedImage, error) {
	return nil, nil
}
```

### templates.go

```go
func InitTemplates() *template.Template {
	return nil
}

func RenderPage(w http.ResponseWriter, v PageView) {
}

func RenderResult(w http.ResponseWriter, v PageView) {
}
```

### handlers.go

```go
type TypeOption struct {
	Value    string
	Label    string
	Selected bool
}

type PageView struct {
	Params       Params
	ImageURL     string
	Types        []TypeOption
	Saved        []SavedImage
	ShowJulia    bool
	ShowExponent bool
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
}

func HandleRender(w http.ResponseWriter, r *http.Request) {
}

func HandleImage(w http.ResponseWriter, r *http.Request) {
}

func HandleSave(w http.ResponseWriter, r *http.Request) {
}

func HandleSelect(w http.ResponseWriter, r *http.Request) {
}
```

### main.go

```go
var templates *template.Template

func main() {
}
```
