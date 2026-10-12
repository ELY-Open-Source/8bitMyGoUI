package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type SelectOptions struct {
	// Width is in dots; zero fits the longest option.
	Width    int
	Disabled bool
}

type SelectResult struct {
	Element *ui.Element
	Changed bool
}

// Select picks one option from a popup list.
func Select(c *ui.Context, text string, value *string, options []string, o SelectOptions) SelectResult {
	k := pixel.K(c)
	off := o.Disabled || pixel.Current(c).Inert
	var sel *ui.SelectParts[string]
	var shell, base *ui.Element
	var st pixel.State
	// Re-keying on each off state forgets an open popup.
	ui.Box(c).Key(off).Children(func() {
		shell, base, st = pixel.Shell(c, func() *ui.Element {
			sel = ui.SelectBase(c, value)
			return sel.Trigger.Label(text).Disabled(o.Disabled)
		})
	})
	shell.Size(k.D(pixel.Or("Select Width", o.Width, widest(options)+22)), k.D(pixel.LineHeight+4)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		ink, dy := face(g, Normal, st, sel.Open() && !off)
		g.Text(6, 2+dy, *value, ink, 1)
		g.Sprite(g.W-10, (g.H-tinyDown.H)/2+dy, tinyDown, 0, ink, 1)
		finish(g, st)
	})
	if !off {
		sel.Popup(func(panel *ui.Element) {
			menu(c, panel, func() {
				for _, opt := range options {
					item := sel.Item(opt).Label(opt)
					row(k, item, opt, item.Highlighted(), opt == *value)
				}
			})
		})
	}
	res := SelectResult{Element: shell, Changed: base.Changed()}
	if res.Changed {
		slog.Debug("8bit: select", "label", text, "value", *value)
	}
	return res
}

// widest is the widest option's width in dots.
func widest(options []string) int {
	w := 0
	for _, opt := range options {
		w = max(w, pixel.Measure(opt))
	}
	return w
}

// menu frames a popup panel as a game window.
func menu(c *ui.Context, panel *ui.Element, fn func()) {
	k := pixel.K(c)
	panel.Padding(k.D(3)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
	}).Children(func() {
		ui.Column(c).AlignItems(ui.Stretch).Children(fn)
	})
}

// row draws an option; a hand marks the lit one.
func row(k pixel.Kit, item *ui.Element, text string, lit, chosen bool) {
	item.Height(k.D(pixel.LineHeight)).MinWidth(k.D(pixel.Measure(text) + 16)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		ink := pixel.Ink
		if chosen {
			ink = pixel.Accent
		}
		if lit {
			g.Rect(0, 0, g.W, g.H, pixel.Lit)
			g.Sprite(1, (g.H-pixel.Hand.H)/2, pixel.Hand, 0, pixel.Edge, 1)
		}
		g.Text(13, 0, text, ink, 1)
	})
}
