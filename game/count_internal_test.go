package game

import (
	"math"
	"testing"
	"time"
)

func TestCountsNeverOverflow(t *testing.T) {
	for _, tc := range [][3]int{{math.MaxInt64, 15, math.MaxInt64}, {100, 7, 46}, {0, 3, 0}, {29, 15, 29}} {
		if got := upTo(tc[0], tc[1], 15); got != tc[2] {
			t.Errorf("upTo(%d, %d, 15) = %d, want %d", tc[0], tc[1], got, tc[2])
		}
	}
	if got := upTo(time.Duration(math.MaxInt64), 15, 15); got != math.MaxInt64 {
		t.Errorf("a long run counts to %d", got)
	}
}
