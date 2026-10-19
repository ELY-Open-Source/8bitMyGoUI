package gallery

import (
	"fmt"
	"slices"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/navigation"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var navigationEntries = []Entry{
	{131, "Tabs", "标签页", "navigation", "Folder tabs; the chosen one joins its panel. Arrows move.", demoTabs},
	{132, "VerticalTabs", "竖排标签", "navigation", "Tabs down the left edge, for settings screens.", demoVerticalTabs},
	{133, "Sidebar", "侧边栏", "navigation", "Sections of items with icons and badges. Up and Down choose; one Tab stop.", demoSidebar},
	{134, "Breadcrumbs", "面包屑", "navigation", "Where you are, split by pointers. Every part but the last walks back.", demoBreadcrumbs},
	{135, "CommandMenu", "指令菜单", "navigation", "FIGHT, MAGIC, ITEM, RUN, with a bobbing ▶. Enter picks; Escape backs out.", demoCommandMenu},
	{136, "MenuBar", "菜单栏", "navigation", "A bar of menus. Click one, then the arrows go across and down.", demoMenuBar},
	{137, "DropdownMenu", "下拉菜单", "navigation", "A button with a menu: icons, rules and shortcut names.", demoDropdownMenu},
	{138, "ContextMenu", "右键菜单", "navigation", "Right click, or the menu key, opens a menu at the pointer.", demoContextMenu},
	{139, "Pagination", "分页", "navigation", "◀ 1 … 4 5 6 … 12 ▶. The current page is pressed in.", demoPagination},
	{140, "PageDots", "页码点", "navigation", "A dot per page; click one, or use the arrows.", demoPageDots},
	{141, "Steps", "关卡阶段", "navigation", "Flags along a trail: cleared, waving where you are, locked ahead.", demoSteps},
	{142, "NavRail", "导航栏", "navigation", "A rail of icons. Hover or focus one and it says its name.", demoNavRail},
	{143, "BottomNav", "底部导航", "navigation", "Icons over labels along the bottom, for small screens.", demoBottomNav},
	{144, "CommandPalette", "命令面板", "navigation", "Cmd+K, type a few letters, press Enter.", demoCommandPalette},
	{145, "BackForward", "前进后退", "navigation", "◀ and ▶ walking a router's history.", demoBackForward},
	{146, "TableOfContents", "目录", "navigation", "A manual's sections; ▶ marks where you are.", demoTableOfContents},
	{147, "LevelSelect", "选关", "navigation", "Stages in a grid: locked, open, or cleared with stars.", demoLevelSelect},
	{148, "WorldMap", "世界地图", "navigation", "Places on dotted trails. The arrows walk the hero; Enter goes in.", demoWorldMap},
	{149, "Toolbar", "工具栏", "navigation", "Tool buttons as one Tab stop; arrows move along.", demoToolbar},
	{150, "DocumentTabs", "文档标签", "navigation", "Tabs you can close with × or Delete.", demoDocumentTabs},
}

type navDemo struct {
	tab, vtab, side, toc, rail, bottom, level, node, doc int
	page, dot                                            int
	crumbs, docs                                         []string
	palette                                              bool
	router                                               *ui.Router
	pen, fill                                            bool
	newDocs                                              int
}

func navState(c *ui.Context) *navDemo {
	return ui.Local(c.Root(), demoKey("navigation"), func() navDemo {
		r := ui.NewRouter("/town")
		r.Push("/field")
		return navDemo{page: 5, dot: 1, level: 2, crumbs: []string{"World", "Cave", "B1", "B2"},
			docs: []string{"map.txt", "save1.dat", "notes.md"}, router: r, pen: true}
	})
}

func demoTabs(c *ui.Context, log func(string)) {
	s := navState(c)
	names := []string{"ITEMS", "MAGIC", "STATUS"}
	body := []string{"Potion ×3, Ether ×1, Tent ×1", "Fire, Cure, Haste", "HP 42/42  MP 12/18  LV 7"}
	if navigation.Tabs(c, "Menu", names, &s.tab, navigation.TabsOptions{}, func() {
		pixel.Text(c, body[s.tab], pixel.TextOptions{})
	}).Changed {
		log(names[s.tab])
	}
}

