package game_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/game"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var bag = []data.Item{{Sprite: icon.Potion, Name: "Potion", Count: 5}, {Sprite: icon.Sword, Name: "Sword"}, {Sprite: icon.Bolt, Name: "Bolt"}, {}, {Sprite: icon.Key, Name: "Key"}}

func TestHotbarKeysWheelAndClicks(t *testing.T) {
	sel, changes, used := 0, 0, -1
	off := false
	wd := &world{clk: newClock()}
	tt := view(400, 100, wd, func(c *ui.Context) {
		r := game.Hotbar(c, "Hotbar", bag, &sel, game.HotbarOptions{Disabled: off})
		if r.Changed {
			changes++
		}
		if r.Used >= 0 {
			used = r.Used
		}
	})
	// Slot i is 22 dots wide, its middle at 22i+11.
	tt.ClickAt(2*(22*2+11), 22)
	if sel != 2 || used != 2 || changes != 1 {
		t.Fatalf("click: slot %d, used %d, %d changes", sel, used, changes)
	}
	tt.Key(0, ui.Key5)
	tt.Key(0, ui.Key9)
	tt.Key(0, ui.KeyLeft)
	if sel != 3 || changes != 3 {
		t.Fatalf("5, 9, Left: slot %d, %d changes", sel, changes)
	}
	// A notch is a 16-dot line; the bar wraps.
	tt.Scroll(22, 22, 0, 32)
	tt.Scroll(22, 22, 0, 32)
	tt.Scroll(22, 22, 0, -16)
	if sel != 0 || changes != 5 {
		t.Fatalf("two notches and a half: slot %d, %d changes", sel, changes)
	}
	tt.Key(0, ui.KeyEnter)
	if used != 0 {
		t.Fatalf("Enter used %d", used)
	}
	keeps(t, tt.Image())
	off = true
	tt.Frame()
	tt.Key(0, ui.Key3)
	tt.Scroll(22, 22, 0, 32)
	if sel != 0 {
		t.Fatalf("a disabled bar moved to %d", sel)
	}
	panics(t, "ten slots", func() {
		view(400, 100, wd, func(c *ui.Context) {
			game.Hotbar(c, "x", make([]data.Item, 10), &sel, game.HotbarOptions{})
		})
	})
}

func TestInertHotbarIgnoresTheWheel(t *testing.T) {
	sel := 0
	wd := &world{clk: newClock(), inert: true}
	tt := view(400, 100, wd, func(c *ui.Context) { game.Hotbar(c, "Hotbar", bag, &sel, game.HotbarOptions{}) })
	tt.Scroll(22, 22, 0, 32)
	tt.ClickAt(2*(22*3+11), 22)
	if sel != 0 {
		t.Fatalf("an inert bar moved to %d", sel)
	}
}

func TestEquipmentWalksTheDoll(t *testing.T) {
	gear := []data.Item{{Sprite: icon.Crown, Name: "Crown"}, {}, {Sprite: icon.Sword, Name: "Iron sword"}, {}, {}}
	sel, picked := game.Weapon, -1
	tt := view(400, 200, &world{clk: newClock()}, func(c *ui.Context) {
		if r := game.Equipment(c, "Gear", gear, &sel, game.EquipmentOptions{}); r.Picked >= 0 {
			picked = r.Picked
		}
	})
	r := find(t, tt, "Gear")
	// The doll is centered under the widest caption.
	w := 3*26 - 4
	for _, s := range []string{"HEAD: Crown", "BODY", "WEAPON: Iron sword", "SHIELD", "RING"} {
		w = max(w, pixel.Measure(s))
	}
	left := (w - (3*26 - 4)) / 2
	click := func(x, y int) { tt.ClickAt(r.X+float32(2*(left+x*26+11)), r.Y+float32(2*(y*26+11))) }
	click(2, 1)
	if sel != game.Shield || picked != game.Shield {
		t.Fatalf("clicking the shield slot chose %d, picked %d", sel, picked)
	}
	// Left of the doll, nothing is hit.
	tt.ClickAt(r.X+float32(2*(left-5)), r.Y+float32(2*(26+11)))
	if sel != game.Shield {
		t.Fatalf("a click beside the doll chose %d", sel)
	}
	for _, step := range []struct {
		key  ui.Key
		want int
	}{{ui.KeyLeft, game.Body}, {ui.KeyUp, game.Head}, {ui.KeyLeft, game.Weapon}, {ui.KeyDown, game.Ring}, {ui.KeyUp, game.Weapon}, {ui.KeyRight, game.Body}, {ui.KeyDown, game.Ring}} {
		tt.Key(0, step.key)
		if sel != step.want {
			t.Fatalf("%v went to %d, not %d", step.key, sel, step.want)
		}
	}
	tt.Key(0, ui.KeySpace)
	if picked != game.Ring {
		t.Fatalf("Space picked %d", picked)
	}
	keeps(t, tt.Image())
	panics(t, "four slots", func() {
		view(400, 200, &world{clk: newClock()}, func(c *ui.Context) {
			game.Equipment(c, "x", gear[:4], &sel, game.EquipmentOptions{})
		})
	})
}

