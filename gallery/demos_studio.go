package gallery

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/media"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var studioEntries = []Entry{
	{356, "PianoKeys", "钢琴键", "media", "Click, or type A–K. W E T Y U play the black keys.", demoPianoKeys},
	{357, "StepSequencer", "步进音序器", "media", "Four tracks, sixteen steps. Press play and nod along.", demoStepSequencer},
	{358, "Tracker", "音轨编辑", "media", "Type A–K to write a note, Delete for a rest. Beats every four rows.", demoTracker},
	{359, "Waveform", "波形", "media", "Every sample a line, up or down from zero.", demoWaveform},
	{360, "VideoFrame", "视频框", "media", "A 16:9 screen with play and a scrubber. The slime does the rest.", demoVideoFrame},
	{361, "ImageCompare", "对比滑块", "media", "Drag the grip: sketch on the left, paint on the right.", demoImageCompare},
	{362, "Gallery", "相册", "media", "Arrows walk the thumbnails; the big picture follows.", demoGallery},
	{363, "Filmstrip", "胶片条", "media", "A sprite's frames on film, sprocket holes and all.", demoFilmstrip},
	{364, "RecordButton", "录制按钮", "media", "The dot blinks while it records, and the time runs.", demoRecordButton},
	{365, "Tuner", "调频器", "media", "Drag the needle across FM. A few stations are on the air.", demoTuner},
}

func demoPianoKeys(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := mediaState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		for _, n := range media.PianoKeys(c, "Piano", media.PianoKeysOptions{}).Struck {
			s.notes = append(s.notes, n)
			log("♪ " + n)
		}
		played := s.notes[max(0, len(s.notes)-8):]
		pixel.Text(c, "♪ "+strings.Join(played, " "), pixel.TextOptions{Slot: pixel.Faint})
	})
}

func demoStepSequencer(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := mediaState(c)
	tracks := []string{"KICK", "SNARE", "HAT", "CLAP"}
	ui.Column(c).Gap(k.D(6)).Children(func() {
		if media.StepSequencer(c, "Beat", tracks, s.pattern, media.StepSequencerOptions{Playing: s.beat}).Changed {
			log("pattern changed")
		}
		if input.Toggle(c, "PLAY", &s.beat, input.ToggleOptions{}).Changed {
			log(fmt.Sprintf("beat playing %v", s.beat))
		}
	})
}

// melody is a pattern with a tune in it.
func melody() []media.Row {
	rows := make([]media.Row, 16)
	for i, n := range map[int]string{0: "C-4", 2: "E-4", 4: "G-4", 6: "E-4", 8: "F-4", 10: "A-4", 12: "G-4", 14: "C-5"} {
		rows[i] = media.Row{Note: n, Instrument: 1 + i/8, Volume: 64 - i*2}
	}
	return rows
}

func demoTracker(c *ui.Context, log func(string)) {
	s := mediaState(c)
	if media.Tracker(c, "Pattern", s.rows, &s.row, media.TrackerOptions{}).Changed {
		log(fmt.Sprintf("pattern edited at row %02d", s.row))
	}
}

func demoWaveform(c *ui.Context, log func(string)) {
	// A plucked string: a fast wiggle fading out.
	samples := make([]float64, 72)
	for i := range samples {
		f := float64(i)
		samples[i] = math.Sin(f*0.7) * math.Exp(-f/24)
	}
	media.Waveform(c, "Pluck", samples, media.WaveformOptions{Height: 40})
}

