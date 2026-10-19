package code

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestDiffTakesOutBeforePuttingIn(t *testing.T) {
	got := diff([]string{"x", "y", "z"}, []string{"x", "q", "z", "w"})
	want := []change{{' ', "x", 1, 1}, {'-', "y", 2, 0}, {'+', "q", 0, 2}, {' ', "z", 3, 3}, {'+', "w", 0, 4}}
	if !slices.Equal(got, want) {
		t.Fatalf("diff = %v, want %v", got, want)
	}
	if got := diff(nil, []string{"a"}); !slices.Equal(got, []change{{'+', "a", 0, 1}}) {
		t.Fatalf("diff from nothing = %v", got)
	}
}

func TestStepsWalkEveryLine(t *testing.T) {
	const w = 7
	for a := range w * w {
		for b := range w * w {
			x0, y0, x1, y1 := a%w, a/w, b%w, b/w
			var path []int
			steps(x0, y0, x1, y1, func(i int) { path = append(path, i) }, w)
			if path[0] != a || path[len(path)-1] != b || len(path) != abs(x1-x0)+abs(y1-y0)+1 {
				t.Fatalf("%d,%d to %d,%d walked %v", x0, y0, x1, y1, path)
			}
			for i := 1; i < len(path); i++ {
				if d := abs(path[i]%w-path[i-1]%w) + abs(path[i]/w-path[i-1]/w); d != 1 {
					t.Fatalf("%d,%d to %d,%d jumps in %v", x0, y0, x1, y1, path)
				}
			}
		}
	}
}

func TestFloodFillsOnlyItsRun(t *testing.T) {
	// A wall of 1s splits the left run off.
	cells := []int{-1, 1, -1, -1, 1, -1, -1, 1, -1}
	if !flood(cells, 3, 3, 0, 5) || !slices.Equal(cells, []int{5, 1, -1, 5, 1, -1, 5, 1, -1}) {
		t.Fatalf("flood = %v", cells)
	}
	if flood(cells, 3, 3, 0, 5) {
		t.Fatal("filling a run with its own slot changed it")
	}
}

func TestPercentileNearestRank(t *testing.T) {
	var ds []time.Duration
	for i := 10; i >= 1; i-- {
		ds = append(ds, time.Duration(i))
	}
	if percentile(ds, 50) != 5 || percentile(ds, 95) != 10 || percentile(ds, 1) != 1 || percentile(nil, 50) != 0 {
		t.Fatal("percentiles are not nearest-rank")
	}
	if ds[0] != 10 {
		t.Fatal("percentile sorted its input")
	}
}

func TestJSONNodesKeepOrderAndTypes(t *testing.T) {
	dec := json.NewDecoder(strings.NewReader(`{"b": 1.50, "a": "中\n", "c": [true, null], "d": {}}`))
	dec.UseNumber()
	n := node(dec, "")
	got := []data.TreeNode{}
	for _, ch := range n.Children {
		got = append(got, data.TreeNode{Label: ch.Label, Value: ch.Value, Slot: ch.Slot})
	}
	want := []data.TreeNode{
		{Label: "b", Value: "1.50", Slot: pixel.Warn}, {Label: "a", Value: `"中\n"`, Slot: pixel.Good},
		{Label: "c", Value: "[2]", Slot: pixel.Faint}, {Label: "d", Value: "{0}", Slot: pixel.Faint},
	}
	if n.Value != "{4}" || !slices.EqualFunc(got, want, func(a, b data.TreeNode) bool { return a.Label == b.Label && a.Value == b.Value && a.Slot == b.Slot }) {
		t.Fatalf("nodes %+v", got)
	}
	arr := n.Children[2].Children
	if len(arr) != 2 || arr[0].Value != "true" || arr[0].Slot != pixel.Info || arr[1].Value != "null" || arr[1].Label != "1" {
		t.Fatalf("array %+v", arr)
	}
	if d := n.Children[3].Children; d == nil || len(d) != 0 {
		t.Fatal("an empty object must still be a branch")
	}
}

func TestVariableKinds(t *testing.T) {
	for _, tc := range []struct {
		v    Var
		want pixel.Slot
	}{
		{Var{Type: "string", Value: `"a"`}, pixel.Good}, {Var{Type: "bool", Value: "true"}, pixel.Info},
		{Var{Type: "int", Value: "1"}, pixel.Warn}, {Var{Type: "uint8", Value: "1"}, pixel.Warn},
		{Var{Type: "float64", Value: "1.5"}, pixel.Warn}, {Var{Type: "rune", Value: "65"}, pixel.Warn},
		{Var{Type: "*Enemy", Value: "nil"}, pixel.Faint}, {Var{Type: "Hero", Value: "{…}"}, pixel.Faint},
	} {
		if got := kind(tc.v); got != tc.want {
			t.Errorf("%+v is %v, want %v", tc.v, got, tc.want)
		}
	}
	ns := nodes([]Var{{Name: "hero", Type: "Hero", Value: "{…}", Fields: []Var{{Name: "HP", Type: "int", Value: "3"}}}})
	if ns[0].Label != "hero: Hero" || ns[0].Children[0].Label != "HP: int" || ns[0].Children[0].Value != "3" {
		t.Fatalf("nodes %+v", ns)
	}
	if glyph('A') != 'A' || glyph(0x7F) != '.' || glyph(0x1F) != '.' || glyph(' ') != ' ' {
		t.Fatal("glyph does not keep printable ASCII only")
	}
}
