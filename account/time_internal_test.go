package account

import (
	"math"
	"testing"
	"time"
)

func TestSecondsLeftRoundsUpWithoutOverflow(t *testing.T) {
	for _, tc := range []struct {
		now    time.Time
		period time.Duration
		want   int64
	}{
		{time.Unix(0, 0), 30 * time.Second, 30},
		{time.Unix(1, 0), 30 * time.Second, 29},
		{time.Unix(29, 5e8), 30 * time.Second, 1},
		{time.Unix(0, 0), math.MaxInt64, 9223372037},
	} {
		if got := secondsLeft(tc.now, tc.period); got != tc.want {
			t.Errorf("%v into %v: %d s left, want %d", tc.now.Unix(), tc.period, got, tc.want)
		}
	}
}
