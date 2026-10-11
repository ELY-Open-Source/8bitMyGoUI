package icon_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type decl struct {
	file string
	expr ast.Expr
}

var upper = regexp.MustCompile(`([a-z])([A-Z])`)

func kebab(s string) string { return strings.ToLower(upper.ReplaceAllString(s, "$1-$2")) }

// decls maps each icon's kebab name to its source.
func decls(t *testing.T) map[string]decl {
	files, _ := filepath.Glob("*.go")
	out := map[string]decl{}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") || f == "icon.go" {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), f, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range file.Decls {
			if g, ok := d.(*ast.GenDecl); ok && g.Tok == token.VAR {
				for _, sp := range g.Specs {
					v := sp.(*ast.ValueSpec)
					for i, n := range v.Names {
						out[kebab(n.Name)] = decl{f, v.Values[i]}
					}
				}
			}
		}
	}
	return out
}

// build rebuilds an icon from its declared art.
func build(t *testing.T, ds map[string]decl, name string) *pixel.Sprite {
	call, ok := ds[name].expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("%s is not a sprite call", name)
	}
	sel := call.Fun.(*ast.SelectorExpr)
	switch sel.Sel.Name {
	case "MustSprite":
		art, err := strconv.Unquote(call.Args[0].(*ast.BasicLit).Value)
		if err != nil {
			t.Fatal(err)
		}
		return pixel.MustSprite(art)
	case "Mirror":
		return build(t, ds, kebab(sel.X.(*ast.Ident).Name)).Mirror()
	}
	t.Fatalf("%s built by %s", name, sel.Sel.Name)
	return nil
}

func render(sp *pixel.Sprite, ink pixel.Slot) []byte {
	return ui.Render(func(c *ui.Context) {
		pixel.Use(c, pixel.Settings{Palette: pixel.SummerDay, Dot: 1})
		pixel.Icon(c, sp, pixel.IconOptions{Slot: ink})
	}, 12, 12, 1).Pix
}

func TestEveryIconIsListedWithItsArt(t *testing.T) {
	ds := decls(t)
	listed := map[string]bool{}
	for _, n := range icon.All {
		if listed[n.Name] {
			t.Errorf("icon %q listed twice", n.Name)
		}
		listed[n.Name] = true
		if _, ok := ds[n.Name]; !ok {
			t.Fatalf("listed icon %q has no declaration", n.Name)
		}
		if string(render(n.Sprite, pixel.Magic)) != string(render(build(t, ds, n.Name), pixel.Magic)) {
			t.Errorf("icon %q is bound to another icon's art", n.Name)
		}
	}
	for name := range ds {
		if !listed[name] {
			t.Errorf("icon %q missing from All", name)
		}
	}
	if len(icon.All) < 100 {
		t.Fatalf("%d icons, want at least 100", len(icon.All))
	}
}

func TestOneBitIconsTakeTheInk(t *testing.T) {
	ds := decls(t)
	back, ink := pixel.SummerDay.C(pixel.Back), pixel.SummerDay.C(pixel.Magic)
	for _, n := range icon.All {
		if ds[n.Name].file == "items.go" {
			continue
		}
		pix, drawn := render(n.Sprite, pixel.Magic), 0
		for i := 0; i < len(pix); i += 4 {
			switch {
			case pix[i] == back.R && pix[i+1] == back.G && pix[i+2] == back.B:
			case pix[i] == ink.R && pix[i+1] == ink.G && pix[i+2] == ink.B:
				drawn++
			default:
				t.Fatalf("%s draws a color other than the ink", n.Name)
			}
		}
		if drawn == 0 {
			t.Fatalf("%s draws nothing in the ink", n.Name)
		}
	}
}

func TestIconsAreDistinctTwelveByTwelve(t *testing.T) {
	seen := map[string]string{}
	for _, n := range icon.All {
		if n.Sprite.W != 12 || n.Sprite.H != 12 || n.Sprite.Frames() != 1 {
			t.Fatalf("%s is %d×%d with %d frames", n.Name, n.Sprite.W, n.Sprite.H, n.Sprite.Frames())
		}
		pix, back, drawn := render(n.Sprite, pixel.Ink), pixel.SummerDay.C(pixel.Back), 0
		for i := 0; i < len(pix); i += 4 {
			if pix[i] != back.R || pix[i+1] != back.G || pix[i+2] != back.B {
				drawn++
			}
		}
		if drawn == 0 {
			t.Fatalf("%s draws nothing", n.Name)
		}
		key := string(pix)
		if other, ok := seen[key]; ok {
			t.Fatalf("%s and %s are the same picture", n.Name, other)
		}
		seen[key] = n.Name
	}
}

func TestItemsAreOutlined(t *testing.T) {
	for name, d := range decls(t) {
		if d.file != "items.go" {
			continue
		}
		art, _ := strconv.Unquote(d.expr.(*ast.CallExpr).Args[0].(*ast.BasicLit).Value)
		var rows []string
		for _, r := range strings.Split(strings.TrimSpace(art), "\n") {
			rows = append(rows, strings.TrimSpace(r))
		}
		for y, row := range rows {
			for x := range len(row) {
				if strings.IndexByte(".e#", row[x]) >= 0 {
					continue
				}
				for _, n := range [][2]int{{x + 1, y}, {x - 1, y}, {x, y + 1}, {x, y - 1}} {
					if n[0] < 0 || n[1] < 0 || n[0] >= len(row) || n[1] >= len(rows) || rows[n[1]][n[0]] == '.' {
						t.Fatalf("%s: dot %d,%d touches the outside without an Edge outline", name, x, y)
					}
				}
			}
		}
	}
}