func demoVerticalTabs(c *ui.Context, log func(string)) {
	s := navState(c)
	names := []string{"Video", "Audio", "Controls"}
	if navigation.VerticalTabs(c, "Settings", names, &s.vtab, navigation.VerticalTabsOptions{}, func() {
		pixel.Text(c, names[s.vtab]+" settings go here.", pixel.TextOptions{Slot: pixel.Faint})
	}).Changed {
		log(names[s.vtab])
	}
}

func demoSidebar(c *ui.Context, log func(string)) {
	s := navState(c)
	if navigation.Sidebar(c, "Places", []navigation.SidebarSection{
		{Title: "WORLD", Items: []navigation.SidebarItem{{Icon: icon.Home, Label: "Town"}, {Icon: icon.Map, Label: "Field"}, {Icon: icon.Skull, Label: "Dungeon"}}},
		{Title: "INBOX", Items: []navigation.SidebarItem{{Icon: icon.Mail, Label: "Letters", Badge: "3"}}},
		{Title: "SYSTEM", Items: []navigation.SidebarItem{{Icon: icon.Settings, Label: "Options"}}},
	}, &s.side, navigation.SidebarOptions{}).Changed {
		log(fmt.Sprintf("item %d", s.side))
	}
}

func demoBreadcrumbs(c *ui.Context, log func(string)) {
	s := navState(c)
	if r := navigation.Breadcrumbs(c, "Path", s.crumbs, navigation.BreadcrumbsOptions{}); r.Clicked >= 0 {
		s.crumbs = s.crumbs[:r.Clicked+1]
		log("back to " + s.crumbs[r.Clicked])
	}
	if len(s.crumbs) < 4 && input.Button(c, "GO DEEPER", input.ButtonOptions{Size: input.Small}).Clicked {
		s.crumbs = append(s.crumbs, fmt.Sprintf("B%d", len(s.crumbs)-1))
		log("deeper")
	}
}

func demoCommandMenu(c *ui.Context, log func(string)) {
	items := []navigation.MenuItem{{Label: "FIGHT"}, {Label: "MAGIC", Disabled: true}, {Label: "ITEM"}, {Label: "RUN"}}
	r := navigation.CommandMenu(c, "Command", items, navigation.CommandMenuOptions{})
	if r.Chosen >= 0 {
		log(items[r.Chosen].Label)
	}
	if r.Canceled {
		log("back")
	}
}

func demoMenuBar(c *ui.Context, log func(string)) {
	menus := []navigation.Menu{
		{Title: "File", Items: []navigation.MenuItem{{Label: "New game"}, {Label: "Load", Icon: icon.FolderOpen}, {Label: "Save", Icon: icon.Save, Keys: "CTRL+S"}, {Separator: true}, {Label: "Quit"}}},
		{Title: "Edit", Items: []navigation.MenuItem{{Label: "Undo", Keys: "CTRL+Z"}, {Label: "Redo", Keys: "CTRL+Y", Disabled: true}}},
		{Title: "Help", Items: []navigation.MenuItem{{Label: "Controls"}, {Label: "About"}}},
	}
	if r := navigation.MenuBar(c, "Main menu", menus, navigation.MenuBarOptions{}); r.Item >= 0 {
		log(menus[r.Menu].Title + " › " + menus[r.Menu].Items[r.Item].Label)
	}
}

func demoDropdownMenu(c *ui.Context, log func(string)) {
	items := []navigation.MenuItem{{Label: "Swap", Icon: icon.Refresh}, {Label: "Equip", Icon: icon.Sword}, {Label: "Status", Icon: icon.User}, {Separator: true}, {Label: "Dismiss", Icon: icon.Trash, Disabled: true}}
	if r := navigation.DropdownMenu(c, "Party", items, navigation.DropdownMenuOptions{}); r.Chosen >= 0 {
		log(items[r.Chosen].Label)
	}
}

