// Package account holds settings, sign-in and player screens.
package account

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// players are the arcade colors of P1 to P4.
var players = [4]pixel.Slot{pixel.Bad, pixel.Info, pixel.Good, pixel.Warn}

// choice panics unless *i points into n choices.
func choice(name string, i *int, n int) {
	if n == 0 || *i < 0 || *i >= n {
		panic(fmt.Sprintf("account: %s choice %d of %d", name, *i, n))
	}
}

// words is a message split on spaces, as option names.
func words(k pixel.Kit, key string) []string { return strings.Fields(k.Msg(key)) }

// caption is a line of micro text named for readers.
func caption(c *ui.Context, k pixel.Kit, s string, slot pixel.Slot) ui.Element {
	return ui.Box(c).Size(k.D(pixel.MicroWidth(s)), k.D(5)).Role(ui.RoleText).Label(s).Draw(func(p *ui.Painter, r ui.Rect) {
		k.Canvas(p, r).Micro(0, 0, s, slot, 1)
	})
}

// mark draws s unnamed, beside the control that names it.
func mark(c *ui.Context, k pixel.Kit, s string, slot pixel.Slot, w int) ui.Element {
	return ui.Box(c).Size(k.D(max(w, pixel.Measure(s))), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		k.Canvas(p, r).Text(0, 0, s, slot, 1)
	})
}
