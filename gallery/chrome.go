package gallery

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

const logo = "8BIT·MYGO·UI"

func (a *App) topBar(c *ui.Context) {
	k := pixel.K(c)
	ui.Row(c).Height(k.D(24)).Padding(0, k.D(6)).Gap(k.D(8)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Face)
		g.HLine(0, g.H-2, g.W, pixel.Shade)
		g.HLine(0, g.H-1, g.W, pixel.Edge)
	}).Children(func() {
		if press(c, pixel.Measure(logo)+2, 18, "home", func(g pixel.Canvas, st pixel.State) {
			ink := pixel.Accent
			if st.Hover {
				ink = pixel.AccentLit
			}
			shadowText(g, 0, 1, logo, ink, 1)
		}) {
			a.Go("/")
		}
		ui.Column(c).Width(k.D(150)).Padding(k.D(2), k.D(4)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
		}).Children(func() {
			a.search = pixel.Editor(c, &a.query, pixel.EditorOptions{Placeholder: "/ search 搜索", Label: "Search components"})
			if m := matches(a.query); a.search.Submitted && len(m) > 0 {
				a.Go(m[0].Path())
			}
		})
		ui.Spacer(c)
		a.paletteControl(c)
		other, mark := pixel.Zh, "中"
		if a.locale == pixel.Zh {
			other, mark = pixel.En, "EN"
		}
		if press(c, 20, 16, "language", func(g pixel.Canvas, st pixel.State) {
			g.Frame(0, 0, g.W, g.H, pixel.Notched, pixel.Face)
			g.Text((g.W-pixel.Measure(mark))/2, 0, mark, pixel.Ink, 1)
		}) {
			a.locale = other
		}
	})
}

func (a *App) paletteControl(c *ui.Context) {
	pal := pixel.Palettes[a.palette]
	if press(c, 12, 16, "previous palette", func(g pixel.Canvas, st pixel.State) {
		g.Sprite(0, 2, icon.ChevronLeft, 0, pixel.Ink, 1)
	}) {
		a.setPalette(a.palette - 1)
	}
	name := strings.ToUpper(pal.Name) + " " + pal.Zh
	press(c, 4*pixel.Slots+6+pixel.Measure(name), 16, "palette "+pal.Name, func(g pixel.Canvas, st pixel.State) {
		g.Rect(0, 2, 4*pixel.Slots+2, 12, pixel.Edge)
		for i := range pixel.Slots {
			g.Rect(1+i*4, 3, 4, 10, pixel.Slot(i))
		}
		g.Text(4*pixel.Slots+6, 0, name, pixel.Ink, 1)
	})
	if press(c, 12, 16, "next palette", func(g pixel.Canvas, st pixel.State) {
		g.Sprite(0, 2, icon.ChevronRight, 0, pixel.Ink, 1)
	}) {
		a.setPalette(a.palette + 1)
	}
}

func (a *App) sidebar(c *ui.Context) {
	k := pixel.K(c)
	here := a.router.Path()
	ui.Scroll(c).Width(k.D(120)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Face)
		g.VLine(g.W-2, 0, g.H, pixel.Shade)
		g.VLine(g.W-1, 0, g.H, pixel.Edge)
	}).Children(func() {
		ui.Column(c).Padding(k.D(4)).Gap(k.D(1)).Children(func() {
			for _, cat := range Categories {
				n := count(cat.Key)
				path := "/cat/" + cat.Key
				on := here == path
				if press(c, 0, 16, cat.Name, func(g pixel.Canvas, st pixel.State) {
					ink := pixel.Ink
					if on {
						g.Rect(0, 0, g.W, g.H, pixel.Accent)
						ink = pixel.AccentInk
					} else if st.Hover {
						g.Rect(0, 0, g.W, g.H, pixel.Lit)
					}
					if n == 0 {
						ink = pixel.Faint
					}
					g.Sprite(2, 2, cat.Icon, 0, ink, 1)
					g.Text(17, 0, cat.Name, ink, 1)
					num := fmt.Sprint(n)
					g.Micro(g.W-pixel.MicroWidth(num)-3, 6, num, ink, 1)
				}) {
					a.Go(path)
				}
			}
		})
	})
}

func (a *App) statusBar(c *ui.Context) {
	k := pixel.K(c)
	hints := "/ search · P palette · [ ] prev next"
	total := fmt.Sprintf("%d ENTRIES", len(Entries()))
	ui.Row(c).Height(k.D(18)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Rect(0, 0, g.W, g.H, pixel.Face)
		g.HLine(0, 0, g.W, pixel.Edge)
		g.HLine(0, 1, g.W, pixel.Lit)
		g.Text(6, 2, hints, pixel.Faint, 1)
		g.Micro(g.W-pixel.MicroWidth(total)-6, 8, total, pixel.Faint, 1)
	})
}

func count(cat string) int {
	n := 0
	for _, e := range Entries() {
		if e.Category == cat {
			n++
		}
	}
	return n
}
