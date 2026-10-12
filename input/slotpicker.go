package input

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SlotPickerOptions struct {
	Disabled bool
}

type SlotPickerResult struct {
	Element *ui.Element
	Changed bool
}

type reelKey struct{}

// reel holds a drag's and the wheel's leftover DIPs.
type reel struct {
	drag, wheel float32
	changed     bool
}

// SlotPicker spins a looping reel; the centered option is chosen.
func SlotPicker(c *ui.Context, text string, value *int, options []string, o SlotPickerOptions) SlotPickerResult {
	k := pixel.K(c)
	n := len(options)
	if *value < 0 || *value >= n {
		panic(fmt.Sprintf("input: SlotPicker value %d outside %d options", *value, n))
	}
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleSlider).Label(text).Description(options[*value]).Disabled(o.Disabled)
	})
	s := ui.Local(base, reelKey{}, func() reel { return reel{} })
	spin := func(d int) {
		*value = ((*value+d)%n + n) % n
		s.changed = true
	}
	row := k.D(pixel.LineHeight)
	base.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputScroll || o.Disabled {
			return false
		}
		s.wheel += ev.DY
		turn(&s.wheel, row, spin)
		return true
	})
	if base.Shortcut(0, ui.KeyDown) {
		spin(1)
	}
	if base.Shortcut(0, ui.KeyUp) {
		spin(-1)
	}
	if _, dy, held := base.Dragged(); held && !o.Disabled {
		s.drag += dy
		turn(&s.drag, row, func(d int) { spin(-d) })
	} else {
		s.drag = 0
	}
	res := SlotPickerResult{Element: shell, Changed: s.changed}
	s.changed = false
	if res.Changed {
		slog.Debug("8bit: slot picker", "label", text, "value", options[*value])
	}
	off := int(s.drag / k.Dot())
	shell.Size(k.D(widest(options)+28), k.D(3*pixel.LineHeight+4)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
		g.Clip(2, 2, g.W-4, g.H-4, func() {
			for i := -2; i <= 2; i++ {
				ink := pixel.Faint
				if i == 0 {
					ink = pixel.Ink
				}
				opt := options[((*value+i)%n+n)%n]
				g.Text((g.W-pixel.Measure(opt))/2, 2+pixel.LineHeight*(1+i)+off, opt, ink, 1)
			}
		})
		g.Veil(2, 2, g.W-4, 6, pixel.Shade, 4)
		g.Veil(2, g.H-8, g.W-4, 6, pixel.Shade, 4)
		mid := 2 + pixel.LineHeight + (pixel.LineHeight-tinyRight.H)/2
		g.Sprite(3, mid, tinyRight, 0, pixel.Accent, 1)
		g.Sprite(g.W-6, mid, tinyLeft, 0, pixel.Accent, 1)
		finish(g, st)
	})
	return res
}

// turn spends whole rows of *left on step.
func turn(left *float32, row float32, step func(int)) {
	for ; *left >= row; *left -= row {
		step(1)
	}
	for ; *left <= -row; *left += row {
		step(-1)
	}
}
