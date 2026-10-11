package pixel

import (
	"slices"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestRuns(t *testing.T) {
	cases := map[string][][2]int{
		"##..###.#": {{0, 2}, {4, 7}, {8, 9}},
		"....":      nil,
		"####":      {{0, 4}},
		".#.#":      {{1, 2}, {3, 4}},
	}
	for pat, want := range cases {
		var got [][2]int
		runs(len(pat), func(i int) bool { return pat[i] == '#' }, func(i, j int) { got = append(got, [2]int{i, j}) })
		if !slices.Equal(got, want) {
			t.Errorf("runs(%q) = %v, want %v", pat, got, want)
		}
	}
}

func TestTextLinesEllipsis(t *testing.T) {
	cases := []struct {
		s        string
		width, n int
		want     []string
	}{
		{"one\ntwo\nthree", 0, 2, []string{"one", "two…"}},
		{"one\ntwo", 0, 2, []string{"one", "two"}},
		{"hello world again", Measure("hello…"), 1, []string{"hello…"}},
		{"hello world again", Measure("hel…"), 1, []string{"hel…"}},
		{"像素风格组件库", 36, 2, []string{"像素风", "格组…"}},
	}
	for _, c := range cases {
		if got := textLines(c.s, c.width, c.n); !slices.Equal(got, c.want) {
			t.Errorf("textLines(%q, %d, %d) = %q, want %q", c.s, c.width, c.n, got, c.want)
		}
	}
}

type fakePainter struct{ fills int }

func (f *fakePainter) Fill(ui.Rect, ui.Color, float32)      { f.fills++ }
func (f *fakePainter) Clip(_ ui.Rect, _ float32, fn func()) { fn() }
func (f *fakePainter) After(time.Duration)                  {}
func (f *fakePainter) Now() time.Time                       { return time.Unix(0, 0) }

func TestDrawingMergesRuns(t *testing.T) {
	sp := MustSprite("aab\n.bb")
	cases := []struct {
		name string
		draw func(g Canvas)
		want int
	}{
		{"text A", func(g Canvas) { g.Text(0, 0, "A", Ink, 1) }, 15},
		{"text A ×2", func(g Canvas) { g.Text(0, 0, "A", Ink, 2) }, 15},
		{"micro 1", func(g Canvas) { g.Micro(0, 0, "1", Ink, 1) }, 5},
		{"sprite", func(g Canvas) { g.Sprite(0, 0, sp, 0, Ink, 1) }, 3},
		{"flat line", func(g Canvas) { g.Line(0, 0, 9, 0, Ink) }, 1},
		{"steep line", func(g Canvas) { g.Line(0, 0, 0, 9, Ink) }, 1},
		{"diagonal", func(g Canvas) { g.Line(0, 0, 9, 9, Ink) }, 10},
		{"circle 3", func(g Canvas) { g.Circle(3, 3, 3, Ink) }, 12},
		{"disc 3", func(g Canvas) { g.Disc(3, 3, 3, Ink) }, 7},
		{"dither 0", func(g Canvas) { g.Dither(0, 0, 4, 4, Face, Ink, 0) }, 1},
		{"dither 8", func(g Canvas) { g.Dither(0, 0, 4, 4, Face, Ink, 8) }, 9},
		{"dither 16", func(g Canvas) { g.Dither(0, 0, 4, 4, Face, Ink, 16) }, 1},
	}
	for _, c := range cases {
		f := &fakePainter{}
		c.draw(Canvas{p: f, k: Kit{pal: SpringDay, dot: 1}, W: 32, H: 32})
		if f.fills != c.want {
			t.Errorf("%s: %d fills, want %d", c.name, f.fills, c.want)
		}
	}
}
