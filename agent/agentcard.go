package agent

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Agent is a party member; its context window is HP.
type Agent struct {
	Name, Class string
	Level       int
	// Used and Window count context tokens.
	Used, Window int
}

type AgentCardOptions struct {
	// Portrait fits 12 dots square; zero means a robot.
	Portrait *pixel.Sprite
}

const cardW, cardH = 136, 48

// AgentCard is a party card: class, level, context as HP.
func AgentCard(c *ui.Context, a Agent, o AgentCardOptions) *ui.Element {
	k := pixel.K(c)
	if a.Name == "" || a.Class == "" || a.Level < 1 || a.Window < 1 || a.Used < 0 || a.Used > a.Window {
		panic(fmt.Sprintf("agent: AgentCard %q, %q, LV %d, %d of %d tokens", a.Name, a.Class, a.Level, a.Used, a.Window))
	}
	face := portrait("AgentCard", o.Portrait)
	left := a.Window - a.Used
	lv, hp := fmt.Sprintf(k.Msg("level"), a.Level), tokens(left)+"/"+tokens(a.Window)
	if pixel.Measure(lv) > cardW-40-28 {
		panic(fmt.Sprintf("agent: AgentCard level %d is too wide", a.Level))
	}
	ink := health(left, a.Window)
	said := lv + " " + a.Class + ", " + k.Msg("hp") + " " + hp
	return ui.Box(c).Size(k.D(cardW), k.D(cardH)).Role(ui.RoleImage).Label(a.Name).Description(said).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, cardW, cardH, pixel.Raised, pixel.Face)
		g.Frame(4, 10, 28, 28, pixel.Sunken, pixel.Well)
		g.Sprite(18-face.W, 24-face.H, face, 0, pixel.Ink, 2)
		x, tw := 36, cardW-40
		g.Text(x, 1, pixel.Wrap(a.Name, tw, 1)[0], pixel.Ink, 1)
		lvW := pixel.Measure(lv)
		g.Text(x, 15, pixel.Wrap(a.Class, tw-lvW-4, 1)[0], pixel.Faint, 1)
		g.Text(cardW-4-lvW, 15, lv, pixel.Ink, 1)
		g.Micro(x, 32, k.Msg("hp"), pixel.Ink, 1)
		g.Micro(cardW-4-pixel.MicroWidth(hp), 32, hp, pixel.Ink, 1)
		gauge(g, x, 38, tw, float64(left)/float64(a.Window), ink)
	})
}
