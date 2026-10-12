package gallery

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/feedback"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/fx"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

func (a *App) content(c *ui.Context) {
	if strings.TrimSpace(a.query) != "" {
		page(c, func() { a.results(c) })
		return
	}
	a.router.View(c, func(r *ui.Route) {
		page(c, func() {
			switch {
			case r.Match("/"):
				r.Title("8bitMyGoUI")
				a.home(c)
			case r.Match("/cat/{key}"):
				a.category(c, r.Param("key"))
			case r.Match("/c/{slug}"), r.Match("/f/{slug}"):
				if e, ok := find(r.Param("slug")); ok && e.Path() == r.Path() {
					r.Title(e.Name)
					a.entry(c, e)
					return
				}
				a.lost(c, r.Path())
			default:
				a.lost(c, r.Path())
			}
		})
	}).Grow(1)
}

// page scrolls a page; a router's column would squeeze it.
func page(c *ui.Context, fn func()) {
	k := pixel.K(c)
	ui.Scroll(c).Grow(1).Padding(k.D(8)).Gap(k.D(6)).Children(fn)
}

func (a *App) home(c *ui.Context) {
	k := pixel.K(c)
	p := k.Palette()
	// The banner lives in the palette's own season and hour.
	layout.Layers(c, layout.LayersOptions{}, func() {
		fx.SeasonScene(c, p.Season, hourOf(p), fx.SeasonSceneOptions{Height: 96, Label: p.Name + " scene"}).FillWidth()
	}, func() {
		ui.Column(c).Padding(k.D(10)).Children(func() {
			text.OutlineText(c, logo, text.OutlineTextOptions{Slot: pixel.Accent, Scale: 3})
		})
	})
	pixel.Text(c, "Pixel components for MyGo. Every dot sits on the grid; every color comes from sixteen.", pixel.TextOptions{Wrap: true})
	pixel.Text(c, "为 MyGo 打造的像素组件。每个点落在点阵上，每种颜色来自十六格调色板。", pixel.TextOptions{Wrap: true, Slot: pixel.Faint})
	comps := len(Entries()) - count("foundations")
	stats := fmt.Sprintf("%d COMPONENTS  %d CATEGORIES  %d PALETTES", comps, len(Categories)-1, len(pixel.Palettes))
	text.MicroText(c, stats, text.MicroTextOptions{Slot: pixel.Accent})
	a.palettes(c)
	// Each grid is one Tab stop; arrows move inside it.
	ui.Grid(c).Columns(3).Gap(k.D(4)).FocusGroup(ui.Horizontal | ui.Vertical).Children(func() {
		for _, cat := range Categories {
			n := count(cat.Key)
			if press(c, 0, 40, cat.Name, func(g pixel.Canvas, st pixel.State) {
				fill := pixel.Face
				if st.Hover {
					fill = pixel.Lit
				}
				g.Frame(0, 0, g.W, g.H, pixel.Window, fill)
				g.Sprite(6, 6, cat.Icon, 0, pixel.Accent, 2)
				g.Text(34, 4, cat.Name, pixel.Ink, 1)
				g.Text(34, 20, cat.Zh, pixel.Faint, 1)
				num := fmt.Sprint(n)
				g.Micro(g.W-pixel.MicroWidth(num)-6, g.H-10, num, pixel.Faint, 1)
			}) {
				a.Go("/cat/" + cat.Key)
			}
		}
	})
}

// palettes shows each palette in its colors; press to pick.
func (a *App) palettes(c *ui.Context) {
	k := pixel.K(c)
	ui.Grid(c).Columns(4).Gap(k.D(4)).FocusGroup(ui.Horizontal | ui.Vertical).Children(func() {
		for i, pal := range pixel.Palettes {
			name := pal.Name
			if i == a.palette {
				name += ", " + strings.ToLower(k.Msg("in_use"))
			}
			set := pixel.Current(c)
			set.Palette = pal
			pixel.Scope(c, set, func() {
				if press(c, 0, 30, name, func(g pixel.Canvas, st pixel.State) {
					fill := pixel.Face
					if i == a.palette || st.Hover {
						fill = pixel.Lit
					}
					g.Frame(0, 0, g.W, g.H, pixel.Window, fill)
					for s := range pixel.Slots {
						g.Rect(5+3*s, 4, 3, 8, pixel.Slot(s))
					}
					g.Text(5, 14, pal.Name+" "+pal.Zh, pixel.Ink, 1)
					if i == a.palette {
						g.Rect(g.W-9, 5, 4, 4, pixel.Accent)
					}
				}) {
					a.setPalette(i)
				}
			})
		}
	})
}

