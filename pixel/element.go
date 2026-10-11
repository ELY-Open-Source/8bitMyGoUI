package pixel

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
)

type TextOptions struct {
	Slot Slot
	// Scale is 1–4; zero means 1.
	Scale int
	// Wrap breaks lines at the width the layout gives.
	Wrap     bool
	MaxLines int
}

// Text is a pixel-font label that screen readers read.
func Text(c *ui.Context, s string, o TextOptions) *ui.Element {
	k := K(c)
	scale := scaleOr1(o.Scale)
	e := ui.Box(c).Label(s).Role(ui.RoleText)
	width, check := 0, func(ui.Rect) {}
	if o.Wrap {
		width, check = WrapWidth(c, e)
	}
	lines := Wrap(s, width/scale, o.MaxLines)
	if o.Wrap {
		e.Height(k.D(len(lines) * LineHeight * scale))
	} else {
		w := 0
		for _, l := range lines {
			w = max(w, Measure(l))
		}
		e.Size(k.D(w*scale), k.D(len(lines)*LineHeight*scale))
	}
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		check(r)
		g := k.Canvas(p, r)
		for i, l := range lines {
			g.Text(0, i*LineHeight*scale, l, o.Slot, scale)
		}
	})
}

// Wrap wraps s to width; past maxLines, an ellipsis.
func Wrap(s string, width, maxLines int) []string {
	if maxLines < 0 {
		panic(fmt.Sprintf("pixel: MaxLines %d", maxLines))
	}
	lines := Lines(s, width)
	if maxLines == 0 || len(lines) <= maxLines {
		return lines
	}
	lines = lines[:maxLines]
	last := []rune(strings.TrimRight(lines[maxLines-1], " "))
	for width > 0 && len(last) > 0 && Measure(string(last)+"…") > width {
		last = last[:len(last)-1]
	}
	lines[maxLines-1] = string(last) + "…"
	return lines
}

type wrapKey struct{}

// WrapWidth is e's last painted width; Draw must call check.
func WrapWidth(c *ui.Context, e *ui.Element) (int, func(r ui.Rect)) {
	dot := K(c).dot
	w := ui.Local(e, wrapKey{}, func() int { return 0 })
	used := *w
	return used, func(r ui.Rect) {
		if now := int(r.W / dot); now != *w {
			*w = now
			c.Invalidate()
		}
	}
}

type IconOptions struct {
	Slot  Slot
	Scale int
	// Label names the icon; empty hides it from assistive technology.
	Label string
}

func Icon(c *ui.Context, sp *Sprite, o IconOptions) *ui.Element {
	k := K(c)
	scale := scaleOr1(o.Scale)
	e := ui.Box(c).Size(k.D(sp.W*scale), k.D(sp.H*scale))
	if o.Label != "" {
		e.Label(o.Label).Role(ui.RoleImage)
	} else {
		e.Role(ui.RoleNone)
	}
	return e.Draw(func(p *ui.Painter, r ui.Rect) {
		k.Canvas(p, r).Sprite(0, 0, sp, 0, o.Slot, scale)
	})
}

func scaleOr1(s int) int {
	if s == 0 {
		return 1
	}
	checkScale(s)
	return s
}

// State is how a control looks this frame.
type State struct{ Hover, Press, Focus, Disabled bool }

// Shell wraps base in an enabled box that paints.
func Shell(c *ui.Context, base func() *ui.Element) (shell, b *ui.Element, st State) {
	shell = ui.Box(c)
	shell.Children(func() { b = base().Fill().FocusRing(false) })
	st = State{Hover: b.Hovered(), Press: b.Pressed(), Focus: b.FocusVisible(), Disabled: b.IsDisabled()}
	return shell, b, st
}

var hand = MustSprite(`
	.eee......
	elllleeee.
	elllllllle
	elllleeee.
	elllle....
	.eeee.....
`)

// FocusMark rings a box with dots and points a hand.
func (g Canvas) FocusMark(x, y, w, h int) {
	rx, ry, rw, rh := x-2, y-2, w+4, h+4
	for i := range rw {
		g.ant(rx+i, ry)
		g.ant(rx+i, ry+rh-1)
	}
	for j := 1; j < rh-1; j++ {
		g.ant(rx, ry+j)
		g.ant(rx+rw-1, ry+j)
	}
	g.Sprite(rx-hand.W-1+g.Step(4)%2, y+(h-hand.H)/2, hand, 0, Edge, 1)
}

func (g Canvas) ant(x, y int) {
	if (g.gx+x+g.gy+y)&1 == 0 {
		g.Dot(x, y, Accent)
	}
}
