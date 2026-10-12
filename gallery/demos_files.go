package gallery

import (
	"fmt"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/files"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var filesEntries = []Entry{
	{406, "FileExplorer", "文件管理器", "files", "Folders on the left, their files on the right. Double click opens.", demoFileExplorer},
	{407, "FileGrid", "文件网格", "files", "Big icons over short names. Arrows walk the grid; Enter opens.", demoFileGrid},
	{408, "FileIcon", "文件图标", "files", "Six kinds, placed dot by dot: folder, text, image, audio, code, archive.", demoFileIcon},
	{409, "PathBar", "路径栏", "files", "Where you are, folder by folder. Click one to climb back up.", demoPathBar},
	{410, "MemoryCard", "记忆卡", "files", "Fifteen blocks, like the grey card of old. DELETE frees one.", demoMemoryCard},
	{411, "DiskUsage", "磁盘占用", "files", "Where the room went, stacked in one bar.", demoDiskUsage},
	{412, "UploadQueue", "上传队列", "files", "Files on their way up. × calls one back.", demoUploadQueue},
	{413, "Trash", "回收站", "files", "A bin that looks full when it is. EMPTY takes it out.", demoTrash},
}

// disk is the demo's drive, from C: down.
var disk = files.Dir{Name: "C:", Files: []files.File{{Name: "AUTOEXEC.BAT", Kind: files.Code, Size: 512}}, Dirs: []files.Dir{
	{Name: "GAMES", Dirs: []files.Dir{{Name: "HERO QUEST", Files: []files.File{
		{Name: "SAVE1.SAV", Kind: files.Archive, Size: 8192},
		{Name: "MAP.PNG", Kind: files.Image, Size: 48_000},
		{Name: "THEME.MID", Kind: files.Audio, Size: 21_504},
		{Name: "CHEATS.TXT", Kind: files.Text, Size: 1_337},
	}}}, Files: []files.File{{Name: "README.TXT", Kind: files.Text, Size: 2_048}}},
	{Name: "MUSIC", Files: []files.File{
		{Name: "BOSS FIGHT.MID", Kind: files.Audio, Size: 30_720},
		{Name: "CREDITS.MID", Kind: files.Audio, Size: 12_288},
	}},
}}

// path is the demo path's folder names, root first.
var path = []string{"C:", "GAMES", "HERO QUEST", "SAVES"}

type filesDemo struct {
	dir                          string
	file, grid, block, trash, at int
	saves                        []files.Save
	ups                          []files.Upload
	started                      time.Time
}

func filesState(c *ui.Context) *filesDemo {
	return ui.Local(c.Root(), demoKey("files"), func() filesDemo {
		return filesDemo{dir: "0/0/0", file: -1, grid: -1, trash: 3, at: len(path) - 1, saves: []files.Save{
			{Title: "HERO QUEST", Icon: icon.Sword},
			{Title: "SPACE RACE", Icon: icon.Star, Slot: pixel.Warn},
			{Title: "MONSTER FARM", Icon: icon.Heart, Slot: pixel.Bad},
			{Title: "PUZZLE TOWER", Icon: icon.Gem, Slot: pixel.Magic},
			{Title: "FISHING DAY", Icon: icon.Drop, Slot: pixel.Info},
		}, ups: []files.Upload{
			{Name: "SPEEDRUN.ZIP", Kind: files.Archive, Size: 9 << 20},
			{Name: "BOSS.PNG", Kind: files.Image, Size: 3 << 20},
			{Name: "LEVEL 9.MAP", Kind: files.Code, Size: 1 << 20, Sent: 1 << 19, Failed: true},
		}}
	})
}

func demoFileExplorer(c *ui.Context, log func(string)) {
	s := filesState(c)
	r := files.FileExplorer(c, "Explorer", disk, &s.dir, &s.file, files.FileExplorerOptions{})
	if r.Changed {
		log(fmt.Sprintf("folder %s, file %d", s.dir, s.file))
	}
	if r.Opened {
		log(fmt.Sprintf("open file %d", s.file))
	}
}

