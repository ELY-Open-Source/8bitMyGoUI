package diagram

import (
	"context"
	"math"
	"regexp"
	"strings"
	"testing"

	pele "github.com/Broderick-Westrope/pele-go"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/egoist/mygo/ui"
)

func TestFlatten(t *testing.T) {
	o, err := flatten("M1 2h3v4H1zm10 0l2,2 -2 2", identity)
	if err != nil {
		t.Fatal(err)
	}
	want := outline{{{1, 2}, {4, 2}, {4, 6}, {1, 6}, {1, 2}}, {{11, 2}, {13, 4}, {11, 6}}}
	if len(o) != 2 || len(o[0]) != 5 || len(o[1]) != 3 {
		t.Fatalf("outline %v, want %v", o, want)
	}
	for i, line := range want {
		for j, p := range line {
			if o[i][j] != p {
				t.Errorf("point %d.%d = %v, want %v", i, j, o[i][j], p)
			}
		}
	}
	turned, err := flatten("M1 0L2 0", rotation(90).then(translation(10, 20)))
	if err != nil || math.Abs(turned[0][1].x-10) > 1e-9 || math.Abs(turned[0][1].y-22) > 1e-9 {
		t.Errorf("a quarter turn then a move puts (2,0) at %v, want (10,22): %v", turned, err)
	}
	for _, d := range []string{"M0 0A5 5 0 1 1 10 0", "M0 0a5,5 0 11 10,0", "M0 0a5 5 0 1110 0"} {
		arc, err := flatten(d, identity)
		if err != nil {
			t.Fatalf("%s: %v", d, err)
		}
		top := 0.0
		for _, p := range arc[0] {
			if r := math.Hypot(p.x-5, p.y); math.Abs(r-5) > 1e-6 {
				t.Fatalf("%s: point %v is %v from the center, want 5", d, p, r)
			}
			top = min(top, p.y)
		}
		if end := arc[0][len(arc[0])-1]; math.Abs(top+5) > 0.05 || math.Abs(end.x-10) > 1e-9 || math.Abs(end.y) > 1e-9 {
			t.Errorf("%s: arc tops at %v and ends at %v, want -5 and (10,0)", d, top, end)
		}
	}
	for d, want := range map[string]string{"M0 0 X": `unknown command 'X'`, "5 5": "unexpected", "M0": "no number", "M0 0A5 5 0 2 1 1 1": "no arc flag", "M0 0z5": "unexpected"} {
		if _, err := flatten(d, identity); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%q: %v, want %q", d, err, want)
		}
	}
}

// TestFillAgreesWithMyGo holds the scanline fill to MyGo's own path test.
func TestFillAgreesWithMyGo(t *testing.T) {
	data := regexp.MustCompile(` d="([^"]*)"`)
	seen := map[string]bool{}
	paths, inside, wrong := 0, 0, 0
	for kind, sources := range corpus(t) {
		r, err := pele.Render(context.Background(), sources[0], pele.Options{})
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range data.FindAllStringSubmatch(r.SVG, -1) {
			ref, err := ui.ParsePath(m[1])
			if err != nil {
				t.Fatalf("%s: MyGo: %q: %v", kind, m[1], err)
			}
			b := ref.Bounds()
			if seen[m[1]] || !strings.ContainsAny(m[1], "Zz") || b.W > 300 || b.H > 300 {
				continue
			}
			seen[m[1]] = true
			o, err := flatten(m[1], translation(float64(8-b.X), float64(8-b.Y)))
			if err != nil {
				t.Fatalf("%s: %q: %v", kind, m[1], err)
			}
			p := plate{w: int(b.W) + 16, h: int(b.H) + 16}
			p.dots = make([]uint8, p.w*p.h)
			p.fill(o, ink{pixel.Ink, 16})
			bad := 0
			for y := range p.h {
				for x := range p.w {
					in := ref.Contains(float32(x)+0.5-8+b.X, float32(y)+0.5-8+b.Y)
					if in {
						inside++
					}
					if in != (p.dots[y*p.w+x] != 0) {
						bad++
					}
				}
			}
			paths++
			// Curves are sampled: a dot on the edge may differ.
			if wrong += bad; bad > 4+int(b.W+b.H)/4 {
				t.Errorf("%s: %d dots differ for %q", kind, bad, m[1])
			}
		}
	}
	if paths < 80 || inside < 100000 || wrong*200 > inside {
		t.Errorf("%d paths, %d dots inside, %d differ: want a wide sample within half a percent", paths, inside, wrong)
	}
}