func demoContextMenu(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	items := []navigation.MenuItem{{Label: "Inspect", Icon: icon.Search}, {Label: "Pet", Icon: icon.Heart}, {Separator: true}, {Label: "Throw", Icon: icon.Bolt}}
	if r := navigation.ContextMenu(c, "Slime", items, navigation.ContextMenuOptions{}, func() {
		ui.Box(c).Size(k.D(160), k.D(40)).Label("Right click the slime").Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Dither(0, 0, g.W, g.H, pixel.Face, pixel.Good, 2)
			g.Text(8, 12, "Right click the slime", pixel.Faint, 1)
		})
	}); r.Chosen >= 0 {
		log(items[r.Chosen].Label)
	}
}

func demoPagination(c *ui.Context, log func(string)) {
	s := navState(c)
	if navigation.Pagination(c, "Pages", &s.page, 12, navigation.PaginationOptions{}).Changed {
		log(fmt.Sprintf("page %d", s.page))
	}
}

func demoPageDots(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := navState(c)
	tips := []string{"Jump with Space.", "Hold Shift to run.", "Talk to everyone.", "Check the pots.", "Save often."}
	ui.Column(c).Gap(k.D(4)).Children(func() {
		pixel.Text(c, tips[s.dot-1], pixel.TextOptions{})
		if navigation.PageDots(c, "Tips", &s.dot, len(tips), navigation.PageDotsOptions{}).Changed {
			log(fmt.Sprintf("tip %d", s.dot))
		}
	})
}

func demoSteps(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := navState(c)
	stops := []string{"Gate", "Forest", "Cave", "Peak", "Castle"}
	ui.Column(c).Gap(k.D(6)).Children(func() {
		navigation.Steps(c, stops, s.level, navigation.StepsOptions{})
		ui.Row(c).Gap(k.D(4)).Children(func() {
			if input.Button(c, "BACK", input.ButtonOptions{Size: input.Small, Disabled: s.level == 0}).Clicked {
				s.level--
				log("at " + stops[s.level])
			}
			if input.Button(c, "ONWARD", input.ButtonOptions{Size: input.Small, Disabled: s.level == len(stops)-1}).Clicked {
				s.level++
				log("at " + stops[s.level])
			}
		})
	})
}

func demoNavRail(c *ui.Context, log func(string)) {
	s := navState(c)
	items := []navigation.NavItem{{Icon: icon.Home, Label: "Home"}, {Icon: icon.Map, Label: "Map"}, {Icon: icon.Gem, Label: "Inventory"}, {Icon: icon.Settings, Label: "Options"}}
	if navigation.NavRail(c, "Rail", items, &s.rail, navigation.NavRailOptions{}).Changed {
		log(items[s.rail].Label)
	}
}

func demoBottomNav(c *ui.Context, log func(string)) {
	s := navState(c)
	items := []navigation.NavItem{{Icon: icon.Home, Label: "Home"}, {Icon: icon.Flag, Label: "Quests"}, {Icon: icon.Chest, Label: "Bag"}, {Icon: icon.User, Label: "Me"}}
	if navigation.BottomNav(c, "Bottom", items, &s.bottom, navigation.BottomNavOptions{}).Changed {
		log(items[s.bottom].Label)
	}
}

func demoCommandPalette(c *ui.Context, log func(string)) {
	s := navState(c)
	cmds := []navigation.MenuItem{{Label: "Save game", Keys: "CTRL+S"}, {Label: "Load game"}, {Label: "Toggle music"}, {Label: "Change palette", Keys: "P"}, {Label: "Rename hero"}, {Label: "Delete save", Disabled: true}}
	if input.Button(c, "OPEN (CMD+K)", input.ButtonOptions{Size: input.Small}).Clicked {
		s.palette = true
		log("palette open")
	}
	if r := navigation.CommandPalette(c, "Commands", cmds, &s.palette, navigation.CommandPaletteOptions{}); r.Ran >= 0 {
		log("ran " + cmds[r.Ran].Label)
	}
}

