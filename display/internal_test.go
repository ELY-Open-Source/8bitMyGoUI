package display

import (
	"image"
	"image/color"
	"testing"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func fill(w, h int, c color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			img.Set(x, y, c)
		}
	}
	return img
}

// dots are a sprite's slots, -1 where clear.
func dots(sp *pixel.Sprite) []int {
	var out []int
	for y := range sp.H {
		for x := range sp.W {
			s, ok := sp.At(0, x, y, pixel.Ink)
			if !ok {
				out = append(out, -1)
				continue
			}
			out = append(out, int(s))
		}
	}
	return out
}

func TestReduce(t *testing.T) {
	pal := pixel.SpringDay
	a := pal.Colors[pixel.Accent]
	for _, s := range dots(Reduce(fill(8, 8, color.RGBA{a.R, a.G, a.B, 255}), 4, 4, pal)) {
		if s != int(pixel.Accent) {
			t.Fatalf("a flat Accent image gave slot %d", s)
		}
	}
	for _, s := range dots(Reduce(fill(4, 4, color.RGBA{}), 4, 4, pal)) {
		if s != -1 {
			t.Fatal("a clear image drew")
		}
	}
	i, b := pal.Colors[pixel.Info], pal.Colors[pixel.Back]
	mid := color.RGBA{uint8((int(i.R) + int(b.R)) / 2), uint8((int(i.G) + int(b.G)) / 2), uint8((int(i.B) + int(b.B)) / 2), 255}
	seen := map[int]int{}
	for _, s := range dots(Reduce(fill(16, 16, mid), 16, 16, pal)) {
		seen[s]++
	}
	if len(seen) < 2 {
		t.Fatalf("a midtone did not dither: %v", seen)
	}
}

func TestSpanCoversSource(t *testing.T) {
	for _, tc := range [][2]int{{96, 48}, {10, 30}, {7, 7}} {
		size, n := tc[0], tc[1]
		next := 0
		for i := range n {
			a, b := span(0, size, i, n)
			if a > next || b <= a || b > size+1 {
				t.Fatalf("span %d of %d over %d is %d..%d after %d", i, n, size, a, b, next)
			}
			next = max(next, b)
		}
		if next < size {
			t.Fatalf("spans over %d stop at %d", size, next)
		}
	}
}

func TestAroundWalksTheEdgeOnce(t *testing.T) {
	w, h := 7, 5
	seen := map[[2]int]bool{}
	for d := range 2*(w+h) - 4 {
		x, y := around(w, h, d)
		if seen[[2]int{x, y}] || (x != 0 && x != w-1 && y != 0 && y != h-1) {
			t.Fatalf("step %d at %d,%d", d, x, y)
		}
		seen[[2]int{x, y}] = true
	}
}

func TestExtremesAndNames(t *testing.T) {
	if d, l := extremes(pixel.SpringDay); d != pixel.Ink || l != pixel.AccentInk {
		t.Fatalf("Spring Day extremes %v %v", d, l)
	}
	if d, l := extremes(pixel.SpringNight); d == l {
		t.Fatal("night extremes collapsed")
	}
	if len(emoji) != 16 || len(EmojiNames) != 16 {
		t.Fatal("emoji count")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("unknown rarity did not panic")
		}
	}()
	Rarity(9).Slot()
}
