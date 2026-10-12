package gallery

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

var dataEntries = []Entry{
	{151, "List", "列表", "data", "Thousands of rows, built as they scroll in. ▶ marks the chosen one.", demoList},
	{152, "Table", "表格", "data", "Click a header to sort ▲, again for ▼. Odd rows are dithered.", demoTable},
	{153, "Tree", "树", "data", "▸ opens, ▾ closes. Right opens, Left climbs back out.", demoTree},
	{154, "Outline", "大纲", "data", "A tree of 2,000 floors that only builds what shows.", demoOutline},
	{155, "GridView", "网格视图", "data", "Whole-dot cells in as many columns as fit; arrows go both ways.", demoGridView},
	{156, "DescriptionList", "描述列表", "data", "Terms on the left, what they mean on the right.", demoDescriptionList},
	{157, "Leaderboard", "排行榜", "data", "Medals for the top three; your row lights up.", demoLeaderboard},
	{158, "Inventory", "物品栏", "data", "Arrows pick a slot. Drag one item onto another to swap them.", demoInventory},
	{159, "ItemSlot", "物品格", "data", "An item with its stack, rarity edge and cooldown veil.", demoItemSlot},
	{160, "Accordion", "手风琴", "data", "Open one answer and the others fold away.", demoAccordion},
	{161, "Collapsible", "折叠区", "data", "One header that shows and hides what is under it.", demoCollapsible},
	{162, "Carousel", "轮播", "data", "◀ ▶ turn the pages; it turns by itself unless you point at it.", demoCarousel},
	{163, "Timeline", "时间线", "data", "Dated events down a line of markers.", demoTimeline},
}

type dataDemo struct {
	list, row, tile, slot, page, size, carousel, hit int
	sort                                             data.Sort
	tree, floor                                      string
	faq                                              []bool
	secret                                           bool
	quests, bench, party                             []string
	bag                                              []data.Item
	rights                                           [][]bool
	filters                                          []bool
	name                                             string
	level                                            float64
	hard                                             bool
	logs                                             int
	loading                                          time.Time
	query                                            string
}

func dataState(c *ui.Context) *dataDemo {
	return ui.Local(c.Root(), demoKey("data"), func() dataDemo {
		return dataDemo{size: 20, tree: "0", floor: "R1", faq: []bool{true, false, false},
			quests: []string{"Find the key", "Save the cat", "Beat the boss", "Buy a potion"},
			bench:  []string{"Dax", "Eko", "Fay"}, party: []string{"Aria", "Brom"},
			bag: []data.Item{{Sprite: icon.Potion, Name: "Potion", Count: 5}, {Sprite: icon.Drop, Name: "Ether", Count: 2, Rarity: display.Uncommon},
				{Sprite: icon.Sword, Name: "Sword", Rarity: display.Rare}, {Sprite: icon.Shield, Name: "Shield", Rarity: display.Epic},
				{Sprite: icon.Crown, Name: "Crown", Rarity: display.Legendary}, {Sprite: icon.Gem, Name: "Gem", Count: 12, Rarity: display.Rare},
				{}, {Sprite: icon.Key, Name: "Key"}, {}, {}, {Sprite: icon.Scroll, Name: "Map scroll"}, {}},
			rights:  [][]bool{{true, false, false}, {false, true, false}, {true, true, true}},
			filters: []bool{false, true, false, false}, name: "Aria", level: 7, logs: 20}
	})
}

var horde = strings.Fields("Slime Bat Goblin Ghost Golem Imp Kobold Lich Mimic Ogre Orc Pixie Rat Skeleton Snake Spider Troll Wisp Wolf Wyrm Zombie Harpy Hydra Kraken Basilisk Djinn Gargoyle Manticore Wyvern Yeti")

func demoList(c *ui.Context, log func(string)) {
	s := dataState(c)
	r := data.List(c, "Monsters", len(horde), &s.list, data.ListOptions{}, func(i int) {
		pixel.Text(c, horde[i], pixel.TextOptions{})
	})
	r.Element.Width(pixel.K(c).D(140))
	if r.Changed {
		log("chose " + horde[s.list])
	}
	if r.Submitted {
		log("opened " + horde[s.list])
	}
}

