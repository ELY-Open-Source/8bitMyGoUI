package datetime

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TimeInputOptions struct {
	Disabled bool
}

type TimeInputResult struct {
	Element *ui.Element
	Changed bool
}

type timeKey struct{}

// timing is the segment edited, and a first digit typed.
type timing struct {
	seg, first int
}

// TimeInput edits HH:MM; Up and Down step, digits type.
func TimeInput(c *ui.Context, text string, tm *time.Time, o TimeInputOptions) TimeInputResult {
	k := pixel.K(c)
	old := *tm
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleTextField).Label(text).Disabled(o.Disabled)
	})
	s := ui.Local(base, timeKey{}, func() timing { return timing{first: -1} })
	h, m := tm.Hour(), tm.Minute()
	set := func(hour, minute int) {
		*tm = time.Date(tm.Year(), tm.Month(), tm.Day(), hour, minute, 0, 0, tm.Location())
		h, m = hour, minute
	}
	step := func(d int) {
		if s.seg == 0 {
			set((h+d+24)%24, m)
		} else {
			set(h, (m+d+60)%60)
		}
		s.first = -1
	}
	if base.Shortcut(0, ui.KeyUp) {
		step(1)
	}
	if base.Shortcut(0, ui.KeyDown) {
		step(-1)
	}
	for _, mv := range []struct {
		key ui.Key
		seg int
	}{{ui.KeyLeft, 0}, {ui.KeyRight, 1}} {
		if base.Shortcut(0, mv.key) {
			s.seg, s.first = mv.seg, -1
		}
	}
	for d := range 10 {
		if !base.Shortcut(0, ui.Key0+ui.Key(d)) {
			continue
		}
		top := 23
		if s.seg == 1 {
			top = 59
		}
		v := d
		if s.first >= 0 && s.first*10+d <= top {
			v = s.first*10 + d
		}
		if s.seg == 0 {
			set(v, m)
		} else {
			set(h, v)
		}
		// A second digit ends a segment; hours pass to minutes.
		if s.first >= 0 {
			s.first, s.seg = -1, 1
		} else {
			s.first = d
		}
	}
	// Enter and Space are taken, so clicks come from pointers.
	base.Shortcut(0, ui.KeyEnter)
	base.Shortcut(0, ui.KeySpace)
	if base.Clicked() {
		x, _, _ := base.PointerPosition()
		s.seg, s.first = 0, -1
		if int(x/k.Dot()) >= 4+pixel.Measure("00:") {
			s.seg = 1
		}
	}
	shown := fmt.Sprintf("%02d:%02d", h, m)
	base.Description(shown)
	seg := s.seg
	shell.Size(k.D(pixel.Measure("00:00")+10), k.D(pixel.LineHeight+6)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
		hh, mm, at := shown[:2], shown[3:], 5
		for i, part := range []string{hh, ":", mm} {
			ink, w := pixel.Ink, pixel.Measure(part)
			if st.Focus && (i == 0 && seg == 0 || i == 2 && seg == 1) {
				g.Rect(at-1, 2, w+2, g.H-4, pixel.Accent)
				ink = pixel.AccentInk
			}
			g.Text(at, 3, part, ink, 1)
			at += w
		}
		g.Finish(st)
	})
	res := TimeInputResult{Element: shell, Changed: !tm.Equal(old)}
	if res.Changed {
		slog.Debug("8bit: time input", "label", text, "time", shown)
	}
	return res
}
