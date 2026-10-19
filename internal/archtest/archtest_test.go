// Package archtest checks the rules of AGENTS.md on every package.
package archtest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

const module = "github.com/ELY-Open-Source/8bitMyGoUI/"

var tiers = map[string]int{
	"pixel": 0, "icon": 1,
	"layout": 2, "text": 2, "display": 2,
	"input": 3, "feedback": 4,
	"navigation": 5, "data": 5, "overlay": 5, "datetime": 5, "chart": 5, "fx": 5,
	"game": 6, "chat": 6, "media": 6, "code": 6, "account": 6, "files": 6,
	"project": 6, "shop": 6, "agent": 6, "desk": 6, "tabletop": 6,
	"gallery": 7, "cmd/gallery": 8, "cmd/web": 8,
}

// Styled ui widgets drawing outside the pixel style.
var bannedFuncs = set("Text Textf RichText Button PrimaryButton Link MenuButton Toggle ToggleGroup Segmented Toolbar " +
	"Checkbox CheckboxGroup Switch Radio Select Slider StepSlider RangeSlider Stepper Rating TextInput TextArea " +
	"NumberInput SearchField Combobox Autocomplete TokenField EditableText FindBar Calendar DateInput TimeInput " +
	"ColorPicker ColorWell Table Tree TreeItem Accordion AccordionItem Collapsible Form Field Fieldset Tabs " +
	"Sidebar SidebarItem SidebarSection Breadcrumbs BackButton ForwardButton Modal AlertDialog Popover Progress " +
	"Spinner Meter Badge Image Icon Avatar Divider Split SplitVertical Outline OutlineTable GridView")

// Element, Context and Painter methods that blend, ease or draw.
var bannedMethods = set("Background Border BorderColor BorderWidth BorderStyle Radius Shadow Gradient " +
	"LinearGradient Opacity Stripes Tooltip Menu ContextMenu Transition Animate AnimateWith Loop Toast " +
	"ToastAction AnimationFrame After FillPath StrokePath RichText Glyphs FillGradient Stroke StrokeDashed " +
	"StrokePathGradient FillPathGradient Material Icon Image Selected")

func set(s string) map[string]bool {
	out := map[string]bool{}
	for _, f := range strings.Fields(s) {
		out[f] = true
	}
	return out
}

type goFile struct {
	path, pkg string
	ast       *ast.File
	src       []byte
}

// files parses every Go file of the module.
func files(t *testing.T) []goFile {
	var out []goFile
	err := filepath.WalkDir("../..", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && p != "../.." && (strings.HasPrefix(d.Name(), ".") || d.Name() == "node_modules") {
			return filepath.SkipDir
		}
		if d.IsDir() || !strings.HasSuffix(p, ".go") {
			return nil
		}
		src, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		f, err := parser.ParseFile(token.NewFileSet(), p, src, parser.ParseComments)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel("../..", filepath.Dir(p))
		out = append(out, goFile{p, filepath.ToSlash(rel), f, src})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) < 20 {
		t.Fatalf("found %d Go files; the walk is broken", len(out))
	}
	return out
}

func TestImportsPointDown(t *testing.T) {
	for _, f := range files(t) {
		if strings.HasSuffix(f.path, "_test.go") || strings.HasPrefix(f.pkg, "tools/") || strings.HasPrefix(f.pkg, "internal/") {
			continue
		}
		tier, ok := tiers[f.pkg]
		if !ok {
			t.Errorf("%s: package %q has no tier in AGENTS.md", f.path, f.pkg)
			continue
		}
		for _, im := range f.ast.Imports {
			path, _ := strconv.Unquote(im.Path.Value)
			if !strings.HasPrefix(path, module) {
				continue
			}
			to := strings.TrimPrefix(path, module)
			if got, ok := tiers[to]; !ok || got >= tier {
				t.Errorf("%s: %s (tier %d) imports %s (tier %d)", f.path, f.pkg, tier, to, got)
			}
		}
	}
}

