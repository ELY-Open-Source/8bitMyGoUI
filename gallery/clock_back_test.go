package gallery

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestDemosSurviveTheClockSteppingBack(t *testing.T) {
	negative := regexp.MustCompile(`(^|[\s(:])-\d`)
	for _, e := range Entries() {
		at := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
		tt := ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }})
			ui.Column(c).Fill().Children(func() { e.Demo(c, func(string) {}) })
		}, 640, 480)
		// Press up to three buttons, run, then fall back.
		for _, b := range buttons(tt, 3) {
			if r, ok := tt.Find(b); ok {
				tt.ClickAt(r.X+r.W/2, r.Y+r.H/2)
			}
		}
		at = at.Add(500 * time.Millisecond)
		tt.Frame()
		// Negative numbers shown before are content, not faults.
		before := map[string]bool{}
		for _, s := range tt.Texts() {
			before[s] = true
		}
		for _, d := range []time.Duration{-time.Hour, 300 * time.Millisecond} {
			at = at.Add(d)
			tt.Frame()
		}
		for _, s := range tt.Texts() {
			if negative.MatchString(s) && !before[s] {
				t.Errorf("%s shows %q after the clock stepped back", e.Slug(), s)
			}
		}
	}
}

// buttons names up to n pressable words on screen.
func buttons(tt *ui.Tester, n int) []string {
	var out []string
	for _, s := range tt.Texts() {
		if len(out) < n && len(s) > 1 && s == upper(s) && strings.ContainsAny(s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			out = append(out, s)
		}
	}
	return out
}

func upper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}

func TestPlayingMediaSurvivesTheClockSteppingBack(t *testing.T) {
	for _, slug := range []string{"346-audioplayer", "353-cassette", "360-videoframe"} {
		e, _ := find(slug)
		at := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
		tt := ui.NewTester(func(c *ui.Context) {
			pixel.Use(c, pixel.Settings{Palette: pixel.SpringDay, Now: func() time.Time { return at }})
			ui.Column(c).Fill().Children(func() { e.Demo(c, func(string) {}) })
		}, 640, 480)
		if err := tt.Click("PLAY"); err != nil {
			t.Fatalf("%s: %v", slug, err)
		}
		for _, d := range []time.Duration{500 * time.Millisecond, -time.Second, 200 * time.Millisecond} {
			at = at.Add(d)
			tt.Frame()
		}
	}
}
