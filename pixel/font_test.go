package pixel

import (
	"encoding/binary"
	"encoding/hex"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

func TestFontHeader(t *testing.T) {
	if font12.n != 24870 || font12.ascent != 13 || font12.descent != 3 {
		t.Fatalf("font: %d glyphs, ascent %d, descent %d", font12.n, font12.ascent, font12.descent)
	}
	if LineHeight != font12.ascent+font12.descent {
		t.Fatalf("LineHeight %d, want %d", LineHeight, font12.ascent+font12.descent)
	}
	prev := rune(-1)
	for i := range font12.n {
		r := rune(binary.LittleEndian.Uint32(font12.index[8*i:]))
		if r <= prev {
			t.Fatalf("index entry %d: %U after %U", i, r, prev)
		}
		prev = r
	}
}

func TestBadFontPanics(t *testing.T) {
	for name, b := range map[string][]byte{
		"magic":     []byte("XXXX\x0d\x03\x00\x00\x00\x00"),
		"truncated": []byte("8BF1\x0d\x03\x05\x00\x00\x00"),
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: loadFont did not panic", name)
				}
			}()
			loadFont(b)
		}()
	}
}

// Rows copied from ark-pixel-12px-proportional-zh_hans.bdf v2026.09.25.
func TestGlyphsMatchBDF(t *testing.T) {
	cases := []struct {
		r                    rune
		adv, w, h, x, y, top int
		hexRows              string
	}{
		{'A', 8, 7, 9, 0, 0, 4, "10 10 28 28 44 7C 44 82 82"},
		{'g', 6, 5, 8, 0, -2, 7, "78 88 88 88 88 78 08 70"},
		{'中', 12, 11, 11, 0, -1, 3, "0400 0400 FFE0 8420 8420 8420 FFE0 0400 0400 0400 0400"},
		{'，', 12, 2, 4, 2, -1, 10, "C0 C0 40 80"},
		{' ', 6, 0, 0, 0, 0, 13, ""},
	}
	for _, c := range cases {
		g, ok := font12.find(c.r)
		if !ok {
			t.Fatalf("%q missing", c.r)
		}
		if g.adv != c.adv || g.w != c.w || g.h != c.h || g.x != c.x || g.y != c.y {
			t.Fatalf("%q metrics adv %d box %d×%d at %d,%d", c.r, g.adv, g.w, g.h, g.x, g.y)
		}
		var want []byte
		for _, row := range strings.Fields(c.hexRows) {
			b, _ := hex.DecodeString(row)
			want = append(want, b...)
		}
		if !slices.Equal(g.rows, want) {
			t.Fatalf("%q rows % X, want % X", c.r, g.rows, want)
		}
		if got := g.top(); got != c.top {
			t.Fatalf("%q top %d, want %d", c.r, got, c.top)
		}
	}
	a := glyphOf('A')
	if !a.on(3, 0) || a.on(0, 0) || !a.on(0, 8) || !a.on(6, 8) {
		t.Fatal("'A' bits read wrong")
	}
}

func TestMissingRuneDrawsTofu(t *testing.T) {
	const r = '\U0001F600'
	if _, ok := font12.find(r); ok {
		t.Fatalf("%U is in the font; pick a missing rune", r)
	}
	if g := glyphOf(r); g.adv != tofu.adv || !slices.Equal(g.rows, tofu.rows) {
		t.Fatalf("%U = adv %d, want tofu", r, g.adv)
	}
}

func TestMeasure(t *testing.T) {
	if got := Measure("Ag中"); got != 8+6+12 {
		t.Fatalf("Measure = %d, want 26", got)
	}
	if Measure("") != 0 {
		t.Fatal("Measure of empty text is not 0")
	}
}

func TestLines(t *testing.T) {
	cases := []struct {
		s     string
		width int
		want  []string
	}{
		{"HELLO WORLD", max(Measure("HELLO"), Measure("WORLD")) + 2, []string{"HELLO", "WORLD"}},
		{"像素风格组件库", 36, []string{"像素风", "格组件", "库"}},
		{"你好，世界", 24, []string{"你", "好，", "世界"}},
		{"「像素」", 24, []string{"「像", "素」"}},
		{"A\nB", 100, []string{"A", "B"}},
		{"a long line", 0, []string{"a long line"}},
		{"x  y", Measure("x"), []string{"x", "y"}},
		{"你 ，好", 24, []string{"你 ，", "好"}},
		{"「 像素", 18, []string{"「 像", "素"}},
		{"［你好］世界", 36, []string{"［你", "好］世", "界"}},
		{"「  像素", 18, []string{"「  像", "素"}},
		{"「（像素）」", 18, []string{"「（像", "素）」"}},
		{"「（像素）」", 24, []string{"「（像", "素）」"}},
		{"AB「（C", Measure("AB") + 4, []string{"AB", "「（C"}},
	}
	for _, c := range cases {
		if got := Lines(c.s, c.width); !slices.Equal(got, c.want) {
			t.Errorf("Lines(%q, %d) = %q, want %q", c.s, c.width, got, c.want)
		}
	}
}

