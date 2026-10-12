package datetime

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

type DigitalClockOptions struct {
	// Scale is 1–4; zero means 2.
	Scale int
}

// DigitalClock shows HH:MM:SS in seven segments; colons blink.
func DigitalClock(c *ui.Context, o DigitalClockOptions) *ui.Element {
	k := pixel.K(c)
	sc := pixel.Or("DigitalClock Scale", o.Scale, 2)
	if sc > 4 {
		panic(fmt.Sprintf("datetime: DigitalClock Scale %d outside 1–4", sc))
	}
	k.Tick(c, 1)
	now := k.Now()
	colon := func() {
		ui.Box(c).Size(k.D(3*sc), k.D(11*sc)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			// The colon blinks each half second, unless reduced.
			if k.Reduce() || g.Step(2)%2 == 0 {
				g.Rect(sc, 3*sc, sc, sc, pixel.Bad)
				g.Rect(sc, 7*sc, sc, sc, pixel.Bad)
			}
		})
	}
	return ui.Row(c).Gap(k.D(2 * sc)).Role(ui.RoleText).Label(now.Format("15:04:05")).Children(func() {
		text.SevenSegment(c, now.Format("15"), text.SevenSegmentOptions{Scale: sc})
		colon()
		text.SevenSegment(c, now.Format("04"), text.SevenSegmentOptions{Scale: sc})
		colon()
		text.SevenSegment(c, now.Format("05"), text.SevenSegmentOptions{Scale: sc})
	})
}
