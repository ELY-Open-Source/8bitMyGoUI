package gallery

import (
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/files"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var storageEntries = []Entry{
	{414, "RecentFiles", "最近文件", "files", "What you touched lately, and how long ago.", demoRecentFiles},
	{415, "FilePreview", "文件预览", "files", "Pictures shrink to fit, text scrolls, and the rest says so.", demoFilePreview},
	{416, "ArchiveViewer", "压缩包", "files", "What is inside the zip, and how much squeezing saved. Headers sort.", demoArchiveViewer},
	{417, "FileDetails", "文件详情", "files", "Name, kind, size to the byte, and when it all happened.", demoFileDetails},
	{418, "SharePanel", "分享面板", "files", "A link to send and how much it lets them do. COPY takes it.", demoSharePanel},
	{419, "SyncStatus", "同步状态", "files", "All saved, saving (the arrows turn), or something went wrong.", demoSyncStatus},
	{420, "QuotaMeter", "配额", "files", "Storage as gems, like hearts in an old adventure. Low turns gold, then red.", demoQuotaMeter},
}

type storageDemo struct {
	recent, entry, kind int
	access              files.Access
}

func storageState(c *ui.Context) *storageDemo {
	return ui.Local(c.Root(), demoKey("storage"), func() storageDemo { return storageDemo{recent: -1, entry: -1} })
}

func demoRecentFiles(c *ui.Context, log func(string)) {
	s := storageState(c)
	now := pixel.K(c).Now()
	recent := []files.File{
		{Name: "SAVE1.SAV", Kind: files.Archive, Modified: now.Add(-20 * time.Second)},
		{Name: "CHEATS.TXT", Kind: files.Text, Modified: now.Add(-12 * time.Minute)},
		{Name: "MAP.PNG", Kind: files.Image, Modified: now.Add(-5 * time.Hour)},
		{Name: "THEME.MID", Kind: files.Audio, Modified: now.Add(-30 * time.Hour)},
		{Name: "MAIN.GO", Kind: files.Code, Modified: now.Add(-9 * 24 * time.Hour)},
	}
	r := files.RecentFiles(c, "Recent", recent, &s.recent, files.RecentFilesOptions{})
	if r.Changed {
		log("chose " + recent[s.recent].Name)
	}
	if r.Opened {
		log("open " + recent[s.recent].Name)
	}
}

const cheats = "UP UP DOWN DOWN LEFT RIGHT LEFT RIGHT B A START.\n\nWorks on the title screen. Thirty lives, or so the playground said. Hold SELECT for the level select, but only on Tuesdays."

func demoFilePreview(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := storageState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		if input.Segmented(c, "File", []string{"CHEATS.TXT", "SUNSET.PNG", "THEME.MID"}, &s.kind, input.SegmentedOptions{}).Changed {
			log([]string{"text", "picture", "no preview"}[s.kind])
		}
		switch s.kind {
		case 0:
			files.FilePreview(c, "CHEATS.TXT", files.Preview{Text: cheats}, files.FilePreviewOptions{H: 80})
		case 1:
			files.FilePreview(c, "SUNSET.PNG", files.Preview{Picture: sunset}, files.FilePreviewOptions{H: 80})
		default:
			files.FilePreview(c, "THEME.MID", files.Preview{}, files.FilePreviewOptions{H: 80})
		}
	})
}

func demoArchiveViewer(c *ui.Context, log func(string)) {
	s := storageState(c)
	entries := []files.Entry{
		{Path: "LEVELS/1-1.MAP", Size: 18_432, Packed: 3_120},
		{Path: "LEVELS/1-2.MAP", Size: 20_480, Packed: 3_604},
		{Path: "MUSIC/THEME.MID", Size: 21_504, Packed: 19_800},
		{Path: "SPRITES.PNG", Size: 65_536, Packed: 64_900},
		{Path: "README.TXT", Size: 2_048, Packed: 900},
	}
	r := files.ArchiveViewer(c, "LOOT.ZIP", entries, &s.entry, files.ArchiveViewerOptions{Height: 80})
	if r.Changed {
		log("chose " + entries[s.entry].Path)
	}
	if r.Sorted {
		log("sorted")
	}
}

func demoFileDetails(c *ui.Context, log func(string)) {
	files.FileDetails(c, files.File{
		Name: "BOSS FIGHT.PNG", Kind: files.Image, Size: 1_234_567,
		Created:  time.Date(1998, 11, 21, 9, 30, 0, 0, time.UTC),
		Modified: time.Date(2026, 4, 1, 23, 59, 0, 0, time.UTC),
	}, files.FileDetailsOptions{})
}

func demoSharePanel(c *ui.Context, log func(string)) {
	s := storageState(c)
	r := files.SharePanel(c, "Share", "https://8bit.example/s/UUDDLRLRBA", &s.access, files.SharePanelOptions{})
	if r.Copied {
		log("link copied")
	}
	if r.Changed {
		log([]string{"view", "comment", "edit"}[s.access])
	}
}

func demoSyncStatus(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		for _, st := range []files.SyncState{files.Synced, files.Syncing, files.SyncFailed} {
			files.SyncStatus(c, st, files.SyncStatusOptions{})
		}
	})
}

func demoQuotaMeter(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(8)).Children(func() {
		for _, used := range []int64{3 << 30, 13 << 30, 16 << 30} {
			files.QuotaMeter(c, "Storage", used, 15<<30, files.QuotaMeterOptions{})
		}
	})
}
