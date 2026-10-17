package account

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// PadButton is a control on the game pad.
type PadButton int

const (
	DPad PadButton = iota
	ButtonA
	ButtonB
	ButtonX
	ButtonY
	ShoulderL
	ShoulderR
	Start
	Select
)

// pad is each button's name, place, label side and row.
var pad = []struct {
	name           string
	x, y, side, ly int
}{
	DPad: {"DPAD", 45, 45, -1, 26}, ButtonA: {"A", 130, 45, 1, 46}, ButtonB: {"B", 120, 55, 1, 74}, ButtonX: {"X", 120, 35, 1, 18},
	ButtonY: {"Y", 110, 45, 1, 32}, ShoulderL: {"L", 40, 17, -1, 4}, ShoulderR: {"R", 120, 17, 1, 4},
	Start: {"START", 87, 50, 1, 60}, Select: {"SELECT", 73, 50, -1, 56},
}

type ControllerMapOptions struct{}

// ControllerMap draws a pad; labels say what each button does.
func ControllerMap(c *ui.Context, text string, labels map[PadButton]string, o ControllerMapOptions) ui.Element {
	k := pixel.K(c)
	var said []string
	for b := DPad; b <= Select; b++ {
		if l, ok := labels[b]; ok {
			said = append(said, pad[b].name+": "+l)
		}
	}
	for b := range labels {
		if b < DPad || b > Select {
			panic(fmt.Sprintf("account: ControllerMap button %d", b))
		}
	}
	// The pad sits mid-picture; labels stand out at the sides.
	const ox, w, h = 70, 310, 96
	return ui.Box(c).Size(k.D(w), k.D(h)).Role(ui.RoleImage).Label(text).Description(strings.Join(said, ", ")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for _, x := range []int{40, 120} {
			g.Rect(ox+x-15, 13, 30, 6, pixel.Shade)
			g.Disc(ox+x, 45, 27, pixel.Edge)
		}
		g.Rect(ox+40, 18, 80, 55, pixel.Edge)
		for _, x := range []int{40, 120} {
			g.Disc(ox+x, 45, 26, pixel.Face)
		}
		g.Rect(ox+40, 19, 80, 53, pixel.Face)
		// Lines first, so buttons sit on their ends.
		for b := DPad; b <= Select; b++ {
			if l, ok := labels[b]; ok {
				lx := ox - 4
				if pad[b].side > 0 {
					lx = ox + 164
				}
				g.Line(ox+pad[b].x, pad[b].y, lx, pad[b].ly+6, pixel.Faint)
				if pad[b].side < 0 {
					g.Text(ox-6-pixel.Measure(l), pad[b].ly, l, pixel.Ink, 1)
				} else {
					g.Text(ox+166, pad[b].ly, l, pixel.Ink, 1)
				}
			}
		}
		g.Rect(ox+38, 43, 15, 5, pixel.Ink)
		g.Rect(ox+43, 38, 5, 15, pixel.Ink)
		for b, s := range map[PadButton]pixel.Slot{ButtonA: pixel.Bad, ButtonB: pixel.Warn, ButtonX: pixel.Info, ButtonY: pixel.Good} {
			g.Disc(ox+pad[b].x, pad[b].y, 4, s)
		}
		for _, b := range []PadButton{Start, Select} {
			g.Rect(ox+pad[b].x-4, pad[b].y-1, 8, 3, pixel.Shade)
		}
	})
}
