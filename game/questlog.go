package game

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

type QuestLogOptions struct {
	Disabled bool
}

type QuestLogResult struct {
	Element ui.Element
	Changed bool
}

// QuestLog lists quests; the chosen one is tracked and detailed.
func QuestLog(c *ui.Context, text string, quests []Quest, tracked *int, o QuestLogOptions) QuestLogResult {
	k := pixel.K(c)
	if len(quests) == 0 || *tracked < 0 || *tracked >= len(quests) {
		panic(fmt.Sprintf("game: QuestLog tracks %d of %d", *tracked, len(quests)))
	}
	for _, q := range quests {
		q.check()
	}
	row := ui.Row(c).Gap(k.D(10)).AlignItems(ui.Start)
	res := QuestLogResult{Element: row}
	row.Children(func() {
		r := data.List(c, text, len(quests), tracked, data.ListOptions{Height: 6 * pixel.LineHeight, Disabled: o.Disabled}, func(i int) string { return quests[i].Title }, func(i int) {
			ui.Row(c).Gap(k.D(3)).AlignItems(ui.Center).Children(func() {
				mark, ink := icon.Flag, pixel.Accent
				switch {
				case i == *tracked:
				case quests[i].done():
					mark, ink = icon.Check, pixel.Good
				default:
					mark = nil
				}
				ui.Box(c).Size(k.D(12), k.D(12)).Draw(func(p *ui.Painter, r ui.Rect) {
					if mark != nil {
						k.Canvas(p, r).Sprite(0, 0, mark, 0, ink, 1)
					}
				})
				slot := pixel.Ink
				if quests[i].done() {
					slot = pixel.Faint
				}
				pixel.Text(c, quests[i].Title, pixel.TextOptions{Slot: slot})
			})
		})
		r.Element.Width(k.D(128))
		res.Changed = r.Changed
		ui.Column(c).Gap(k.D(2)).Children(func() {
			pixel.Text(c, k.Msg("tracking"), pixel.TextOptions{Slot: pixel.Faint})
			QuestTracker(c, quests[*tracked], QuestTrackerOptions{})
		})
	})
	if res.Changed {
		slog.Debug("8bit: quest tracked", "label", text, "quest", quests[*tracked].Title)
	}
	return res
}
