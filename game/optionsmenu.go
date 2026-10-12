package game

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Setting is an options row: a choice, or a level.
type Setting struct {
	Name    string
	Choices []string
	// Max makes a level row, 0 to Max, not Choices.
	Max   int
	Value *int
}

// top is the setting's highest value.
func (s Setting) top() int {
	if s.Max > 0 {
		return s.Max
	}
	return len(s.Choices) - 1
}

// shown is the setting's value as words.
func (s Setting) shown() string {
	if s.Max > 0 {
		return fmt.Sprintf("%d/%d", *s.Value, s.Max)
	}
	return s.Choices[*s.Value]
}

type OptionsMenuOptions struct {
	Disabled bool
}

type OptionsMenuResult struct {
	Element *ui.Element
	Changed bool
}

// rowH is an options row's height in dots.
const rowH = pixel.LineHeight + 2

// OptionsMenu rows settings: Up and Down choose, ◀ ▶ change.
func OptionsMenu(c *ui.Context, text string, settings []Setting, row *int, o OptionsMenuOptions) OptionsMenuResult {
	k := pixel.K(c)
	if len(settings) == 0 || *row < 0 || *row >= len(settings) {
		panic(fmt.Sprintf("game: OptionsMenu row %d of %d", *row, len(settings)))
	}
	name, value := 0, 0
	for _, s := range settings {
		if (len(s.Choices) > 0) == (s.Max > 0) || s.Value == nil || *s.Value < 0 || *s.Value > s.top() {
			panic(fmt.Sprintf("game: setting %q needs Choices or Max, and a value in range", s.Name))
		}
		name = max(name, pixel.Measure(s.Name))
		value = max(value, 4*s.Max)
		for _, ch := range s.Choices {
			value = max(value, pixel.Measure(ch))
		}
	}
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleList).Label(text).Disabled(o.Disabled)
	})
	res, was := OptionsMenuResult{Element: shell}, *row
	step := func(d int) {
		s := settings[*row]
		if v := *s.Value + d; v >= 0 && v <= s.top() {
			*s.Value, res.Changed = v, true
		}
	}
	if base.Shortcut(0, ui.KeyUp) && *row > 0 {
		*row--
	}
	if base.Shortcut(0, ui.KeyDown) && *row < len(settings)-1 {
		*row++
	}
	if base.Shortcut(0, ui.KeyLeft) {
		step(-1)
	}
	if base.Shortcut(0, ui.KeyRight) {
		step(1)
	}
	// Enter and Space are taken, so clicks come from pointers.
	base.Shortcut(0, ui.KeyEnter)
	base.Shortcut(0, ui.KeySpace)
	aw := pixel.Measure("◀")
	ax := pixel.Hand.W + 4 + name + 10
	vx := ax + aw + 3
	bx := vx + value + 3
	// A click picks its row; an arrow click steps too.
	if base.Clicked() {
		x, y, _ := base.PointerPosition()
		dx, dy := int(x/k.Dot()), int(y/k.Dot())
		if i := dy / rowH; dy >= 0 && i < len(settings) {
			*row = i
			switch {
			case dx >= ax && dx < ax+aw:
				step(-1)
			case dx >= bx && dx < bx+aw:
				step(1)
			}
		}
	}
	base.Description(settings[*row].Name + ": " + settings[*row].shown())
	if res.Changed || *row != was {
		slog.Debug("8bit: option", "label", text, "setting", settings[*row].Name, "value", settings[*row].shown())
	}
	cur := *row
	values := make([]int, len(settings))
	for i, s := range settings {
		values[i] = *s.Value
	}
	shell.Size(k.D(bx+aw), k.D(len(settings)*rowH)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, s := range settings {
			y, v := i*rowH, values[i]
			if i == cur {
				g.Sprite(0, y+(rowH-pixel.Hand.H)/2, pixel.Hand, 0, pixel.Edge, 1)
			}
			g.Text(pixel.Hand.W+4, y, s.Name, pixel.Ink, 1)
			arrow := func(x int, glyph string, on bool) {
				ink := pixel.Faint
				if on {
					ink = pixel.Accent
				}
				g.Text(x, y, glyph, ink, 1)
			}
			arrow(ax, "◀", v > 0)
			arrow(bx, "▶", v < s.top())
			if s.Max > 0 {
				for j := range s.Max {
					if j < v {
						g.Rect(vx+j*4, y+5, 3, 7, pixel.Good)
					} else {
						g.Checker(vx+j*4, y+5, 3, 7, pixel.Shade)
					}
				}
				continue
			}
			ch := s.Choices[v]
			g.Text(vx+(value-pixel.Measure(ch))/2, y, ch, pixel.Ink, 1)
		}
		g.Finish(st)
	})
	return res
}
