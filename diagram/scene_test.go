package diagram

import (
	"strings"
	"testing"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// grid is a scene's dots: slot plus one, or 0.
func grid(t *testing.T, sc *scene) *plate {
	t.Helper()
	p := blank(sc.w, sc.h)
	for _, c := range sc.cells() {
		if c.x < 0 || c.y < 0 || c.w < 1 || c.h < 1 || c.x+c.w > sc.w || c.y+c.h > sc.h || c.slot >= pixel.Slots {
			t.Fatalf("cell %+v leaves a scene of %d by %d dots", c, sc.w, sc.h)
		}
		for y := c.y; y < c.y+c.h; y++ {
			for x := c.x; x < c.x+c.w; x++ {
				p.dots[y*p.w+x] = uint8(c.slot) + 1
			}
		}
	}
	return p
}

func draw(t *testing.T, source string) *plate {
	t.Helper()
	sc, err := render(source, tokens(pixel.SpringDay), &pixel.SpringDay.Colors)
	if err != nil {
		t.Fatalf("%q: %v", source, err)
	}
	return grid(t, sc)
}

// sketch draws an SVG body on 100 by 100 dots.
func sketch(t *testing.T, body string) *plate {
	t.Helper()
	sc, err := parse(`<svg width="100" height="100">`+body+`</svg>`, tokens(pixel.SpringDay), &pixel.SpringDay.Colors)
	if err != nil {
		t.Fatalf("%s: %v", body, err)
	}
	return grid(t, sc)
}

func (p *plate) slot(x, y int) pixel.Slot { return pixel.Slot(p.at(x, y)) - 1 }

// box bounds slot s and counts its dots.
func (p *plate) box(s pixel.Slot) (x0, y0, x1, y1, n int) {
	x0, y0, x1, y1 = p.w, p.h, -1, -1
	for i, d := range p.dots {
		if d == uint8(s)+1 {
			x0, y0, x1, y1, n = min(x0, i%p.w), min(y0, i/p.w), max(x1, i%p.w), max(y1, i/p.w), n+1
		}
	}
	return
}

func TestCorpusDraws(t *testing.T) {
	docs := corpus(t)
	if len(docs) != 33 {
		t.Fatalf("%d diagram types, want 33", len(docs))
	}
	for _, p := range pixel.Palettes[:2] {
		for kind, sources := range docs {
			for i, source := range sources {
				sc, err := render(source, tokens(p), &p.Colors)
				if err != nil {
					t.Errorf("%s %d: %v", kind, i, err)
					continue
				}
				grid(t, sc)
				// Gantt 6 holds settings and no task.
				if blank := kind == "gantt" && i == 6; (len(sc.cells()) == 0) != blank {
					t.Errorf("%s %d: %d cells on %d by %d dots", kind, i, len(sc.cells()), sc.w, sc.h)
				}
			}
		}
	}
}

// only keeps the dots of a rectangle.
func (p *plate) only(x0, y0, x1, y1 int) *plate {
	in := blank(p.w, p.h)
	for y := y0; y <= y1; y++ {
		copy(in.dots[y*p.w+x0:y*p.w+x1+1], p.dots[y*p.w+x0:y*p.w+x1+1])
	}
	return in
}

func TestFlowchartLandsOnDots(t *testing.T) {
	// Pele: Kiln's rect 8,8–65,50; edge 65–106 on 29; head 113.
	p := draw(t, "flowchart LR\n  A[Kiln] --> B[Glaze]")
	for x := 8; x <= 65; x++ {
		if p.slot(x, 8) != pixel.Edge || p.slot(x, 50) != pixel.Edge {
			t.Fatalf("column %d of Kiln's frame is not Edge", x)
		}
	}
	for y := 8; y <= 50; y++ {
		if p.slot(8, y) != pixel.Edge || p.slot(65, y) != pixel.Edge {
			t.Fatalf("row %d of Kiln's frame is not Edge", y)
		}
	}
	if p.slot(12, 12) != pixel.Face || p.slot(61, 46) != pixel.Face || p.at(4, 4) != 0 || p.at(67, 8) != 0 {
		t.Error("want Face inside Kiln's box and nothing outside it")
	}
	for x := 67; x <= 103; x++ {
		if p.slot(x, 29) != pixel.Ink || p.at(x, 27) != 0 {
			t.Fatalf("column %d: want a one-dot edge on row 29", x)
		}
	}
	if p.slot(106, 27) != pixel.Ink || p.slot(106, 31) != pixel.Ink || p.slot(111, 29) != pixel.Ink || p.at(111, 27) != 0 {
		t.Error("want an arrowhead from column 105 narrowing to its tip at 113")
	}
	x0, y0, x1, y1, n := p.only(9, 9, 64, 49).box(pixel.Ink)
	if w := pixel.Measure("Kiln"); n == 0 || x1-x0 >= w || x0+x1 < 71 || x0+x1 > 75 || y0 < 22 || y1 > 36 {
		t.Errorf("label dots span %d..%d by %d..%d (%d dots), want Kiln centered on column 36.5", x0, x1, y0, y1, n)
	}
}

func TestInkTakesTheNearestSlot(t *testing.T) {
	p := draw(t, "flowchart LR\n  A[Kiln] --> B[Glaze]\n  style A fill:#c83b3b,stroke:#4f9a40,color:#ffffff\n  style B fill:#3f7fcb,fill-opacity:50%")
	if p.slot(12, 12) != pixel.Bad || p.slot(8, 20) != pixel.Good {
		t.Errorf("Kiln is %v framed %v, want the slots nearest its colors: Bad and Good", p.slot(12, 12), p.slot(8, 20))
	}
	if _, _, _, _, n := p.box(pixel.AccentInk); n == 0 {
		t.Error("no dot of the white label")
	}
	info, face := 0, 0
	// Glaze's rect runs 113,8 to 178,50.
	for y := 12; y < 46; y++ {
		for x := 117; x < 175; x++ {
			if p.at(x, y) == 0 {
				face++
			} else if p.slot(x, y) == pixel.Info {
				info++
			}
		}
	}
	if info == 0 || face == 0 || info*10 < face*8 || info*10 > face*12 {
		t.Errorf("Glaze holds %d Info dots and %d bare ones, want about half and half", info, face)
	}
}

func TestTextAnchorsScaleAndTurns(t *testing.T) {
	w := pixel.Measure("Kiln")
	for anchor, left := range map[string]int{"start": 50, "middle": 50 - w/2, "end": 50 - w} {
		x0, _, x1, y1, _ := sketch(t, `<text x="50" y="40" text-anchor="`+anchor+`">Kiln</text>`).box(pixel.Ink)
		if x0 < left || x0 > left+2 || x1 >= left+w || y1 != 39 {
			t.Errorf("%s: dots span columns %d..%d down to row %d, want from %d on the baseline at 40", anchor, x0, x1, y1, left)
		}
	}
	_, y0, _, y1, _ := sketch(t, `<text x="10" y="40">K</text>`).box(pixel.Ink)
	_, by0, _, by1, _ := sketch(t, `<text x="10" y="40" font-size="24">K</text>`).box(pixel.Ink)
	if by1-by0+1 != 2*(y1-y0+1) {
		t.Errorf("a 24 px K is %d dots tall against %d, want twice", by1-by0+1, y1-y0+1)
	}
	// A quarter turn stands the line up; others stay level.
	x0, y0, x1, y1, _ := sketch(t, `<g transform="translate(50,60) rotate(-90)"><text x="0" y="0">Kiln</text></g>`).box(pixel.Ink)
	if x1-x0 >= y1-y0 || y1 > 60 || y0 < 60-w {
		t.Errorf("turned text spans %d..%d by %d..%d, want it to run up from 50,60", x0, x1, y0, y1)
	}
	x0, y0, x1, y1, _ = sketch(t, `<g transform="translate(50,60) rotate(-45)"><text x="0" y="0">Kiln</text></g>`).box(pixel.Ink)
	if x1-x0 <= y1-y0 || x0 < 50 || y1 != 59 {
		t.Errorf("text turned 45 degrees spans %d..%d by %d..%d, want it upright at 50,60", x0, x1, y0, y1)
	}
	under := sketch(t, `<text x="10" y="40" style="text-decoration:underline">Kiln</text>`)
	for x := 10; x < 10+w; x++ {
		if under.slot(x, 42) != pixel.Ink {
			t.Fatalf("no underline at column %d", x)
		}
	}
	halo := sketch(t, `<text x="10" y="40" fill="none" stroke="var(--_bg)" stroke-width="3">K</text><text x="10" y="40">K</text>`)
	if _, _, _, _, n := halo.box(pixel.Back); n == 0 {
		t.Error("an outlined text left no halo")
	}
}

func TestUnknownMarkupFails(t *testing.T) {
	for in, want := range map[string]string{
		`<svg width="10" height="10"><line/></svg>`:                             "unsupported element <line>",
		`<svg width="10" height="10"><g clip-path="url(#a)"/></svg>`:            `unsupported attribute "clip-path"`,
		`<svg width="10" height="10"><g transform="scale(2)"/></svg>`:           `unsupported transform "scale(2)"`,
		`<svg width="10" height="10"><rect fill="var(--_zz)"/></svg>`:           `unknown token "--_zz"`,
		`<svg width="10" height="10"><rect fill="var(--pele-new,#123)"/></svg>`: `unknown token "--pele-new"`,
		`<svg width="10" height="10"><rect width="wide"/></svg>`:                `width="wide" is not a number`,
		`<svg width="10" height="10"><svg width="4" height="4"/></svg>`:         "unsupported viewBox",
		`<svg width="10" height="10"><path d="M0 0 L"/></svg>`:                  "<path>",
		`<svg width="10" height="10"><g>`:                                       "diagram:",
		``:                                                                      "no size",
		`<svg width="0" height="10"></svg>`:                                     "dots a side",
		`<svg width="99999" height="10"></svg>`:                                 "dots a side",
	} {
		if _, err := parse(in, tokens(pixel.SpringDay), &pixel.SpringDay.Colors); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v, want %q", in, err, want)
		}
	}
}

func TestCornersAreSquareUnlessTheShapeRounds(t *testing.T) {
	// Glaze is a stadium: its rect brings its own rx.
	p := draw(t, "flowchart LR\n  A[Kiln] --> B([Glaze])")
	if p.slot(8, 8) != pixel.Edge || p.slot(65, 50) != pixel.Edge {
		t.Error("Kiln's frame lacks its corner dots")
	}
	left := func(y int) int {
		for x := 100; x < p.w; x++ {
			if p.at(x, y) != 0 {
				return x
			}
		}
		return p.w
	}
	// Row 36 clears the arrowhead on the middle row.
	if side, top := left(36), left(9); side < 112 || side > 115 || top < side+8 {
		t.Errorf("Glaze starts at column %d by its middle and %d near its top, want about 113 and a round corner", side, top)
	}
}
