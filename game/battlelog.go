package game

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Line is a battle log line; Slot zero is Ink.
type Line struct {
	Text string
	Slot pixel.Slot
}

type BattleLogOptions struct {
	// Width and Height are dots; zeros mean 180 by 68.
	Width, Height int
}

// BattleLog keeps the newest lines last; old ones roll off.
func BattleLog(c *ui.Context, text string, lines []Line, o BattleLogOptions) *ui.Element {
	k := pixel.K(c)
	w, h := pixel.Or("BattleLog Width", o.Width, 180), pixel.Or("BattleLog Height", o.Height, 68)
	n := (h - 4) / pixel.LineHeight
	if n < 1 {
		panic("game: BattleLog is shorter than a line")
	}
	shown := lines[max(0, len(lines)-n):]
	e := ui.Box(c).Role(ui.RoleStatus).Label(text)
	if len(lines) > 0 {
		e.Description(lines[len(lines)-1].Text)
	}
	return e.Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, w, h, pixel.Sunken, pixel.Well)
		g.Clip(2, 2, w-4, h-4, func() {
			for i, l := range shown {
				g.Text(4, h-2-(len(shown)-i)*pixel.LineHeight, l.Text, l.Slot, 1)
			}
		})
	})
}
