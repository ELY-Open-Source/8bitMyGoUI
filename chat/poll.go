package chat

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// PollOption is a choice and its votes.
type PollOption struct {
	Text  string
	Votes int
}

type PollOptions struct {
	// Width is in dots; zero means 180.
	Width    int
	Disabled bool
}

type PollResult struct {
	Element *ui.Element
	// Voted is the option voted for this frame, else -1.
	Voted int
}

// Poll shows vote bars; a click moves your vote there.
func Poll(c *ui.Context, question string, opts []PollOption, mine *int, o PollOptions) PollResult {
	k := pixel.K(c)
	w := pixel.Or("Poll Width", o.Width, 180)
	if len(opts) < 2 || *mine < -1 || *mine >= len(opts) {
		panic(fmt.Sprintf("chat: Poll %q has %d options and vote %d", question, len(opts), *mine))
	}
	total := 0
	for i, op := range opts {
		if op.Votes < 0 || i == *mine && op.Votes == 0 {
			panic(fmt.Sprintf("chat: Poll option %q has %d votes", op.Text, op.Votes))
		}
		total += op.Votes
	}
	res := PollResult{Voted: -1}
	res.Element = ui.Column(c).Width(k.D(w)).Gap(k.D(3)).Children(func() {
		para(c, k, question, w, 0, pixel.Ink)
		ui.Column(c).Gap(k.D(2)).FocusGroup(ui.Vertical).Role(ui.RoleGroup).Label(question).Children(func() {
			for i := range opts {
				if bar(c, k, opts[i], total, i == *mine, w, o.Disabled) {
					if *mine >= 0 {
						opts[*mine].Votes--
					}
					opts[i].Votes++
					*mine, res.Voted = i, i
					slog.Debug("8bit: poll vote", "poll", question, "option", opts[i].Text)
				}
			}
		})
		pixel.Text(c, fmt.Sprintf(k.Msg("votes"), total), pixel.TextOptions{Slot: pixel.Faint})
	})
	return res
}

// bar fills an option's share; ▶ marks your vote.
func bar(c *ui.Context, k pixel.Kit, op PollOption, total int, mine bool, w int, disabled bool) bool {
	label := fmt.Sprintf("%s, %d", op.Text, op.Votes)
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.ButtonBase(c).Label(label).Disabled(disabled).FocusRing(false)
	})
	n := strconv.Itoa(op.Votes)
	shell.Size(k.D(w), k.D(pixel.LineHeight+4)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, w, g.H, pixel.Sunken, pixel.Well)
		fill := 0
		if total > 0 {
			fill = (w - 4) * op.Votes / total
		}
		g.Rect(2, 2, fill, g.H-4, pixel.Accent)
		// Over the bar the words take the bar's own ink.
		words := func(ink, mark pixel.Slot) {
			if mine {
				g.Text(3, 2, "▶", mark, 1)
			}
			g.Text(14, 2, line(op.Text, w-30), ink, 1)
			g.Micro(w-4-pixel.MicroWidth(n), 7, n, ink, 1)
		}
		g.Clip(2, 0, fill, g.H, func() { words(pixel.AccentInk, pixel.AccentInk) })
		g.Clip(2+fill, 0, w-4-fill, g.H, func() { words(pixel.Ink, pixel.Accent) })
		g.Finish(st)
	})
	return base.Clicked() && !mine
}
