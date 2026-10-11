package pixel_test

import (
	"testing"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestRectAndLines(t *testing.T) {
	d := canvas(t, 8, 6, func(g pixel.Canvas) {
		g.Rect(1, 1, 3, 2, pixel.Accent)
		g.HLine(0, 4, 5, pixel.Good)
		g.VLine(7, 0, 6, pixel.Bad)
		g.Dot(6, 5, pixel.Magic)
		g.Rect(5, 1, 0, 3, pixel.Ink)
	})
	if d.count(1, 1, 3, 2, pixel.Accent) != 6 || d.at(0, 1) != pixel.Back || d.at(4, 2) != pixel.Back {
		t.Fatal("Rect covers the wrong dots")
	}
	if d.count(0, 4, 5, 1, pixel.Good) != 5 || d.at(5, 4) != pixel.Back {
		t.Fatal("HLine covers the wrong dots")
	}
	if d.count(7, 0, 1, 6, pixel.Bad) != 6 || d.at(6, 5) != pixel.Magic {
		t.Fatal("VLine or Dot wrong")
	}
	if d.count(5, 1, 1, 3, pixel.Ink) != 0 {
		t.Fatal("zero-width Rect drew")
	}
}

func TestLineBresenham(t *testing.T) {
	cases := [][4]int{{0, 0, 9, 3}, {9, 0, 0, 3}, {2, 0, 4, 9}, {0, 9, 6, 0}, {3, 3, 3, 3}, {0, 5, 9, 5}}
	for _, c := range cases {
		d := canvas(t, 10, 10, func(g pixel.Canvas) { g.Line(c[0], c[1], c[2], c[3], pixel.Ink) })
		if d.at(c[0], c[1]) != pixel.Ink || d.at(c[2], c[3]) != pixel.Ink {
			t.Fatalf("line %v misses an end", c)
		}
		dx, dy := c[2]-c[0], c[3]-c[1]
		major := max(dx, -dx, dy, -dy)
		if got := d.count(0, 0, 10, 10, pixel.Ink); got != major+1 {
			t.Fatalf("line %v has %d dots, want %d", c, got, major+1)
		}
		xMajor := max(dx, -dx) >= max(dy, -dy)
		for i := range 10 {
			var n int
			if xMajor {
				n = d.count(i, 0, 1, 10, pixel.Ink)
			} else {
				n = d.count(0, i, 10, 1, pixel.Ink)
			}
			inside := i >= min(c[0], c[2]) && i <= max(c[0], c[2])
			if !xMajor {
				inside = i >= min(c[1], c[3]) && i <= max(c[1], c[3])
			}
			if inside != (n == 1) || n > 1 {
				t.Fatalf("line %v: %d dots across major step %d", c, n, i)
			}
		}
	}
}

func TestLineTieBreak(t *testing.T) {
	d := canvas(t, 3, 2, func(g pixel.Canvas) { g.Line(0, 0, 2, 1, pixel.Ink) })
	for _, p := range [][2]int{{0, 0}, {1, 1}, {2, 1}} {
		if d.at(p[0], p[1]) != pixel.Ink {
			t.Fatalf("Line(0,0,2,1) misses %v", p)
		}
	}
	if d.at(1, 0) != pixel.Back {
		t.Fatal("Line(0,0,2,1) took the wrong tie")
	}
}

func TestCircleAndDisc(t *testing.T) {
	const cx, cy, r = 7, 7, 5
	disc := canvas(t, 15, 15, func(g pixel.Canvas) { g.Disc(cx, cy, r, pixel.Info) })
	ring := canvas(t, 15, 15, func(g pixel.Canvas) { g.Circle(cx, cy, r, pixel.Ink) })
	for _, p := range [][2]int{{cx + r, cy}, {cx - r, cy}, {cx, cy + r}, {cx, cy - r}} {
		if disc.at(p[0], p[1]) != pixel.Info || ring.at(p[0], p[1]) != pixel.Ink {
			t.Fatalf("radius point %v missing", p)
		}
	}
	if disc.at(cx, cy) != pixel.Info || ring.at(cx, cy) != pixel.Back || disc.at(cx+r, cy+r) != pixel.Back {
		t.Fatal("disc must be solid, circle hollow, corners empty")
	}
	for y := range 15 {
		for x := range 15 {
			if disc.at(x, y) != disc.at(2*cx-x, y) || disc.at(x, y) != disc.at(x, 2*cy-y) || disc.at(x, y) != disc.at(y, x) {
				t.Fatalf("disc not symmetric at %d,%d", x, y)
			}
			if ring.at(x, y) == pixel.Ink && disc.at(x, y) != pixel.Info {
				t.Fatalf("circle dot %d,%d outside the disc", x, y)
			}
		}
	}
}

func TestCircleExact(t *testing.T) {
	ring := canvas(t, 7, 7, func(g pixel.Canvas) { g.Circle(3, 3, 3, pixel.Ink) })
	disc := canvas(t, 7, 7, func(g pixel.Canvas) { g.Disc(3, 3, 3, pixel.Ink) })
	wantRing := []string{"..###..", ".#...#.", "#.....#", "#.....#", "#.....#", ".#...#.", "..###.."}
	wantDisc := []string{"..###..", ".#####.", "#######", "#######", "#######", ".#####.", "..###.."}
	for y := range 7 {
		for x := range 7 {
			if (ring.at(x, y) == pixel.Ink) != (wantRing[y][x] == '#') || (disc.at(x, y) == pixel.Ink) != (wantDisc[y][x] == '#') {
				t.Fatalf("radius 3 dot %d,%d differs from the midpoint shape", x, y)
			}
		}
	}
}

func TestCircleRadiusSix(t *testing.T) {
	octant := [][2]int{{6, 0}, {6, 1}, {6, 2}, {5, 3}, {4, 4}}
	want := map[[2]int]bool{}
	for _, p := range octant {
		for _, q := range [][2]int{{p[0], p[1]}, {p[1], p[0]}} {
			for _, s := range [][2]int{{1, 1}, {-1, 1}, {1, -1}, {-1, -1}} {
				want[[2]int{6 + s[0]*q[0], 6 + s[1]*q[1]}] = true
			}
		}
	}
	d := canvas(t, 13, 13, func(g pixel.Canvas) { g.Circle(6, 6, 6, pixel.Ink) })
	for y := range 13 {
		for x := range 13 {
			if (d.at(x, y) == pixel.Ink) != want[[2]int{x, y}] {
				t.Fatalf("radius 6 dot %d,%d differs from the midpoint shape", x, y)
			}
		}
	}
}

func TestBayerMatrix(t *testing.T) {
	want := [4][4]int{{0, 8, 2, 10}, {12, 4, 14, 6}, {3, 11, 1, 9}, {15, 7, 13, 5}}
	first := [4][4]int{}
	for level := 16; level >= 1; level-- {
		d := canvas(t, 4, 4, func(g pixel.Canvas) { g.Veil(0, 0, 4, 4, pixel.Ink, level) })
		for y := range 4 {
			for x := range 4 {
				if d.at(x, y) == pixel.Ink {
					first[y][x] = level - 1
				}
			}
		}
	}
	if first != want {
		t.Fatalf("Bayer thresholds %v, want %v", first, want)
	}
}

func TestFrameStyles(t *testing.T) {
	const w, h = 8, 7
	want := map[pixel.FrameStyle][][3]int{
		pixel.Flat:    {{0, 0, int(pixel.Edge)}, {1, 1, int(pixel.Face)}, {w - 1, h - 1, int(pixel.Edge)}},
		pixel.Raised:  {{0, 0, int(pixel.Edge)}, {1, 1, int(pixel.Lit)}, {w - 2, h - 2, int(pixel.Shade)}, {2, 2, int(pixel.Face)}},
		pixel.Sunken:  {{1, 1, int(pixel.Shade)}, {w - 2, h - 2, int(pixel.Lit)}, {2, 2, int(pixel.Face)}},
		pixel.Double:  {{0, 0, int(pixel.Edge)}, {1, 1, int(pixel.Face)}, {2, 2, int(pixel.Edge)}, {3, 3, int(pixel.Face)}},
		pixel.Window:  {{0, 0, int(pixel.Back)}, {1, 0, int(pixel.Edge)}, {1, 1, int(pixel.Lit)}, {2, 2, int(pixel.Shade)}, {3, 3, int(pixel.Face)}},
		pixel.Notched: {{0, 0, int(pixel.Back)}, {1, 0, int(pixel.Edge)}, {0, 1, int(pixel.Edge)}, {1, 1, int(pixel.Face)}},
	}
	for style, checks := range want {
		d := canvas(t, w, h, func(g pixel.Canvas) { g.Frame(0, 0, w, h, style, pixel.Face) })
		for _, ck := range checks {
			if got := d.at(ck[0], ck[1]); got != pixel.Slot(ck[2]) {
				t.Errorf("style %d dot %d,%d = slot %d, want %d", style, ck[0], ck[1], got, ck[2])
			}
		}
	}
	mustPanic(t, "2×2 frame", func() {
		canvas(t, 4, 4, func(g pixel.Canvas) { g.Frame(0, 0, 2, 2, pixel.Flat, pixel.Face) })
	})
}

func TestDitherLevels(t *testing.T) {
	for level := 0; level <= 16; level++ {
		d := canvas(t, 8, 8, func(g pixel.Canvas) { g.Dither(0, 0, 8, 8, pixel.Face, pixel.Ink, level) })
		for _, b := range [][2]int{{0, 0}, {4, 0}, {0, 4}, {4, 4}} {
			if n := d.count(b[0], b[1], 4, 4, pixel.Ink); n != level {
				t.Fatalf("level %d: 4×4 block at %v has %d ink dots", level, b, n)
			}
		}
	}
	d := canvas(t, 6, 6, func(g pixel.Canvas) { g.Checker(0, 0, 6, 6, pixel.Ink) })
	for y := range 6 {
		for x := range 6 {
			if want := (x+y)%2 == 0; (d.at(x, y) == pixel.Ink) != want {
				t.Fatalf("checker dot %d,%d", x, y)
			}
		}
	}
}

func TestDitherAnchorsToWindow(t *testing.T) {
	d := canvas(t, 8, 4, func(g pixel.Canvas) {
		g.Rect(0, 0, 8, 4, pixel.Accent)
		g.Veil(1, 0, 3, 4, pixel.Ink, 8)
		g.Veil(4, 1, 4, 3, pixel.Ink, 8)
	})
	for y := range 4 {
		for x := 1; x < 8; x++ {
			if x >= 4 && y == 0 {
				continue
			}
			want := pixel.Accent
			if (x+y)%2 == 0 {
				want = pixel.Ink
			}
			if got := d.at(x, y); got != want {
				t.Fatalf("veil dot %d,%d = %d, want %d", x, y, got, want)
			}
		}
	}
	mustPanic(t, "level 17", func() {
		canvas(t, 2, 2, func(g pixel.Canvas) { g.Veil(0, 0, 1, 1, pixel.Ink, 17) })
	})
}

func TestTextAndMicroDraw(t *testing.T) {
	a := []string{"...#...", "...#...", "..#.#..", "..#.#..", ".#...#.", ".#####.", ".#...#.", "#.....#", "#.....#"}
	for _, scale := range []int{1, 2} {
		d := canvas(t, 8*scale, 16*scale, func(g pixel.Canvas) { g.Text(0, 0, "A", pixel.Ink, scale) })
		for row, line := range a {
			for col, ch := range line {
				if got := d.at(col*scale, (4+row)*scale) == pixel.Ink; got != (ch == '#') {
					t.Fatalf("×%d 'A' dot %d,%d", scale, col, row)
				}
			}
		}
		if d.count(0, 0, 8*scale, 4*scale, pixel.Ink) != 0 {
			t.Fatalf("×%d 'A' drew above its top", scale)
		}
	}
	one := []string{".#.", "##.", ".#.", ".#.", "###"}
	d := canvas(t, 8, 5, func(g pixel.Canvas) { g.Micro(0, 0, "11", pixel.Good, 1) })
	for row, line := range one {
		for col, ch := range line {
			if (d.at(col, row) == pixel.Good) != (ch == '#') || (d.at(col+4, row) == pixel.Good) != (ch == '#') {
				t.Fatalf("micro '1' dot %d,%d", col, row)
			}
		}
	}
}
