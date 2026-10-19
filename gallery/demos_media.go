package gallery

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/media"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var mediaEntries = []Entry{
	{346, "AudioPlayer", "音乐播放器", "media", "A whole deck in one box. Stop rewinds to the top.", demoAudioPlayer},
	{347, "Transport", "播放控制", "media", "Back, play or pause, stop, skip: the four buttons every deck agrees on.", demoTransport},
	{348, "Scrubber", "进度定位", "media", "Drag to seek, or nudge with arrows. Time at both ends.", demoScrubber},
	{349, "Volume", "音量", "media", "Eight bars; click one or use arrows. The speaker mutes.", demoVolume},
	{350, "Equalizer", "均衡器", "media", "Five faders, ±12 dB. Bass boost at your own risk.", demoEqualizer},
	{351, "Visualizer", "频谱", "media", "Bars bounce out from the middle, hottest at the tips.", demoVisualizer},
	{352, "Playlist", "播放列表", "media", "♪ marks the track that plays. Enter plays the one you pick.", demoPlaylist},
	{353, "Cassette", "磁带", "media", "Reels turn while it plays; watch the tape cross over.", demoCassette},
	{354, "Cartridge", "卡带", "media", "A cart with a sticker. Blow on the contacts before use.", demoCartridge},
	{355, "Vinyl", "黑胶", "media", "Turns in eight frames; the arm swings onto the groove.", demoVinyl},
}

type mediaDemo struct {
	player             media.Player
	track, pick        int
	deck               string
	seek               time.Duration
	level              int
	muted              bool
	gains              []float64
	tape, spin, video  bool
	tapeAt, videoAt    time.Duration
	notes              []string
	beat               bool
	pattern            [][16]bool
	rows               []media.Row
	row, chosen        int
	split, mhz         float64
	rec                bool
	ran, taped, filmed time.Time
}

var tracks = []media.Track{
	{Title: "Overworld Theme", Artist: "The Square Waves", Length: 2*time.Minute + 47*time.Second},
	{Title: "Boss Battle", Artist: "Chip Wizards", Length: 3*time.Minute + 12*time.Second},
	{Title: "Save Point Lullaby", Artist: "Triangle Wave", Length: 1*time.Minute + 58*time.Second},
	{Title: "Credits Roll", Artist: "Pulse & Noise", Length: 4*time.Minute + 5*time.Second},
}

func mediaState(c *ui.Context) *mediaDemo {
	return ui.Local(c.Root(), demoKey("media"), func() mediaDemo {
		return mediaDemo{
			player: media.Player{Level: 6}, pick: 1, deck: "STOPPED", seek: 72 * time.Second, level: 5,
			gains:   []float64{6, 2, 0, 3, 5},
			pattern: [][16]bool{beat("x...x...x...x..."), beat("....x.......x..."), beat("x.x.x.x.x.x.x.xx"), beat("............x...")},
			rows:    melody(), split: 0.5, mhz: 98.3,
		}
	})
}

// beat reads a pattern: x hits, . rests.
func beat(s string) [16]bool {
	var b [16]bool
	for i, r := range s {
		b[i] = r == 'x'
	}
	return b
}

// spin moves *at while playing, as an audio engine would.
func spin(c *ui.Context, last *time.Time, playing bool, at *time.Duration, length time.Duration) {
	k := pixel.K(c)
	if playing && !last.IsZero() {
		// A clock stepping back plays nothing backward.
		*at = (*at + max(0, k.Now().Sub(*last))) % length
	}
	*last = time.Time{}
	if playing {
		*last = k.Now()
		k.Tick(c, 4)
	}
}

func demoAudioPlayer(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := mediaState(c)
	t := tracks[s.track]
	spin(c, &s.ran, s.player.Playing, &s.player.At, t.Length)
	ui.Box(c).Width(k.D(240)).Children(func() {
		res := media.AudioPlayer(c, "Now playing", t, &s.player, media.AudioPlayerOptions{})
		switch {
		case res.Previous:
			s.track = (s.track + len(tracks) - 1) % len(tracks)
		case res.Next:
			s.track = (s.track + 1) % len(tracks)
		}
		if res.Previous || res.Next {
			s.player.At = 0
			log("♪ " + tracks[s.track].Title)
		}
		if res.Changed {
			log(fmt.Sprintf("playing %v at %s, volume %d, muted %v", s.player.Playing, s.player.At.Truncate(time.Second), s.player.Level, s.player.Muted))
		}
	})
}

