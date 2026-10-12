package datetime

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SeasonWheelOptions struct {
	Disabled bool
}

type SeasonWheelResult struct {
	Element *ui.Element
	Changed bool
}

// wheel is the season wheel's width in dots.
const wheel = 64

// seasonArt draws each season, clockwise from the top right.
var seasonArt = [4]*pixel.Sprite{icon.Flower, icon.Sun, icon.Leaf, icon.Snowflake}

// SeasonWheel chooses one of four seasons on a wheel.
func SeasonWheel(c *ui.Context, text string, season *pixel.Season, o SeasonWheelOptions) SeasonWheelResult {
	k := pixel.K(c)
	if *season > pixel.Winter {
		panic(fmt.Sprintf("datetime: unknown Season %d", *season))
	}
	old := *season
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleRadioGroup).Label(text).Disabled(o.Disabled)
	})
	if pixel.Shortcut(base, 0, ui.KeyRight) || pixel.Shortcut(base, 0, ui.KeyDown) {
		*season = (*season + 1) % 4
	}
	if pixel.Shortcut(base, 0, ui.KeyLeft) || pixel.Shortcut(base, 0, ui.KeyUp) {
		*season = (*season + 3) % 4
	}
	// Enter and Space are taken, so clicks come from pointers.
	pixel.Shortcut(base, 0, ui.KeyEnter)
	pixel.Shortcut(base, 0, ui.KeySpace)
	if base.Clicked() {
		x, y, _ := base.PointerPosition()
		*season = pixel.Season(quarter(int(x/k.Dot())-wheel/2, int(y/k.Dot())-wheel/2))
	}
	names := []string{k.Msg("spring"), k.Msg("summer"), k.Msg("autumn"), k.Msg("winter")}
	base.Description(names[*season])
	at := int(*season)
	shell.Size(k.D(wheel), k.D(wheel)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		mid, rad := wheel/2, wheel/2-1
		g.Disc(mid, mid, rad, pixel.Face)
		// The chosen quarter lights dot by dot inside the rim.
		for y := range wheel {
			for x := range wheel {
				dx, dy := x-mid, y-mid
				if quarter(dx, dy) == at && dx*dx+dy*dy <= (rad-1)*(rad-1) {
					g.Dot(x, y, pixel.AccentLit)
				}
			}
		}
		g.Circle(mid, mid, rad, pixel.Edge)
		g.VLine(mid, 1, wheel-2, pixel.Edge)
		g.HLine(1, mid, wheel-2, pixel.Edge)
		for i, sp := range seasonArt {
			x, y := mid+6, mid-6-sp.H
			if i == 1 || i == 2 {
				y = mid + 6
			}
			if i >= 2 {
				x = mid - 6 - sp.W
			}
			ink := pixel.Faint
			if i == at {
				ink = pixel.Accent
			}
			g.Sprite(x, y, sp, 0, ink, 1)
		}
		g.Finish(st)
	})
	res := SeasonWheelResult{Element: shell, Changed: *season != old}
	if res.Changed {
		slog.Debug("8bit: season wheel", "label", text, "season", names[*season])
	}
	return res
}

// quarter is an offset's quarter, clockwise from the top right.
func quarter(dx, dy int) int {
	switch {
	case dx >= 0 && dy < 0:
		return 0
	case dx >= 0:
		return 1
	case dy >= 0:
		return 2
	}
	return 3
}
