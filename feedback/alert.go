package feedback

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

// Kind is a message's tone.
type Kind uint8

const (
	Info Kind = iota
	Success
	Warning
	Danger
)

// look is a kind's slot and icon.
func (k Kind) look() (pixel.Slot, *pixel.Sprite) {
	switch k {
	case Info:
		return pixel.Info, icon.Info
	case Success:
		return pixel.Good, icon.Check
	case Warning:
		return pixel.Warn, icon.Warning
	case Danger:
		return pixel.Bad, icon.Error
	}
	panic(fmt.Sprintf("feedback: unknown Kind %d", k))
}

type AlertOptions struct {
	Kind Kind
	// Closable adds a × that reports Closed.
	Closable bool
	Disabled bool
}

type AlertResult struct {
	Element *ui.Element
	Closed  bool
}

// Alert is a framed callout: icon, title and text.
func Alert(c *ui.Context, title, message string, o AlertOptions) AlertResult {
	k := pixel.K(c)
	slot, sp := o.Kind.look()
	res := AlertResult{}
	res.Element = ui.Row(c).Gap(k.D(6)).Padding(k.D(6), k.D(8), k.D(6), k.D(10)).AlignItems(ui.Start).Role(ui.RoleStatus).Label(title).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
		g.Rect(3, 3, 3, g.H-6, slot)
	}).Children(func() {
		display.Icon(c, sp, display.IconOptions{Slot: slot})
		ui.Column(c).Grow(1).Gap(k.D(2)).Children(func() {
			pixel.Text(c, title, pixel.TextOptions{Slot: slot})
			text.Paragraph(c, message, text.ParagraphOptions{})
		})
		if o.Closable {
			res.Closed = input.IconButton(c, icon.Close, k.Msg("close"), input.IconButtonOptions{Variant: input.Ghost, Size: input.Small, Disabled: o.Disabled}).Clicked
		}
	})
	if res.Closed {
		slog.Debug("8bit: alert closed", "title", title)
	}
	return res
}
