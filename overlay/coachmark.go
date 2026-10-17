package overlay

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

type CoachmarkOptions struct {
	Disabled bool
}

type CoachmarkResult struct {
	Element ui.Element
	Next    bool
}

// coachWidth is the bubble's text width in dots.
const coachWidth = 140

// Coachmark points a bubble at target: step count, NEXT ▶.
func Coachmark(c *ui.Context, target ui.Element, message string, step, steps int, o CoachmarkOptions) CoachmarkResult {
	k := pixel.K(c)
	if step < 1 || step > steps {
		panic(fmt.Sprintf("overlay: Coachmark step %d of %d", step, steps))
	}
	var res CoachmarkResult
	res.Element = bubble(c, target.Bounds(), func() {
		ui.Column(c).Width(k.D(coachWidth)).Gap(k.D(4)).Role(ui.RoleGroup).Label(message).Children(func() {
			text.Paragraph(c, message, text.ParagraphOptions{})
			ui.Row(c).AlignItems(ui.Center).Children(func() {
				text.MicroText(c, fmt.Sprintf("%d/%d", step, steps), text.MicroTextOptions{Slot: pixel.Faint})
				ui.Box(c).Grow(1)
				res.Next = input.Button(c, k.Msg("next_step"), input.ButtonOptions{Size: input.Small, Variant: input.Primary, Disabled: o.Disabled}).Clicked
			})
		})
	}, coachWidth+16)
	if res.Next {
		slog.Debug("8bit: coachmark next", "step", step, "steps", steps)
	}
	return res
}
