package game_test

import (
	"testing"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/data"
	"github.com/ELY-Open-Source/8bitMyGoUI/game"
	"github.com/ELY-Open-Source/8bitMyGoUI/icon"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

func TestPressStartBlinksAndReports(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		pressed := 0
		tt := view(300, 60, wd, func(c *ui.Context) {
			if game.PressStart(c, game.PressStartOptions{}).Pressed {
				pressed++
			}
		})
		var ink []int
		for range 2 {
			ink = append(ink, count(t, tt.Image(), 0, 0, 100, 20, pixel.Ink))
			wd.clk.t = wd.clk.t.Add(500 * time.Millisecond)
			tt.Frame()
		}
		if blinked := (ink[0] == 0) != (ink[1] == 0); blinked == reduce {
			t.Fatalf("reduced %v: PRESS START showed %v", reduce, ink)
		}
		tt.Click("PRESS START")
		tt.Key(0, ui.KeyEnter)
		if pressed != 2 {
			t.Fatalf("a click and Enter pressed %d times", pressed)
		}
	}
}

func TestTitleScreenStartsThenChooses(t *testing.T) {
	started, chosen := false, -1
	tt := view(600, 300, &world{clk: newClock()}, func(c *ui.Context) {
		if r := game.TitleScreen(c, "SLIME QUEST", &started, game.TitleScreenOptions{}); r.Chosen >= 0 {
			chosen = r.Chosen
		}
	})
	tt.Click("PRESS START")
	if !started || tt.HasText("PRESS START") {
		t.Fatalf("PRESS START left %q", tt.Texts())
	}
	// The menu takes focus; no save turns CONTINUE off.
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if chosen != -1 {
		t.Fatalf("CONTINUE with no save chose %d", chosen)
	}
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	if chosen != 2 {
		t.Fatalf("OPTIONS chose %d", chosen)
	}
	tt.Key(0, ui.KeyEscape)
	if started || !tt.HasText("PRESS START") {
		t.Fatal("Escape did not go back to PRESS START")
	}
	keeps(t, tt.Image())
	// Shown first thing, the menu still takes the keys.
	again, picked := true, -1
	tt = view(600, 300, &world{clk: newClock()}, func(c *ui.Context) {
		if r := game.TitleScreen(c, "SLIME QUEST", &again, game.TitleScreenOptions{}); r.Chosen >= 0 {
			picked = r.Chosen
		}
	})
	tt.Key(0, ui.KeyEnter)
	if picked != 0 {
		t.Fatalf("Enter on a fresh menu chose %d", picked)
	}
}

func TestSaveSlotsChooseByClickAndKeys(t *testing.T) {
	saves := []game.Save{{Place: "Woods", Played: 90 * time.Minute, Portrait: game.Hero}, {}, {Place: "Castle", Played: time.Minute}}
	sel, chosen := 0, -1
	tt := view(400, 300, &world{clk: newClock()}, func(c *ui.Context) {
		if r := game.SaveSlots(c, "Saves", saves, &sel, game.SaveSlotsOptions{}); r.Chosen >= 0 {
			chosen = r.Chosen
		}
	})
	r := find(t, tt, "Saves")
	// Cards are 32 dots tall with a 3-dot gap.
	tt.ClickAt(r.X+40, r.Y+float32(2*(2*35+10)))
	if sel != 2 || chosen != 2 {
		t.Fatalf("clicking card 3 chose %d, %d", sel, chosen)
	}
	tt.Key(0, ui.KeyUp)
	tt.Key(0, ui.KeyEnter)
	if chosen != 1 {
		t.Fatalf("Up, Enter chose %d", chosen)
	}
	keeps(t, tt.Image())
	// A long place name stays on its card.
	long := []game.Save{{Place: "The Very Long Forest of Many Wobbly Slimes", Played: time.Hour, Portrait: game.Hero}}
	first := 0
	tt = view(500, 100, &world{clk: newClock()}, func(c *ui.Context) {
		game.SaveSlots(c, "Long", long, &first, game.SaveSlotsOptions{})
	})
	if n := count(t, tt.Image(), 180, 0, 70, 40, pixel.Ink); n != 0 {
		t.Fatalf("%d place dots past the card", n)
	}
	panics(t, "an empty save with play time", func() {
		view(400, 300, &world{clk: newClock()}, func(c *ui.Context) {
			game.SaveSlots(c, "x", []game.Save{{Played: time.Hour}}, &first, game.SaveSlotsOptions{})
		})
	})
}

