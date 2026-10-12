package files

import (
	"image"
	"math"
	"slices"
	"testing"
)

func TestSizeSaysUnits(t *testing.T) {
	for n, want := range map[int64]string{
		0: "0 B", 1023: "1023 B", 1024: "1.0 KB", 1536: "1.5 KB", 10 << 10: "10 KB",
		1023<<10 + 600: "1.0 MB", 1 << 30: "1.0 GB", 41 << 30: "41 GB", 1 << 50: "1.0 PB", math.MaxInt64: "8192 PB",
	} {
		if got := size(n); got != want {
			t.Errorf("size(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestGroupedPutsCommasInThousands(t *testing.T) {
	for n, want := range map[int64]string{0: "0", 999: "999", 1000: "1,000", 1234567: "1,234,567"} {
		if got := grouped(n); got != want {
			t.Errorf("grouped(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestLabelKeepsTheExtension(t *testing.T) {
	for name, want := range map[string][]string{
		"MAP.PNG": {"MAP.PNG"}, "CHEATS.TXT": {"CHEATS", ".TXT"}, "SECRET LEVEL": {"SECRET", "LEVEL"},
	} {
		if got := label(name, 60); !slices.Equal(got, want) {
			t.Errorf("label(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestFitScalesWholeUpAndAnyDown(t *testing.T) {
	for _, tc := range []struct{ iw, ih, w, h, ww, wh int }{
		{64, 40, 160, 80, 128, 80}, {640, 400, 160, 120, 160, 100}, {10, 10, 160, 120, 120, 120},
	} {
		if w, h := fit(image.Rect(0, 0, tc.iw, tc.ih), tc.w, tc.h); w != tc.ww || h != tc.wh {
			t.Errorf("%d×%d in %d×%d is %d×%d, want %d×%d", tc.iw, tc.ih, tc.w, tc.h, w, h, tc.ww, tc.wh)
		}
	}
}
