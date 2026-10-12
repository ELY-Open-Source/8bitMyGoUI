package chat

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/display"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type AttachmentChipOptions struct {
	// Width caps the name in dots; zero means 100.
	Width    int
	Disabled bool
}

type AttachmentChipResult struct {
	Element *ui.Element
	Removed bool
}

// kinds pick a file's icon by its extension.
var kinds = map[string]*pixel.Sprite{
	".png": icon.Image, ".jpg": icon.Image, ".gif": icon.Image, ".mp3": icon.Music, ".wav": icon.Music, ".ogg": icon.Music,
	".mp4": icon.Video, ".mov": icon.Video, ".go": icon.Code, ".js": icon.Code, ".py": icon.Code,
	".txt": icon.FileText, ".md": icon.FileText, ".pdf": icon.FileText,
}

// AttachmentChip shows a file's icon, name, size and ×.
func AttachmentChip(c *ui.Context, name string, size int64, o AttachmentChipOptions) AttachmentChipResult {
	k := pixel.K(c)
	if size < 0 {
		panic(fmt.Sprintf("chat: attachment %q has size %d", name, size))
	}
	sp, ok := kinds[strings.ToLower(filepath.Ext(name))]
	if !ok {
		sp = icon.File
	}
	var res AttachmentChipResult
	res.Element = ui.Row(c).Gap(k.D(3)).Padding(k.D(2), k.D(2), k.D(2), k.D(4)).AlignItems(ui.Center).Draw(func(p *ui.Painter, r ui.Rect) {
		g := k.Canvas(p, r)
		g.Frame(0, 0, g.W, g.H, pixel.Notched, pixel.Face)
	}).Children(func() {
		display.Icon(c, sp, display.IconOptions{Slot: pixel.Faint})
		pixel.Text(c, line(name, pixel.Or("AttachmentChip Width", o.Width, 100)), pixel.TextOptions{})
		pixel.Text(c, bytes(size), pixel.TextOptions{Slot: pixel.Faint})
		if input.IconButton(c, icon.Close, k.Msg("remove")+" "+name, input.IconButtonOptions{Variant: input.Ghost, Size: input.Small, Disabled: o.Disabled}).Clicked {
			res.Removed = true
			slog.Debug("8bit: attachment removed", "name", name)
		}
	})
	return res
}

// bytes prints n in binary units, tenths below ten.
func bytes(n int64) string {
	v, units := float64(n), []string{"B", "KB", "MB", "GB", "TB"}
	i := 0
	for v >= 1024 && i < len(units)-1 {
		v, i = v/1024, i+1
	}
	if i > 0 && v < 10 {
		return fmt.Sprintf("%.1f %s", v, units[i])
	}
	return fmt.Sprintf("%.0f %s", v, units[i])
}