func demoFileGrid(c *ui.Context, log func(string)) {
	s := filesState(c)
	all := slices.Concat(disk.Dirs[0].Dirs[0].Files, []files.File{{Name: "SECRET LEVEL", Kind: files.Folder}})
	r := files.FileGrid(c, "Files", all, &s.grid, files.FileGridOptions{Height: 150})
	if r.Changed {
		log("chose " + all[s.grid].Name)
	}
	if r.Opened {
		log("open " + all[s.grid].Name)
	}
}

func demoFileIcon(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	kinds := []files.Kind{files.Folder, files.Text, files.Image, files.Audio, files.Code, files.Archive}
	ui.Column(c).Gap(k.D(6)).Children(func() {
		for _, scale := range []int{1, 2} {
			ui.Row(c).Gap(k.D(6)).Children(func() {
				for _, kind := range kinds {
					files.FileIcon(c, kind, files.FileIconOptions{Scale: scale})
				}
			})
		}
	})
}

func demoPathBar(c *ui.Context, log func(string)) {
	s := filesState(c)
	// Climbing up trims the path; the reset walks back down.
	r := files.PathBar(c, "Path", path[:s.at+1], files.PathBarOptions{})
	if r.Clicked >= 0 {
		s.at = r.Clicked
		log("up to " + path[s.at])
	}
	if s.at < len(path)-1 && input.Button(c, "BACK DOWN", input.ButtonOptions{Size: input.Small}).Clicked {
		s.at = len(path) - 1
		log("down to " + path[s.at])
	}
}

func demoMemoryCard(c *ui.Context, log func(string)) {
	s := filesState(c)
	r := files.MemoryCard(c, "Memory card 1", s.saves, &s.block, files.MemoryCardOptions{})
	if r.Changed {
		log(fmt.Sprintf("block %d", s.block+1))
	}
	if r.Deleted >= 0 {
		log("deleted " + s.saves[r.Deleted].Title)
		s.saves = append(s.saves[:r.Deleted:r.Deleted], s.saves[r.Deleted+1:]...)
	}
}

func demoDiskUsage(c *ui.Context, log func(string)) {
	files.DiskUsage(c, "Disk", []files.Use{
		{Name: "GAMES", Bytes: 41 << 30}, {Name: "SCREENSHOTS", Bytes: 6 << 30},
		{Name: "MUSIC", Bytes: 12 << 30}, {Name: "SAVES", Bytes: 300 << 20},
	}, 128<<30, files.DiskUsageOptions{Width: 220})
}

func demoUploadQueue(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := filesState(c)
	if s.started.IsZero() {
		s.started = k.Now()
	}
	// The demo sends a megabyte a second, row by row.
	sent := int64(k.Since(c, s.started, 4, 60)) << 18
	ups := make([]files.Upload, len(s.ups))
	for i, u := range s.ups {
		if !u.Failed {
			u.Sent = min(u.Size, sent)
			sent = max(0, sent-u.Size)
		}
		ups[i] = u
	}
	if r := files.UploadQueue(c, "Uploads", ups, files.UploadQueueOptions{}); r.Canceled >= 0 {
		log("canceled " + s.ups[r.Canceled].Name)
		s.ups = append(s.ups[:r.Canceled:r.Canceled], s.ups[r.Canceled+1:]...)
	}
}

func demoTrash(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := filesState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		if files.Trash(c, "Trash", s.trash, files.TrashOptions{}).Emptied {
			log(fmt.Sprintf("emptied %d", s.trash))
			s.trash = 0
		}
		if input.Button(c, "THROW AWAY", input.ButtonOptions{Size: input.Small}).Clicked {
			s.trash++
			log(fmt.Sprintf("threw one away: %d", s.trash))
		}
	})
}
