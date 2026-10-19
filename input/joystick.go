package input

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type JoystickOptions struct {
	Disabled bool
}

type JoystickResult struct {
	Element ui.Element
	// X and Y are in −1…1; Y grows downward.
	X, Y  float64
	Moved bool
}

type stickKey struct{}

// stick is the knob and its last key nudge.
type stick struct {
	x, y float64
	key  time.Time
}

const (
	stickSize = 40
	// stickReach is how far the knob travels, in dots.
	stickReach = 12
)

// Joystick drags a knob in a ring; arrows nudge it.
func Joystick(c *ui.Context, text string, o JoystickOptions) JoystickResult {
	k := pixel.K(c)
	shell, base, st := pixel.Shell(c, func() ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Disabled(o.Disabled)
	})
	s := ui.Local(base, stickKey{}, func() stick { return stick{} })
	ox, oy, now, on := s.x, s.y, k.Now(), !base.IsDisabled()
	for _, m := range []struct {
		key    ui.Key
		dx, dy float64
	}{{ui.KeyLeft, -0.25, 0}, {ui.KeyRight, 0.25, 0}, {ui.KeyUp, 0, -0.25}, {ui.KeyDown, 0, 0.25}} {
		if pixel.Shortcut(base, 0, m.key) {
			s.x, s.y, s.key = s.x+m.dx, s.y+m.dy, now
		}
	}
	_, _, held := base.Dragged()
	switch {
	case held && on:
		px, py, _ := base.PointerPosition()
		s.x = (float64(px/k.Dot()) - stickSize/2) / stickReach
		s.y = (float64(py/k.Dot()) - stickSize/2) / stickReach
		s.key = time.Time{}
	case k.Elapsed(&s.key) > repeatGap:
		s.x, s.y = 0, 0
	default:
		k.Tick(c, 10)
	}
	if l := math.Hypot(s.x, s.y); l > 1 {
		s.x, s.y = s.x/l, s.y/l
	}
	base.Description(fmt.Sprintf("%.2f, %.2f", s.x, s.y))
	x, y := stickSize/2+int(math.Round(s.x*stickReach)), stickSize/2+int(math.Round(s.y*stickReach))
	shell.Size(k.D(stickSize), k.D(stickSize)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		const m = stickSize / 2
		g.Disc(m, m, m-1, pixel.Edge)
		g.Disc(m, m, m-2, pixel.Well)
		g.Dither(m-2, 4, 4, stickSize-8, pixel.Well, pixel.Shade, 4)
		g.Dither(4, m-2, stickSize-8, 4, pixel.Well, pixel.Shade, 4)
		g.Disc(x+1, y+1, 7, pixel.Shade)
		g.Disc(x, y, 7, pixel.Edge)
		g.Disc(x, y, 6, pixel.Accent)
		g.Disc(x-2, y-2, 2, pixel.AccentLit)
		g.Finish(st)
	})
	return JoystickResult{Element: shell, X: s.x, Y: s.y, Moved: s.x != ox || s.y != oy}
}
