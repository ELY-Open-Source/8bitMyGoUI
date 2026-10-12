package files_test

import (
	"image"
	"image/color"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/files"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

func TestMemoryCardChoosesAndDeletes(t *testing.T) {
	saves := []files.Save{{Title: "QUEST", Icon: icon.Sword}, {Title: "RACE", Icon: icon.Star, Slot: pixel.Magic}}
	block, deleted := 0, -1
	tt := view(400, 400, newWorld(), func(c *ui.Context) {
		if r := files.MemoryCard(c, "Card", saves, &block, files.MemoryCardOptions{}); r.Deleted >= 0 {
			deleted = r.Deleted
		}
	})
	x, y, _, _ := dots(find(t, tt, "Card"))
	// The sword keeps its colors; Slot inks the star.
	if count(tt.Image(), x+1, y+1, 24, 24, pixel.Warn) == 0 || count(tt.Image(), x+27, y+1, 24, 24, pixel.Magic) == 0 {
		t.Fatal("the sword or the Magic star is missing")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	click(t, tt, "DELETE RACE")
	if block != 1 || deleted != 1 {
		t.Fatalf("Right, DELETE chose %d, deleted %d", block, deleted)
	}
	// A free block has nothing to delete.
	deleted = -1
	tt.Key(ui.Shift, ui.KeyTab)
	tt.Key(0, ui.KeyDown)
	click(t, tt, "DELETE FREE")
	if block != 6 || deleted != -1 || !tt.HasText("2/15 BLOCKS") {
		t.Fatalf("on a free block DELETE gave %d at %d", deleted, block)
	}
	mustPanic(t, "sixteen saves", func() {
		b, full := 0, make([]files.Save, 16)
		for i := range full {
			full[i] = files.Save{Title: "S", Icon: icon.Star}
		}
		view(400, 400, newWorld(), func(c *ui.Context) { files.MemoryCard(c, "C", full, &b, files.MemoryCardOptions{}) })
	})
}

func TestDiskUsageSharesTheBar(t *testing.T) {
	tt := view(500, 400, newWorld(), func(c *ui.Context) {
		files.DiskUsage(c, "Disk", []files.Use{{Name: "GAMES", Bytes: 1 << 30}, {Name: "MUSIC", Bytes: 1 << 30, Slot: pixel.Magic}}, 4<<30, files.DiskUsageOptions{Width: 204})
	})
	x, y, w, _ := dots(find(t, tt, "Disk"))
	img := tt.Image()
	// A quarter each, out of a 200-dot inside.
	if g, m := count(img, x, y+5, w, 1, pixel.Info), count(img, x, y+5, w, 1, pixel.Magic); g != 50 || m != 50 {
		t.Fatalf("GAMES fills %d, MUSIC %d, want 50 each", g, m)
	}
	for _, s := range []string{"GAMES", "1.0 GB", "FREE", "2.0 GB"} {
		find(t, tt, s)
	}
	mustPanic(t, "more than the disk", func() {
		view(500, 400, newWorld(), func(c *ui.Context) {
			files.DiskUsage(c, "D", []files.Use{{Name: "A", Bytes: 3}, {Name: "B", Bytes: 2}}, 4, files.DiskUsageOptions{})
		})
	})
}

func TestUploadQueueCancelsAndFinishes(t *testing.T) {
	ups := []files.Upload{{Name: "A.ZIP", Kind: files.Archive, Sent: 5, Size: 10}, {Name: "B.PNG", Kind: files.Image, Sent: 10, Size: 10}, {Name: "C.GO", Kind: files.Code, Sent: 1, Size: 10, Failed: true}}
	canceled := -1
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		if r := files.UploadQueue(c, "Up", ups, files.UploadQueueOptions{}); r.Canceled >= 0 {
			canceled = r.Canceled
		}
	})
	click(t, tt, "CANCEL C.GO")
	if canceled != 2 {
		t.Fatalf("CANCEL C.GO gave %d", canceled)
	}
	// Done has no cancel; failed says so.
	if _, ok := tt.Find("CANCEL B.PNG"); ok || !tt.HasText("DONE") || !tt.HasText("FAILED") {
		t.Fatalf("texts %q", tt.Texts())
	}
	mustPanic(t, "more sent than the size", func() {
		view(600, 400, newWorld(), func(c *ui.Context) {
			files.UploadQueue(c, "U", []files.Upload{{Name: "A", Sent: 11, Size: 10}}, files.UploadQueueOptions{})
		})
	})
}

