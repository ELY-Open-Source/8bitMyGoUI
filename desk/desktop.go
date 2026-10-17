package desk

import (
	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type DesktopOptions struct {
	Wallpaper Wallpaper
	// Width and Height are dots; zeros mean 320 by 200.
	Width, Height int
}

// Desktop lays icons down a dithered wallpaper, column by column.
func Desktop(c *ui.Context, text string, o DesktopOptions, icons func()) ui.Element {
	k := pixel.K(c)
	o.Wallpaper.name(k)
	w, h := pixel.Or("Desktop Width", o.Width, 320), pixel.Or("Desktop Height", o.Height, 200)
	return ui.Column(c).Size(k.D(w), k.D(h)).Wrap().AlignContent(ui.Start).Padding(k.D(4)).Gap(k.D(2)).Clip().
		FocusGroup(ui.Horizontal | ui.Vertical).Role(ui.RoleGroup).Label(text).Draw(func(p *ui.Painter, r ui.Rect) {
		paper(k.Canvas(p, r), 0, 0, w, h, o.Wallpaper)
	}).Children(icons)
}
