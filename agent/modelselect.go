package agent

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Model is a fighter to pick; stats run 1–5.
type Model struct {
	Name string
	// Portrait fits 12 dots square; zero means a robot.
	Portrait           *pixel.Sprite
	Speed, Power, Cost int
}

type ModelSelectOptions struct {
	Disabled bool
}

type ModelSelectResult struct {
	Element *ui.Element
	Changed bool
}

type pickKey struct{}

const fighter = 30

// ModelSelect picks a model the way fighters are picked.
func ModelSelect(c *ui.Context, text string, models []Model, chosen *int, o ModelSelectOptions) ModelSelectResult {
	k := pixel.K(c)
	if len(models) == 0 || *chosen < 0 || *chosen >= len(models) {
		panic(fmt.Sprintf("agent: ModelSelect chose %d of %d", *chosen, len(models)))
	}
	faces := make([]*pixel.Sprite, len(models))
	for i, m := range models {
		for _, s := range []int{m.Speed, m.Power, m.Cost} {
			if m.Name == "" || s < 1 || s > 5 {
				panic(fmt.Sprintf("agent: ModelSelect model %q has a stat of %d", m.Name, s))
			}
		}
		faces[i] = portrait("ModelSelect", m.Portrait)
	}
	cols, old := min(len(models), 4), *chosen
	row := ui.Row(c).Gap(k.D(8)).AlignItems(ui.Start)
	cur := ui.Local(row, pickKey{}, func() int { return *chosen })
	*cur = min(*cur, len(models)-1)
	res := ModelSelectResult{}
	res.Element = row.Children(func() {
		shell, base, st, picked := pixel.Pad(c, text, cols, len(models), fighter, fighter, cur, o.Disabled)
		if picked >= 0 {
			*chosen = picked
		}
		at, pick := *cur, *chosen
		m := models[at]
		stats := []struct {
			name  string
			value int
			ink   pixel.Slot
		}{{k.Msg("speed"), m.Speed, pixel.Info}, {k.Msg("power"), m.Power, pixel.Bad}, {k.Msg("cost"), m.Cost, pixel.Warn}}
		said := m.Name
		for _, s := range stats {
			said += fmt.Sprintf(", %s %d", s.name, s.value)
		}
		if at == pick {
			said += ", " + k.Msg("chosen")
		}
		base.Description(said)
		shell.Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			for i, f := range faces {
				x, y := i%cols*fighter, i/cols*fighter
				// The cursor lights its cell; the portrait shows on it.
				if i == at && st.Focus {
					g.Cursor(x, y, fighter, fighter)
				} else {
					g.Frame(x+1, y+1, 28, 28, pixel.Sunken, pixel.Well)
				}
				g.Sprite(x+15-f.W, y+15-f.H, f, 0, pixel.Ink, 2)
				// The chosen fighter wears an Accent ring and corner.
				if i == pick {
					g.Rect(x+1, y+1, 28, 1, pixel.Accent)
					g.Rect(x+1, y+28, 28, 1, pixel.Accent)
					g.Rect(x+1, y+1, 1, 28, pixel.Accent)
					g.Rect(x+28, y+1, 1, 28, pixel.Accent)
					g.Rect(x+2, y+2, 6, 2, pixel.Accent)
					g.Rect(x+2, y+4, 4, 2, pixel.Accent)
					g.Rect(x+2, y+6, 2, 2, pixel.Accent)
				}
			}
			g.Finish(st)
		})
		// The panel shows the fighter under the cursor, bobbing.
		face := faces[at]
		ui.Box(c).Size(k.D(120), k.D(98)).Role(ui.RoleNone).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Frame(0, 0, g.W, g.H, pixel.Double, pixel.Face)
			g.Sprite(6+3*(12-face.W)/2, 6-g.Step(4)%2+3*(12-face.H)/2, face, 0, pixel.Ink, 3)
			g.Text(48, 6, pixel.Wrap(m.Name, g.W-52, 1)[0], pixel.Ink, 1)
			if at == pick {
				g.Text(48, 22, pixel.Wrap(k.Msg("chosen"), g.W-52, 1)[0], pixel.Good, 1)
			}
			for i, s := range stats {
				y := 46 + i*pixel.LineHeight
				g.Text(6, y, pixel.Wrap(s.name, 40, 1)[0], pixel.Faint, 1)
				for j := range 5 {
					ink := pixel.Shade
					if j < s.value {
						ink = s.ink
					}
					g.Rect(50+j*12, y+5, 10, 6, ink)
				}
			}
		})
	})
	res.Changed = *chosen != old
	if res.Changed {
		slog.Debug("8bit: model chosen", "label", text, "model", models[*chosen].Name)
	}
	return res
}
