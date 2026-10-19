package pixel_test

import (
	"math"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func luminance(c ui.Color) float64 {
	lin := func(v uint8) float64 {
		s := float64(v) / 255
		if s <= 0.04045 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	return 0.2126*lin(c.R) + 0.7152*lin(c.G) + 0.0722*lin(c.B)
}

func contrast(a, b ui.Color) float64 {
	la, lb := luminance(a), luminance(b)
	return (max(la, lb) + 0.05) / (min(la, lb) + 0.05)
}

func TestPaletteContrast(t *testing.T) {
	rules := []struct {
		fg, bg pixel.Slot
		min    float64
	}{
		{pixel.Ink, pixel.Face, 4.5},
		{pixel.Ink, pixel.Back, 4.5},
		{pixel.Ink, pixel.Well, 4.5},
		{pixel.AccentInk, pixel.Accent, 4.5},
		{pixel.Faint, pixel.Face, 3},
		{pixel.Edge, pixel.Back, 3},
	}
	for _, p := range pixel.Palettes {
		for _, r := range rules {
			if got := contrast(p.C(r.fg), p.C(r.bg)); got < r.min {
				t.Errorf("%s: slot %d on %d contrast %.2f, want ≥ %.1f", p.Name, r.fg, r.bg, got, r.min)
			}
		}
	}
}

func TestPaletteColorsDistinctAndOpaque(t *testing.T) {
	if len(pixel.Palettes) != 8 {
		t.Fatalf("got %d palettes, want 8", len(pixel.Palettes))
	}
	for _, p := range pixel.Palettes {
		seen := map[ui.Color]pixel.Slot{}
		for s, c := range p.Colors {
			if c.A != 255 {
				t.Errorf("%s slot %d alpha %d", p.Name, s, c.A)
			}
			if prev, ok := seen[c]; ok {
				t.Errorf("%s slots %d and %d share %v", p.Name, prev, s, c)
			}
			seen[c] = pixel.Slot(s)
		}
	}
}

func TestPaletteBevelOrder(t *testing.T) {
	for _, p := range pixel.Palettes {
		face := luminance(p.C(pixel.Face))
		if luminance(p.C(pixel.Lit)) <= face || luminance(p.C(pixel.Shade)) >= face {
			t.Errorf("%s: Lit must be lighter and Shade darker than Face", p.Name)
		}
		if p.Night != (luminance(p.C(pixel.Back)) < luminance(p.C(pixel.Ink))) {
			t.Errorf("%s: night palettes are dark, day palettes light", p.Name)
		}
	}
}

func TestFind(t *testing.T) {
	seen := map[*pixel.Palette]bool{}
	for s := pixel.Spring; s <= pixel.Winter; s++ {
		for _, night := range []bool{false, true} {
			p := pixel.Find(s, night)
			if p.Season != s || p.Night != night {
				t.Errorf("Find(%d, %v) = %s", s, night, p.Name)
			}
			seen[p] = true
		}
	}
	if len(seen) != 8 {
		t.Fatalf("Find reached %d palettes, want 8", len(seen))
	}
}

func TestCycle(t *testing.T) {
	band := []pixel.Slot{pixel.Info, pixel.Accent, pixel.Lit}
	cases := []struct {
		step int
		want []pixel.Slot
	}{
		{0, []pixel.Slot{pixel.Info, pixel.Accent, pixel.Lit}},
		{1, []pixel.Slot{pixel.Accent, pixel.Lit, pixel.Info}},
		{-1, []pixel.Slot{pixel.Lit, pixel.Info, pixel.Accent}},
		{7, []pixel.Slot{pixel.Accent, pixel.Lit, pixel.Info}},
		{math.MaxInt, []pixel.Slot{pixel.Accent, pixel.Lit, pixel.Info}},
		{math.MinInt, []pixel.Slot{pixel.Accent, pixel.Lit, pixel.Info}},
	}
	for _, c := range cases {
		if got := pixel.Cycle(band, c.step); !slices.Equal(got, c.want) {
			t.Errorf("Cycle(step %d) = %v, want %v", c.step, got, c.want)
		}
	}
	if !slices.Equal(band, []pixel.Slot{pixel.Info, pixel.Accent, pixel.Lit}) {
		t.Fatalf("Cycle changed its input: %v", band)
	}
}
