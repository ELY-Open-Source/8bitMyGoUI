package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type PinInputOptions struct {
	// Length is the digit count; zero means 4.
	Length   int
	Disabled bool
}

type PinInputResult struct {
	Element ui.Element
	// Complete reports the last box filled.
	Changed, Complete bool
}

type pinKey struct{}

// pin: a change since last frame, and the prior value.
type pin struct {
	changed bool
	kept    string
}

const pinBox, pinGap = 12, 4

// PinInput types digits into boxes; Backspace goes back, paste fills.
func PinInput(c *ui.Context, text string, value *string, o PinInputOptions) PinInputResult {
	k := pixel.K(c)
	n := pixel.Or("PinInput Length", o.Length, 4)
	shell, base, st := pixel.Shell(c, func() ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleTextField).Label(text).Description(*value).Disabled(o.Disabled)
	})
	p := ui.Local(base, pinKey{}, func() pin { return pin{} })
	changed := &p.changed
	// The handler runs before this frame knows it is off.
	if p.changed && base.IsDisabled() {
		*value, p.changed = p.kept, false
	}
	put := func(s string) {
		for _, r := range s {
			if r >= '0' && r <= '9' && len(*value) < n {
				*value += string(r)
				*changed = true
			}
		}
	}
	base.HandleInput(func(ev ui.InputEvent) bool {
		if !p.changed {
			p.kept = *value
		}
		switch {
		case ev.Kind == ui.InputText:
			put(ev.Text)
		// Browsers paste by command, native hosts by Cmd+V.
		case ev.Kind == ui.InputCommand && ev.Text == "paste", ev.Kind == ui.InputKeyDown && ev.Key == ui.KeyV && ev.Mods == ui.Cmd:
			put(c.Services().ReadClipboard())
		case ev.Kind == ui.InputKeyDown && ev.Key == ui.KeyBackspace && ev.Mods == 0:
			if len(*value) > 0 {
				*value, *changed = (*value)[:len(*value)-1], true
			}
		default:
			return false
		}
		return true
	})
	h := pixel.LineHeight + 4
	base.TextCaret(ui.Rect{X: k.D(len(*value) * (pinBox + pinGap)), W: k.D(1), H: k.D(h)})
	res := PinInputResult{Element: shell, Changed: *changed, Complete: *changed && len(*value) == n}
	*changed = false
	focused := base.Focused()
	shell.Size(k.D(n*(pinBox+pinGap)-pinGap), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i := range n {
			x := i * (pinBox + pinGap)
			g.Frame(x, 0, pinBox, h, pixel.Sunken, pixel.Well)
			switch {
			case i < len(*value):
				d := (*value)[i : i+1]
				g.Text(x+(pinBox-pixel.Measure(d))/2, 2, d, pixel.Ink, 1)
			case i == len(*value) && focused:
				g.HLine(x+2, h-4, pinBox-4, pixel.Accent)
			}
		}
		g.Finish(st)
	})
	if res.Complete {
		slog.Debug("8bit: pin complete", "label", text)
	}
	return res
}
