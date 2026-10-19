package files

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/feedback"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

// Upload is a file on its way up.
type Upload struct {
	// Name is unique in the queue; it keys the row.
	Name string
	Kind Kind
	// Sent and Size are bytes; Sent reaches Size when done.
	Sent, Size int64
	Failed     bool
}

type UploadQueueOptions struct {
	// Width is in dots; zero means 240.
	Width    int
	Disabled bool
}

type UploadQueueResult struct {
	Element ui.Element
	// Canceled is the upload whose ✕ was pressed, else -1.
	Canceled int
}

// UploadQueue lists uploads with progress; ✕ cancels one.
func UploadQueue(c *ui.Context, text string, ups []Upload, o UploadQueueOptions) UploadQueueResult {
	k := pixel.K(c)
	w := pixel.Or("UploadQueue Width", o.Width, 240)
	// ✕, icon and gaps take 40 dots, plus the percent.
	bar := w - 44 - pixel.MicroWidth("100%")
	if bar < 8 {
		panic(fmt.Sprintf("files: UploadQueue Width %d leaves a %d-dot bar", w, bar))
	}
	names := map[string]bool{}
	for _, u := range ups {
		look(u.Kind)
		if u.Name == "" || names[u.Name] || u.Size <= 0 || u.Sent < 0 || u.Sent > u.Size {
			panic(fmt.Sprintf("files: UploadQueue %q, twice %v, sent %d of %d", u.Name, names[u.Name], u.Sent, u.Size))
		}
		names[u.Name] = true
	}
	res := UploadQueueResult{Canceled: -1}
	res.Element = ui.Column(c).Width(k.D(w)).Gap(k.D(4)).Role(ui.RoleList).Label(text).Children(func() {
		for i, u := range ups {
			done := u.Sent == u.Size && !u.Failed
			ui.Row(c).Key(u.Name).FillWidth().Gap(k.D(4)).AlignItems(ui.Center).Role(ui.RoleListItem).Children(func() {
				FileIcon(c, u.Kind, FileIconOptions{})
				ui.Column(c).Grow(1).Gap(k.D(1)).Children(func() {
					ui.Row(c).FillWidth().Gap(k.D(4)).Children(func() {
						pixel.Text(c, u.Name, pixel.TextOptions{Wrap: true, MaxLines: 1}).Grow(1)
						switch {
						case u.Failed:
							pixel.Text(c, k.Msg("failed"), pixel.TextOptions{Slot: pixel.Bad})
						case done:
							pixel.Text(c, k.Msg("done"), pixel.TextOptions{Slot: pixel.Good})
						}
					})
					feedback.Progress(c, u.Name, float64(u.Sent)/float64(u.Size), feedback.ProgressOptions{Width: bar, Percent: !done && !u.Failed})
				})
				// A finished upload has nothing left to cancel.
				if done {
					ui.Box(c).Size(k.D(16), k.D(16)).Draw(func(p *ui.Painter, r ui.Rect) {
						k.Canvas(p, r).Sprite(2, 2, icon.Check, 0, pixel.Good, 1)
					})
					return
				}
				if input.IconButton(c, icon.Close, k.Msg("cancel")+" "+u.Name, input.IconButtonOptions{Size: input.Small, Disabled: o.Disabled}).Clicked {
					res.Canceled = i
				}
			})
		}
	})
	if res.Canceled >= 0 {
		slog.Debug("8bit: upload canceled", "label", text, "upload", res.Canceled)
	}
	return res
}
