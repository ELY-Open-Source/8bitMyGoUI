package agent

import (
	"fmt"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Member is a summoned agent and those it summoned.
type Member struct {
	Name     string
	Status   Status
	Summoned []Member
}

type SubAgentPartyOptions struct {
	// Width is in dots; zero means 200.
	Width int
}

// partyRow: a member, and per depth whether its branch ends.
type partyRow struct {
	m    Member
	last []bool
}

// SubAgentParty draws the summoning tree, each agent with its status.
func SubAgentParty(c *ui.Context, text string, party []Member, o SubAgentPartyOptions) ui.Element {
	k := pixel.K(c)
	w := pixel.Or("SubAgentParty Width", o.Width, 200)
	var rows []partyRow
	var flat func(ms []Member, last []bool)
	flat = func(ms []Member, last []bool) {
		for i, m := range ms {
			l := append(slices.Clone(last), i == len(ms)-1)
			rows = append(rows, partyRow{m, l})
			flat(m.Summoned, l)
		}
	}
	flat(party, nil)
	if len(rows) == 0 {
		panic("agent: SubAgentParty has nobody")
	}
	said := make([]string, len(rows))
	for i, r := range rows {
		word, _ := r.m.Status.look(k)
		depth := len(r.last) - 1
		if need := depth*12 + 18 + 24 + pixel.Measure(word); r.m.Name == "" || w < need {
			panic(fmt.Sprintf("agent: SubAgentParty member %q needs Width %d", r.m.Name, need))
		}
		said[i] = r.m.Name + ", " + word
	}
	return ui.Box(c).Size(k.D(w), k.D(len(rows)*pixel.LineHeight)).Role(ui.RoleImage).Label(text).Description(strings.Join(said, "; ")).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		for i, row := range rows {
			y, depth := i*pixel.LineHeight, len(row.last)-1
			// Lines run from each summoner down to its summoned.
			for a := 1; a < depth; a++ {
				if !row.last[a] {
					g.VLine((a-1)*12+5, y, pixel.LineHeight, pixel.Edge)
				}
			}
			if depth > 0 {
				x := (depth-1)*12 + 5
				g.VLine(x, y, 8, pixel.Edge)
				if !row.last[depth] {
					g.VLine(x, y+8, pixel.LineHeight-8, pixel.Edge)
				}
				g.HLine(x, y+7, 7, pixel.Edge)
			}
			word, slot := row.m.Status.look(k)
			x := depth * 12
			g.Sprite(x, y+2, icon.Robot, 0, slot, 1)
			if len(row.m.Summoned) > 0 {
				g.VLine(x+5, y+14, 2, pixel.Edge)
			}
			ww := pixel.Measure(word)
			g.Text(x+14, y, pixel.Wrap(row.m.Name, w-x-18-ww, 1)[0], pixel.Ink, 1)
			g.Text(w-ww, y, word, slot, 1)
		}
	})
}