func TestLinesFitAndKeepText(t *testing.T) {
	texts := []string{
		"The quick brown fox jumps over the lazy dog near the riverbank.",
		"像素风格的组件库，在限制中创造丰富。每个点都经过认真摆放！",
		"Mixed 中英文 text with ABCDEFGHIJKLMNOPQRSTUVWXYZ and 標點，結束。",
		"Supercalifragilisticexpialidocious",
		"「像素」风格 （测试） 文本 ，结束 。再见",
		"［测试］｛内容｝〔括号〕《书名》“引号”‘单引’",
		"「（像素）」  「  空格」 『《嵌套》』",
	}
	for _, s := range texts {
		for width := 6; width <= 120; width += 7 {
			lines := Lines(s, width)
			for _, l := range lines {
				core := strings.TrimLeftFunc(strings.TrimRightFunc(l, hanging(closer)), hanging(opener))
				if Measure(core) > width && len([]rune(core)) > 1 {
					t.Fatalf("Lines(%q, %d): %q is %d wide", s, width, l, Measure(l))
				}
				if r := []rune(l); closer(r[0]) {
					t.Fatalf("Lines(%q, %d): %q starts with a closer", s, width, l)
				}
				if r := []rune(l); opener(r[len(r)-1]) {
					t.Fatalf("Lines(%q, %d): %q ends with an opener", s, width, l)
				}
			}
			if got := strings.ReplaceAll(strings.Join(lines, ""), " ", ""); got != strings.ReplaceAll(s, " ", "") {
				t.Fatalf("Lines(%q, %d) lost text: %q", s, width, got)
			}
		}
	}
}

func TestLinesRandom(t *testing.T) {
	alphabet := []rune("ab  像素，。「」（）［］《》“”!?.WM")
	rng := rand.New(rand.NewPCG(8, 16))
	for n := range 3000 {
		rs := make([]rune, 1+rng.IntN(24))
		for i := range rs {
			rs[i] = alphabet[rng.IntN(len(alphabet))]
		}
		s, width := string(rs), 6+rng.IntN(60)
		lines := Lines(s, width)
		for j, l := range lines {
			r := []rune(strings.TrimRight(l, " "))
			if len(r) == 0 {
				continue
			}
			core := strings.TrimLeftFunc(strings.TrimRightFunc(l, hanging(closer)), hanging(opener))
			if j > 0 && closer(r[0]) || j < len(lines)-1 && opener(r[len(r)-1]) || Measure(core) > width && len([]rune(core)) > 1 {
				t.Fatalf("case %d: Lines(%q, %d) = %q", n, s, width, lines)
			}
		}
		if got := strings.ReplaceAll(strings.Join(lines, ""), " ", ""); got != strings.ReplaceAll(s, " ", "") {
			t.Fatalf("case %d: Lines(%q, %d) lost text: %q", n, s, width, lines)
		}
	}
}

func hanging(punct func(rune) bool) func(rune) bool {
	return func(r rune) bool { return r == ' ' || punct(r) }
}

func TestMicroFont(t *testing.T) {
	set := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ+-.:/%! "
	seen := map[[MicroHeight]uint8]rune{}
	for _, r := range set {
		g := microGlyph(r)
		if prev, ok := seen[g]; ok && r != ' ' {
			t.Errorf("micro %q and %q are the same glyph", prev, r)
		}
		seen[g] = r
	}
	if microGlyph('a') != microGlyph('A') {
		t.Fatal("micro font does not upper-case")
	}
	if got := microGlyph('1'); got != [MicroHeight]uint8{2, 6, 2, 2, 7} {
		t.Fatalf("micro '1' = %v", got)
	}
	if MicroWidth("HP 99") != 19 || MicroWidth("") != 0 {
		t.Fatalf("MicroWidth = %d, %d", MicroWidth("HP 99"), MicroWidth(""))
	}
	defer func() {
		if recover() == nil {
			t.Fatal("micro '#' did not panic")
		}
	}()
	MicroWidth("#")
}
