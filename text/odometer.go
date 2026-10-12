package text

import (
	"fmt"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type OdometerOptions struct {
	// Digits is zero for as many as value has.
	Digits int
	Slot   pixel.Slot
}

type odoKey struct{}

type odoState struct {
	shown float64
	last  int
}

const odoCell = 10

// Odometer rolls its digits toward value, step by step.
func Odometer(c *ui.Context, value int, o OdometerOptions) *ui.Element {
	if value < 0 {
		panic(fmt.Sprintf("text: Odometer value %d is negative", value))
	}
	k := pixel.K(c)
	digits := pixel.Or("Odometer Digits", o.Digits, len(strconv.Itoa(value)))
	if value >= int(math.Pow10(digits)) {
		panic(fmt.Sprintf("text: Odometer value %d needs more than %d digits", value, digits))
	}
	h := pixel.LineHeight + 2
	e := ui.Box(c).Size(k.D(digits*odoCell), k.D(h)).Label(strconv.Itoa(value)).Role(ui.RoleText)
	st := ui.Local(e, odoKey{}, func() odoState { return odoState{shown: float64(value), last: -1} })
	target := float64(value)
	switch {
	case k.Reduce() || st.shown == target:
		st.shown, st.last = target, -1
	default:
		tick := k.Tick(c, 30)
		if st.last >= 0 {
			for range tick - st.last {
				st.shown = roll(st.shown, target)
			}
		}
		st.last = tick
	}
	shown := st.shown
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i := range digits {
			x := (digits - 1 - i) * odoCell
			d, frac := wheel(shown, i)
			g.Frame(x, 0, odoCell-1, h, pixel.Sunken, pixel.Well)
			g.Clip(x+2, 2, odoCell-5, h-4, func() {
				y := 1 - int(frac*float64(h))
				g.Text(x+2, y, strconv.Itoa(d), o.Slot, 1)
				g.Text(x+2, y+h, strconv.Itoa((d+1)%10), o.Slot, 1)
			})
		}
	})
}

// roll moves shown a step toward target, faster when far.
func roll(shown, target float64) float64 {
	diff := target - shown
	stride := max(0.25, math.Abs(diff)/6)
	if math.Abs(diff) <= stride {
		return target
	}
	return shown + math.Copysign(stride, diff)
}

// wheel is digit i of shown and its roll.
func wheel(shown float64, i int) (d int, frac float64) {
	place := math.Pow10(i)
	v := shown / place
	d = int(math.Floor(v)) % 10
	if lower := math.Mod(shown, place); lower > place-1 {
		frac = lower - (place - 1)
	}
	return d, frac
}
