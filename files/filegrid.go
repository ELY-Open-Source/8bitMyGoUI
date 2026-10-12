package files

import (
	"log/slog"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

type FileGridOptions struct {
	// Height is in dots; zero means 140.
	Height   int
	Disabled bool
}

type FileGridResult struct {
	Element *ui.Element
	// Opened is Enter or a double click on a file.
	Changed, Opened bool
}

const gridCellW, gridCellH = 64, 64

// FileGrid shows files as big icons over their names.
func FileGrid(c *ui.Context, text string, files []File, selected *int, o FileGridOptions) FileGridResult {
	k := pixel.K(c)
	pick("FileGrid", selected, len(files))
	for _, f := range files {
		check("FileGrid", f)
	}
	g := data.GridView(c, text, len(files), selected, data.GridViewOptions{CellWidth: gridCellW, CellHeight: gridCellH, Height: o.Height, Disabled: o.Disabled}, func(i int) {
		ui.Column(c).FillWidth().Padding(k.D(4), 0, 0, 0).Gap(k.D(2)).AlignItems(ui.Center).Children(func() {
			FileIcon(c, files[i].Kind, FileIconOptions{Scale: 2})
			lines := label(files[i].Name, gridCellW-4)
			ui.Box(c).Size(k.D(gridCellW), k.D(2*pixel.LineHeight)).Role(ui.RoleText).Label(files[i].Name).Draw(func(p *ui.Painter, r ui.Rect) {
				g := k.Canvas(p, r)
				for j, l := range lines {
					g.Text((gridCellW-pixel.Measure(l))/2, j*pixel.LineHeight, l, pixel.Ink, 1)
				}
			})
		})
	})
	res := FileGridResult{Element: g.Element, Changed: g.Changed, Opened: g.Submitted && *selected >= 0}
	if res.Opened {
		slog.Debug("8bit: file grid open", "label", text, "file", *selected)
	}
	return res
}

// label fits a name in two lines, extension whole.
func label(name string, w int) []string {
	if pixel.Measure(name) <= w {
		return []string{name}
	}
	if i := strings.LastIndex(name, "."); i > 0 {
		return []string{pixel.Wrap(name[:i], w, 1)[0], pixel.Wrap(name[i:], w, 1)[0]}
	}
	return pixel.Wrap(name, w, 2)
}