func TestOptionsMenuStepsAndClamps(t *testing.T) {
	speed, bgm, row, changes := 0, 9, 0, 0
	tt := view(500, 100, &world{clk: newClock()}, func(c *ui.Context) {
		if game.OptionsMenu(c, "Options", []game.Setting{{Name: "TEXT SPEED", Choices: []string{"SLOW", "FAST"}, Value: &speed}, {Name: "BGM", Max: 10, Value: &bgm}}, &row, game.OptionsMenuOptions{}).Changed {
			changes++
		}
	})
	// The focusing click lands on the second row.
	tt.Click("Options")
	for _, k := range []ui.Key{ui.KeyUp, ui.KeyLeft, ui.KeyRight, ui.KeyRight, ui.KeyDown, ui.KeyRight, ui.KeyRight, ui.KeyDown, ui.KeyLeft} {
		tt.Key(0, k)
	}
	// SLOW stays, one step to FAST, BGM tops at 10.
	if speed != 1 || bgm != 9 || row != 1 || changes != 3 {
		t.Fatalf("speed %d, bgm %d, row %d, %d changes", speed, bgm, row, changes)
	}
	keeps(t, tt.Image())
	panics(t, "a setting with both kinds", func() {
		view(500, 100, &world{clk: newClock()}, func(c *ui.Context) {
			game.OptionsMenu(c, "x", []game.Setting{{Name: "x", Choices: []string{"A"}, Max: 3, Value: &speed}}, &row, game.OptionsMenuOptions{})
		})
	})
}

func TestOptionsMenuKeysDoNotReplayClicks(t *testing.T) {
	level, row := 5, 0
	tt := view(500, 100, &world{clk: newClock()}, func(c *ui.Context) {
		game.OptionsMenu(c, "Options", []game.Setting{{Name: "BGM", Max: 10, Value: &level}}, &row, game.OptionsMenuOptions{})
	})
	r := find(t, tt, "Options")
	// ▶ follows the hand, name, ◀ and ten blocks.
	bx := pixel.Hand.W + 4 + pixel.Measure("BGM") + 10 + pixel.Measure("◀") + 3 + 40 + 3
	tt.ClickAt(r.X+float32(2*(bx+2)), r.Y+16)
	tt.Key(0, ui.KeyEnter)
	tt.Key(0, ui.KeySpace)
	if level != 6 {
		t.Fatalf("a click on ▶, Enter and Space left BGM at %d", level)
	}
}

func TestGameOverDithersInAndAnswers(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		quit := false
		tt := view(600, 300, wd, func(c *ui.Context) {
			quit = quit || game.GameOver(c, game.GameOverOptions{}).Quit
		})
		r := find(t, tt, "GAME OVER")
		x, y := dots(r)
		first := count(t, tt.Image(), x, y, int(r.W)/2, int(r.H)/2, pixel.Bad)
		wd.clk.t = wd.clk.t.Add(3 * time.Second)
		tt.Frame()
		full := count(t, tt.Image(), x, y, int(r.W)/2, int(r.H)/2, pixel.Bad)
		if full == 0 || (first == full) != reduce {
			t.Fatalf("reduced %v: GAME OVER drew %d dots, then %d", reduce, first, full)
		}
		tt.Key(0, ui.KeyDown)
		tt.Key(0, ui.KeyEnter)
		if !quit {
			t.Fatalf("reduced %v: Down, Enter did not quit", reduce)
		}
	}
}

