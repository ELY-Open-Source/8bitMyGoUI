package chat

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type UserCardOptions struct {
	Title  string
	Status display.Status
	// Actions label the buttons along the bottom.
	Actions []string
	// Width is in dots; zero means 160.
	Width    int
	Disabled bool
}

type UserCardResult struct {
	Element ui.Element
	// Clicked is the action clicked this frame, else -1.
	Clicked int
}

// UserCard shows an avatar, name, title and action buttons.
func UserCard(c *ui.Context, name string, o UserCardOptions) UserCardResult {
	k := pixel.K(c)
	w := pixel.Or("UserCard Width", o.Width, 160)
	res := UserCardResult{Clicked: -1}
	res.Element = ui.Column(c).Width(k.D(w)).Padding(k.D(6)).Gap(k.D(5)).Role(ui.RoleGroup).Label(name).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
	}).Children(func() {
		ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
			display.Avatar(c, name, display.AvatarOptions{Scale: 2})
			ui.Column(c).Gap(k.D(1)).Children(func() {
				ui.Row(c).Gap(k.D(3)).AlignItems(ui.Center).Children(func() {
					pixel.Text(c, line(name, w-62), pixel.TextOptions{})
					display.StatusDot(c, o.Status, display.StatusDotOptions{})
				})
				if o.Title != "" {
					pixel.Text(c, line(o.Title, w-50), pixel.TextOptions{Slot: pixel.Faint})
				}
			})
		})
		if len(o.Actions) == 0 {
			return
		}
		ui.Row(c).Gap(k.D(3)).FocusGroup(ui.Horizontal).Children(func() {
			for i, a := range o.Actions {
				if input.Button(c, a, input.ButtonOptions{Size: input.Small, Disabled: o.Disabled}).Clicked {
					res.Clicked = i
					slog.Debug("8bit: usercard action", "name", name, "action", a)
				}
			}
		})
	})
	return res
}
