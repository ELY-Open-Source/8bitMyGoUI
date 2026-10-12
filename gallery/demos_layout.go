package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var layoutEntries = []Entry{
	{1, "Screen", "屏幕", "layout", "The root of a screen: Back fills the window, with an optional dithered vignette.", demoScreen},
	{2, "Window", "框体", "layout", "A beveled window with a title tab, the shape of every RPG menu.", demoWindow},
	{3, "Panel", "面板", "layout", "A box framed raised, sunken, flat, double, window or notched.", demoPanel},
	{4, "Card", "卡片", "layout", "A raised panel with a header strip, a body and a footer.", demoCard},
	{5, "Section", "分节", "layout", "A titled group inside a dotted rule.", demoSection},
	{6, "Divider", "分隔线", "layout", "Rules across or down: solid, dotted, dashed, double, ornament.", demoDivider},
	{7, "ScrollArea", "卷动区", "layout", "A scroll view with a pixel scrollbar: arrows step, the track pages, the thumb drags.", demoScrollArea},
	{8, "Split", "分栏", "layout", "Two panes and a grip. Drag it, or focus it and use the arrow keys.", demoSplit},
	{9, "Grid", "网格", "layout", "Equal columns with solid or dotted lines between the cells.", demoGrid},
	{10, "Tiles", "平铺", "layout", "Equal tiles that wrap into rows as the width allows.", demoTiles},
	{11, "Masonry", "瀑布流", "layout", "Each item drops into the shortest column.", demoMasonry},
	{12, "AspectBox", "比例框", "layout", "Keeps a ratio in whole dots at whatever width it gets.", demoAspectBox},
	{13, "Layers", "图层", "layout", "Layers stacked in one box, each over the one before.", demoLayers},
	{14, "Viewport", "视口", "layout", "A fixed virtual resolution, drawn at the largest Dot that fits and letterboxed.", demoViewport},
	{15, "TitleBar", "标题栏", "layout", "Stripes, a title, a close box and a minimize box.", demoTitleBar},
	{16, "StatusBar", "状态栏", "layout", "Segments split by sunken dividers; the first one grows.", demoStatusBar},
	{17, "Resizable", "可调尺寸", "layout", "A box with a corner grip. Drag it, or focus it and use the arrow keys.", demoResizable},
	{18, "Handheld", "掌机", "layout", "A pocket console: a screen, a power LED, a D-pad, A and B.", demoHandheld},
	{19, "CRTFrame", "显像管", "layout", "A monitor bezel with scanlines over the picture.", demoCRTFrame},
	{20, "AppFrame", "应用框架", "layout", "Header, sidebar, main and footer; on a narrow screen the sidebar folds above.", demoAppFrame},
}

var monsters = []string{"Slime", "Bat", "Goblin", "Mimic", "Ghost", "Golem", "Drake", "Imp"}

func demoScreen(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	on := ui.Local(c.Root(), demoKey("screen"), func() bool { return true })
	ui.Column(c).Gap(k.D(6)).Children(func() {
		ui.Box(c).FillWidth().Height(k.D(96)).Children(func() {
			layout.Screen(c, layout.ScreenOptions{Vignette: *on, Padding: 12}, func() {
				pixel.Text(c, "A NEW SAVE FILE", pixel.TextOptions{})
				pixel.Text(c, "Day 1. The cartridge starts.", pixel.TextOptions{Slot: pixel.Faint})
			})
		})
		if input.Button(c, "VIGNETTE", input.ButtonOptions{}).Clicked {
			*on = !*on
			log(fmt.Sprintf("vignette %v", *on))
		}
	})
}

func demoWindow(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Start).Children(func() {
		layout.Window(c, layout.WindowOptions{Title: "ITEMS"}, func() {
			for _, s := range []string{"Potion x3", "Ether  x1", "Key    x2"} {
				pixel.Text(c, s, pixel.TextOptions{})
			}
		})
		layout.Window(c, layout.WindowOptions{}, func() {
			pixel.Text(c, "HP 32/40", pixel.TextOptions{})
			pixel.Text(c, "MP 12/18", pixel.TextOptions{})
		})
	})
}

func demoPanel(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Wrap().Gap(k.D(6)).Children(func() {
		for i, name := range []string{"RAISED", "SUNKEN", "FLAT", "DOUBLE", "WINDOW", "NOTCHED"} {
			layout.Panel(c, layout.PanelOptions{Style: pixel.FrameStyle(i)}, func() { pixel.Text(c, name, pixel.TextOptions{}) })
		}
	})
}

func demoCard(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Box(c).Width(k.D(200)).Children(func() {
		layout.Card(c, layout.CardOptions{Title: "QUEST", Icon: icon.Scroll, Footer: func() {
			if input.Button(c, "ACCEPT", input.ButtonOptions{}).Clicked {
				log("quest accepted")
			}
			if input.Button(c, "LATER", input.ButtonOptions{}).Clicked {
				log("quest saved for later")
			}
		}}, func() {
			pixel.Text(c, "Find the lost cartridge.", pixel.TextOptions{})
			pixel.Text(c, "Reward: 300 G", pixel.TextOptions{Slot: pixel.Faint})
		})
	})
}

