package code

import (
	"log/slog"
	"unicode/utf8"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TerminalOptions struct {
	// Prompt starts the input line; zero means "$".
	Prompt string
	// Width and Height are dots; zeros mean 240 by 120.
	Width, Height int
	Disabled      bool
}

type TerminalResult struct {
	Element *ui.Element
	// Ran is the command Enter sent, else "".
	Ran     string
	Changed bool
}

type termKey struct{}

// term is the scroll, history, how far back, the draft.
type term struct {
	scroll  ui.ScrollState
	history []string
	back    int
	draft   string
}

// Terminal prints lines over a prompt; Enter runs, arrows recall.
func Terminal(c *ui.Context, text string, lines []string, input *string, o TerminalOptions) TerminalResult {
	k := pixel.K(c)
	w, h := pixel.Or("Terminal Width", o.Width, 240), pixel.Or("Terminal Height", o.Height, 120)
	prompt := o.Prompt
	if prompt == "" {
		prompt = "$"
	}
	var res TerminalResult
	res.Element = layout.Panel(c, layout.PanelOptions{Style: pixel.Sunken, Fill: pixel.Well, Padding: 2}, func() {
		box := ui.Column(c).Size(k.D(w), k.D(h)).Gap(k.D(2))
		tm := ui.Local(box, termKey{}, func() term { return term{} })
		box.Children(func() {
			layout.ScrollArea(c, &tm.scroll, layout.ScrollAreaOptions{Label: text, Follow: true, Disabled: o.Disabled}, func() {
				ui.Column(c).FillWidth().Children(func() {
					for _, l := range lines {
						pixel.Text(c, l, pixel.TextOptions{Wrap: true})
					}
				})
			}).Element.Grow(1).FillWidth()
			row, off := ui.Row(c).FillWidth().Gap(k.D(4)), false
			row.Children(func() {
				pixel.Text(c, prompt, pixel.TextOptions{Slot: pixel.Good})
				ed := pixel.Editor(c, input, pixel.EditorOptions{Label: k.Msg("command"), Disabled: o.Disabled})
				ed.Element.Grow(1)
				res.Changed, off = ed.Changed, ed.Input.IsDisabled()
				if ed.Submitted && !off {
					// The editor drew the old line; draw the empty one.
					c.Invalidate()
					res.Ran, *input, tm.back = *input, "", 0
					if n := len(tm.history); res.Ran != "" && (n == 0 || tm.history[n-1] != res.Ran) {
						tm.history = append(tm.history, res.Ran)
					}
				}
			})
			// Up and Down are history here, taken from the editor.
			if up, down := pixel.Shortcut(row, 0, ui.KeyUp), pixel.Shortcut(row, 0, ui.KeyDown); (up || down) && !off {
				res.Changed = recall(tm, input, up) || res.Changed
			}
		})
	})
	if res.Ran != "" || res.Changed {
		slog.Debug("8bit: terminal", "label", text, "ran", res.Ran != "", "runes", utf8.RuneCountInString(*input))
	}
	return res
}

// recall steps through history; past the newest is the draft.
func recall(tm *term, input *string, up bool) bool {
	n := len(tm.history)
	switch {
	case up && tm.back < n:
		if tm.back == 0 {
			tm.draft = *input
		}
		tm.back++
	case !up && tm.back > 0:
		tm.back--
	default:
		return false
	}
	*input = tm.draft
	if tm.back > 0 {
		*input = tm.history[n-tm.back]
	}
	return true
}
