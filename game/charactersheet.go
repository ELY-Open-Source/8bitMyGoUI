package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

// Stat is a named number on a sheet.
type Stat struct {
	Name  string
	Value int
}

// Character is who a character sheet describes.
type Character struct {
	Name, Class string
	Level       int
	Portrait    *pixel.Sprite
	Stats       []Stat
}

type CharacterSheetOptions struct{}

// CharacterSheet shows a portrait, name, class, level and stats.
func CharacterSheet(c *ui.Context, h Character, o CharacterSheetOptions) *ui.Element {
	k := pixel.K(c)
	if h.Portrait == nil || h.Level < 1 {
		panic(fmt.Sprintf("game: character %q needs a portrait and a level (%d)", h.Name, h.Level))
	}
	name, value := 0, 0
	stats := make([]string, len(h.Stats))
	for i, s := range h.Stats {
		name, value = max(name, pixel.Measure(s.Name)), max(value, pixel.Measure(strconv.Itoa(s.Value)))
		stats[i] = s.Name + " " + strconv.Itoa(s.Value)
	}
	return ui.Row(c).Gap(k.D(8)).AlignItems(ui.Start).Role(ui.RoleGroup).Label(h.Name).Children(func() {
		ui.Box(c).Size(k.D(40), k.D(40)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, 40, 40, pixel.Sunken, pixel.Well)
			g.Sprite(2, 2, h.Portrait, 0, pixel.Ink, 3)
		})
		ui.Column(c).Gap(k.D(3)).Children(func() {
			text.ShadowText(c, h.Name, text.ShadowTextOptions{Scale: 2})
			pixel.Text(c, h.Class+"  "+fmt.Sprintf(k.Msg("level"), h.Level), pixel.TextOptions{Slot: pixel.Faint})
			// Stats stand in a column: names Faint, values right-aligned.
			w := name + 8 + value
			ui.Box(c).Role(ui.RoleText).Label(strings.Join(stats, ", ")).Size(k.D(w), k.D(len(h.Stats)*pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				for i, s := range h.Stats {
					n := strconv.Itoa(s.Value)
					g.Text(0, i*pixel.LineHeight, s.Name, pixel.Faint, 1)
					g.Text(w-pixel.Measure(n), i*pixel.LineHeight, n, pixel.Ink, 1)
				}
			})
		})
	})
}
