package gallery

import (
	"fmt"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/layout"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
	"github.com/ELY-Open-Source/8bitMyGoUI/text"
)

var recordEntries = []Entry{
	{164, "StatCompare", "属性对比", "data", "What a new sword does to you: ▲ in green, ▼ in red.", demoStatCompare},
	{165, "Bestiary", "怪物图录", "data", "Monsters you met, and ??? shadows of those you have not.", demoBestiary},
	{166, "PropertyGrid", "属性表", "data", "Edit text, numbers and switches row by row.", demoPropertyGrid},
	{167, "Pager", "分页器", "data", "Pick a page size; ◀ ▶ step through 312 items.", demoPager},
	{168, "Feed", "动态流", "data", "What happened, newest first, with icons and times.", demoFeed},
	{169, "Matrix", "矩阵", "data", "Who knows which spell. Arrows move; Space flips.", demoMatrix},
	{170, "Sortable", "可排序列表", "data", "Drag a quest by its grip, or move it with Alt and the arrows.", demoSortable},
	{171, "Transfer", "双栏选择", "data", "Pick a hero, then → or ← to swap benches.", demoTransfer},
	{172, "InfiniteList", "无限列表", "data", "Scroll to the end; it loads twenty more.", demoInfiniteList},
	{173, "FilterBar", "筛选栏", "data", "Chips turn filters on; CLEAR turns them all off.", demoFilterBar},
	{174, "SearchResults", "搜索结果", "data", "Type a spell name; matches light up.", demoSearchResults},
	{175, "CardGrid", "卡片网格", "data", "As many columns of cards as the width holds.", demoCardGrid},
}

func demoStatCompare(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		pixel.Text(c, "Iron sword → Flame sword", pixel.TextOptions{Slot: pixel.Faint})
		data.StatCompare(c, []data.Change{{Label: "ATK", Old: 12, New: 18}, {Label: "DEF", Old: 9, New: 7}, {Label: "SPD", Old: 5, New: 5}, {Label: "MAG", Old: 3, New: 6}}, data.StatCompareOptions{})
	})
}

func demoBestiary(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(6)).Children(func() {
		ui.Box(c).Width(k.D(130)).Children(func() {
			data.Bestiary(c, data.Monster{Sprite: slime, Name: "Slime", Seen: true, Stats: [][2]string{{"HP", "8"}, {"ATK", "2"}},
				Notes: "Wobbles happily. Mostly harmless."}, data.BestiaryOptions{})
		})
		ui.Box(c).Width(k.D(130)).Children(func() {
			data.Bestiary(c, data.Monster{Sprite: icon.Bug, Name: "Cave beetle", Stats: [][2]string{{"HP", "30"}, {"ATK", "9"}},
				Notes: "Clicks in the dark."}, data.BestiaryOptions{Scale: 3})
		})
	})
}

func demoPropertyGrid(c *ui.Context, log func(string)) {
	s := dataState(c)
	props := []data.Property{{Name: "Name", Text: &s.name}, {Name: "Level", Number: &s.level, Min: 1, Max: 99}, {Name: "Hard mode", Switch: &s.hard}}
	if data.PropertyGrid(c, "Hero", props, data.PropertyGridOptions{NameWidth: 56}).Changed {
		log(fmt.Sprintf("%s, level %g, hard %v", s.name, s.level, s.hard))
	}
}

func demoPager(c *ui.Context, log func(string)) {
	s := dataState(c)
	if data.Pager(c, "Items", 312, &s.page, &s.size, data.PagerOptions{}).Changed {
		log(fmt.Sprintf("page %d, %d per page", s.page+1, s.size))
	}
}

func demoFeed(c *ui.Context, log func(string)) {
	ui.Box(c).Width(pixel.K(c).D(200)).Children(func() {
		data.Feed(c, []data.Activity{{Icon: icon.Star, Text: "Aria found a star shard.", Time: "1 min ago"}, {Icon: icon.Coin, Text: "Brom sold three slime jellies.", Time: "4 min ago"},
			{Icon: icon.Trophy, Text: "You cleared World 1-2.", Time: "9 min ago"}, {Icon: icon.Heart, Text: "Cyra healed the party.", Time: "12 min ago"}}, data.FeedOptions{})
	})
}

