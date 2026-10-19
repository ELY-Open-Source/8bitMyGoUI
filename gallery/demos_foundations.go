package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/layout"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var foundationEntries = []Entry{
	{0, "Palettes", "调色板", "foundations", "Eight palettes of sixteen slots: four seasons, by day and by night. Components name slots; rotating the palette repaints everything.", demoPalettes},
	{0, "Canvas", "画布", "foundations", "Rectangles, Bresenham lines, midpoint circles and six frame styles, all snapped to the dot grid.", demoCanvas},
	{0, "Dither", "抖动", "foundations", "Ordered Bayer dither stands in for gradients, shadows and translucency: seventeen levels between two slots.", demoDither},
	{0, "Pixel Font", "像素字体", "foundations", "Ark Pixel 12px with CJK, at integer scales, plus a hand-placed 3×5 micro font.", demoFont},
	{0, "Icons", "图标", "foundations", "One hundred and nine icons placed dot by dot. Click one to name it.", demoIcons},
	{0, "Sprites", "精灵", "foundations", "Sprites are string art in slot letters, with frames stepped by the clock. Click to turn around.", demoSprites},
	{0, "Editor", "编辑器", "foundations", "One text editing core drawn in the pixel font: caret, selection, clipboard, undo and input methods.", demoEditor},
	{0, "Focus", "焦点", "foundations", "Tab to see the focus ring and the pointing hand. Disabled controls dither out.", demoFocus},
}

func demoPalettes(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Grid(c).Columns(2).Gap(k.D(6)).Children(func() {
		for _, pal := range pixel.Palettes {
			s := pixel.Current(c)
			s.Palette = pal
			pixel.Scope(c, s, func() {
				kp := pixel.K(c)
				shell, base, st := pixel.Shell(c, func() ui.Element { return ui.ButtonBase(c).Label(pal.Name) })
				shell.Height(kp.D(52)).Draw(func(p *ui.Painter, r ui.Rect) {
					g := kp.Canvas(p, r)
					g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
					g.Text(6, 4, pal.Name+" "+pal.Zh, pixel.Ink, 1)
					g.Rect(5, 23, pixel.Slots*5+2, 10, pixel.Edge)
					for i := range pixel.Slots {
						g.Rect(6+i*5, 24, 5, 8, pixel.Slot(i))
					}
					g.Sprite(g.W-46, 4, icon.Heart, 0, pixel.Bad, 1)
					g.Sprite(g.W-32, 4, icon.Star, 0, pixel.Warn, 1)
					g.Sprite(g.W-18, 4, icon.Coin, 0, pixel.Ink, 1)
					g.Text(6, 34, "HP 99 · 你好", pixel.Accent, 1)
					if st.Focus {
						g.FocusMark(0, 0, g.W, g.H)
					}
				})
				if base.Clicked() {
					log("palette " + pal.Name)
				}
			})
		}
	})
}

func demoCanvas(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	styles := []string{"RAISED", "SUNKEN", "FLAT", "DOUBLE", "WINDOW", "NOTCHED"}
	ui.Box(c).Height(k.D(132)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, name := range styles {
			x := (i % 3) * 64
			y := (i / 3) * 34
			g.Frame(x, y, 58, 28, pixel.FrameStyle(i), pixel.Face)
			g.Micro(x+(58-pixel.MicroWidth(name))/2, y+12, name, pixel.Ink, 1)
		}
		for i := range 8 {
			g.Line(220, 100, 220+(i-4)*12, 72, pixel.Slot(8+i))
		}
		g.Disc(30, 104, 20, pixel.Info)
		g.Circle(30, 104, 24, pixel.Accent)
		g.Disc(90, 104, 12, pixel.Warn)
		g.Circle(90, 104, 16, pixel.Bad)
		g.Bevel(120, 84, 40, 40, pixel.Lit, pixel.Shade)
	})
}

func demoDither(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Box(c).Height(k.D(118)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for level := range 17 {
			g.Dither(level*12, 0, 12, 24, pixel.Face, pixel.Accent, level)
			g.Micro(level*12+2, 27, fmt.Sprint(level), pixel.Faint, 1)
		}
		for band := range 8 {
			g.Dither(0, 40+band*6, 204, 6, pixel.Info, pixel.Back, band*2)
		}
		g.Sprite(220, 44, icon.Chest, 0, pixel.Ink, 3)
		g.Checker(220, 44, 36, 36, pixel.Face)
		g.Micro(220, 84, "VEIL", pixel.Faint, 1)
	})
}

