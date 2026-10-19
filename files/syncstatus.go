package files

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// SyncState is where a sync stands.
type SyncState int

const (
	Synced SyncState = iota
	Syncing
	SyncFailed
)

type SyncStatusOptions struct{}

var (
	cloud = pixel.MustSprite(`
		....eeee....
		..eellllee..
		.ellllllllle
		elllllllllle
		elllllllllle
		elllllllllle
		.eeeeeeeeee.`)
	tick = pixel.MustSprite(`
		.......
		......g
		.....gg
		g...gg.
		gg.gg..
		.ggg...
		..g....`)
	// Four quarter turns of a round arrow.
	turning = pixel.MustSprite(`
		..ccc..
		.c...c.
		c....cc
		c...ccc
		c......
		.c...c.
		..ccc..

		..ccc..
		.c...c.
		c.....c
		c.....c
		c..c..c
		.c.ccc.
		...cc..

		..ccc..
		.c...c.
		......c
		ccc...c
		cc....c
		.c...c.
		..ccc..

		..cc...
		.ccc.c.
		c..c..c
		c.....c
		c.....c
		.c...c.
		..ccc..`)
	bang = pixel.MustSprite(`
		..rrr..
		.rrlrr.
		rrrlrrr
		rrrlrrr
		rrrrrrr
		.rrlrr.
		..rrr..`)
)

var syncs = [...]struct {
	msg  string
	ink  pixel.Slot
	mark *pixel.Sprite
}{
	{"synced", pixel.Faint, tick}, {"syncing", pixel.Ink, turning}, {"sync_failed", pixel.Bad, bang},
}

// SyncStatus is a cloud: synced, syncing (arrows turning) or failed.
func SyncStatus(c *ui.Context, state SyncState, o SyncStatusOptions) ui.Element {
	k := pixel.K(c)
	if state < Synced || state > SyncFailed {
		panic(fmt.Sprintf("files: SyncStatus state %d", state))
	}
	s := syncs[state]
	say := k.Msg(s.msg)
	return ui.Box(c).Size(k.D(16+pixel.Measure(say)), k.D(16)).Role(ui.RoleStatus).Label(say).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Sprite(0, 3, cloud, 0, pixel.Ink, 1)
		// Syncing is work, so it turns under Reduce Motion.
		frame := 0
		if state == Syncing {
			frame = g.Busy(4) % 4
		}
		g.Sprite(7, 8, s.mark, frame, pixel.Ink, 1)
		g.Text(16, 0, say, s.ink, 1)
	})
}
