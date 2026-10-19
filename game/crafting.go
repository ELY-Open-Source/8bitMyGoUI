package game

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Recipe makes Result from item names in a 3×3 grid.
type Recipe struct {
	Grid   [9]string
	Result data.Item
}

type CraftingOptions struct {
	Disabled bool
}

type CraftingResult struct {
	Element ui.Element
	Changed bool
	// Picked is the cell chosen this frame, or -1.
	Picked  int
	Crafted bool
}

// Crafting is a 3×3 grid showing a known recipe's result.
func Crafting(c *ui.Context, text string, grid []data.Item, recipes []Recipe, selected *int, o CraftingOptions) CraftingResult {
	k := pixel.K(c)
	if len(grid) != 9 || *selected < 0 || *selected > 8 {
		panic(fmt.Sprintf("game: Crafting has %d cells, cell %d; it takes 9", len(grid), *selected))
	}
	var made data.Item
	for _, r := range recipes {
		if r.makes(grid) {
			made = r.Result
			break
		}
	}
	row := ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center)
	res := CraftingResult{Element: row, Picked: -1}
	row.Children(func() {
		old := *selected
		shell, base, st, picked := pixel.Pad(c, text, 3, 9, data.SlotSize, data.SlotSize, selected, o.Disabled)
		res.Changed, res.Picked = *selected != old, picked
		base.Description(grid[*selected].String())
		sel, shown := *selected, append([]data.Item(nil), grid...)
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, it := range shown {
				data.DrawSlot(g, i%3*data.SlotSize, i/3*data.SlotSize, it, data.ItemSlotOptions{Selected: i == sel})
			}
			g.Finish(st)
		})
		pixel.Text(c, "→", pixel.TextOptions{Slot: pixel.Faint})
		data.ItemSlot(c, made, data.ItemSlotOptions{})
		res.Crafted = input.Button(c, k.Msg("craft"), input.ButtonOptions{Variant: input.Primary, Disabled: o.Disabled || made.Sprite == nil}).Clicked
	})
	if res.Changed || res.Crafted {
		slog.Debug("8bit: crafting", "label", text, "cell", *selected, "crafted", res.Crafted, "item", made.Name)
	}
	return res
}

// makes is whether the grid lays out the recipe's names.
func (r Recipe) makes(grid []data.Item) bool {
	for i, name := range r.Grid {
		if grid[i].Name != name {
			return false
		}
	}
	return true
}
