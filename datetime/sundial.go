package datetime

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SunDialOptions struct {
	Disabled bool
}

type SunDialResult struct {
	Element ui.Element
	Changed bool
}

// Dial size in dots: the arc's middle, radius and box.
const dialW, dialH, dialR = 96, 56, 40

// SunDial arcs the sun or moon; a drag sets *hour.
func SunDial(c *ui.Context, text string, hour *int, o SunDialOptions) SunDialResult {
	k := pixel.K(c)
	if *hour < 0 || *hour > 23 {
		panic(fmt.Sprintf("datetime: SunDial hour %d outside 0–23", *hour))
	}
	old := *hour
	shell, base, st := pixel.Shell(c, func() ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleSlider).Label(text).Disabled(o.Disabled)
	})
	if pixel.Shortcut(base, 0, ui.KeyLeft) {
		*hour = (*hour + 23) % 24
	}
	if pixel.Shortcut(base, 0, ui.KeyRight) {
		*hour = (*hour + 1) % 24
	}
	// The width spans the whole day, midnight to midnight.
	if _, _, held := base.Dragged(); held && !st.Disabled && !pixel.Current(c).Inert {
		x, _, _ := base.PointerPosition()
		*hour = max(0, min(23, int(x/k.D(dialW)*24)))
	}
	h := *hour
	label := fmt.Sprintf("%02d:00", h)
	base.Description(label)
	day := h >= 6 && h < 18
	shell.Size(k.D(dialW), k.D(dialH)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		sky, body := pixel.Lit, icon.Sun
		if !day {
			sky, body = pixel.Shade, icon.Moon
		}
		g.Frame(0, 0, g.W, g.H, pixel.Sunken, sky)
		cx, cy := dialW/2, dialH-6
		g.HLine(2, cy, dialW-4, pixel.Edge)
		for i := 0; i <= 48; i += 2 {
			a := math.Pi * float64(i) / 48
			g.Dot(cx+int(math.Round(dialR*math.Cos(a))), cy-int(math.Round(dialR*math.Sin(a))), pixel.Faint)
		}
		t := float64((h+18)%12) / 12
		a := math.Pi * (1 - t)
		x, y := cx+int(math.Round(dialR*math.Cos(a))), cy-int(math.Round(dialR*math.Sin(a)))
		g.Sprite(x-body.W/2, y-body.H/2, body, 0, pixel.Ink, 1)
		g.Micro(3, 3, label, pixel.Ink, 1)
		g.Finish(st)
	})
	res := SunDialResult{Element: shell, Changed: *hour != old}
	if res.Changed {
		slog.Debug("8bit: sundial", "label", text, "hour", *hour)
	}
	return res
}
