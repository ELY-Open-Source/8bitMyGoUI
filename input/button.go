package input

import (
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type ButtonOptions struct {
	Variant Variant
	Size    Size
	Icon    *pixel.Sprite
	// Loading shows a spinner and takes no clicks.
	Loading  bool
	Disabled bool
	// Label names it for assistive technology; zero uses text.
	Label string
}

type ButtonResult struct {
	Element *ui.Element
	// Base is the focusable inside Element.
	Base    *ui.Element
	Clicked bool
}

// Button is a beveled button whose face sinks when pressed.
func Button(c *ui.Context, text string, o ButtonOptions) ButtonResult {
	k := pixel.K(c)
	h, pad, top, sc := o.Size.metrics()
	name := label(o.Label, text)
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.ButtonBase(c).Label(name).Disabled(o.Disabled || o.Loading)
	})
	w := pixel.Measure(text) * sc
	lead := 0
	if o.Icon != nil || o.Loading {
		lead = 3
		if o.Icon != nil {
			lead = o.Icon.W
		}
		if text != "" {
			lead += 3
		}
	}
	if text == "" && o.Icon != nil {
		w, pad = 0, (h-o.Icon.W)/2
	}
	shell.Size(k.D(w+lead+2*pad), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		ink, dy := face(g, o.Variant, st, false)
		x := pad
		switch {
		case o.Loading:
			cy := (h-3)/2 + dy
			for i, d := range [4][2]int{{1, 0}, {2, 1}, {1, 2}, {0, 1}} {
				if g.Step(8)%4 != i {
					g.Dot(x+d[0], cy+d[1], ink)
				}
			}
			x += lead
		case o.Icon != nil:
			g.Sprite(x, (h-o.Icon.H)/2+dy, o.Icon, 0, ink, 1)
			x += lead
		}
		if text != "" {
			g.Text(x, top+dy, text, ink, sc)
		}
		g.Finish(st)
	})
	res := ButtonResult{Element: shell, Base: base, Clicked: base.Clicked()}
	if res.Clicked {
		slog.Debug("8bit: button", "label", name)
	}
	return res
}

type IconButtonOptions struct {
	Variant  Variant
	Size     Size
	Disabled bool
}

type IconButtonResult = ButtonResult

// IconButton is a square button named for assistive technology.
func IconButton(c *ui.Context, icon *pixel.Sprite, name string, o IconButtonOptions) IconButtonResult {
	if name == "" {
		panic("input: IconButton needs a name")
	}
	return Button(c, "", ButtonOptions{Variant: o.Variant, Size: o.Size, Icon: icon, Disabled: o.Disabled, Label: name})
}
