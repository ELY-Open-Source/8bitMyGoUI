package overlay

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/layout"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var more = pixel.MustSprite(`
	#####
	.###.
	..#..
`)

type DialogueBoxOptions struct {
	// Name labels the speaker's tab.
	Name string
	// Portrait draws the speaker; nil shows none.
	Portrait *pixel.Sprite
	// CPS is the typing speed; zero means 30.
	CPS      int
	Disabled bool
}

type DialogueBoxResult struct {
	Element *ui.Element
	// Changed turns a page; Done passes the last one.
	Changed, Done bool
}

type dialogueKey struct{}

// typing is the page, when it began, and if out.
type typing struct {
	page  int
	since time.Time
	out   bool
}

// DialogueBox types pages; a press shows all, then turns.
func DialogueBox(c *ui.Context, pages []string, page *int, o DialogueBoxOptions) DialogueBoxResult {
	k := pixel.K(c)
	if *page < 0 || *page >= len(pages) {
		panic(fmt.Sprintf("overlay: DialogueBox page %d of %d", *page, len(pages)))
	}
	cps := pixel.Or("DialogueBox CPS", o.CPS, 30)
	off := o.Disabled || pixel.Current(c).Inert
	box := ui.Box(c).Role(ui.RoleText).Label(strings.TrimSpace(o.Name + " " + pages[*page])).FocusRing(false)
	if !off {
		box.Focusable()
	}
	pixel.Park(c, box, off)
	st := ui.Local(box, dialogueKey{}, func() typing { return typing{page: -1} })
	if st.page != *page {
		*st = typing{page: *page, since: k.Now()}
	}
	res := DialogueBoxResult{Element: box}
	// A press shows the whole page, then turns it.
	if !off && box.Clicked() {
		switch {
		case !st.out:
			st.out = true
		case *page < len(pages)-1:
			*page, res.Changed = *page+1, true
			*st = typing{page: *page, since: k.Now()}
		default:
			res.Done = true
		}
	}
	focus := box.FocusVisible()
	box.DrawOver(func(p *ui.Painter, r ui.Rect) {
		if focus {
			g := k.Canvas(p, r)
			g.FocusMark(0, 0, g.W, g.H)
		}
	}).Children(func() {
		layout.Window(c, layout.WindowOptions{Title: o.Name}, func() {
			ui.Row(c).Gap(k.D(8)).Children(func() {
				if sp := o.Portrait; sp != nil {
					ui.Box(c).Size(k.D(sp.W*2+6), k.D(sp.H*2+6)).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
						g := k.Canvas(p, r)
						g.Frame(0, 0, g.W, g.H, pixel.Sunken, pixel.Well)
						g.Sprite(3, 3, sp, 0, pixel.Ink, 2)
					})
				}
				words := ui.Box(c).Grow(1)
				w, check := pixel.WrapWidth(c, words)
				lines := pixel.Wrap(pages[*page], max(0, w-8), 0)
				total := 0
				for _, l := range lines {
					total += len([]rune(l))
				}
				shown := total
				if !st.out && !k.Reduce() {
					shown = k.Since(c, st.since, cps, total)
				}
				st.out = shown == total
				words.Height(k.D(max(2, len(lines)) * pixel.LineHeight)).Draw(func(p *ui.Painter, r ui.Rect) {
					check(r)
					g := k.Canvas(p, r)
					left, end := shown, 0
					for i, l := range lines {
						lr := []rune(l)
						n := min(len(lr), max(0, left))
						g.Text(0, i*pixel.LineHeight, string(lr[:n]), pixel.Ink, 1)
						left -= len(lr)
						end = pixel.Measure(string(lr[:n]))
					}
					if shown == total && g.Step(2)%2 == 0 {
						g.Sprite(end+3, len(lines)*pixel.LineHeight-6, more, 0, pixel.Accent, 1)
					}
				})
			})
		})
	})
	if res.Changed || res.Done {
		slog.Debug("8bit: dialogue", "name", o.Name, "page", *page, "done", res.Done)
	}
	return res
}
