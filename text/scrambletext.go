package text

import (
	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

const scrambleRunes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789#$%&*+?"

type ScrambleTextOptions struct {
	Slot  pixel.Slot
	Scale int
	// FPS is 4–15, zero for 15; two steps settle one.
	FPS int
}

type scrambleKey struct{}

type scrambleState struct {
	text        string
	start, step int
}

// ScrambleText settles from random glyphs to s, left to right.
func ScrambleText(c *ui.Context, s string, o ScrambleTextOptions) ui.Element {
	k, sc := pixel.K(c), scale(o.Scale)
	fps := pixel.Or("ScrambleText FPS", o.FPS, 15)
	rs := []rune(s)
	e := ui.Box(c).Size(k.D(pixel.Measure(s)*sc), k.D(pixel.LineHeight*sc)).Label(s).Role(ui.RoleText)
	st := ui.Local(e, scrambleKey{}, func() scrambleState { return scrambleState{start: -1} })
	if st.text != s {
		*st = scrambleState{text: s, start: -1}
	}
	settled := len(rs)
	if !k.Reduce() && (st.start < 0 || (st.step-st.start)/2 < len(rs)) {
		st.step = k.Tick(c, fps)
		if st.start < 0 {
			st.start = st.step
		}
		settled = min(len(rs), (st.step-st.start)/2)
	}
	step := st.step
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		each(s, sc, func(i int, ch rune, x int) {
			if i >= settled && ch != ' ' {
				ch = rune(scrambleRunes[noise(len(rs), i, step)%uint32(len(scrambleRunes))])
			}
			g.Text(x, 0, string(ch), o.Slot, sc)
		})
	})
}
