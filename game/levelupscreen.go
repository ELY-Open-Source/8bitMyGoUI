package game

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

// Attribute is a stat raised at a level up.
type Attribute struct {
	Name string
	Base int
	// Added is the points put in it this time.
	Added int
}

type LevelUpScreenOptions struct {
	Disabled bool
}

type LevelUpScreenResult struct {
	Element *ui.Element
	Changed bool
	Done    bool
}

// back is ◀, the triangle icon turned around.
var back = icon.TriangleRight.Mirror()

// LevelUpScreen spends points on attributes with ◀ ▶, then DONE.
func LevelUpScreen(c *ui.Context, title string, attrs []Attribute, points *int, o LevelUpScreenOptions) LevelUpScreenResult {
	k := pixel.K(c)
	if *points < 0 {
		panic(fmt.Sprintf("game: LevelUpScreen has %d points", *points))
	}
	name := 0
	for _, a := range attrs {
		if a.Base < 0 || a.Added < 0 {
			panic(fmt.Sprintf("game: attribute %q is %d%+d", a.Name, a.Base, a.Added))
		}
		name = max(name, pixel.Measure(a.Name))
	}
	col := ui.Column(c).Gap(k.D(4)).Role(ui.RoleGroup).Label(title)
	res := LevelUpScreenResult{Element: col}
	col.Children(func() {
		ui.Row(c).Gap(k.D(10)).AlignItems(ui.End).Children(func() {
			text.ShadowText(c, k.Msg("level_up"), text.ShadowTextOptions{Scale: 2, Slot: pixel.Warn})
			pixel.Text(c, fmt.Sprintf(k.Msg("points"), *points), pixel.TextOptions{Slot: pixel.Accent})
		})
		for i := range attrs {
			a := &attrs[i]
			ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Width(k.D(name + 6)).Children(func() { pixel.Text(c, a.Name, pixel.TextOptions{}) })
				if input.IconButton(c, back, a.Name+" "+k.Msg("down"), input.IconButtonOptions{Size: input.Small, Disabled: o.Disabled || a.Added == 0}).Clicked {
					a.Added, *points, res.Changed = a.Added-1, *points+1, true
				}
				ui.Box(c).Width(k.D(40)).Children(func() {
					ui.Row(c).Children(func() {
						pixel.Text(c, strconv.Itoa(a.Base+a.Added), pixel.TextOptions{})
						if a.Added > 0 {
							pixel.Text(c, "+"+strconv.Itoa(a.Added), pixel.TextOptions{Slot: pixel.Good})
						}
					})
				})
				if input.IconButton(c, icon.TriangleRight, a.Name+" "+k.Msg("up"), input.IconButtonOptions{Size: input.Small, Disabled: o.Disabled || *points == 0}).Clicked {
					a.Added, *points, res.Changed = a.Added+1, *points-1, true
				}
			})
		}
		res.Done = input.Button(c, k.Msg("done"), input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled}).Clicked
	})
	if res.Changed {
		slog.Debug("8bit: level up", "label", title, "points", *points)
	}
	return res
}