func (a *App) category(c *ui.Context, key string) {
	var cat Category
	for _, x := range Categories {
		if x.Key == key {
			cat = x
		}
	}
	if cat.Key == "" {
		a.lost(c, "/cat/"+key)
		return
	}
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		display.Icon(c, cat.Icon, display.IconOptions{Slot: pixel.Accent, Scale: 2})
		text.ShadowText(c, cat.Name, text.ShadowTextOptions{Scale: 2})
		pixel.Text(c, cat.Zh, pixel.TextOptions{Slot: pixel.Faint})
	})
	ui.Grid(c).Columns(3).Gap(k.D(4)).FocusGroup(ui.Horizontal | ui.Vertical).Children(func() {
		for _, e := range Entries() {
			if e.Category == key {
				a.card(c, e)
			}
		}
	})
}

func (a *App) card(c *ui.Context, e Entry) {
	if press(c, 0, 36, e.Name, func(g pixel.Canvas, st pixel.State) {
		fill := pixel.Face
		if st.Hover {
			fill = pixel.Lit
		}
		g.Frame(0, 0, g.W, g.H, pixel.Raised, fill)
		if e.ID > 0 {
			g.Micro(5, 4, fmt.Sprintf("%03d", e.ID), pixel.Accent, 1)
		}
		g.Text(5, 10, e.Name, pixel.Ink, 1)
		g.Text(g.W-pixel.Measure(e.Zh)-5, 10, e.Zh, pixel.Faint, 1)
	}) {
		a.Go(e.Path())
	}
}

func (a *App) entry(c *ui.Context, e Entry) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		if e.ID > 0 {
			text.MicroText(c, fmt.Sprintf("%03d", e.ID), text.MicroTextOptions{Slot: pixel.Accent, Scale: 2})
		}
		text.ShadowText(c, e.Name, text.ShadowTextOptions{Scale: 2})
		pixel.Text(c, e.Zh, pixel.TextOptions{Slot: pixel.Faint})
	})
	pixel.Text(c, e.Summary, pixel.TextOptions{Wrap: true})
	layout.Panel(c, layout.PanelOptions{Style: pixel.Window, Padding: 8}, func() { e.Demo(c, a.logger(e)) }).MinHeight(k.D(80))
	code := strings.ReplaceAll(e.Code(), "\t", "  ")
	ui.Row(c).Gap(k.D(4)).AlignItems(ui.Start).Children(func() {
		layout.Panel(c, layout.PanelOptions{Style: pixel.Sunken, Fill: pixel.Well, Padding: 2}, func() {
			ui.ScrollBoth(c).Height(k.D(160)).Children(func() { codeBlock(c, code) })
		}).Grow(1)
		ui.Column(c).Width(k.D(140)).Gap(k.D(4)).Children(func() {
			if input.Button(c, "COPY CODE", input.ButtonOptions{}).Clicked {
				c.WriteClipboard(code)
				a.logger(e)("code copied")
			}
			layout.Panel(c, layout.PanelOptions{Style: pixel.Sunken, Fill: pixel.Well, Padding: 2}, func() {
				for _, l := range a.logs[e.Slug()] {
					pixel.Text(c, "> "+l, pixel.TextOptions{Slot: pixel.Faint, MaxLines: 1})
				}
			}).Height(k.D(140))
		})
	})
	ui.Row(c).Gap(k.D(4)).Children(func() {
		if input.Button(c, "◀ PREV", input.ButtonOptions{}).Clicked {
			a.step(-1)
		}
		if input.Button(c, "NEXT ▶", input.ButtonOptions{}).Clicked {
			a.step(1)
		}
	})
}

func codeBlock(c *ui.Context, code string) {
	k := pixel.K(c)
	lines := strings.Split(strings.TrimRight(code, "\n"), "\n")
	w := 0
	for _, l := range lines {
		w = max(w, pixel.Measure(l))
	}
	ui.Box(c).Size(k.D(w), k.D(len(lines)*pixel.LineHeight)).Label(code).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, l := range lines {
			ink := pixel.Ink
			if strings.HasPrefix(strings.TrimSpace(l), "//") {
				ink = pixel.Faint
			}
			g.Text(0, i*pixel.LineHeight, l, ink, 1)
		}
	})
}

func (a *App) results(c *ui.Context) {
	found := matches(a.query)
	pixel.Text(c, fmt.Sprintf("%d results for %q", len(found), strings.TrimSpace(a.query)), pixel.TextOptions{Slot: pixel.Faint})
	k := pixel.K(c)
	ui.Grid(c).Columns(3).Gap(k.D(4)).FocusGroup(ui.Horizontal | ui.Vertical).Children(func() {
		for _, e := range found {
			a.card(c, e)
		}
	})
}

// lost answers unknown paths with GAME OVER.
func (a *App) lost(c *ui.Context, path string) {
	if feedback.ErrorState(c, "No page at "+path, feedback.ErrorStateOptions{}).Continue {
		a.Go("/")
	}
}
