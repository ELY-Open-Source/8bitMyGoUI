package game_test

import (
	"testing"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/game"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
)

func TestSpaceAfterEnterPicksTheCursor(t *testing.T) {
	gear := []data.Item{{Sprite: icon.Crown, Name: "Crown"}, {}, {Sprite: icon.Sword, Name: "Sword"}, {}, {}}
	sel, picks := game.Head, []int{}
	tt := view(400, 200, &world{clk: newClock()}, func(c *ui.Context) {
		if r := game.Equipment(c, "Gear", gear, &sel, game.EquipmentOptions{}); r.Picked >= 0 {
			picks = append(picks, r.Picked)
		}
	})
	r := find(t, tt, "Gear")
	// The pointer rests over the doll's lower row.
	tt.Move(r.X+r.W/2, r.Y+2*(26+11))
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeySpace)
	if len(picks) != 2 || picks[0] != game.Head || picks[1] != game.Head || sel != game.Head {
		t.Fatalf("Enter then Space picked %v, cursor %d", picks, sel)
	}
}

func TestSpaceLearnsTheCursorNotThePointer(t *testing.T) {
	skills := []game.Skill{
		{Name: "Root", Sprite: icon.Star},
		{Name: "Slash", Sprite: icon.Sword, At: game.Spot{X: 1}, Needs: []int{0}},
	}
	cur, points, learned := 0, 2, -1
	tt := view(400, 200, &world{clk: newClock()}, func(c *ui.Context) {
		if r := game.SkillTree(c, "Skills", skills, &cur, &points, game.SkillTreeOptions{}); r.Learned >= 0 {
			learned = r.Learned
		}
	})
	r := find(t, tt, "Skills")
	// The pointer rests on Slash; the cursor is on Root.
	tt.Move(r.X+2*(28+10), r.Y+2*(20+10))
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeySpace)
	if learned != 0 || cur != 0 {
		t.Fatalf("Space learned %d, cursor %d", learned, cur)
	}
}

func TestHotbarWheelAsAVeilFallsMovesNothing(t *testing.T) {
	wd, slot := &world{clk: newClock()}, 0
	tt := view(400, 200, wd, func(c *ui.Context) {
		game.Hotbar(c, "Slots", []data.Item{{Name: "A"}, {Name: "B"}}, &slot, game.HotbarOptions{})
	})
	r := find(t, tt, "Slots")
	wd.inert = true
	tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 40)
	if slot != 0 {
		t.Fatalf("the wheel moved a veiled bar to %d", slot)
	}
	wd.inert = false
	tt.Frame()
	tt.Scroll(r.X+r.W/2, r.Y+r.H/2, 0, 40)
	if slot != 1 {
		t.Fatalf("the wheel did not move the bar: %d", slot)
	}
}