var roster = []struct {
	name, class string
	lv          int
}{{"Aria", "Knight", 7}, {"Brom", "Monk", 5}, {"Cyra", "Mage", 9}, {"Dax", "Thief", 4}, {"Eko", "Bard", 6},
	{"Fay", "Cleric", 8}, {"Gus", "Ranger", 3}, {"Hal", "Knight", 10}, {"Ivy", "Mage", 2}, {"Jun", "Monk", 11}}

func demoTable(c *ui.Context, log func(string)) {
	s := dataState(c)
	order := make([]int, len(roster))
	for i := range order {
		order[i] = i
	}
	slices.SortStableFunc(order, func(a, b int) int {
		x, y := roster[a], roster[b]
		d := [3]int{cmp.Compare(x.name, y.name), cmp.Compare(x.class, y.class), cmp.Compare(x.lv, y.lv)}[s.sort.Column]
		if s.sort.Desc {
			return -d
		}
		return d
	})
	cols := []data.Column{{Title: "Name", Width: 50}, {Title: "Class", Width: 56}, {Title: "LV", Width: 30}}
	r := data.Table(c, "Roster", cols, len(roster), &s.sort, &s.row, data.TableOptions{}, func(row, col int) string {
		h := roster[order[row]]
		return [3]string{h.name, h.class, fmt.Sprint(h.lv)}[col]
	})
	if r.Sorted {
		log(fmt.Sprintf("sorted by %s, desc %v", cols[s.sort.Column].Title, s.sort.Desc))
	}
	if r.Changed {
		log("chose row " + fmt.Sprint(s.row))
	}
	if r.Submitted {
		log("opened " + roster[order[s.row]].name)
	}
}

func demoTree(c *ui.Context, log func(string)) {
	s := dataState(c)
	roots := []data.TreeNode{
		{Label: "Bag", Children: []data.TreeNode{{Label: "Potions", Children: []data.TreeNode{{Label: "Potion ×5"}, {Label: "Ether ×2"}}}, {Label: "Tent"}}},
		{Label: "Key items", Children: []data.TreeNode{{Label: "Old key"}, {Label: "Map scroll"}}},
		{Label: "Gold: 120"},
	}
	r := data.Tree(c, "Items", roots, &s.tree, data.TreeOptions{})
	if r.Changed {
		log("chose " + s.tree)
	}
	if r.Toggled {
		log("opened or closed " + s.tree)
	}
}

func demoOutline(c *ui.Context, log func(string)) {
	s := dataState(c)
	roots := make([]string, 8)
	for i := range roots {
		roots[i] = fmt.Sprintf("R%d", i+1)
	}
	children := func(k string) []string {
		n := map[int]int{0: 12, 1: 20}[strings.Count(k, "/")]
		var out []string
		for i := range n {
			out = append(out, fmt.Sprintf("%s/%d", k, i+1))
		}
		return out
	}
	label := func(k string) string {
		parts := strings.Split(k, "/")
		name := []string{"Region " + parts[0][1:], "Dungeon ", "Floor B"}[len(parts)-1]
		if len(parts) > 1 {
			name += parts[len(parts)-1]
		}
		return name
	}
	r := data.Outline(c, "World", roots, children, label, &s.floor, data.OutlineOptions{})
	r.Element.Width(pixel.K(c).D(160))
	if r.Changed {
		log("chose " + label(s.floor))
	}
	if r.Submitted {
		log("entered " + label(s.floor))
	}
}

var tiles = []*pixel.Sprite{icon.Coin, icon.Gem, icon.Potion, icon.Sword, icon.Shield, icon.Chest, icon.Crown, icon.Scroll, icon.Flame, icon.Leaf}

func demoGridView(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := dataState(c)
	r := data.GridView(c, "Loot", 200, &s.tile, data.GridViewOptions{CellWidth: 24, CellHeight: 24, Height: 100}, func(i int) {
		ui.Box(c).Grow(1).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
			display.Icon(c, tiles[i%len(tiles)], display.IconOptions{Label: fmt.Sprintf("loot %d", i)})
		})
	})
	r.Element.Width(k.D(200))
	if r.Changed {
		log(fmt.Sprintf("chose loot %d", s.tile))
	}
	if r.Submitted {
		log(fmt.Sprintf("took loot %d", s.tile))
	}
}

