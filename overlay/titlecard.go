package overlay

import (
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type TitleCardOptions struct {
	// Cue shows the card each time it changes.
	Cue int
	// Hold keeps the title up; zero means two seconds.
	Hold time.Duration
}

// TitleCard dithers in, holds title, dithers out; fills its parent.
func TitleCard(c *ui.Context, title string, o TitleCardOptions) ui.Element {
	k := pixel.K(c)
	if o.Hold < 0 {
		panic("overlay: TitleCard Hold is negative")
	}
	hold := o.Hold
	if hold == 0 {
		hold = 2 * time.Second
	}
	e := ui.Box(c).Fill()
	// Eight steps a second: four in, the hold, four out.
	held := int(hold * 8 / time.Second)
	step := k.Since(c, k.Cue(e, o.Cue), 8, held+8)
	if step >= held+8 {
		return e
	}
	level := 16
	switch {
	case k.Reduce():
	case step < 4:
		level = 4 * (step + 1)
	case step >= held+4:
		level = 4 * (held + 7 - step)
	}
	deep, light := dark(k)
	return e.Role(ui.RoleStatus).Label(title).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Veil(0, 0, g.W, g.H, deep, level)
		if level == 16 {
			g.Text((g.W-pixel.Measure(title)*2)/2, (g.H-pixel.LineHeight*2)/2, title, light, 2)
		}
	})
}
