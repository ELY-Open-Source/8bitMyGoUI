package chat

import (
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ComposerOptions struct {
	Placeholder string
	// Width is in dots, zero 160; Lines zero means 2.
	Width, Lines int
	Disabled     bool
}

type ComposerResult struct {
	Element *ui.Element
	Changed bool
	// Sent is the text sent this frame, else empty.
	Sent string
}

// Composer: Enter or SEND ▶ sends; Shift+Enter breaks lines.
func Composer(c *ui.Context, text string, value *string, o ComposerOptions) ComposerResult {
	k := pixel.K(c)
	lines := pixel.Or("Composer Lines", o.Lines, 2)
	var res ComposerResult
	res.Element = ui.Row(c).Gap(k.D(4)).AlignItems(ui.End).Children(func() {
		f := input.TextArea(c, text, value, input.TextAreaOptions{Placeholder: o.Placeholder, Width: o.Width, Lines: lines, Disabled: o.Disabled, EnterSends: true})
		empty := strings.TrimSpace(*value) == ""
		b := input.Button(c, k.Msg("send"), input.ButtonOptions{Disabled: o.Disabled || empty})
		res.Changed = f.Changed
		if f.Changed {
			slog.Debug("8bit: composer changed", "label", text, "runes", utf8.RuneCountInString(*value))
		}
		if (f.Submitted || b.Clicked) && !empty && !o.Disabled {
			res.Sent, *value = *value, ""
			// The cleared field shows next frame and keeps focus.
			f.Input.Focus()
			c.Invalidate()
			slog.Debug("8bit: composer sent", "label", text, "runes", utf8.RuneCountInString(res.Sent))
		}
	})
	return res
}