func TestNoBannedCalls(t *testing.T) {
	for _, f := range files(t) {
		// cmd/web drives a Tester and builds no views.
		if strings.HasSuffix(f.path, "_test.go") || f.pkg == "pixel" || f.pkg == "cmd/web" || strings.HasPrefix(f.pkg, "tools/") || strings.HasPrefix(f.pkg, "internal/") {
			continue
		}
		imported := map[string]string{}
		for _, im := range f.ast.Imports {
			path, _ := strconv.Unquote(im.Path.Value)
			name := path[strings.LastIndex(path, "/")+1:]
			if im.Name != nil {
				name = im.Name.Name
			}
			imported[name] = path
		}
		painters := map[*ast.Ident]bool{}
		var stack []ast.Node
		ast.Inspect(f.ast, func(n ast.Node) bool {
			if n == nil {
				stack = stack[:len(stack)-1]
				return true
			}
			var parent ast.Node
			if len(stack) > 0 {
				parent = stack[len(stack)-1]
			}
			stack = append(stack, n)
			switch n := n.(type) {
			case *ast.Field:
				if star, ok := n.Type.(*ast.StarExpr); ok && isUI(star.X, "Painter", imported) {
					for _, name := range n.Names {
						painters[name] = true
					}
				}
			case *ast.Ident:
				if n.Obj != nil && n.Obj.Kind == ast.Var && !painters[n] && isPainter(n, painters) && !intoCanvas(n, parent) {
					t.Errorf("%s: painter %s used outside k.Canvas; only pixel draws", f.path, n.Name)
				}
			case *ast.CallExpr:
				sel, ok := n.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				pkg := false
				if id, ok := sel.X.(*ast.Ident); ok {
					_, pkg = imported[id.Name]
				}
				if pkg && isUI(sel.X, "", imported) && bannedFuncs[sel.Sel.Name] || !pkg && bannedMethods[sel.Sel.Name] {
					t.Errorf("%s: banned call %s", f.path, sel.Sel.Name)
				}
				// An element's keys go through pixel.Shortcut; c's are the window's.
				if id, ok := sel.X.(*ast.Ident); sel.Sel.Name == "Shortcut" && !pkg && !(ok && id.Name == "c") {
					t.Errorf("%s: %s.Shortcut skips the enabled check; use pixel.Shortcut", f.path, receiver(sel.X))
				}
				// It moves the focus before an overlay notes its opener.
				if sel.Sel.Name == "AutoFocus" && !pkg {
					t.Errorf("%s: %s.AutoFocus loses an overlay's opener; use pixel.AutoFocus", f.path, receiver(sel.X))
				}
			}
			return true
		})
	}
}

// isPainter reports whether id refers to a *ui.Painter parameter.
func isPainter(id *ast.Ident, painters map[*ast.Ident]bool) bool {
	for p := range painters {
		if id.Obj == p.Obj {
			return true
		}
	}
	return false
}

// intoCanvas reports whether id is a Canvas call argument.
func intoCanvas(id *ast.Ident, parent ast.Node) bool {
	call, ok := parent.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "Canvas" && slices.Contains(call.Args, ast.Expr(id))
}

// isUI reports whether x names package ui, or ui.name.
func isUI(x ast.Expr, name string, imported map[string]string) bool {
	if sel, ok := x.(*ast.SelectorExpr); ok && name != "" {
		return sel.Sel.Name == name && isUI(sel.X, "", imported)
	}
	id, ok := x.(*ast.Ident)
	return ok && name == "" && imported[id.Name] == "github.com/egoist/mygo/ui"
}

var directive = regexp.MustCompile(`^//(go:|line )`)