func demoMatrix(c *ui.Context, log func(string)) {
	s := dataState(c)
	heroes, spells := []string{"Aria", "Brom", "Cyra"}, []string{"FIRE", "ICE", "BOLT"}
	if data.Matrix(c, "Spells known", heroes, spells, s.rights, data.MatrixOptions{}).Changed {
		log(fmt.Sprintf("%v", s.rights))
	}
}

func demoSortable(c *ui.Context, log func(string)) {
	s := dataState(c)
	if data.Sortable(c, "Quest order", &s.quests, data.SortableOptions{}).Changed {
		log(strings.Join(s.quests, " · "))
	}
}

func demoTransfer(c *ui.Context, log func(string)) {
	s := dataState(c)
	if data.Transfer(c, "Party", &s.bench, &s.party, data.TransferOptions{Height: 70}).Changed {
		log("party: " + strings.Join(s.party, ", "))
	}
}

func demoInfiniteList(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := dataState(c)
	// A pretend server answers after one second.
	if !s.loading.IsZero() && k.Since(c, s.loading, 4, 4) == 4 {
		s.logs, s.loading = s.logs+20, time.Time{}
		log(fmt.Sprintf("loaded %d rows", s.logs))
	}
	r := data.InfiniteList(c, "Battle log", s.logs, data.InfiniteListOptions{Height: 100, Done: s.logs >= 100}, func(i int) {
		pixel.Text(c, fmt.Sprintf("Turn %d: the slime wobbles.", i+1), pixel.TextOptions{})
	})
	r.Element.Width(k.D(180))
	if r.More {
		s.loading = k.Now()
		log("asked for more")
	}
}

func demoFilterBar(c *ui.Context, log func(string)) {
	s := dataState(c)
	labels := []string{"FIRE", "ICE", "BOLT", "HOLY"}
	if data.FilterBar(c, "Elements", labels, s.filters, data.FilterBarOptions{}).Changed {
		log(fmt.Sprintf("%v", s.filters))
	}
}

var spellBook = []data.Hit{{Title: "Fire", Text: "A small burst of fire."}, {Title: "Fira", Text: "A bigger fire. Warm."}, {Title: "Blizzard", Text: "Ice falls from nowhere."},
	{Title: "Thunder", Text: "Bolts strike one foe."}, {Title: "Cure", Text: "Heals a little."}, {Title: "Firaga", Text: "All the fire at once."}}

func demoSearchResults(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := dataState(c)
	var hits []data.Hit
	for _, h := range spellBook {
		if strings.Contains(strings.ToLower(h.Title+" "+h.Text), strings.ToLower(s.query)) {
			hits = append(hits, h)
		}
	}
	s.hit = min(s.hit, len(hits)-1)
	ui.Column(c).Gap(k.D(4)).Width(k.D(180)).Children(func() {
		if input.SearchField(c, "Search spells", &s.query, input.SearchFieldOptions{Placeholder: "fire"}).Changed {
			log("query " + s.query)
		}
		r := data.SearchResults(c, "Spells", s.query, hits, &s.hit, data.SearchResultsOptions{Height: 90})
		if r.Changed {
			log("chose " + hits[s.hit].Title)
		}
		if r.Submitted {
			log("cast " + hits[s.hit].Title)
		}
	})
}

func demoCardGrid(c *ui.Context, log func(string)) {
	cards := []struct {
		art        *pixel.Sprite
		name, line string
	}{{icon.Sword, "Sword", "Pointy end first."}, {icon.Shield, "Shield", "Blocks most things."}, {icon.Potion, "Potion", "Tastes like cherry."},
		{icon.Crown, "Crown", "Heavy is the head."}, {icon.Map, "Map", "X marks the spot."}, {icon.Key, "Key", "Opens one door."}}
	data.CardGrid(c, len(cards), data.CardGridOptions{MinWidth: 90}, func(i int) {
		layout.Card(c, layout.CardOptions{Title: cards[i].name, Icon: cards[i].art}, func() {
			text.Paragraph(c, cards[i].line, text.ParagraphOptions{})
		})
	})
}