func demoDescriptionList(c *ui.Context, log func(string)) {
	ui.Box(c).Width(pixel.K(c).D(220)).Children(func() {
		data.DescriptionList(c, [][2]string{{"Class", "Knight"}, {"Weapon", "Rusty sword"},
			{"Home", "Castle Town, two screens west of the cave"}, {"Motto", "It is dangerous to go alone."}}, data.DescriptionListOptions{TermWidth: 50})
	})
}

func demoLeaderboard(c *ui.Context, log func(string)) {
	scores := []data.Score{{Name: "ZAK", Score: 98400}, {Name: "MOM", Score: 87150}, {Name: "AAA", Score: 80000}, {Name: "BOB", Score: 64200}, {Name: "YOU", Score: 51300}, {Name: "KEN", Score: 32000}}
	ui.Box(c).Width(pixel.K(c).D(160)).Children(func() {
		data.Leaderboard(c, scores, data.LeaderboardOptions{Player: 4})
	})
}

func demoInventory(c *ui.Context, log func(string)) {
	s := dataState(c)
	r := data.Inventory(c, "Bag", s.bag, &s.slot, data.InventoryOptions{})
	if r.Swapped {
		log("swapped slots")
	}
	if r.Changed {
		log(fmt.Sprintf("slot %d: %s", s.slot, s.bag[s.slot].Name))
	}
}

func demoItemSlot(c *ui.Context, log func(string)) {
	ui.Row(c).Gap(pixel.K(c).D(6)).Children(func() {
		data.ItemSlot(c, data.Item{Sprite: icon.Potion, Name: "Potion", Count: 3}, data.ItemSlotOptions{})
		data.ItemSlot(c, data.Item{Sprite: icon.Crown, Name: "Crown", Rarity: display.Legendary}, data.ItemSlotOptions{Selected: true})
		data.ItemSlot(c, data.Item{Sprite: icon.Bolt, Name: "Bolt scroll", Rarity: display.Rare}, data.ItemSlotOptions{Cooldown: 0.5})
		data.ItemSlot(c, data.Item{Name: "Empty"}, data.ItemSlotOptions{})
	})
}

func demoAccordion(c *ui.Context, log func(string)) {
	s := dataState(c)
	qs := []string{"How do I save?", "Where is the key?", "Can I pet the dog?"}
	as := []string{"Visit any inn, or press START.", "Under the third torch.", "Yes. Always yes."}
	ui.Box(c).Width(pixel.K(c).D(180)).Children(func() {
		if data.Accordion(c, "Help", qs, s.faq, data.AccordionOptions{}, func(i int) {
			text.Paragraph(c, as[i], text.ParagraphOptions{})
		}).Changed {
			log(fmt.Sprintf("open %v", s.faq))
		}
	})
}

func demoCollapsible(c *ui.Context, log func(string)) {
	s := dataState(c)
	if data.Collapsible(c, "Secret moves", &s.secret, data.CollapsibleOptions{}, func() {
		pixel.Text(c, "↑ ↑ ↓ ↓ ← → ← → B A", pixel.TextOptions{})
	}).Changed {
		log(fmt.Sprintf("open %v", s.secret))
	}
}

func demoCarousel(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := dataState(c)
	tips := []struct {
		art  *pixel.Sprite
		line string
	}{{icon.Heart, "Hearts heal you."}, {icon.Key, "Keys open doors."}, {icon.Coin, "Coins buy things."}}
	if data.Carousel(c, "Tips", len(tips), &s.carousel, data.CarouselOptions{Every: 4}, func(i int) {
		ui.Column(c).Width(k.D(120)).AlignItems(ui.Center).Gap(k.D(4)).Children(func() {
			display.Icon(c, tips[i].art, display.IconOptions{Scale: 2})
			pixel.Text(c, tips[i].line, pixel.TextOptions{})
		})
	}).Changed {
		log(fmt.Sprintf("page %d", s.carousel+1))
	}
}

func demoTimeline(c *ui.Context, log func(string)) {
	ui.Box(c).Width(pixel.K(c).D(200)).Children(func() {
		data.Timeline(c, []data.Event{{Date: "1985", Text: "A plumber jumps into a pipe."}, {Date: "1986", Text: "A boy in green finds a sword."},
			{Date: "1987", Text: "A last fantasy begins."}, {Date: "1991", Text: "A blue hedgehog runs very fast."}}, data.TimelineOptions{})
	})
}
