package agent

import (
	"math"
	"testing"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestTokensPrintShort(t *testing.T) {
	for n, want := range map[int]string{
		0: "0", 950: "950", 1000: "1K", 1500: "1.5K", 9999: "9.9K", 12345: "12K", 999999: "999K",
		1_000_000: "1M", 2_500_000: "2.5M", 1_000_000_000: "1B", math.MaxInt: "9.2E",
	} {
		if got := tokens(n); got != want {
			t.Errorf("tokens(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestRingRunsClockwiseOnce(t *testing.T) {
	pts := ring(17)
	seen := map[[2]int]bool{}
	for i, p := range pts {
		if seen[p] {
			t.Fatalf("dot %v twice", p)
		}
		seen[p] = true
		if d := math.Hypot(float64(p[0]), float64(p[1])); math.Abs(d-17) > 1 {
			t.Fatalf("dot %v is %.1f from the middle", p, d)
		}
		// Twelve o'clock first, then right before left.
		if i == 0 && p != [2]int{0, -17} {
			t.Fatalf("the ring starts at %v", p)
		}
	}
	if q := pts[len(pts)/4]; q[0] < 10 || q[1] < -10 {
		t.Fatalf("a quarter round is at %v, not near three o'clock", q)
	}
}

func TestTailKeepsTheFileName(t *testing.T) {
	long := "pixel/very/long/path/canvas.go"
	if got := tail(long, pixel.Measure("…/canvas.go")); got != "…/canvas.go" {
		t.Errorf("tail = %q", got)
	}
	if got := tail("a.go", 100); got != "a.go" {
		t.Errorf("a short path became %q", got)
	}
}
