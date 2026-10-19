// Package feedback tells the player what is going on.
package feedback

import (
	"strconv"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// pingpong runs n back and forth over 0…span.
func pingpong(n, span int) int {
	if span <= 0 {
		return 0
	}
	n %= 2 * span
	if n > span {
		return 2*span - n
	}
	return n
}

// level is Good, then Warn past warn, Bad past bad.
func level(v, warn, bad float64) pixel.Slot {
	switch {
	case bad != 0 && v > bad:
		return pixel.Bad
	case warn != 0 && v > warn:
		return pixel.Warn
	}
	return pixel.Good
}

// fmtNum prints a number without trailing zeros.
func fmtNum(v float64) string { return strconv.FormatFloat(v, 'g', 4, 64) }