func demoFont(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		pixel.Text(c, "The quick brown fox jumps over the lazy dog.", pixel.TextOptions{})
		pixel.Text(c, "在限制中创造丰富。每个点都经过认真摆放。", pixel.TextOptions{})
		pixel.Text(c, "LEVEL 99", pixel.TextOptions{Scale: 2, Slot: pixel.Accent})
		pixel.Text(c, "像素", pixel.TextOptions{Scale: 3, Slot: pixel.Magic})
		ui.Box(c).Height(k.D(7)).Draw(func(p *ui.Painter, r ui.Rect) {
			k.Canvas(p, r).Micro(0, 0, "0123456789 ABCDEFGHIJKLMNOPQRSTUVWXYZ +-.:/%!", pixel.Ink, 1)
		})
		pixel.Text(c, "Long text wraps at the width the layout gives it, and Chinese 中文也会按字断行，标点不会出现在行首。", pixel.TextOptions{Wrap: true, Slot: pixel.Faint})
	})
}

func demoIcons(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Wrap().Gap(k.D(2)).Children(func() {
		for _, n := range icon.All {
			if input.IconButton(c, n.Sprite, n.Name, input.IconButtonOptions{Variant: input.Ghost}).Clicked {
				log("icon " + n.Name)
			}
		}
	})
}

var slime = pixel.MustSprite(`
	................
	................
	......eeee......
	....eeggggee....
	...egllgggggge..
	..egllgggggggge.
	..eglggkggkggge.
	..egggggkkgggge.
	..eggggggggggse.
	..egggggggggsse.
	...essssssssse..
	....eeeeeeeee...

	................
	................
	................
	................
	.....eeeeee.....
	...eegllggggee..
	..egllggkggkgge.
	.eglggggkkgggsge
	.egggggggggggsse
	.eggggggggggssse
	..esssssssssssse
	...eeeeeeeeeeee.
`)

var slimeLeft = slime.Mirror()

type demoKey string

func demoSprites(c *ui.Context, log func(string)) {
	left := ui.Local(c.Root(), demoKey("sprites"), func() bool { return false })
	sp := slime
	if *left {
		sp = slimeLeft
	}
	k := pixel.K(c)
	shell, base, st := pixel.Shell(c, func() ui.Element { return ui.ButtonBase(c).Label("slime") })
	shell.Size(k.D(120), k.D(48)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 40, g.W, 8, pixel.Good)
		g.Dither(0, 40, g.W, 8, pixel.Good, pixel.Shade, 4)
		step := g.Step(4)
		g.Sprite(8+step%8*2, 16, sp, step/2, pixel.Ink, 2)
		if st.Focus {
			g.FocusMark(0, 0, g.W, g.H)
		}
	})
	if base.Clicked() {
		*left = !*left
		log(fmt.Sprintf("turned, facing left %v", *left))
	}
}

type editorDemo struct{ name, secret, notes string }

func demoEditor(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	st := ui.Local(c.Root(), demoKey("editor"), func() editorDemo { return editorDemo{notes: "Two lines\nof notes"} })
	field := func(v *string, o pixel.EditorOptions) {
		layout.Panel(c, layout.PanelOptions{Style: pixel.Sunken, Fill: pixel.Well, Padding: 1}, func() {
			e := pixel.Editor(c, v, o)
			if e.Changed {
				log(o.Label + " changed")
			}
			if e.Submitted {
				log(o.Label + " submitted")
			}
		})
	}
	ui.Column(c).Gap(k.D(6)).Children(func() {
		field(&st.name, pixel.EditorOptions{Label: "Name", Placeholder: "Your name 你的名字"})
		field(&st.secret, pixel.EditorOptions{Label: "Password", Placeholder: "Password", Mask: '★'})
		field(&st.notes, pixel.EditorOptions{Label: "Notes", Multiline: true})
	})
}

func demoFocus(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(16)).Padding(0, k.D(14)).Children(func() {
		for _, name := range []string{"ATTACK", "MAGIC", "ITEM"} {
			if input.Button(c, name, input.ButtonOptions{}).Clicked {
				log(name)
			}
		}
		input.Button(c, "RUN", input.ButtonOptions{Disabled: true})
	})
}
