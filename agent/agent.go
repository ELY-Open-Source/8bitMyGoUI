// Package agent dresses AI agent work as an RPG party.
package agent

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Status is where an agent stands.
type Status int

const (
	Idle Status = iota
	Working
	Blocked
	Done
)

// look is a status's word and slot; unknown ones panic.
func (s Status) look(k pixel.Kit) (string, pixel.Slot) {
	if s < Idle || s > Done {
		panic(fmt.Sprintf("agent: unknown Status %d", s))
	}
	return k.Msg([...]string{"idle", "working", "blocked", "done"}[s]), [...]pixel.Slot{pixel.Faint, pixel.Info, pixel.Warn, pixel.Good}[s]
}

// tokens prints n short: 950, 1.5K, 128K, 2M.
func tokens(n int) string {
	if n < 1000 {
		return strconv.Itoa(n)
	}
	div, unit := 1000, 0
	for n/div >= 1000 {
		div, unit = div*1000, unit+1
	}
	s := strconv.Itoa(n / div)
	if t := n % div / (div / 10); n/div < 10 && t > 0 {
		s += "." + strconv.Itoa(t)
	}
	return s + "KMBTPE"[unit:unit+1]
}

// gauge frames a bar eight dots high and fills part.
func gauge(g pixel.Canvas, x, y, w int, part float64, ink pixel.Slot) {
	g.Frame(x, y, w, 8, pixel.Sunken, pixel.Well)
	g.Rect(x+2, y+2, int(part*float64(w-4)), 4, ink)
}

// health: Good past half, Warn past a quarter, else Bad.
func health(left, full int) pixel.Slot {
	switch {
	case left <= full/4:
		return pixel.Bad
	case left <= full/2:
		return pixel.Warn
	}
	return pixel.Good
}

// portrait is p or the robot; over 12 dots panics.
func portrait(name string, p *pixel.Sprite) *pixel.Sprite {
	if p == nil {
		return icon.Robot
	}
	if p.W > 12 || p.H > 12 {
		panic(fmt.Sprintf("agent: %s portrait is %dx%d, over 12", name, p.W, p.H))
	}
	return p
}
