package files_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/files"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var drive = files.Dir{Name: "C:", Files: []files.File{{Name: "BOOT.INI", Kind: files.Text, Size: 10}}, Dirs: []files.Dir{
	{Name: "GAMES", Files: []files.File{{Name: "README.TXT", Kind: files.Text, Size: 2048}}, Dirs: []files.Dir{
		{Name: "QUEST", Files: []files.File{{Name: "SAVE.SAV", Kind: files.Archive, Size: 8192}, {Name: "MAP.PNG", Kind: files.Image, Size: 4096}}},
	}},
	{Name: "MUSIC"},
}}

func TestFileExplorerOpensAtThePathAndWalksFolders(t *testing.T) {
	path, file, opened := "0/0/0", 1, 0
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		if files.FileExplorer(c, "Explorer", drive, &path, &file, files.FileExplorerOptions{}).Opened {
			opened++
		}
	})
	// The tree opens to QUEST; the list shows its files.
	if find(t, tt, "QUEST"); !tt.HasText("8.0 KB") {
		t.Fatalf("no sizes in %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	if path != "0/0" || file != -1 {
		t.Fatalf("Up chose %q, file %d", path, file)
	}
	find(t, tt, "GAMES")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if file != 0 || opened != 1 {
		t.Fatalf("Down, Enter in the list chose %d, opened %d", file, opened)
	}
	for _, bad := range []string{"0/7", "1", "0/x", "0/00", "0/+0", "0/-0"} {
		mustPanic(t, "the path "+bad, func() {
			p, f := bad, -1
			view(600, 400, newWorld(), func(c *ui.Context) {
				files.FileExplorer(c, "E", drive, &p, &f, files.FileExplorerOptions{})
			})
		})
	}
}

