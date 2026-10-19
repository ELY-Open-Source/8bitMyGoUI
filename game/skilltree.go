package game

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Skill is a node of a skill tree.
type Skill struct {
	Name   string
	Sprite *pixel.Sprite
	// At places it on the tree's grid.
	At Spot
	// Needs lists the skills that open it.
	Needs   []int
	Learned bool
}

type SkillTreeOptions struct {
	Disabled bool
}

type SkillTreeResult struct {
	Element ui.Element
	Changed bool
	// Learned is the skill learned this frame, or -1.
	Learned int
}

// node is a node's side; cellPitch adds the gap.
const node, cellPitch = 20, 28

// SkillTree links skills; points learn the cursor's skill.
func SkillTree(c *ui.Context, text string, skills []Skill, cursor, points *int, o SkillTreeOptions) SkillTreeResult {
	k := pixel.K(c)
	if *cursor < 0 || *cursor >= len(skills) || *points < 0 {
		panic(fmt.Sprintf("game: SkillTree cursor %d of %d, %d points", *cursor, len(skills), *points))
	}
	spots, cols, rows := make([]Spot, len(skills)), 0, 0
	for i, s := range skills {
		if s.Sprite == nil || s.At.X < 0 || s.At.Y < 0 {
			panic(fmt.Sprintf("game: skill %q needs a sprite and a place on the grid, not %v", s.Name, s.At))
		}
		for _, n := range s.Needs {
			if n < 0 || n >= len(skills) || n == i {
				panic(fmt.Sprintf("game: skill %q needs %d", s.Name, n))
			}
		}
		for j := range i {
			if skills[j].At == s.At {
				panic(fmt.Sprintf("game: skills %q and %q share %v", skills[j].Name, s.Name, s.At))
			}
		}
		spots[i], cols, rows = s.At, max(cols, s.At.X+1), max(rows, s.At.Y+1)
	}
	if loop := cycle(skills); loop >= 0 {
		panic(fmt.Sprintf("game: skill %q needs itself through others", skills[loop].Name))
	}
	open := func(i int) bool {
		for _, n := range skills[i].Needs {
			if !skills[n].Learned {
				return false
			}
		}
		return !skills[i].Learned
	}
	old := *cursor
	shell, base, st := pixel.Shell(c, func() ui.Element {
		return ui.Box(c).Focusable().Role(ui.RoleGroup).Label(text).Disabled(o.Disabled)
	})
	arrows(base, spots, cursor)
	learn := pixel.Shortcut(base, 0, ui.KeyEnter) || pixel.Shortcut(base, 0, ui.KeySpace)
	top := pixel.LineHeight + 4
	// A click moves the cursor; clicking it again learns.
	if base.Clicked() {
		x, y, _ := base.PointerPosition()
		dx, dy := int(x/k.Dot()), int(y/k.Dot())-top
		for i, s := range spots {
			if dx >= s.X*cellPitch && dx < s.X*cellPitch+node && dy >= s.Y*cellPitch && dy < s.Y*cellPitch+node {
				learn, *cursor = learn || i == *cursor, i
			}
		}
	}
	res := SkillTreeResult{Element: shell, Changed: *cursor != old, Learned: -1}
	if res.Changed {
		slog.Debug("8bit: skill cursor", "label", text, "skill", skills[*cursor].Name)
	}
	if learn && open(*cursor) && *points > 0 {
		skills[*cursor].Learned, *points, res.Learned = true, *points-1, *cursor
		slog.Debug("8bit: skill learned", "label", text, "skill", skills[*cursor].Name, "points", *points)
	}
	states := make([]string, len(skills))
	for i := range skills {
		switch {
		case skills[i].Learned:
			states[i] = k.Msg("learned")
		case open(i):
			states[i] = k.Msg("can_learn")
		default:
			states[i] = k.Msg("locked")
		}
	}
	caption := skills[*cursor].Name + " · " + states[*cursor]
	cur, shown, left := *cursor, append([]Skill(nil), skills...), fmt.Sprintf(k.Msg("points"), *points)
	base.Description(caption + " · " + left)
	opens := make([]bool, len(skills))
	for i := range skills {
		opens[i] = open(i)
	}
	w := max(cols*cellPitch-(cellPitch-node), pixel.Measure(left)+8+pixel.Measure(caption))
	shell.Size(k.D(w), k.D(top+rows*cellPitch-(cellPitch-node)+4)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Text(0, 0, left, pixel.Accent, 1)
		g.Text(w-pixel.Measure(caption), 0, caption, pixel.Faint, 1)
		mid := func(s Spot) (int, int) { return s.X*cellPitch + node/2, top + s.Y*cellPitch + node/2 }
		for _, s := range shown {
			for _, n := range s.Needs {
				x0, y0 := mid(shown[n].At)
				x1, y1 := mid(s.At)
				ink := pixel.Shade
				if shown[n].Learned {
					ink = pixel.Accent
				}
				g.Line(x0, y0, x1, y1, ink)
			}
		}
		for i, s := range shown {
			x, y := s.At.X*cellPitch, top+s.At.Y*cellPitch
			sx, sy := x+(node-s.Sprite.W)/2, y+(node-s.Sprite.H)/2
			switch {
			case s.Learned:
				g.Frame(x, y, node, node, pixel.Raised, pixel.Lit)
				g.Sprite(sx, sy, s.Sprite, 0, pixel.Ink, 1)
			case opens[i]:
				g.Frame(x, y, node, node, pixel.Flat, pixel.Face)
				g.Sprite(sx, sy, s.Sprite, 0, pixel.Ink, 1)
			default:
				g.Frame(x, y, node, node, pixel.Sunken, pixel.Well)
				g.Silhouette(sx, sy, s.Sprite, 0, pixel.Shade, 1)
			}
			if i == cur {
				g.HLine(x-2, y-2, node+4, pixel.Accent)
				g.HLine(x-2, y+node+1, node+4, pixel.Accent)
				g.VLine(x-2, y-2, node+4, pixel.Accent)
				g.VLine(x+node+1, y-2, node+4, pixel.Accent)
			}
		}
		g.Finish(st)
	})
	return res
}

// cycle finds a skill needing itself, or -1.
func cycle(skills []Skill) int {
	state := make([]int, len(skills))
	var visit func(i int) bool
	// visit is false on a loop; state 1 is mid-walk.
	visit = func(i int) bool {
		if state[i] != 0 {
			return state[i] == 2
		}
		state[i] = 1
		for _, n := range skills[i].Needs {
			if !visit(n) {
				return false
			}
		}
		state[i] = 2
		return true
	}
	for i := range skills {
		if !visit(i) {
			return i
		}
	}
	return -1
}