func demoBackForward(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := navState(c)
	ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		if navigation.BackForward(c, s.router, navigation.BackForwardOptions{}).Moved {
			log("now at " + s.router.Path())
		}
		pixel.Text(c, s.router.Path(), pixel.TextOptions{Slot: pixel.Accent})
		for _, place := range []string{"/castle", "/shop"} {
			if input.Button(c, place, input.ButtonOptions{Size: input.Small}).Clicked {
				s.router.Push(place)
				log("went to " + place)
			}
		}
	})
}

func demoTableOfContents(c *ui.Context, log func(string)) {
	s := navState(c)
	items := []navigation.TocItem{{Title: "Getting started"}, {Title: "Controls", Level: 1}, {Title: "Saving", Level: 1}, {Title: "Combat"}, {Title: "Magic", Level: 1}, {Title: "Credits"}}
	if navigation.TableOfContents(c, "Manual", items, &s.toc, navigation.TableOfContentsOptions{}).Changed {
		log(items[s.toc].Title)
	}
}

func demoLevelSelect(c *ui.Context, log func(string)) {
	s := navState(c)
	stages := []navigation.Stage{
		{Name: "1-1", Cleared: true, Stars: 3}, {Name: "1-2", Cleared: true, Stars: 2}, {Name: "1-3", Cleared: true, Stars: 1},
		{Name: "1-4"}, {Name: "1-5", Locked: true}, {Name: "2-1", Locked: true}, {Name: "2-2", Locked: true},
		{Name: "2-3", Locked: true}, {Name: "2-4", Locked: true}, {Name: "BOSS", Locked: true},
	}
	r := navigation.LevelSelect(c, "Stages", stages, &s.level, navigation.LevelSelectOptions{})
	if r.Changed {
		log("on " + stages[s.level].Name)
	}
	if r.Chosen {
		log("enter " + stages[s.level].Name)
	}
}

func demoWorldMap(c *ui.Context, log func(string)) {
	s := navState(c)
	nodes := []navigation.MapNode{{Name: "Home", X: 24, Y: 100}, {Name: "Woods", X: 84, Y: 70}, {Name: "Lake", X: 150, Y: 100}, {Name: "Tower", X: 120, Y: 28}, {Name: "Lair", X: 206, Y: 40, Locked: true}}
	paths := [][2]int{{0, 1}, {1, 2}, {1, 3}, {3, 4}, {2, 4}}
	r := navigation.WorldMap(c, "World", nodes, paths, &s.node, navigation.WorldMapOptions{})
	if r.Changed {
		log("walked to " + nodes[s.node].Name)
	}
	if r.Entered {
		log("entered " + nodes[s.node].Name)
	}
}

func demoToolbar(c *ui.Context, log func(string)) {
	s := navState(c)
	navigation.Toolbar(c, "Tools", navigation.ToolbarOptions{}, func() {
		if input.Toggle(c, "", &s.pen, input.ToggleOptions{Icon: icon.Edit, Label: "Pen"}).Changed {
			log(fmt.Sprintf("pen %v", s.pen))
		}
		if input.Toggle(c, "", &s.fill, input.ToggleOptions{Icon: icon.Drop, Label: "Fill"}).Changed {
			log(fmt.Sprintf("fill %v", s.fill))
		}
		if input.IconButton(c, icon.Prev, "Undo", input.IconButtonOptions{}).Clicked {
			log("undo")
		}
		if input.IconButton(c, icon.Trash, "Clear", input.IconButtonOptions{Variant: input.Danger}).Clicked {
			log("clear")
		}
	})
}

func demoDocumentTabs(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := navState(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		r := navigation.DocumentTabs(c, "Files", s.docs, &s.doc, navigation.DocumentTabsOptions{})
		if r.Changed {
			log("open " + s.docs[s.doc])
		}
		if r.Closed >= 0 {
			log("closed " + s.docs[r.Closed])
			s.docs = slices.Delete(s.docs, r.Closed, r.Closed+1)
			s.doc = max(0, min(s.doc, len(s.docs)-1))
		}
		if input.Button(c, "NEW TAB", input.ButtonOptions{Size: input.Small}).Clicked {
			s.newDocs++
			s.docs = append(s.docs, fmt.Sprintf("untitled%d.txt", s.newDocs))
			log("new tab")
		}
	})
}