func demoVideoFrame(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := mediaState(c)
	const length = 20 * time.Second
	spin(c, &s.filmed, s.video, &s.videoAt, length)
	at := s.videoAt
	screen := func() {
		ui.Box(c).Size(k.D(160), k.D(90)).Draw(func(p *ui.Painter, r ui.Rect) {
			g := k.Canvas(p, r)
			g.Rect(0, 0, g.W, g.H, pixel.Info)
			g.Rect(0, 72, g.W, 18, pixel.Good)
			g.Dither(0, 72, g.W, 18, pixel.Good, pixel.Shade, 4)
			// The slime hops right with the film's own clock.
			hop := int(at / (time.Second / 4))
			x := int(float64(at) / float64(length) * 128)
			g.Sprite(x, 48-hop%2*4, slime, hop%2, pixel.Ink, 2)
		})
	}
	res := media.VideoFrame(c, "Slime documentary", &s.video, &s.videoAt, length, media.VideoFrameOptions{}, screen)
	if res.Changed {
		log(fmt.Sprintf("video playing %v at %s", s.video, s.videoAt.Truncate(time.Second)))
	}
}

func demoImageCompare(c *ui.Context, log func(string)) {
	s := mediaState(c)
	after := scene(pixel.Warn, pixel.Info, pixel.Good, 24, 5)
	// The sketch is the painting in one faint ink.
	before := pixel.SpriteOf(after.W, after.H, func(x, y int) (pixel.Slot, bool) {
		if sl, _ := after.At(0, x, y, pixel.Ink); sl == pixel.Info {
			return pixel.Lit, true
		}
		return pixel.Faint, true
	})
	if media.ImageCompare(c, "Sketch and painting", before, after, &s.split, media.ImageCompareOptions{Scale: 4}).Changed {
		log(fmt.Sprintf("split at %.0f%%", s.split*100))
	}
}

var views = []struct {
	name string
	pic  *pixel.Sprite
}{
	{"Morning", scene(pixel.Warn, pixel.AccentLit, pixel.Good, 4, 10)},
	{"Noon", scene(pixel.Warn, pixel.Info, pixel.Good, 16, 3)},
	{"Dusk", scene(pixel.Bad, pixel.Accent, pixel.Magic, 27, 11)},
	{"Night", scene(pixel.Lit, pixel.Well, pixel.Shade, 22, 4)},
}

func demoGallery(c *ui.Context, log func(string)) {
	s := mediaState(c)
	pics, names := make([]*pixel.Sprite, len(views)), make([]string, len(views))
	for i, v := range views {
		pics[i], names[i] = v.pic, v.name
	}
	if media.Gallery(c, "Views", pics, names, &s.chosen, media.GalleryOptions{}).Changed {
		log("showing " + names[s.chosen])
	}
}

var coin = pixel.MustSprite(`
	..eeee..
	.eyyyye.
	eyylyyye
	eyylyyye
	eyylyyye
	eyylyyye
	.eyyyye.
	..eeee..

	...ee...
	..eyye..
	..eyle..
	..eyle..
	..eyle..
	..eyle..
	..eyye..
	...ee...

	...ee...
	...ee...
	...ee...
	...ee...
	...ee...
	...ee...
	...ee...
	...ee...

	...ee...
	..eyye..
	..elye..
	..elye..
	..elye..
	..elye..
	..eyye..
	...ee...
`)

func demoFilmstrip(c *ui.Context, log func(string)) {
	media.Filmstrip(c, "Coin spin", coin, media.FilmstripOptions{Scale: 3})
}

func demoRecordButton(c *ui.Context, log func(string)) {
	s := mediaState(c)
	if media.RecordButton(c, "Record", &s.rec, media.RecordButtonOptions{}).Changed {
		log(fmt.Sprintf("recording %v", s.rec))
	}
}

var stations = map[float64]string{88.1: "PIXEL FM", 95.5: "CHIPTUNE RADIO", 101.1: "8-BIT HITS", 104.3: "BOSS BATTLE FM"}

func demoTuner(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := mediaState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		if media.Tuner(c, "Radio", &s.mhz, media.TunerOptions{}).Changed {
			log(fmt.Sprintf("%.1f MHz", s.mhz))
		}
		on := "· · · static · · ·"
		if name, ok := stations[s.mhz]; ok {
			on = "♪ " + name
		}
		pixel.Text(c, on, pixel.TextOptions{Slot: pixel.Faint})
	})
}
