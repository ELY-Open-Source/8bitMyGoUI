package text

import (
	"reflect"
	"testing"
)

func TestWaveAndShake(t *testing.T) {
	want := []int{0, 1, 2, 1, 0, -1, -2, -1}
	for n, w := range want {
		if got := wave(2, n); got != w {
			t.Fatalf("wave(2, %d) = %d, want %d", n, got, w)
		}
	}
	moved := false
	for step := range 20 {
		dx, dy := shake(false, 7, 3, step)
		if dx < -1 || dx > 1 || dy < -1 || dy > 1 {
			t.Fatalf("shake %d,%d outside one dot", dx, dy)
		}
		moved = moved || dx != 0 || dy != 0
		if rx, ry := shake(true, 7, 3, step); rx != 0 || ry != 0 {
			t.Fatal("shake moved under Reduce Motion")
		}
	}
	if !moved {
		t.Fatal("shake never moved")
	}
	if a, b := noise(1, 2, 3), noise(1, 2, 3); a != b {
		t.Fatal("noise is not seeded")
	}
}

func TestRollAndWheel(t *testing.T) {
	v, steps := 0.0, 0
	for v != 100 {
		v = roll(v, 100)
		if steps++; steps > 100 {
			t.Fatal("roll never reached the target")
		}
	}
	if got := roll(0, 100); got != 100.0/6 {
		t.Fatalf("first roll step %v", got)
	}
	for _, tc := range []struct {
		shown float64
		i, d  int
		frac  float64
	}{{129.5, 0, 9, 0.5}, {129.5, 1, 2, 0.5}, {129.5, 2, 1, 0}, {125, 1, 2, 0}} {
		if d, frac := wheel(tc.shown, tc.i); d != tc.d || frac != tc.frac {
			t.Errorf("wheel(%v, %d) = %d %v, want %d %v", tc.shown, tc.i, d, frac, tc.d, tc.frac)
		}
	}
}

func TestMatches(t *testing.T) {
	for _, tc := range []struct {
		s, q string
		want [][2]int
	}{
		{"Slime, King Slime", "slime", [][2]int{{0, 5}, {12, 17}}},
		{"aaa", "aa", [][2]int{{0, 2}}},
		{"水晶水晶", "晶", [][2]int{{1, 2}, {3, 4}}},
		{"abc", "", nil},
	} {
		if got := matches([]rune(tc.s), []rune(tc.q)); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("matches(%q, %q) = %v, want %v", tc.s, tc.q, got, tc.want)
		}
	}
}

func TestCellWidths(t *testing.T) {
	if cellW('8') != 7 || cellW(':') != 3 || cellW('-') != 7 {
		t.Fatal("cell widths")
	}
	defer func() {
		if recover() == nil {
			t.Fatal("SevenSegment took an x")
		}
	}()
	cellW('x')
}
