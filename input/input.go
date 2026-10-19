// Package input takes the player's choices: buttons, boxes, fields.
package input

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type Variant uint8

const (
	Normal Variant = iota
	Primary
	Ghost
	Danger
)

type Size uint8

const (
	Medium Size = iota
	Small
	Large
)

// metrics are a size's height, padding, text top and scale.
func (s Size) metrics() (h, pad, top, scale int) {
	switch s {
	case Small:
		return pixel.LineHeight, 4, 0, 1
	case Medium:
		return pixel.LineHeight + 4, 8, 2, 1
	case Large:
		return 2*pixel.LineHeight + 4, 12, 2, 2
	}
	panic(fmt.Sprintf("input: unknown Size %d", s))
}

// face paints a button body, returning ink and press shift.
func face(g pixel.Canvas, v Variant, st pixel.State, on bool) (ink pixel.Slot, dy int) {
	fill, ink := pixel.Face, pixel.Ink
	switch v {
	case Primary:
		fill, ink = pixel.Accent, pixel.AccentInk
	case Danger:
		fill, ink = pixel.Bad, pixel.AccentInk
	case Normal, Ghost:
	default:
		panic(fmt.Sprintf("input: unknown Variant %d", v))
	}
	if st.Press || on {
		dy = 1
	}
	if g.W < 3 || g.H < 3 {
		return ink, dy
	}
	if v == Ghost && !st.Hover && !st.Press && !on {
		return ink, dy
	}
	if v == Ghost {
		fill = pixel.Lit
	}
	if on {
		fill = pixel.Shade
	}
	g.Frame(0, dy, g.W, g.H-dy, pixel.Notched, fill)
	switch {
	case on || st.Press:
		g.HLine(1, dy+1, g.W-2, pixel.Edge)
	default:
		g.Bevel(1, 1, g.W-2, g.H-2, pixel.Lit, pixel.Shade)
		if st.Hover {
			g.Veil(2, 2, g.W-4, g.H-4, pixel.Lit, 4)
		}
	}
	return ink, dy
}

// label is s, or the fallback when s is empty.
func label(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

// fmtValue prints a number without trailing zeros.
func fmtValue(v float64) string { return strconv.FormatFloat(v, 'g', 4, 64) }

var (
	tinyUp = pixel.MustSprite(`
		..#..
		.###.
		#####
	`)
	tinyDown = pixel.MustSprite(`
		#####
		.###.
		..#..
	`)
	tinyLeft = pixel.MustSprite(`
		..#
		.##
		###
		.##
		..#
	`)
	tinyRight = pixel.MustSprite(`
		#..
		##.
		###
		##.
		#..
	`)
)