func TestTrashEmptiesOnlyWhenFull(t *testing.T) {
	n, emptied := 3, 0
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		if files.Trash(c, "Bin", n, files.TrashOptions{}).Emptied {
			emptied++
			n = 0
		}
	})
	x, y, _, h := dots(find(t, tt, "Bin"))
	// The bin sits mid-row; paper adds Lit dots while full.
	y += (h - 32) / 2
	full := count(tt.Image(), x, y, 32, 32, pixel.Lit)
	click(t, tt, "EMPTY")
	click(t, tt, "EMPTY")
	if bare := count(tt.Image(), x, y, 32, 32, pixel.Lit); emptied != 1 || !tt.HasText("NO ITEMS") || full != bare+8 {
		t.Fatalf("emptied %d times; Lit %d full, %d bare; texts %q", emptied, full, bare, tt.Texts())
	}
	mustPanic(t, "-1 in the bin", func() {
		view(400, 200, newWorld(), func(c *ui.Context) { files.Trash(c, "B", -1, files.TrashOptions{}) })
	})
}

func TestFilePreviewShowsPictureTextOrNothing(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 40, 20))
	for i := range img.Pix {
		img.Pix[i] = 255
	}
	img.SetRGBA(0, 0, color.RGBA{A: 255})
	pv := files.Preview{Picture: img}
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		files.FilePreview(c, "A.PNG", pv, files.FilePreviewOptions{W: 100, H: 60})
	})
	// 40×20 fits 100×60 twice over, by whole steps.
	if r := find(t, tt, "PREVIEW"); r.W != 2*80 || r.H != 2*40 {
		t.Fatalf("the picture is %v, want 80×40 dots", r)
	}
	pv = files.Preview{Text: "LINE ONE"}
	tt.Frame()
	// The text scrolls in the pane's full height.
	if !tt.HasText("LINE ONE") || find(t, tt, "PREVIEW").H != 2*60 {
		t.Fatalf("the text pane %v: %q", find(t, tt, "PREVIEW"), tt.Texts())
	}
	pv = files.Preview{}
	tt.Frame()
	if !tt.HasText("NO PREVIEW") {
		t.Fatalf("texts %q", tt.Texts())
	}
}

func TestFileDetailsReadsEveryField(t *testing.T) {
	f := files.File{Name: "BOSS.PNG", Kind: files.Image, Size: 1_234_567, Modified: time.Date(2026, 4, 1, 23, 59, 0, 0, time.UTC)}
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		files.FileDetails(c, f, files.FileDetailsOptions{})
	})
	// No creation date reads as a lone dash.
	for _, s := range []string{"1.2 MB (1,234,567 B)", "-", "2026-04-01 23:59"} {
		find(t, tt, s)
	}
	// IMAGE names the icon, then fills the KIND row.
	if n := strings.Count(strings.Join(tt.Texts(), "|"), "IMAGE"); n != 2 {
		t.Fatalf("IMAGE shows %d times in %q", n, tt.Texts())
	}
}

