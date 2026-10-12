package files

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

// Dir is a named folder of folders and files.
type Dir struct {
	Name  string
	Dirs  []Dir
	Files []File
}

type FileExplorerOptions struct {
	// Height is the file list's in dots; zero means 120.
	Height   int
	Disabled bool
}

type FileExplorerResult struct {
	Element *ui.Element
	// Opened is Enter or a double click on a file.
	Changed, Opened bool
}

// FileExplorer sets the folder tree beside the chosen folder's files.
func FileExplorer(c *ui.Context, text string, root Dir, path *string, file *int, o FileExplorerOptions) FileExplorerResult {
	k := pixel.K(c)
	dir, ok := folderAt(root, *path)
	if !ok {
		panic(fmt.Sprintf("files: FileExplorer has no folder at %q", *path))
	}
	pick("FileExplorer", file, len(dir.Files))
	h := pixel.Or("FileExplorer Height", o.Height, 120)
	tree, tw := node(root, 0)
	var res FileExplorerResult
	res.Element = ui.Row(c).FillWidth().Gap(k.D(4)).AlignItems(ui.Start).Role(ui.RoleGroup).Label(text).Children(func() {
		old := *path
		data.Tree(c, k.Msg("folders"), []data.TreeNode{tree}, path, data.TreeOptions{Width: tw, Disabled: o.Disabled})
		if *path != old {
			// A new folder starts with no file chosen.
			dir, _ = folderAt(root, *path)
			*file, res.Changed = -1, true
			slog.Debug("8bit: file explorer folder", "label", text, "path", *path)
		}
		files := dir.Files
		for _, f := range files {
			check("FileExplorer", f)
		}
		l := data.List(c, dir.Name, len(files), file, data.ListOptions{Height: h, Disabled: o.Disabled}, func(i int) string { return files[i].Name }, func(i int) {
			row(c, k, files[i], func() { pixel.Text(c, size(files[i].Size), pixel.TextOptions{Slot: pixel.Faint}) })
		})
		l.Element.Grow(1)
		res.Changed = res.Changed || l.Changed
		res.Opened = l.Submitted && *file >= 0
	})
	if res.Opened {
		slog.Debug("8bit: file explorer open", "label", text, "file", *file)
	}
	return res
}

// node makes a folder's tree and the width it needs.
func node(f Dir, depth int) (data.TreeNode, int) {
	if f.Name == "" {
		panic("files: FileExplorer has a folder with no name")
	}
	// data.Tree indents ten dots a level; the arrow takes twelve.
	n, w := data.TreeNode{Label: f.Name}, depth*10+16+pixel.Measure(f.Name)
	for _, sub := range f.Dirs {
		child, cw := node(sub, depth+1)
		n.Children, w = append(n.Children, child), max(w, cw)
	}
	return n, w
}

// folderAt follows a tree path, "0/1", down from root.
func folderAt(root Dir, path string) (Dir, bool) {
	parts := strings.Split(path, "/")
	if parts[0] != "0" {
		return Dir{}, false
	}
	f := root
	for _, p := range parts[1:] {
		// Only the plain form, as data.Tree writes it, is known.
		i, err := strconv.Atoi(p)
		if err != nil || strconv.Itoa(i) != p || i < 0 || i >= len(f.Dirs) {
			return Dir{}, false
		}
		f = f.Dirs[i]
	}
	return f, true
}
