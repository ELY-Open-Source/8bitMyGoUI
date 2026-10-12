// Package project holds task, planning and team components.
package project

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Level is how urgent a task is.
type Level int

const (
	Low Level = iota
	Medium
	High
	Critical
)

var levels = [...]string{"low", "medium", "high", "critical"}

// word is a level in the locale; unknown ones panic.
func (l Level) word(k pixel.Kit) string {
	if l < Low || l > Critical {
		panic(fmt.Sprintf("project: unknown level %d", l))
	}
	return k.Msg(levels[l])
}

// Label is a named tag in a palette slot.
type Label struct {
	Name string
	Slot pixel.Slot
}

// month is a month's short name in the frame's locale.
func month(k pixel.Kit, m time.Month) string {
	return strings.Fields(k.Msg("months"))[m-1]
}

// hms prints d as H:MM:SS, or MM:SS under an hour.
func hms(d time.Duration) string {
	d = d.Truncate(time.Second)
	h, m, s := int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

// show draws s unnamed, beside an element that names it.
func show(c *ui.Context, k pixel.Kit, s string, slot pixel.Slot, w int) *ui.Element {
	return ui.Box(c).Size(k.D(w), k.D(pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
		k.Canvas(p, r).Text(0, 0, pixel.Wrap(s, w, 1)[0], slot, 1)
	})
}