func TestSharePanelCopiesAndKeepsTheLinkOutOfTheLog(t *testing.T) {
	wd := newWorld()
	access := files.CanView
	const link = "https://x.example/s/SECRET"
	tt := view(600, 300, wd, func(c *ui.Context) {
		files.SharePanel(c, "Share", link, &access, files.SharePanelOptions{})
	})
	log := logged(func() {
		click(t, tt, "COPY")
		click(t, tt, "EDIT")
	})
	if tt.Clipboard() != link || access != files.CanEdit || !tt.HasText("COPIED!") || strings.Contains(log, "SECRET") {
		t.Fatalf("clipboard %q, access %d; log:\n%s", tt.Clipboard(), access, log)
	}
	later(tt, wd, 3*time.Second)
	if !tt.HasText("COPY") {
		t.Fatalf("COPIED! stayed: %q", tt.Texts())
	}
	mustPanic(t, "an access past EDIT", func() {
		a := files.Access(3)
		view(600, 300, newWorld(), func(c *ui.Context) { files.SharePanel(c, "S", link, &a, files.SharePanelOptions{}) })
	})
}

func TestSyncStatusTurnsWhileSyncing(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := newWorld()
		wd.reduce = reduce
		tt := view(300, 100, wd, func(c *ui.Context) {
			files.SyncStatus(c, files.Syncing, files.SyncStatusOptions{})
		})
		x, y, _, _ := dots(find(t, tt, "SYNCING"))
		before := badge(tt, x+7, y+8)
		later(tt, wd, 250*time.Millisecond)
		// Work under way keeps turning, even under Reduce Motion.
		if badge(tt, x+7, y+8) == before {
			t.Errorf("reduce %v: the arrows held still", reduce)
		}
	}
	tt := view(300, 100, newWorld(), func(c *ui.Context) {
		files.SyncStatus(c, files.SyncFailed, files.SyncStatusOptions{})
	})
	x, y, w, h := dots(find(t, tt, "SYNC FAILED"))
	if count(tt.Image(), x+16, y, w-16, h, pixel.Bad) == 0 {
		t.Fatal("SYNC FAILED is not written in Bad")
	}
	mustPanic(t, "a fourth state", func() {
		view(300, 100, newWorld(), func(c *ui.Context) { files.SyncStatus(c, 3, files.SyncStatusOptions{}) })
	})
}

func TestQuotaMeterFillsGemsByHalves(t *testing.T) {
	for _, tc := range []struct {
		used  int64
		slot  pixel.Slot
		fills int
	}{{3, pixel.Info, 9}, {9, pixel.Warn, 28}, {12, pixel.Bad, 38}, {math.MaxInt64, pixel.Bad, 38}} {
		tt := view(400, 200, newWorld(), func(c *ui.Context) {
			files.QuotaMeter(c, "Space", tc.used, 10, files.QuotaMeterOptions{Gems: 2})
		})
		x, y, _, _ := dots(find(t, tt, "Space"))
		// A gem inks nineteen dots, its left half nine.
		if n := count(tt.Image(), x, y, 17, 7, tc.slot); n != tc.fills {
			t.Errorf("%d of 10: %d dots of slot %d, want %d", tc.used, n, tc.slot, tc.fills)
		}
	}
	mustPanic(t, "41 gems", func() {
		view(400, 200, newWorld(), func(c *ui.Context) { files.QuotaMeter(c, "S", 1, 2, files.QuotaMeterOptions{Gems: 41}) })
	})
}

func TestFileIconsDiffer(t *testing.T) {
	kinds := []files.Kind{files.Folder, files.Text, files.Image, files.Audio, files.Code, files.Archive}
	tt := view(400, 200, newWorld(), func(c *ui.Context) {
		ui.Row(c).Children(func() {
			for _, k := range kinds {
				files.FileIcon(c, k, files.FileIconOptions{Scale: 2})
			}
		})
	})
	// S7: icons are 12 dots square, here at scale 2.
	if r := find(t, tt, "FOLDER"); r.W != 2*24 || r.H != 2*24 {
		t.Fatalf("a FileIcon at scale 2 is %v", r)
	}
	seen := map[string]bool{}
	for _, name := range []string{"FOLDER", "TEXT", "IMAGE", "AUDIO", "CODE", "ARCHIVE"} {
		x, y, w, h := dots(find(t, tt, name))
		var b strings.Builder
		for j := y; j < y+h; j += 2 {
			for i := x; i < x+w; i += 2 {
				b.WriteByte(byte(at(tt.Image(), i, j)))
			}
		}
		seen[b.String()] = true
	}
	if len(seen) != len(kinds) {
		t.Fatalf("%d kinds look alike", len(kinds)-len(seen))
	}
	mustPanic(t, "scale 5", func() {
		view(400, 200, newWorld(), func(c *ui.Context) { files.FileIcon(c, files.Text, files.FileIconOptions{Scale: 5}) })
	})
	mustPanic(t, "a kind past the six", func() {
		view(400, 200, newWorld(), func(c *ui.Context) { files.FileIcon(c, 6, files.FileIconOptions{}) })
	})
}