func TestCreditsRollOnceAndAgainOnCue(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		wd := &world{clk: newClock(), reduce: reduce}
		cue, done := 0, 0
		tt := view(400, 400, wd, func(c *ui.Context) {
			if game.Credits(c, "Credits", []game.Credit{{Heading: "DIRECTOR", Names: []string{"Aria"}}, {Heading: "CAST", Names: []string{"A", "B", "C", "D", "E", "F"}}}, game.CreditsOptions{Cue: cue}).Done {
				done++
			}
		})
		// Rolling, it starts empty; reduced, it shows all ten lines.
		if lit := count(t, tt.Image(), 0, 0, 160, 96, pixel.Warn) > 0; lit != reduce || (done == 1) != reduce {
			t.Fatalf("reduced %v: first frame lit %v, done %d", reduce, lit, done)
		}
		if h := find(t, tt, "Credits").H; (h == 2*160) != reduce {
			t.Fatalf("reduced %v: the window is %v px tall", reduce, h)
		}
		for range 3 {
			wd.clk.t = wd.clk.t.Add(5 * time.Second)
			tt.Frame()
		}
		if done != 1 || count(t, tt.Image(), 0, 0, 160, 96, pixel.Ink) == 0 {
			t.Fatalf("reduced %v: done %d times; the last name not resting", reduce, done)
		}
		cue = 1
		tt.Frame()
		for range 3 {
			wd.clk.t = wd.clk.t.Add(5 * time.Second)
			tt.Frame()
		}
		if done != 2 {
			t.Fatalf("reduced %v: a new cue rolled %d times in all", reduce, done)
		}
		// An old cue mid-roll rolls anew; reduced, both end.
		cue = 2
		tt.Frame()
		cue = 1
		tt.Frame()
		for range 3 {
			wd.clk.t = wd.clk.t.Add(5 * time.Second)
			tt.Frame()
		}
		if want := map[bool]int{false: 3, true: 4}[reduce]; done != want {
			t.Fatalf("reduced %v: going back to cue 1 rolled %d times in all", reduce, done)
		}
	}
}

func TestEndingCardDithersIn(t *testing.T) {
	wd := &world{clk: newClock()}
	tt := view(400, 200, wd, func(c *ui.Context) { game.EndingCard(c, game.EndingCardOptions{}) })
	first := count(t, tt.Image(), 0, 0, 200, 80, pixel.Back)
	wd.clk.t = wd.clk.t.Add(time.Second)
	tt.Frame()
	half := count(t, tt.Image(), 0, 0, 200, 80, pixel.Back)
	wd.clk.t = wd.clk.t.Add(2 * time.Second)
	tt.Frame()
	if full := count(t, tt.Image(), 0, 0, 200, 80, pixel.Back); !(first < half && half < full) {
		t.Fatalf("THE END dots %d, %d, %d", first, half, full)
	}
}

func TestEndingCardRefusesNoRoom(t *testing.T) {
	panics(t, "a card narrower than THE END", func() {
		view(400, 200, &world{clk: newClock()}, func(c *ui.Context) { game.EndingCard(c, game.EndingCardOptions{Width: 40}) })
	})
}

func TestLoadingTipsTurnEveryFourSeconds(t *testing.T) {
	wd := &world{clk: newClock()}
	tt := view(500, 100, wd, func(c *ui.Context) {
		game.LoadingTips(c, []string{"Tip one.", "Tip two."}, 0.5, game.LoadingTipsOptions{})
	})
	a := tt.HasText("TIP: Tip one.") || tt.HasText("TIP: Tip two.")
	before := tt.HasText("TIP: Tip one.")
	wd.clk.t = wd.clk.t.Add(4 * time.Second)
	tt.Frame()
	if !a || tt.HasText("TIP: Tip one.") == before {
		t.Fatalf("tips did not turn: %q", tt.Texts())
	}
}

