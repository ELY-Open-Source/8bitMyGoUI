package project

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Issue is a bug or a wish, under an id.
type Issue struct {
	ID       string
	Title    string
	Labels   []Label
	Assignee string
	Priority Level
}

type IssueCardOptions struct {
	// Width is in dots; zero means 160.
	Width    int
	Disabled bool
}

type IssueCardResult struct {
	Element *ui.Element
	Clicked bool
}

// IssueCard shows an issue's id, title, labels and assignee.
func IssueCard(c *ui.Context, is Issue, o IssueCardOptions) IssueCardResult {
	k := pixel.K(c)
	if is.ID == "" || is.Title == "" {
		panic("project: IssueCard needs an id and a title")
	}
	said := make([]string, 0, len(is.Labels)+2)
	for _, l := range is.Labels {
		if l.Name == "" {
			panic("project: IssueCard has a label with no name")
		}
		said = append(said, l.Name)
	}
	who := is.Assignee
	if who == "" {
		who = k.Msg("nobody")
	}
	said = append(said, is.Priority.word(k), who)
	w := pixel.Or("IssueCard Width", o.Width, 160)
	lines := pixel.Wrap(is.Title, w-8, 2)
	// The shell draws the card; the button stays bare.
	shell, base, st := pixel.Shell(c, func() *ui.Element {
		return ui.ButtonBase(c).Label(is.ID + " " + is.Title).Disabled(o.Disabled)
	})
	base.Description(strings.Join(said, ", "))
	res := IssueCardResult{Element: shell, Clicked: base.Clicked()}
	h := 3 + pixel.LineHeight*(len(lines)+2) + 6
	shell.Size(k.D(w), k.D(h)).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		fill := pixel.Face
		if st.Hover {
			fill = pixel.Lit
		}
		g.Frame(0, 0, w, h, pixel.Raised, fill)
		g.Text(4, 3, is.ID, pixel.Faint, 1)
		mark(g, w-16, 5, is.Priority, 1)
		for i, l := range lines {
			g.Text(4, 3+pixel.LineHeight*(i+1), l, pixel.Ink, 1)
		}
		// Labels run until the assignee's seat; the rest wait.
		y, x := 3+pixel.LineHeight*(len(lines)+1)+2, 4
		for _, l := range is.Labels {
			cw := 10 + pixel.Measure(l.Name)
			if x+cw > w-24 {
				break
			}
			g.Rect(x, y+5, 7, 7, pixel.Edge)
			g.Rect(x+1, y+6, 5, 5, l.Slot)
			g.Text(x+10, y, l.Name, pixel.Faint, 1)
			x += cw + 6
		}
		if is.Assignee == "" {
			g.Veil(w-20, y, 16, 16, pixel.Shade, 8)
		} else {
			display.Face(g, w-20, y, is.Assignee, 1)
		}
		g.Finish(st)
	})
	if res.Clicked {
		slog.Debug("8bit: issue card", "id", is.ID)
	}
	return res
}
