package desk

import (
	"slices"
	"testing"
	"time"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestCalcPrintsAndPresses(t *testing.T) {
	for keys, want := range map[string]string{
		"100": "100", "007": "7", "0": "0", "1234567890123": "123456789012", "5+": "5", "5-9=": "-4", "2÷4=": "0.5",
	} {
		var cl calc
		for _, r := range keys {
			cl.press(string(r))
		}
		if got := cl.display(); got != want {
			t.Errorf("%s shows %q, want %q", keys, got, want)
		}
	}
	cl := calc{acc: -0.0000000000001}
	if got := cl.display(); got != "0" {
		t.Errorf("a tiny minus shows %q", got)
	}
}

func TestFieldLaysFairAndOpens(t *testing.T) {
	f := newField(25)
	f.lay(5, 5, 10, 12, 1, time.Unix(1, 0))
	for _, i := range append(around(5, 5, 12), 12) {
		if f.mine[i] {
			t.Fatalf("a mine at %d, beside the first click", i)
		}
	}
	// A crowded field keeps only the first click safe.
	g := newField(9)
	g.lay(3, 3, 8, 4, 1, time.Unix(1, 0))
	if g.mine[4] || slices.Index(g.mine, false) != 4 {
		t.Fatalf("crowded mines %v", g.mine)
	}
	// With mines on 7 and 8, the flood misses 6.
	h := newField(9)
	h.mine[7], h.mine[8], h.laid = true, true, true
	h.reveal(3, 3, 0, time.Unix(2, 0))
	if h.over || !h.open[4] || !h.open[5] || h.open[6] {
		t.Fatalf("a flood from a corner opened %v", h.open)
	}
	h.reveal(3, 3, 8, time.Unix(3, 0))
	if !h.over || h.won || h.boom != 8 {
		t.Fatal("stepping on the mine did not end it")
	}
}

func TestPitGrowsAndCrashes(t *testing.T) {
	var s pit
	s.reset(6, 4, 1, 1)
	s.food = [2]int{4, 2}
	if s.step(6, 4) || len(s.body) != 4 || s.score != 1 {
		t.Fatalf("eating gave %v, score %d", s.body, s.score)
	}
	// Into the cell the tail leaves is no crash.
	s.body, s.next, s.food = [][2]int{{1, 1}, {1, 2}, {2, 2}, {2, 1}}, [2]int{1, 0}, [2]int{5, 3}
	if s.step(6, 4) {
		t.Fatal("chasing its tail crashed")
	}
	s.next = [2]int{0, -1}
	s.step(6, 4)
	if !s.step(6, 4) || !s.dead {
		t.Fatal("the top wall did not stop it")
	}
}

func TestSwapsSortTheDisk(t *testing.T) {
	from := scatter(100, 2)
	cells := slices.Clone(from)
	for _, s := range swaps(from) {
		cells[s[0]], cells[s[1]] = cells[s[1]], cells[s[0]]
	}
	if !slices.IsSortedFunc(cells, func(a, b int) int { return order(a) - order(b) }) {
		t.Fatalf("not sorted: %v", cells)
	}
}

func TestFieldTimeNeverGoesBelowZero(t *testing.T) {
	f := field{laid: true, over: true, start: time.Unix(100, 0), end: time.Unix(95, 0)}
	if s := f.seconds(nil, pixel.Kit{}); s != 0 {
		t.Fatalf("an end before the start shows %d", s)
	}
}