func TestInteractPromptBobsOnceASecond(t *testing.T) {
	wd := &world{clk: newClock()}
	tt := view(300, 200, wd, func(c *ui.Context) {
		var target ui.Element
		ui.Column(c).Padding(60).Children(func() { target = ui.Box(c).Size(20, 20) })
		game.InteractPrompt(c, target, "OPEN", game.InteractPromptOptions{})
	})
	start := wd.clk.t
	var ys []float32
	for _, at := range []time.Duration{100, 350, 600} {
		wd.clk.t = start.Add(at * time.Millisecond)
		tt.Frame()
		// The bottom edge bobs; the top may sit clipped.
		r := find(t, tt, "A OPEN")
		ys = append(ys, r.Y+r.H)
	}
	if ys[0] != ys[1] || ys[1] == ys[2] {
		t.Fatalf("the prompt sat at %v at .1, .35 and .6 s", ys)
	}
}

func TestHealthBarTrailStepsAtFifteen(t *testing.T) {
	wd, hp := &world{clk: newClock()}, 100
	tt := view(300, 100, wd, func(c *ui.Context) { game.HealthBar(c, "HP", hp, 100, game.HealthBarOptions{}) })
	hp = 50
	tt.Frame()
	wd.clk.t = wd.clk.t.Add(350 * time.Millisecond)
	// Frames each 10 ms show two steps at most.
	seen := map[string]bool{}
	for range 6 {
		wd.clk.t = wd.clk.t.Add(10 * time.Millisecond)
		tt.Frame()
		seen[string(tt.Image().Pix)] = true
	}
	if len(seen) > 2 {
		t.Fatalf("the trail drew %d pictures in 60 ms", len(seen))
	}
}

func TestDamageNumbersRiseInWholeFrames(t *testing.T) {
	wd := &world{clk: newClock()}
	hits := []game.Hit{{Value: 999, X: 60, Y: 60}}
	tt := view(300, 200, wd, func(c *ui.Context) { game.DamageNumbers(c, "Damage", &hits, game.DamageNumbersOptions{}) })
	wd.clk.t = wd.clk.t.Add(500 * time.Millisecond)
	tt.Frame()
	seen := map[string]bool{}
	for range 7 {
		seen[string(tt.Image().Pix)] = true
		wd.clk.t = wd.clk.t.Add(10 * time.Millisecond)
		tt.Frame()
	}
	if len(seen) > 2 {
		t.Fatalf("%d pictures in 60 ms", len(seen))
	}
}

func TestSpinsAndHitsRestartWhenTheClockStepsBack(t *testing.T) {
	wd := &world{clk: newClock()}
	faces := []data.Item{{Sprite: icon.Star, Name: "Star"}, {Sprite: icon.Coin, Name: "Coin"}}
	spin, done := 0, 0
	tt := view(300, 100, wd, func(c *ui.Context) {
		if game.GachaReel(c, "Reels", faces, [3]int{1, 1, 1}, spin, game.GachaReelOptions{}).Done {
			done++
		}
	})
	spin = 1
	tt.Frame()
	// An hour back, the spin times from now and lands.
	wd.clk.t = wd.clk.t.Add(-time.Hour)
	tt.Frame()
	for range 30 {
		wd.clk.t = wd.clk.t.Add(200 * time.Millisecond)
		tt.Frame()
	}
	if done != 1 {
		t.Fatalf("the reels landed %d times", done)
	}
	hits := []game.Hit{{Value: 7, X: 60, Y: 40}}
	tt = view(300, 200, wd, func(c *ui.Context) { game.DamageNumbers(c, "Damage", &hits, game.DamageNumbersOptions{}) })
	fresh := string(tt.Image().Pix)
	wd.clk.t = wd.clk.t.Add(-time.Second)
	tt.Frame()
	if string(tt.Image().Pix) != fresh {
		t.Fatal("a hit moved as the clock stepped back")
	}
}
