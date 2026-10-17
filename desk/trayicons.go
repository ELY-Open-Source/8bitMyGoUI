package desk

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/overlay"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// TrayIcon is a small icon by the clock, tipped.
type TrayIcon struct {
	Icon *pixel.Sprite
	Tip  string
}

type TrayIconsOptions struct {
	Disabled bool
}

type TrayIconsResult struct {
	Element ui.Element
	// Clicked is the icon pressed this frame, else -1.
	Clicked int
}

// TrayIcons are small icons; hover or focus shows tips.
func TrayIcons(c *ui.Context, text string, icons []TrayIcon, o TrayIconsOptions) TrayIconsResult {
	k := pixel.K(c)
	seen := map[string]bool{}
	for i, t := range icons {
		sprite("TrayIcons", t.Icon)
		if t.Tip == "" || seen[t.Tip] {
			panic(fmt.Sprintf("desk: TrayIcons icon %d tip %q is blank or twice", i, t.Tip))
		}
		seen[t.Tip] = true
	}
	res := TrayIconsResult{Clicked: -1}
	res.Element = ui.Row(c).Gap(k.D(1)).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).Label(text).Children(func() {
		for i, t := range icons {
			// Keyed by tip: icons come and go.
			ui.Box(c).Key(t.Tip).Children(func() {
				overlay.Tooltip(c, t.Tip, overlay.TooltipOptions{}, func() {
					shell, base, st := pixel.Shell(c, func() ui.Element {
						return ui.Box(c).Focusable().Role(ui.RoleButton).Label(t.Tip).Disabled(o.Disabled)
					})
					if !base.IsDisabled() && base.Clicked() {
						res.Clicked = i
					}
					shell.Size(k.D(14), k.D(14)).Draw(func(p *ui.Painter, r ui.Rect) {
						g := k.Canvas(p, r)
						dy := 0
						if st.Press {
							dy = 1
						}
						g.Sprite(1, 1+dy, t.Icon, 0, pixel.Ink, 1)
						g.Finish(st)
					})
				})
			})
		}
	})
	if res.Clicked >= 0 {
		slog.Debug("8bit: tray icon clicked", "label", text, "tip", icons[res.Clicked].Tip)
	}
	return res
}
