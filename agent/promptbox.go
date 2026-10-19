package agent

import (
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type PromptBoxOptions struct {
	Placeholder string
	// Width is in dots, zero 240; Lines zero means 3.
	Width, Lines int
	Disabled     bool
}

type PromptBoxResult struct {
	Element ui.Element
	Changed bool
	// Cast is the prompt cast this frame, else empty.
	Cast string
}

// PromptBox: Enter or CAST casts the prompt; Shift+Enter breaks lines.
func PromptBox(c *ui.Context, text string, value *string, o PromptBoxOptions) PromptBoxResult {
	k := pixel.K(c)
	w, lines := pixel.Or("PromptBox Width", o.Width, 240), pixel.Or("PromptBox Lines", o.Lines, 3)
	var res PromptBoxResult
	res.Element = ui.Column(c).Gap(k.D(4)).AlignItems(ui.End).Children(func() {
		f := input.TextArea(c, text, value, input.TextAreaOptions{Placeholder: o.Placeholder, Width: w, Lines: lines, Disabled: o.Disabled, EnterSends: true})
		empty := strings.TrimSpace(*value) == ""
		b := input.Button(c, k.Msg("cast"), input.ButtonOptions{Variant: input.Primary, Icon: icon.Sparkle, Disabled: o.Disabled || empty})
		res.Changed = f.Changed
		// Off or veiled, a key sent this frame casts nothing.
		if (f.Submitted || b.Clicked) && !empty && !f.Input.IsDisabled() {
			res.Cast, *value = *value, ""
			// The cleared field shows next frame and keeps focus.
			f.Input.Focus()
			c.Invalidate()
			slog.Debug("8bit: prompt cast", "label", text, "runes", utf8.RuneCountInString(res.Cast))
		}
	})
	return res
}
