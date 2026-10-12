package project

import (
	"testing"
	"time"
)

func TestHmsDropsHoursUnderAnHour(t *testing.T) {
	for d, want := range map[time.Duration]string{0: "00:00", 59 * time.Second: "00:59", 25 * time.Minute: "25:00", 61 * time.Minute: "1:01:00", 1500 * time.Millisecond: "00:01"} {
		if got := hms(d); got != want {
			t.Errorf("hms(%v) = %q, want %q", d, got, want)
		}
	}
}