func TestCommentsAndFileLength(t *testing.T) {
	for _, f := range files(t) {
		if n := strings.Count(string(f.src), "\n"); n >= 500 {
			t.Errorf("%s: %d lines, keep files under 500", f.path, n)
		}
		for _, g := range f.ast.Comments {
			lines := 0
			for _, c := range g.List {
				if directive.MatchString(c.Text) || c.Text == "//" {
					continue
				}
				lines++
				if words := len(strings.Fields(strings.TrimPrefix(c.Text, "//"))); words > 10 || strings.HasPrefix(c.Text, "/*") {
					t.Errorf("%s: comment of %d words: %s", f.path, words, c.Text)
				}
			}
			if lines > 1 {
				t.Errorf("%s: comment of %d lines near %q", f.path, lines, g.List[0].Text)
			}
		}
	}
}

// TestWebScripts applies the comment and length rules to scripts.
func TestWebScripts(t *testing.T) {
	marks := map[string]*regexp.Regexp{".js": regexp.MustCompile(`^\s*//\s?(.*)$`), ".sh": regexp.MustCompile(`^\s*#\s(.*)$`)}
	seen := 0
	err := filepath.WalkDir("../../web", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && slices.Contains([]string{"node_modules", "dist", "test-results", "playwright-report"}, d.Name()) {
			return filepath.SkipDir
		}
		mark := marks[filepath.Ext(p)]
		if d.IsDir() || mark == nil {
			return nil
		}
		seen++
		src, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		lines := strings.Split(string(src), "\n")
		if len(lines) > 500 {
			t.Errorf("%s: %d lines, keep files under 500", p, len(lines))
		}
		run := 0
		for i, l := range lines {
			m := mark.FindStringSubmatch(l)
			if m == nil {
				run = 0
				continue
			}
			if run++; run > 1 {
				t.Errorf("%s:%d: comment of %d lines", p, i+1, run)
			}
			if words := len(strings.Fields(m[1])); words > 10 {
				t.Errorf("%s:%d: comment of %d words: %s", p, i+1, words, strings.TrimSpace(l))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if seen < 4 {
		t.Fatalf("found %d web scripts; the walk is broken", seen)
	}
}

// TestComponentsTakeTheirOptions holds PRD 5.1's shape, one per file.
func TestComponentsTakeTheirOptions(t *testing.T) {
	for _, f := range files(t) {
		tier, ok := tiers[f.pkg]
		if !ok || tier < 2 || tier > 6 || strings.HasSuffix(f.path, "_test.go") {
			continue
		}
		var found []string
		for _, d := range f.ast.Decls {
			fn, ok := d.(*ast.FuncDecl)
			// Undrawn bases, as MyGo's own, are parts, not components.
			if !ok || fn.Recv != nil || !fn.Name.IsExported() || strings.HasSuffix(fn.Name.Name, "Base") || len(fn.Type.Params.List) == 0 || !isContext(fn.Type.Params.List[0].Type) {
				continue
			}
			// PRD G1: one component per file.
			if found = append(found, fn.Name.Name); len(found) > 1 {
				t.Errorf("%s holds %v; one component per file", f.path, found)
			}
			want := fn.Name.Name + "Options"
			if !slices.ContainsFunc(fn.Type.Params.List, func(p *ast.Field) bool { return typeName(p.Type) == want }) {
				t.Errorf("%s: %s takes no %s", f.path, fn.Name.Name, want)
			}
			if r := fn.Type.Results; r != nil && len(r.List) == 1 && strings.HasSuffix(typeName(r.List[0].Type), "Result") && typeName(r.List[0].Type) != fn.Name.Name+"Result" {
				t.Errorf("%s: %s returns %s", f.path, fn.Name.Name, typeName(r.List[0].Type))
			}
		}
	}
}

// isContext reports whether x is *ui.Context.
func isContext(x ast.Expr) bool {
	star, ok := x.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "Context"
}

// typeName is a plain or generic type's name, else "".
func typeName(x ast.Expr) string {
	switch x := x.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.IndexExpr:
		return typeName(x.X)
	case *ast.IndexListExpr:
		return typeName(x.X)
	}
	return ""
}

// receiver names a selector's receiver for a message.
func receiver(x ast.Expr) string {
	switch x := x.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return receiver(x.X) + "." + x.Sel.Name
	}
	return "an element"
}
