package input

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Direction is a d-pad arm.
type Direction uint8

const (
	None Direction = iota
	Up
	Down
	Left
	Right
)

func (d Direction) String() string { return [...]string{"none", "up", "down", "left", "right"}[d] }

type DPadOptions struct {
	Disabled bool
}

type DPadResult struct {
	Element ui.Element
	// Pressed is the arm pressed this frame.
	Pressed Direction
}

type dpadKey struct{}

// flash is the arm a key lit, and when.
type flash struct {
	arm Direction
	at  time.Time
}

// DPad is a cross pad reporting arrows and pressed arms.
func DPad(c *ui.Context, text string, o DPadOptions) DPadResult {
	k := pixel.K(c)
	shell, base, st := pixel.Shell(c, func() ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Disabled(o.Disabled)
	})
	lit := ui.Local(base, dpadKey{}, func() flash { return flash{} })
	res := DPadResult{Element: shell}
	for _, m := range []struct {
		key ui.Key
		arm Direction
	}{{ui.KeyUp, Up}, {ui.KeyDown, Down}, {ui.KeyLeft, Left}, {ui.KeyRight, Right}} {
		if pixel.Shortcut(base, 0, m.key) {
			res.Pressed = m.arm
		}
	}
	// Taking Enter and Space keeps keys from clicking.
	pixel.Shortcut(base, 0, ui.KeyEnter)
	pixel.Shortcut(base, 0, ui.KeySpace)
	at := func() Direction {
		x, y, _ := base.PointerPosition()
		return armAt(int(x/k.Dot()), int(y/k.Dot()))
	}
	if base.Clicked() {
		res.Pressed = at()
	}
	if res.Pressed != None {
		*lit = flash{res.Pressed, k.Now()}
		slog.Debug("8bit: dpad", "label", text, "arm", res.Pressed)
	}
	down := None
	if k.Elapsed(&lit.at) < 150*time.Millisecond {
		down = lit.arm
		k.Tick(c, 10)
	}
	if st.Press {
		down = at()
	}
	shell.Size(k.D(36), k.D(36)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(12, 0, 12, 36, pixel.Edge)
		g.Rect(0, 12, 36, 12, pixel.Edge)
		for _, a := range []struct {
			arm        Direction
			x, y, w, h int
			sp         *pixel.Sprite
		}{{Up, 13, 1, 10, 11, tinyUp}, {Down, 13, 24, 10, 11, tinyDown}, {Left, 1, 13, 11, 10, tinyLeft}, {Right, 24, 13, 11, 10, tinyRight}} {
			fill, dy := pixel.Face, 0
			if down == a.arm {
				fill, dy = pixel.Shade, 1
			}
			g.Rect(a.x, a.y, a.w, a.h, fill)
			g.Sprite(a.x+(a.w-a.sp.W)/2, a.y+(a.h-a.sp.H)/2+dy, a.sp, 0, pixel.Ink, 1)
		}
		g.Rect(13, 13, 10, 10, pixel.Shade)
		g.Finish(st)
	})
	return res
}

// armAt is the arm under a dot of the pad.
func armAt(x, y int) Direction {
	switch {
	case x >= 12 && x < 24 && y < 12:
		return Up
	case x >= 12 && x < 24 && y >= 24:
		return Down
	case y >= 12 && y < 24 && x < 12:
		return Left
	case y >= 12 && y < 24 && x >= 24:
		return Right
	}
	return None
}
