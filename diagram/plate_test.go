package diagram

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func blank(w, h int) *plate { return &plate{w: w, h: h, dots: make([]uint8, w*h)} }

func (p *plate) at(x, y int) uint8 { return p.dots[y*p.w+x] }

func TestStrokeLeavesNoGap(t *testing.T) {
	// Sides of 42 and 98 once lost a dot.
	for _, size := range [][2]float64{{98, 42}, {69, 42}, {63.5, 45.25}, {7, 3}} {
		p := blank(120, 60)
		o, err := flatten("M8 8h"+num(float32(size[0]))+"v"+num(float32(size[1]))+"h"+num(float32(-size[0]))+"Z", identity)
		if err != nil {
			t.Fatal(err)
		}
		p.stroke(o, 1, nil, 0, ink{pixel.Edge, 16})
		x1, y1 := dot(8+size[0]), dot(8+size[1])
		for x := 8; x <= x1; x++ {
			if p.at(x, 8) == 0 || p.at(x, y1) == 0 {
				t.Errorf("%v: no dot at column %d of the top or bottom side", size, x)
			}
		}
		for y := 8; y <= y1; y++ {
			if p.at(8, y) == 0 || p.at(x1, y) == 0 {
				t.Errorf("%v: no dot at row %d of the left or right side", size, y)
			}
		}
		if p.at(7, 8) != 0 || p.at(x1+1, y1) != 0 || p.at(9, 9) != 0 {
			t.Errorf("%v: dots outside the one-dot frame", size)
		}
	}
}

func TestStrokeWidthDashAndSlope(t *testing.T) {
	p := blank(40, 20)
	o, _ := flatten("M2 10H32", identity)
	p.stroke(o, 3, []float64{4, 2}, 0, ink{pixel.Ink, 16})
	row := ""
	for x := range 36 {
		if p.at(x, 10) != 0 {
			row += "#"
		} else {
			row += "."
		}
	}
	if want := "..####..####..####..####..####..#..."; row != want {
		t.Errorf("dashes %s, want %s: four on, two off, however wide the brush", row, want)
	}
	if p.at(3, 9) == 0 || p.at(3, 11) == 0 || p.at(3, 8) != 0 || p.at(3, 12) != 0 {
		t.Error("want a brush three dots tall")
	}
	q := blank(12, 12)
	o, _ = flatten("M1 1L10 7", identity)
	q.stroke(o, 1, nil, 0, ink{pixel.Ink, 16})
	for x := 1; x <= 10; x++ {
		n := 0
		for y := range 12 {
			if q.at(x, y) != 0 {
				n++
			}
		}
		if n != 1 {
			t.Errorf("column %d of a shallow line holds %d dots, want 1", x, n)
		}
	}
}

func TestFillWindsAndDithers(t *testing.T) {
	p := blank(20, 20)
	o, _ := flatten("M2 2h12v12h-12zM5 5v6h6v-6z", identity)
	p.fill(o, ink{pixel.Face, 16})
	if p.at(3, 3) == 0 || p.at(13, 13) == 0 || p.at(7, 7) != 0 || p.at(14, 14) != 0 || p.at(1, 1) != 0 {
		t.Error("want a square of 12 dots with a hole wound the other way")
	}
	half := blank(8, 8)
	o, _ = flatten("M0 0h8v8h-8z", identity)
	half.fill(o, ink{pixel.Face, 8})
	on := 0
	for _, d := range half.dots {
		if d != 0 {
			on++
		}
	}
	if on != 32 {
		t.Errorf("level 8 inks %d of 64 dots, want half", on)
	}
}

func TestCellsCoverEachDotOnce(t *testing.T) {
	p := blank(16, 12)
	o, _ := flatten("M2 2h10v6h-10z", identity)
	p.fill(o, ink{pixel.Face, 16})
	p.stroke(o, 1, nil, 0, ink{pixel.Edge, 16})
	p.fill(o, ink{pixel.Accent, 4})
	got := blank(16, 12)
	for _, c := range p.cells() {
		for y := c.y; y < c.y+c.h; y++ {
			for x := c.x; x < c.x+c.w; x++ {
				if got.at(x, y) != 0 {
					t.Fatalf("dot %d,%d sits in two cells", x, y)
				}
				got.dots[y*got.w+x] = uint8(c.slot) + 1
			}
		}
	}
	for i := range p.dots {
		if got.dots[i] != p.dots[i] {
			t.Fatalf("dot %d,%d: cells give %d, the plate holds %d", i%16, i/16, got.dots[i], p.dots[i])
		}
	}
	if cells := blank(8, 8); len(cells.cells()) != 0 {
		t.Error("an empty plate has cells")
	}
	solid := blank(8, 8)
	solid.fill(o, ink{pixel.Face, 16})
	if n := len(solid.cells()); n != 1 {
		t.Errorf("a plain rectangle makes %d cells, want 1", n)
	}
}

// stamp is the plain brush: a square at every dot.
func stamp(p *plate, o outline, width int, k ink) {
	for _, line := range o {
		for i := 1; i < len(line); i++ {
			x, y, x1, y1 := dot(line[i-1].x), dot(line[i-1].y), dot(line[i].x), dot(line[i].y)
			dx, dy := abs(x1-x), -abs(y1-y)
			for lean := dx + dy; ; {
				p.block(x-width/2, y-width/2, width, width, k)
				if x == x1 && y == y1 {
					break
				}
				twice := 2 * lean
				if twice >= dy {
					lean, x = lean+dy, x+sign(x1-x)
				}
				if twice <= dx {
					lean, y = lean+dx, y+sign(y1-y)
				}
			}
		}
	}
}

func TestBrushEdgesInkWhatStampingInks(t *testing.T) {
	o, err := flatten("M5 30L40 4L44 38L9 9L-6 20L20 50M30 30L30 30M2 2C30 0 0 30 38 38", identity)
	if err != nil {
		t.Fatal(err)
	}
	for _, width := range []int{1, 2, 3, 7, 40, 2 * limit} {
		for _, level := range []int{16, 5} {
			got, want := blank(48, 44), blank(48, 44)
			got.stroke(o, width, nil, 0, ink{pixel.Edge, level})
			stamp(want, o, width, ink{pixel.Edge, level})
			if string(got.dots) != string(want.dots) {
				t.Errorf("width %d at level %d: the edge brush and the stamp differ", width, level)
			}
		}
	}
}

func TestAbsurdLengthsFinish(t *testing.T) {
	done := make(chan *plate, 1)
	go func() {
		wide := draw(t, "flowchart LR\n  A[Kiln] --> B[Glaze]\n  style A stroke:#4f9a40,stroke-width:1e9px,stroke-dasharray:1e9 1e9\n"+
			"  style B stroke-width:1e9em,letter-spacing:1e9px,word-spacing:1e9px,text-decoration:underline")
		done <- wide
	}()
	select {
	case p := <-done:
		// Glaze's frame floods the plate; only its label shows through.
		_, _, _, _, edge := p.box(pixel.Edge)
		_, _, _, _, label := p.box(pixel.Ink)
		if label == 0 || edge+label != p.w*p.h {
			t.Errorf("%d Edge and %d Ink dots of %d, want a frame wider than the plate under its label", edge, label, p.w*p.h)
		}
	case <-time.After(20 * time.Second):
		t.Fatal("a style's length kept the walker busy")
	}
}
