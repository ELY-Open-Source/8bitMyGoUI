package overlay

import (
	"log/slog"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/text"
)

// Notification is a card in the queue Notifications shows.
type Notification struct {
	Icon        *pixel.Sprite
	Title, Text string
	// For is how long it shows; zero means five seconds.
	For time.Duration
	// id and shown are set as it first shows.
	id    uint64
	shown time.Time
}

type NotificationsOptions struct {
	// Width is each card's in dots; zero means 180.
	Width    int
	Disabled bool
}

type NotificationsResult struct {
	// Element holds the cards on screen.
	Element ui.Element
	// Dismissed is the card × closed this frame, or nil.
	Dismissed *Notification
}

type notesKey struct{}

// Notifications stacks timed cards at the top right; × dismisses.
func Notifications(c *ui.Context, queue *[]Notification, o NotificationsOptions) NotificationsResult {
	k := pixel.K(c)
	w := pixel.Or("Notifications Width", o.Width, 180)
	now := k.Now()
	kept := (*queue)[:0]
	for _, n := range *queue {
		if n.For < 0 {
			panic("overlay: Notification For is negative")
		}
		life := n.For
		if life == 0 {
			life = 5 * time.Second
		}
		if n.shown.IsZero() || k.Elapsed(&n.shown) < life {
			kept = append(kept, n)
		}
	}
	*queue = kept
	shown := min(4, len(*queue))
	last := ui.Local(c.Root(), notesKey{}, func() uint64 { return 0 })
	for i := range shown {
		if (*queue)[i].id == 0 {
			*last++
			(*queue)[i].id, (*queue)[i].shown = *last, now
		}
	}
	if shown > 0 {
		k.Tick(c, 4)
	}
	var res NotificationsResult
	ui.Overlay(c, func() {
		res.Element = ui.Column(c).Absolute().Right(k.D(8)).Top(k.D(8)).Width(k.D(w)).Gap(k.D(4)).AlignItems(ui.Stretch).Children(func() {
			for i := range shown {
				n := (*queue)[i]
				ui.Row(c).Key(n.id).Padding(k.D(5), k.D(6)).Gap(k.D(6)).Role(ui.RoleStatus).Label(n.Title).Description(n.Text).Draw(func(p *ui.Painter, r ui.Rect) {
					g := k.Canvas(p, r)
					g.Frame(0, 0, g.W, g.H, pixel.Window, pixel.Face)
				}).Children(func() {
					art := n.Icon
					if art == nil {
						art = icon.Bell
					}
					ui.Box(c).Size(k.D(art.W), k.D(art.H)).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
						k.Canvas(p, r).Sprite(0, 0, art, 0, pixel.Accent, 1)
					})
					ui.Column(c).Grow(1).Children(func() {
						pixel.Text(c, n.Title, pixel.TextOptions{})
						text.Paragraph(c, n.Text, text.ParagraphOptions{Slot: pixel.Faint})
					})
					if input.IconButton(c, icon.Close, k.Msg("close"), input.IconButtonOptions{Size: input.Small, Variant: input.Ghost, Disabled: o.Disabled}).Clicked {
						res.Dismissed = &n
					}
				})
			}
		})
	})
	if res.Dismissed != nil {
		*queue = dropNote(*queue, res.Dismissed.id)
		slog.Debug("8bit: notification dismissed", "title", res.Dismissed.Title)
	}
	return res
}

// dropNote removes the card with id from q.
func dropNote(q []Notification, id uint64) []Notification {
	for i := range q {
		if q[i].id == id {
			return append(q[:i], q[i+1:]...)
		}
	}
	return q
}