func demoTransport(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := mediaState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		playing := s.deck == "PLAYING"
		res := media.Transport(c, "Deck", &playing, media.TransportOptions{})
		switch {
		case res.Stopped:
			s.deck = "STOPPED"
		case res.Changed && playing:
			s.deck = "PLAYING"
		case res.Changed:
			s.deck = "PAUSED"
		case res.Previous:
			log("PREVIOUS")
		case res.Next:
			log("NEXT")
		}
		if res.Changed || res.Stopped {
			log(s.deck)
		}
		pixel.Text(c, s.deck, pixel.TextOptions{Slot: pixel.Faint})
	})
}

func demoScrubber(c *ui.Context, log func(string)) {
	s := mediaState(c)
	if media.Scrubber(c, "Seek", &s.seek, tracks[1].Length, media.ScrubberOptions{Width: 120}).Changed {
		log("seek " + s.seek.String())
	}
}

func demoVolume(c *ui.Context, log func(string)) {
	s := mediaState(c)
	if media.Volume(c, "Volume", &s.level, &s.muted, media.VolumeOptions{}).Changed {
		log(fmt.Sprintf("volume %d of 8, muted %v", s.level, s.muted))
	}
}

func demoEqualizer(c *ui.Context, log func(string)) {
	s := mediaState(c)
	bands := []string{"60", "250", "1K", "4K", "12K"}
	if b := media.Equalizer(c, "Equalizer", bands, s.gains, media.EqualizerOptions{}).Band; b >= 0 {
		log(fmt.Sprintf("%s Hz %+.1f dB", bands[b], s.gains[b]))
	}
}

func demoVisualizer(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	// A made-up song: each bar sways at its own pace.
	t := 0.0
	if !k.Reduce() {
		t = float64(k.Tick(c, 12)) / 12
	}
	levels := make([]float64, 24)
	for i := range levels {
		f := float64(i)
		levels[i] = (0.55 + 0.45*math.Sin(t*(2+f/5)+f)) * (1 - f/32)
	}
	media.Visualizer(c, "Visualizer", levels, media.VisualizerOptions{Height: 40})
}

func demoPlaylist(c *ui.Context, log func(string)) {
	s := mediaState(c)
	res := media.Playlist(c, "Playlist", tracks, s.track, &s.pick, media.PlaylistOptions{Height: 80})
	if res.Submitted {
		s.track, s.player.At = s.pick, 0
		log("♪ " + tracks[s.track].Title)
	} else if res.Changed {
		log("picked " + tracks[s.pick].Title)
	}
}

func demoCassette(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := mediaState(c)
	const side = 45 * time.Second
	spin(c, &s.taped, s.tape, &s.tapeAt, side)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		media.Cassette(c, "Mixtape, side A", media.CassetteOptions{Title: "Mixtape 86", Playing: s.tape, Progress: float64(s.tapeAt) / float64(side)})
		if input.Toggle(c, "PLAY", &s.tape, input.ToggleOptions{}).Changed {
			log(fmt.Sprintf("tape playing %v", s.tape))
		}
	})
}

func demoCartridge(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(8)).Children(func() {
		media.Cartridge(c, "Pixel Quest cartridge", media.CartridgeOptions{Title: "Pixel Quest", Art: scene(pixel.Warn, pixel.Info, pixel.Good, 8, 4)})
		media.Cartridge(c, "Night Run cartridge", media.CartridgeOptions{Title: "Night Run 2", Art: scene(pixel.Lit, pixel.Magic, pixel.Shade, 26, 3)})
	})
}

func demoVinyl(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := mediaState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		media.Vinyl(c, "Record", media.VinylOptions{Playing: s.spin})
		if input.Toggle(c, "PLAY", &s.spin, input.ToggleOptions{}).Changed {
			log(fmt.Sprintf("record playing %v", s.spin))
		}
	})
}

// scene is hills under a sky, sun at x, y.
func scene(sun, sky, hill pixel.Slot, x, y int) *pixel.Sprite {
	return pixel.SpriteOf(32, 18, func(i, j int) (pixel.Slot, bool) {
		switch {
		case (i-x)*(i-x)+(j-y)*(j-y) <= 5:
			return sun, true
		case float64(j) > 12+2*math.Sin(float64(i)/4):
			return hill, true
		}
		return sky, true
	})
}