func TestFileGridMovesByRowsAndOpens(t *testing.T) {
	grid := make([]files.File, 7)
	for i := range grid {
		grid[i] = files.File{Name: string(rune('A' + i)), Kind: files.Kind(i % 6), Size: 1}
	}
	sel, opened := -1, 0
	// 400 pixels is 200 dots: three 64-dot cells a row.
	tt := view(400, 400, newWorld(), func(c *ui.Context) {
		if files.FileGrid(c, "Files", grid, &sel, files.FileGridOptions{}).Opened {
			opened++
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if sel != 4 || opened != 1 {
		t.Fatalf("Right, Right, Down, Enter chose %d, opened %d", sel, opened)
	}
	find(t, tt, "E")
	mustPanic(t, "a choice past the files", func() {
		s := 7
		view(400, 400, newWorld(), func(c *ui.Context) { files.FileGrid(c, "F", grid, &s, files.FileGridOptions{}) })
	})
	mustPanic(t, "a negative size", func() {
		s := -1
		view(400, 400, newWorld(), func(c *ui.Context) {
			files.FileGrid(c, "F", []files.File{{Name: "A", Size: -1}}, &s, files.FileGridOptions{})
		})
	})
}

func TestRecentFilesSayWhenAndOpen(t *testing.T) {
	wd := newWorld()
	recent := []files.File{
		{Name: "NOTES.TXT", Kind: files.Text, Modified: wd.clk.t.Add(-12 * time.Minute)},
		{Name: "SONG.MID", Kind: files.Audio, Modified: wd.clk.t.Add(-30 * time.Hour)},
	}
	sel, opened := -1, false
	tt := view(500, 300, wd, func(c *ui.Context) {
		opened = files.RecentFiles(c, "Recent", recent, &sel, files.RecentFilesOptions{}).Opened || opened
	})
	if !tt.HasText("12 MIN AGO") || !tt.HasText("YESTERDAY") {
		t.Fatalf("no times in %q", tt.Texts())
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if sel != 1 || !opened {
		t.Fatalf("Down, Down, Enter chose %d, opened %v", sel, opened)
	}
}

func TestArchiveViewerSortsAndKeepsTheChoice(t *testing.T) {
	entries := []files.Entry{{Path: "B.TXT", Size: 300, Packed: 100}, {Path: "A.TXT", Size: 100, Packed: 50}, {Path: "C.TXT", Size: 200, Packed: 90}}
	sel, sorted := 0, false
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		sorted = files.ArchiveViewer(c, "Zip", entries, &sel, files.ArchiveViewerOptions{}).Sorted || sorted
	})
	// By name B is the second row; Down reaches C.
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if sel != 2 {
		t.Fatalf("Down from B chose entry %d, want C", sel)
	}
	// By size C sits in the middle; Down reaches B.
	sorted = false
	click(t, tt, "SIZE")
	if !sorted {
		t.Fatal("a header click did not report Sorted")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if sel != 0 || !tt.HasText("3 FILES: 600 B → 240 B") {
		t.Fatalf("by size, Down from C chose %d; texts %q", sel, tt.Texts())
	}
	// Biggest first, B leads; Down reaches C.
	click(t, tt, "SIZE")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	if sel != 2 {
		t.Fatalf("by size falling, Down from B chose %d", sel)
	}
	mustPanic(t, "sizes past int64", func() {
		s := -1
		view(600, 400, newWorld(), func(c *ui.Context) {
			files.ArchiveViewer(c, "Z", []files.Entry{{Path: "A", Size: 1 << 62}, {Path: "B", Size: 1 << 62}, {Path: "C", Size: 1 << 62}}, &s, files.ArchiveViewerOptions{})
		})
	})
}

func TestPathBarClimbsToParents(t *testing.T) {
	got := -1
	tt := view(600, 200, newWorld(), func(c *ui.Context) {
		if r := files.PathBar(c, "Path", []string{"C:", "GAMES", "QUEST"}, files.PathBarOptions{}); r.Clicked >= 0 {
			got = r.Clicked
		}
	})
	click(t, tt, "GAMES")
	if got != 1 {
		t.Fatalf("clicking GAMES gave %d", got)
	}
	// The folders are one Tab stop; arrows walk them.
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyEnter)
	if got != 0 || !tt.Focused("C:") {
		t.Fatalf("Left, Enter gave %d", got)
	}
	got = -1
	click(t, tt, "QUEST")
	if got != -1 {
		t.Fatalf("the folder you are in clicked as %d", got)
	}
}

func TestBrowsersRefuseNamelessThings(t *testing.T) {
	cases := map[string]func(c *ui.Context){
		"a nameless folder": func(c *ui.Context) {
			p, f := "0", -1
			files.FileExplorer(c, "E", files.Dir{Name: "C:", Dirs: []files.Dir{{}}}, &p, &f, files.FileExplorerOptions{})
		},
		"an empty path part": func(c *ui.Context) {
			files.PathBar(c, "P", []string{"C:", ""}, files.PathBarOptions{})
		},
		"a recent file with no time": func(c *ui.Context) {
			s := -1
			files.RecentFiles(c, "R", []files.File{{Name: "A", Kind: files.Text}}, &s, files.RecentFilesOptions{})
		},
	}
	for name, build := range cases {
		mustPanic(t, name, func() { view(600, 300, newWorld(), build) })
	}
}

func TestFileExplorerTreeFitsLongNames(t *testing.T) {
	root := files.Dir{Name: "C:", Dirs: []files.Dir{{Name: "GAMES", Dirs: []files.Dir{{Name: "THE VERY LONG QUEST"}}}}}
	path, file := "0/0/0", -1
	tt := view(900, 300, newWorld(), func(c *ui.Context) {
		files.FileExplorer(c, "E", root, &path, &file, files.FileExplorerOptions{})
	})
	// Depth 2 indents 20 dots; the arrow takes 12.
	if w := find(t, tt, "FOLDERS").W / 2; int(w) < 20+12+pixel.Measure("THE VERY LONG QUEST") {
		t.Fatalf("the tree is %v dots for a %d-dot name", w, pixel.Measure("THE VERY LONG QUEST"))
	}
}