func demoSection(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Box(c).Width(k.D(220)).Children(func() {
		layout.Section(c, "Sound", layout.SectionOptions{}, func() {
			pixel.Text(c, "Music   8", pixel.TextOptions{})
			pixel.Text(c, "Effects 5", pixel.TextOptions{})
		})
	})
}

func demoDivider(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	names := []string{"SOLID", "DOTTED", "DASHED", "DOUBLE", "ORNAMENT"}
	ui.Row(c).Gap(k.D(12)).AlignItems(ui.Stretch).Children(func() {
		ui.Column(c).Grow(1).Gap(k.D(4)).Children(func() {
			for i, n := range names {
				pixel.Text(c, n, pixel.TextOptions{Slot: pixel.Faint})
				layout.Divider(c, layout.DividerOptions{Style: layout.DividerStyle(i)})
			}
		})
		for i := range names {
			layout.Divider(c, layout.DividerOptions{Vertical: true, Style: layout.DividerStyle(i), Slot: pixel.Accent})
		}
	})
}

func demoScrollArea(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := ui.Local(c.Root(), demoKey("scroll"), func() ui.ScrollState { return ui.ScrollState{} })
	ui.Column(c).Gap(k.D(4)).Children(func() {
		r := layout.ScrollArea(c, s, layout.ScrollAreaOptions{Label: "Bestiary"}, func() {
			for i := range 40 {
				pixel.Text(c, fmt.Sprintf("No.%03d %s", i+1, monsters[i%len(monsters)]), pixel.TextOptions{})
			}
		})
		r.Element.Size(k.D(200), k.D(112))
		if r.Changed {
			log(fmt.Sprintf("scrolled to %d dots", int(s.Y/k.Dot())))
		}
		pixel.Text(c, fmt.Sprintf("Y %d of %d dots", int(s.Y/k.Dot()), int(s.MaxY/k.Dot())), pixel.TextOptions{Slot: pixel.Faint})
	})
}

func demoSplit(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	size := ui.Local(c.Root(), demoKey("split"), func() int { return 120 })
	pane := func(name string) func() {
		return func() {
			layout.Panel(c, layout.PanelOptions{Style: pixel.Sunken, Fill: pixel.Well}, func() { pixel.Text(c, name, pixel.TextOptions{}) }).Fill()
		}
	}
	ui.Box(c).FillWidth().Height(k.D(96)).Children(func() {
		r := layout.Split(c, size, layout.SplitOptions{Label: "Map and log"}, pane("MAP"), pane("LOG"))
		r.Element.Fill()
		if r.Changed {
			log(fmt.Sprintf("split %d dots", *size))
		}
	})
}

func demoGrid(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	items := []*pixel.Sprite{icon.Sword, icon.Shield, icon.Potion, icon.Gem, icon.Key, icon.Coin}
	ui.Row(c).Gap(k.D(12)).AlignItems(ui.Start).Children(func() {
		for _, lines := range []layout.GridLines{layout.SolidLines, layout.DottedLines} {
			layout.Grid(c, len(items), layout.GridOptions{Columns: 3, Lines: lines}, func(i int) {
				display.Icon(c, items[i], display.IconOptions{Scale: 2})
			})
		}
	})
}

func demoTiles(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	items := []*pixel.Sprite{icon.Sun, icon.Moon, icon.Cloud, icon.Snowflake, icon.Leaf, icon.Flower, icon.Flame, icon.Drop, icon.Bolt, icon.Star}
	layout.Tiles(c, len(items), layout.TilesOptions{W: 28}, func(i int) {
		ui.Box(c).Fill().Center().Draw(func(p *ui.Painter, r ui.Rect) {
			k.Canvas(p, r).Frame(0, 0, 28, 28, pixel.Raised, pixel.Face)
		}).Children(func() { display.Icon(c, items[i], display.IconOptions{}) })
	}).Width(k.D(160))
}

func demoMasonry(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	heights := []int{40, 24, 56, 32, 24, 48, 32, 40}
	slots := []pixel.Slot{pixel.Accent, pixel.Good, pixel.Info, pixel.Warn, pixel.Magic}
	layout.Masonry(c, heights, layout.MasonryOptions{}, func(i int) {
		ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Flat, slots[i%len(slots)])
			g.Micro(3, 3, fmt.Sprint(i+1), pixel.Ink, 1)
		})
	}).Width(k.D(180))
}

func demoAspectBox(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Box(c).Width(k.D(160)).Children(func() {
		layout.AspectBox(c, layout.AspectBoxOptions{W: 16, H: 9}, func() {
			ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				for i := range 17 {
					g.Dither(0, i*g.H/17, g.W, g.H/17+1, pixel.Back, pixel.Info, i)
				}
				g.Text(4, 2, fmt.Sprintf("16:9  %d×%d", g.W, g.H), pixel.Ink, 1)
			})
		})
	})
}

