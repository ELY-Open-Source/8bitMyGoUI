package files_test

import (
	"image"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/files"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
)

// offPalette counts dots no palette color matches.
func offPalette(img *image.RGBA) int {
	b := img.Bounds()
	return count(img, 0, 0, b.Dx()/2, b.Dy()/2, none)
}

func TestFilesStayOnThePaletteWhenOff(t *testing.T) {
	list := []files.File{{Name: "A.TXT", Kind: files.Text, Size: 1, Modified: newWorld().clk.t}}
	builds := map[string]func(c *ui.Context, off bool){
		"explorer": func(c *ui.Context, off bool) {
			p, f := "0/0/0", 0
			files.FileExplorer(c, "X", drive, &p, &f, files.FileExplorerOptions{Disabled: off})
		},
		"grid": func(c *ui.Context, off bool) {
			i := 0
			files.FileGrid(c, "X", list, &i, files.FileGridOptions{Disabled: off})
		},
		"path": func(c *ui.Context, off bool) {
			files.PathBar(c, "X", []string{"C:", "A"}, files.PathBarOptions{Disabled: off})
		},
		"card": func(c *ui.Context, off bool) {
			i := 0
			files.MemoryCard(c, "X", []files.Save{{Title: "A", Icon: icon.Sword}}, &i, files.MemoryCardOptions{Disabled: off})
		},
		"uploads": func(c *ui.Context, off bool) {
			files.UploadQueue(c, "X", []files.Upload{{Name: "A", Sent: 1, Size: 2}}, files.UploadQueueOptions{Disabled: off})
		},
		"trash": func(c *ui.Context, off bool) {
			files.Trash(c, "X", 2, files.TrashOptions{Disabled: off})
		},
		"recent": func(c *ui.Context, off bool) {
			i := 0
			files.RecentFiles(c, "X", list, &i, files.RecentFilesOptions{Disabled: off})
		},
		"archive": func(c *ui.Context, off bool) {
			i := 0
			files.ArchiveViewer(c, "X", []files.Entry{{Path: "A", Size: 2, Packed: 1}}, &i, files.ArchiveViewerOptions{Disabled: off})
		},
		"share": func(c *ui.Context, off bool) {
			a := files.CanView
			files.SharePanel(c, "X", "https://x.example", &a, files.SharePanelOptions{Disabled: off})
		},
	}
	for name, build := range builds {
		for _, inert := range []bool{false, true} {
			wd := newWorld()
			wd.inert = inert
			tt := view(700, 400, wd, func(c *ui.Context) { build(c, !inert) })
			if n := offPalette(tt.Image()); n != 0 {
				t.Errorf("%s, inert %v: %d pixels off the palette", name, inert, n)
			}
		}
	}
}