// badge is the slots of the 7×7 mark there.
func badge(tt *ui.Tester, x, y int) string {
	var b strings.Builder
	for j := y; j < y+7; j++ {
		for i := x; i < x+7; i++ {
			b.WriteByte(byte(at(tt.Image(), i, j)))
		}
	}
	return b.String()
}

func TestUploadQueueRowsFitTheirWidth(t *testing.T) {
	ups := []files.Upload{{Name: "A.ZIP", Kind: files.Archive, Sent: 5, Size: 10}}
	tt := view(600, 300, newWorld(), func(c *ui.Context) {
		files.UploadQueue(c, "Up", ups, files.UploadQueueOptions{})
	})
	q, b := find(t, tt, "Up"), find(t, tt, "CANCEL A.ZIP")
	// The ✕ keeps its 16 dots inside the 240-dot queue.
	if b.W != 2*16 || b.X+b.W > q.X+q.W {
		t.Fatalf("the ✕ is %v in a queue of %v", b, q)
	}
}

func TestMemoryCardFitsAFullCard(t *testing.T) {
	full, b := make([]files.Save, 15), 0
	for i := range full {
		full[i] = files.Save{Title: "HERO QUEST", Icon: icon.Star}
	}
	tt := view(500, 400, newWorld(), func(c *ui.Context) {
		files.MemoryCard(c, "Card", full, &b, files.MemoryCardOptions{})
	})
	card := find(t, tt, "Card")
	// Every line stays within the grid's 130 dots.
	for _, s := range []string{"15/15 BLOCKS", "HERO QUEST", "DELETE HERO QUEST"} {
		if r := find(t, tt, s); r.X < card.X || r.X+r.W > card.X+card.W {
			t.Errorf("%q at %v spills out of %v", s, r, card)
		}
	}
}

func TestUploadQueueFocusFollowsItsRow(t *testing.T) {
	ups := []files.Upload{{Name: "A.ZIP", Sent: 1, Size: 9}, {Name: "B.ZIP", Sent: 1, Size: 9}, {Name: "C.ZIP", Sent: 1, Size: 9}}
	gone := ""
	tt := view(600, 400, newWorld(), func(c *ui.Context) {
		if r := files.UploadQueue(c, "Up", ups, files.UploadQueueOptions{}); r.Canceled >= 0 {
			gone = ups[r.Canceled].Name
		}
	})
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyTab)
	// A leaves; the focus stays on B's ✕, not C's.
	ups = ups[1:]
	tt.Frame()
	tt.Key(0, ui.KeyEnter)
	if gone != "B.ZIP" {
		t.Fatalf("Enter after A left canceled %q", gone)
	}
	for name, bad := range map[string][]files.Upload{
		"a bar of no dots": {{Name: "A", Sent: 1, Size: 2}},
		"two A uploads":    {{Name: "A", Sent: 1, Size: 2}, {Name: "A", Sent: 1, Size: 2}},
	} {
		w := 59
		if len(bad) == 2 {
			w = 0
		}
		mustPanic(t, name, func() {
			view(600, 400, newWorld(), func(c *ui.Context) { files.UploadQueue(c, "U", bad, files.UploadQueueOptions{Width: w}) })
		})
	}
}
