package tabletop

import "testing"

func TestAnswerWinsBlocksThenFills(t *testing.T) {
	for _, tc := range []struct {
		b    [9]int
		want int
	}{
		{[9]int{2, 2, 0, 1, 1, 0, 0, 0, 0}, 2},
		{[9]int{1, 1, 0, 0, 2, 0, 0, 0, 0}, 2},
		{[9]int{1, 0, 0, 0, 0, 0, 0, 0, 0}, 4},
		{[9]int{0, 0, 0, 0, 1, 0, 0, 0, 0}, 0},
		{[9]int{1, 2, 1, 2, 2, 1, 0, 1, 0}, 8}, {[9]int{1, 2, 1, 1, 2, 2, 2, 1, 0}, 8},
	} {
		if got := answer(tc.b); got != tc.want {
			t.Errorf("answer(%v) = %d, want %d", tc.b, got, tc.want)
		}
	}
}

func TestFullSeesEveryLine(t *testing.T) {
	lines := [][]int{{0, 1, 2, 3, 4}, {2, 7, 12, 17, 22}, {0, 6, 12, 18, 24}, {4, 8, 12, 16, 20}}
	for _, l := range lines {
		var m [25]bool
		for _, i := range l[:4] {
			m[i] = true
		}
		if full(&m) {
			t.Fatalf("four of %v is full", l)
		}
		m[l[4]] = true
		if !full(&m) {
			t.Fatalf("%v is not full", l)
		}
	}
}