func TestSkillTreeLearnsOpenSkillsWithPoints(t *testing.T) {
	skills := []game.Skill{
		{Name: "Slash", Sprite: icon.Sword, At: game.Spot{X: 1}, Learned: true},
		{Name: "Guard", Sprite: icon.Shield, At: game.Spot{Y: 1}, Needs: []int{0}},
		{Name: "Fire", Sprite: icon.Flame, At: game.Spot{X: 2, Y: 1}, Needs: []int{0}, Learned: true},
		{Name: "Meteor", Sprite: icon.Star, At: game.Spot{X: 2, Y: 2}, Needs: []int{2, 1}},
	}
	cur, points, learned := 0, 1, -1
	tt := view(400, 200, &world{clk: newClock()}, func(c *ui.Context) {
		if r := game.SkillTree(c, "Skills", skills, &cur, &points, game.SkillTreeOptions{}); r.Learned >= 0 {
			learned = r.Learned
		}
	})
	r := find(t, tt, "Skills")
	// Nodes sit 28 dots apart under a 20-dot header.
	node := func(x, y int) { tt.ClickAt(r.X+float32(2*(x*28+10)), r.Y+float32(2*(20+y*28+10))) }
	node(1, 0)
	tt.Key(0, ui.KeyEnter)
	if points != 1 {
		t.Fatal("a learned skill took a point")
	}
	node(2, 2)
	tt.Key(0, ui.KeyEnter)
	if cur != 3 || skills[3].Learned || points != 1 {
		t.Fatal("a locked skill was learned")
	}
	node(0, 1)
	if cur != 1 || skills[1].Learned {
		t.Fatal("the first click on a skill learned it")
	}
	node(0, 1)
	if !skills[1].Learned || points != 0 || learned != 1 {
		t.Fatalf("the second click: learned %d, %d points", learned, points)
	}
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if cur != 3 || skills[3].Learned {
		t.Fatal("Meteor was learned with no points")
	}
	points = 1
	tt.Key(0, ui.KeyEnter)
	if !skills[3].Learned || points != 0 || learned != 3 {
		t.Fatal("Meteor was not learned with a point")
	}
	keeps(t, tt.Image())
	first := 0
	panics(t, "skills needing each other", func() {
		view(400, 200, &world{clk: newClock()}, func(c *ui.Context) {
			game.SkillTree(c, "x", []game.Skill{{Name: "a", Sprite: icon.Star, Needs: []int{1}}, {Name: "b", Sprite: icon.Star, At: game.Spot{X: 1}, Needs: []int{0}}}, &first, &points, game.SkillTreeOptions{})
		})
	})
	panics(t, "a skill off the grid", func() {
		view(400, 200, &world{clk: newClock()}, func(c *ui.Context) {
			game.SkillTree(c, "x", []game.Skill{{Name: "a", Sprite: icon.Star, At: game.Spot{X: -1}}}, &first, &points, game.SkillTreeOptions{})
		})
	})
	panics(t, "a skill needing itself", func() {
		view(400, 200, &world{clk: newClock()}, func(c *ui.Context) {
			game.SkillTree(c, "x", []game.Skill{{Name: "a", Sprite: icon.Star, Needs: []int{0}}}, &cur, &points, game.SkillTreeOptions{})
		})
	})
}

func TestSkillButtonFiresOnlyWhenReady(t *testing.T) {
	cool, clicks := 0.5, 0
	tt := view(200, 100, &world{clk: newClock()}, func(c *ui.Context) {
		if game.SkillButton(c, "Fireball", icon.Flame, "Q", cool, game.SkillButtonOptions{}).Clicked {
			clicks++
		}
	})
	// The wipe darkens the left half for half a cooldown.
	wiped := count(t, tt.Image(), 2, 2, 9, 20, pixel.Edge)
	if err := tt.Click("Fireball"); err != nil || clicks != 0 {
		t.Fatalf("a cooling skill fired: %v", err)
	}
	cool = 0
	tt.Frame()
	if n := count(t, tt.Image(), 2, 2, 9, 20, pixel.Edge); wiped < n+50 {
		t.Fatalf("wipe dots %d, ready %d", wiped, n)
	}
	tt.Click("Fireball")
	if clicks != 1 {
		t.Fatal("a ready skill did not fire")
	}
}
