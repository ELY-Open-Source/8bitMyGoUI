package overlay

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

type ChoicePromptOptions struct {
	Disabled bool
}

type ChoicePromptResult struct {
	Element ui.Element
	// Chosen is the choice taken this frame, or -1.
	Chosen int
	// Canceled is Escape or a click outside.
	Canceled bool
}

type choiceKey struct{}

// ChoicePrompt asks a question; a hand picks, or keys 1–9.
func ChoicePrompt(c *ui.Context, question string, choices []string, open *bool, o ChoicePromptOptions) ChoicePromptResult {
	k := pixel.K(c)
	if len(choices) == 0 {
		panic("overlay: ChoicePrompt has no choices")
	}
	home := ui.Box(c)
	cur := ui.Local(home, choiceKey{}, func() int { return 0 })
	was := *open
	res := ChoicePromptResult{Chosen: -1}
	res.Element = ui.DialogBase(c, open, func(back, panel ui.Element) {
		dim(k, back, pixel.Back, 4)
		panel.Width(k.D(180)).Label(question).Focusable().FocusRing(false).Children(func() {
			pixel.Harbor(c, panel, func() {
				layout.Window(c, layout.WindowOptions{}, func() {
					text.Paragraph(c, question, text.ParagraphOptions{})
					ui.Box(c).Height(k.D(6))
					if i := choose(c, question, choices, cur, true, o.Disabled); i >= 0 {
						res.Chosen, *open = i, false
					}
				})
			})
		})
	})
	res.Canceled = was && !*open && res.Chosen < 0
	if !*open {
		*cur = 0
	}
	if res.Chosen >= 0 || res.Canceled {
		slog.Debug("8bit: choice", "question", question, "chosen", res.Chosen)
	}
	return res
}
