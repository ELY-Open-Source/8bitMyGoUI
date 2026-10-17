package layout

import (
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SectionOptions struct {
	// Gap is between children; zero means 4.
	Gap int
}

// Section titles a group inside a dotted rule.
func Section(c *ui.Context, title string, o SectionOptions, children func()) ui.Element {
	k := pixel.K(c)
	up := strings.ToUpper(title)
	return ui.Column(c).Gap(k.D(pixel.Or("Section Gap", o.Gap, 4))).Role(ui.RoleGroup).Label(title).Children(func() {
		ui.Box(c).FillWidth().Height(k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			mid, tw := pixel.LineHeight/2, pixel.Measure(up)
			g.Checker(0, mid, 6, 1, pixel.Faint)
			g.Text(10, 0, up, pixel.Ink, 1)
			g.Checker(14+tw, mid, g.W-14-tw, 1, pixel.Faint)
		})
		children()
	})
}