func demoLayers(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	layout.Layers(c, layout.LayersOptions{}, func() {
		ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Dither(0, 0, g.W, g.H, pixel.Well, pixel.Info, 6)
		})
	}, func() {
		ui.Box(c).Fill().Center().Children(func() { display.Icon(c, icon.Chest, display.IconOptions{Scale: 3}) })
	}, func() {
		ui.Box(c).Padding(k.D(3)).Children(func() { pixel.Text(c, "LAYER 3", pixel.TextOptions{}) })
	}).Size(k.D(120), k.D(64))
}

func demoViewport(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Box(c).FillWidth().Height(k.D(100)).Children(func() {
		layout.Viewport(c, layout.ViewportOptions{W: 64, H: 40}, func() {
			kv := pixel.K(c)
			ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
				g := kv.Canvas(p, r)
				g.Rect(0, 0, g.W, g.H, pixel.Info)
				g.Rect(0, 30, g.W, 10, pixel.Good)
				g.Disc(50, 10, 5, pixel.Warn)
				g.Micro(3, 3, "64X40", pixel.Ink, 1)
			})
		})
	})
}

func demoTitleBar(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Box(c).Width(k.D(220)).Children(func() {
		r := layout.TitleBar(c, "Untitled", layout.TitleBarOptions{})
		if r.Close {
			log("close")
		}
		if r.Minimize {
			log("minimize")
		}
	})
}

func demoStatusBar(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Box(c).Width(k.D(260)).Children(func() {
		layout.StatusBar(c, []string{"Ready", "Ln 1, Col 1", "UTF-8"}, layout.StatusBarOptions{})
	})
}

func demoResizable(c *ui.Context, log func(string)) {
	w := ui.Local(c.Root(), demoKey("resize-w"), func() int { return 120 })
	h := ui.Local(c.Root(), demoKey("resize-h"), func() int { return 56 })
	r := layout.Resizable(c, w, h, layout.ResizableOptions{MinW: 48, MinH: 32, MaxW: 240, MaxH: 120}, func() {
		layout.Panel(c, layout.PanelOptions{Style: pixel.Sunken, Fill: pixel.Well}, func() {
			pixel.Text(c, fmt.Sprintf("%d×%d", *w, *h), pixel.TextOptions{})
		}).Fill()
	})
	if r.Changed {
		log(fmt.Sprintf("size %d×%d", *w, *h))
	}
}

func demoHandheld(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	at := ui.Local(c.Root(), demoKey("handheld"), func() [2]int { return [2]int{36, 20} })
	r := layout.Handheld(c, layout.HandheldOptions{ScreenW: 80, ScreenH: 56, Brand: "8BIT POCKET", Label: "Pocket console"}, func() {
		ui.Box(c).Fill().Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Rect(0, 0, g.W, g.H, pixel.Well)
			g.Sprite(at[0], at[1], slime, g.Step(4)/2, pixel.Ink, 1)
		})
	})
	moves := map[string][2]int{"up": {0, -4}, "down": {0, 4}, "left": {-4, 0}, "right": {4, 0}}
	for name, on := range map[string]bool{"up": r.Up, "down": r.Down, "left": r.Left, "right": r.Right} {
		if on {
			at[0] = max(0, min(at[0]+moves[name][0], 80-slime.W))
			at[1] = max(0, min(at[1]+moves[name][1], 56-slime.H))
			log(name)
		}
	}
	if r.A {
		log("A: jump")
	}
	if r.B {
		log("B: run")
	}
}

func demoCRTFrame(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	layout.CRTFrame(c, layout.CRTFrameOptions{Brand: "8BIT CRT"}, func() {
		ui.Column(c).Size(k.D(180), k.D(72)).Padding(k.D(4)).Children(func() {
			pixel.Text(c, "READY.", pixel.TextOptions{Slot: pixel.Good})
			pixel.Text(c, "10 PRINT \"HI\"", pixel.TextOptions{Slot: pixel.Good})
		})
	}).AlignSelf(ui.Start)
}

func demoAppFrame(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Box(c).FillWidth().Height(k.D(150)).Children(func() {
		layout.AppFrame(c, layout.AppFrameOptions{SidebarWidth: 72}, func() {
			pixel.Text(c, "8BIT MAIL", pixel.TextOptions{})
		}, func() {
			for _, s := range []string{"Inbox", "Sent", "Spam"} {
				if input.Button(c, s, input.ButtonOptions{}).Clicked {
					log("open " + s)
				}
			}
		}, func() {
			pixel.Text(c, "No new mail. Go outside.", pixel.TextOptions{})
		}, func() {
			pixel.Text(c, "3 folders", pixel.TextOptions{Slot: pixel.Faint})
		})
	})
}
