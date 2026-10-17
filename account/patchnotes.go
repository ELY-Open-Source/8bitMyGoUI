package account

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

// Patch is a version's changes, by kind.
type Patch struct {
	Version               string
	Added, Fixed, Changed []string
}

type PatchNotesOptions struct {
	// Width is in dots; zero means 240.
	Width int
}

// PatchNotes lists versions, each line tagged ADDED, FIXED or CHANGED.
func PatchNotes(c *ui.Context, title string, patches []Patch, o PatchNotesOptions) ui.Element {
	k := pixel.K(c)
	w := pixel.Or("PatchNotes Width", o.Width, 240)
	// Tags draw in the main font, so any locale reads.
	tw := 0
	for _, key := range []string{"added", "fixed", "changed"} {
		tw = max(tw, pixel.Measure(k.Msg(key))+6)
	}
	return ui.Column(c).Width(k.D(w)).Gap(k.D(3)).Role(ui.RoleGroup).Label(title).Children(func() {
		for i, pt := range patches {
			if i > 0 {
				ui.Box(c).Height(k.D(4))
			}
			text.ShadowText(c, pt.Version, text.ShadowTextOptions{Slot: pixel.Accent})
			for _, kind := range []struct {
				key   string
				slot  pixel.Slot
				lines []string
			}{{"added", pixel.Good, pt.Added}, {"fixed", pixel.Info, pt.Fixed}, {"changed", pixel.Warn, pt.Changed}} {
				tag := k.Msg(kind.key)
				for _, l := range kind.lines {
					ui.Row(c).FillWidth().Gap(k.D(4)).Role(ui.RoleText).Label(tag + ": " + l).Children(func() {
						ui.Box(c).Size(k.D(tw), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
							g := k.Canvas(p, r)
							g.Rect(0, 0, pixel.Measure(tag)+6, pixel.LineHeight, kind.slot)
							g.Text(3, 0, tag, pixel.AccentInk, 1)
						})
						pixel.Text(c, l, pixel.TextOptions{Wrap: true}).Grow(1)
					})
				}
			}
		}
	})
}
