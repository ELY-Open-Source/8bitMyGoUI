package text

import (
	"fmt"
	"math"
	"strconv"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
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
func Odometer(c *ui.Context, value int, o OdometerOptions) ui.Element {
	if value < 0 {
		panic(fmt.Sprintf("text: Odometer value %d is negative", value))
	}
	k := pixel.K(c)
	// Counted as text, so the largest int needs no float.
	digits := pixel.Or("Odometer Digits", o.Digits, len(strconv.Itoa(value)))
	if len(strconv.Itoa(value)) > digits {
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
		// S5 caps steps at 15 fps; two rolls a step.
		tick := k.Tick(c, 15)
		if st.last >= 0 {
			for range 2 * (tick - st.last) {
				st.shown = roll(st.shown, target)
			}
		}
		st.last = tick
	}
	// At rest digits come from the int; floats would round.
	shown, settled, exact := st.shown, st.last < 0, fmt.Sprintf("%0*d", digits, value)
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i := range digits {
			x := (digits - 1 - i) * odoCell
			d, frac := int(exact[digits-1-i]-'0'), 0.0
			if !settled {
				d, frac = wheel(shown, i)
			}
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
	// Near 2^53 a step may not move; then land.
	next := shown + math.Copysign(stride, diff)
	if math.Abs(diff) <= stride || next == shown {
		return target
	}
	return next
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
