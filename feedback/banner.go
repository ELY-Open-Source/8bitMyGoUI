package feedback

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

type BannerOptions struct {
	Kind Kind
	// Actions label buttons; Action reports the one clicked.
	Actions []string
	// Dismissable adds a × that reports Dismissed.
	Dismissable bool
	Disabled    bool
}

type BannerResult struct {
	Element   ui.Element
	Action    string
	Dismissed bool
}

// Banner is a full-width strip with a message and actions.
func Banner(c *ui.Context, message string, o BannerOptions) BannerResult {
	k := pixel.K(c)
	slot, sp := o.Kind.look()
	res := BannerResult{}
	res.Element = ui.Row(c).Gap(k.D(6)).Padding(k.D(4), k.D(8)).AlignItems(ui.Center).Role(ui.RoleStatus).Label(message).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Dither(0, 0, g.W, g.H, slot, pixel.Face, 12)
		g.HLine(0, g.H-1, g.W, slot)
	}).Children(func() {
		display.Icon(c, sp, display.IconOptions{Slot: slot})
		ui.Column(c).Grow(1).Children(func() { text.Paragraph(c, message, text.ParagraphOptions{}) })
		for _, a := range o.Actions {
			if input.Button(c, a, input.ButtonOptions{Size: input.Small, Disabled: o.Disabled}).Clicked {
				res.Action = a
			}
		}
		if o.Dismissable {
			res.Dismissed = input.IconButton(c, icon.Close, k.Msg("close"), input.IconButtonOptions{Variant: input.Ghost, Size: input.Small, Disabled: o.Disabled}).Clicked
		}
	})
	if res.Action != "" || res.Dismissed {
		slog.Debug("8bit: banner", "action", res.Action, "dismissed", res.Dismissed)
	}
	return res
}
