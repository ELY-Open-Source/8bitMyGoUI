package gallery

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegistry(t *testing.T) {
	slugs, ids := map[string]bool{}, map[int]bool{}
	known := map[string]bool{}
	for _, c := range Categories {
		known[c.Key] = true
	}
	for _, e := range Entries() {
		if e.Name == "" || e.Zh == "" || e.Summary == "" || e.Demo == nil {
			t.Fatalf("entry %+v is incomplete", e)
		}
		if !known[e.Category] {
			t.Fatalf("%s: unknown category %q", e.Name, e.Category)
		}
		if (e.ID == 0) != (e.Category == "foundations") || e.ID < 0 || e.ID > 500 {
			t.Fatalf("%s: id %d in %s", e.Name, e.ID, e.Category)
		}
		if slugs[e.Slug()] || e.ID > 0 && ids[e.ID] {
			t.Fatalf("%s: duplicate slug or id", e.Slug())
		}
		slugs[e.Slug()], ids[e.ID] = true, true
		if !strings.HasPrefix(e.Code(), "func "+funcName(e.Demo)+"(") {
			t.Fatalf("%s: code snippet %q", e.Slug(), e.Code()[:min(40, len(e.Code()))])
		}
	}
}

// Every exported component function appears in a demo.
func TestEveryComponentHasADemo(t *testing.T) {
	var demos strings.Builder
	for _, e := range Entries() {
		demos.WriteString(e.Code())
	}
	for _, cat := range Categories[1:] {
		files, _ := filepath.Glob(filepath.Join("..", cat.Key, "*.go"))
		for _, f := range files {
			if strings.HasSuffix(f, "_test.go") {
				continue
			}
			src, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			file, err := parser.ParseFile(token.NewFileSet(), f, src, 0)
			if err != nil {
				t.Fatal(err)
			}
			for _, d := range file.Decls {
				fn, ok := d.(*ast.FuncDecl)
				if !ok || fn.Recv != nil || !fn.Name.IsExported() || fn.Type.Params.NumFields() == 0 {
					continue
				}
				if p, ok := fn.Type.Params.List[0].Type.(*ast.StarExpr); ok {
					if sel, ok := p.X.(*ast.SelectorExpr); ok && sel.Sel.Name == "Context" {
						if call := cat.Key + "." + fn.Name.Name + "("; !strings.Contains(demos.String(), call) {
							t.Errorf("%s has no demo", call)
						}
					}
				}
			}
		}
	}
}
