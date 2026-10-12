package account

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type OptionCyclerOptions struct {
	// Width is the value's room; zero fits the longest.
	Width    int
	Disabled bool
}

type OptionCyclerResult struct {
	Element *ui.Element
	Changed bool
}

// OptionCycler steps ◀ VALUE ▶ through options, wrapping around.
func OptionCycler(c *ui.Context, text string, options []string, selected *int, o OptionCyclerOptions) OptionCyclerResult {
	k := pixel.K(c)
	choice("OptionCycler", selected, len(options))
	longest := 0
	for _, s := range options {
		longest = max(longest, pixel.Measure(s))
	}
	w := pixel.Or("OptionCycler Width", o.Width, longest)
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleStepper).Label(text).Disabled(o.Disabled)
	})
	n, old := len(options), *selected
	if !base.IsDisabled() {
		step := 0
		if base.Shortcut(0, ui.KeyLeft) {
			step--
		}
		// Enter and Space step on; taken, clicks come from pointers.
		right, enter, space := base.Shortcut(0, ui.KeyRight), base.Shortcut(0, ui.KeyEnter), base.Shortcut(0, ui.KeySpace)
		if right || enter || space {
			step++
		}
		// A click on the left half steps back, else on.
		if base.Clicked() {
			if x, _, _ := base.PointerPosition(); x < shell.Bounds().W/2 {
				step--
			} else {
				step++
			}
		}
		*selected = ((*selected+step)%n + n) % n
	}
	base.Description(options[*selected])
	value := options[*selected]
	total := w + 32
	shell.Size(k.D(total), k.D(pixel.LineHeight+4)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Text(2, 2, "◀", pixel.Accent, 1)
		g.Text(16+(w-pixel.Measure(value))/2, 2, value, pixel.Ink, 1)
		g.Text(total-2-pixel.Measure("▶"), 2, "▶", pixel.Accent, 1)
		g.Finish(st)
	})
	res := OptionCyclerResult{Element: shell, Changed: *selected != old}
	if res.Changed {
		slog.Debug("8bit: option cycler", "label", text, "selected", *selected)
	}
	return res
}
