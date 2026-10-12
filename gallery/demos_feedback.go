package gallery

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/feedback"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var feedbackEntries = []Entry{
	{106, "Progress", "进度条", "feedback", "A bar that fills chunk by chunk. Below zero, a block runs back and forth.", demoProgress},
	{107, "Spinner", "转圈器", "feedback", "Eight frames of beads going round. It keeps turning under Reduce Motion: it means work.", demoSpinner},
	{108, "Meter", "计量条", "feedback", "A level that turns Warn past a mark and Bad past the next.", demoMeter},
	{109, "Skeleton", "骨架屏", "feedback", "Dithered stand-ins with a shine sweeping across in steps.", demoSkeleton},
	{110, "Alert", "提示框", "feedback", "A framed callout: info, success, warning or danger, with an optional ×.", demoAlert},
	{111, "Banner", "横幅", "feedback", "A full-width strip with actions and a ×.", demoBanner},
	{112, "EmptyState", "空状态", "feedback", "NO ITEMS over an empty chest, and something to do about it.", demoEmptyState},
	{113, "ErrorState", "错误态", "feedback", "GAME OVER, a CONTINUE? countdown and RETRY. The gallery's own errors use it.", demoErrorState},
	{114, "SuccessState", "通关", "feedback", "STAGE CLEAR! The stars pop in one by one.", demoSuccessState},
	{115, "LoadingScreen", "加载画面", "feedback", "NOW LOADING with stepping dots and a tip every four seconds.", demoLoadingScreen},
	{116, "DotsLoader", "点点加载", "feedback", "Three dots hopping in turn.", demoDotsLoader},
	{117, "BarsLoader", "柱状加载", "feedback", "Five bars rising and falling, a step apart.", demoBarsLoader},
	{118, "Hourglass", "沙漏", "feedback", "Sand falls in steps, then the glass flips.", demoHourglass},
	{119, "Toast", "吐司", "feedback", "Short messages at the bottom, three at a time, each for a few seconds. Some can be undone.", demoToast},
	{120, "InlineMessage", "行内提示", "feedback", "A small warning or check note under a field.", demoInlineMessage},
	{121, "CountUp", "计数", "feedback", "Numbers count up in steps, the way scores should.", demoCountUp},
	{122, "Pulse", "脉冲", "feedback", "A dithered ring grows around what wants attention.", demoPulse},
	{123, "LevelUp", "升级", "feedback", "LEVEL UP! slides in with sparkles when the level rises.", demoLevelUp},
	{124, "Achievement", "成就", "feedback", "A trophy card slides in, waits, and slides away.", demoAchievement},
	{125, "ComboCounter", "连击", "feedback", "×12 COMBO! It pops a dot bigger with every hit.", demoComboCounter},
	{126, "ContinuePrompt", "继续？", "feedback", "CONTINUE? YES or NO before the count runs out.", demoContinuePrompt},
	{127, "SaveIndicator", "存档指示", "feedback", "A floppy spins while saving, then SAVED with a check.", demoSaveIndicator},
	{128, "BusyOverlay", "加载遮罩", "feedback", "A checker veil and spinner over busy content; it takes the pointer and the focus.", demoBusyOverlay},
	{129, "ReadyGo", "预备开始", "feedback", "3, 2, 1, GO! in big digits, a second each.", demoReadyGo},
	{130, "Shake", "抖动", "feedback", "Jolts its child side to side, for a wrong answer.", demoShake},
}

type feedbackDemo struct {
	progress, hp                    float64
	closed                          bool
	tries, replays, prompts, score  int
	level, cue, combo, start, shake int
	toasts                          []feedback.Toast
	hero, secret, motto             string
	saving, busy                    time.Time
}

func feedbackState(c *ui.Context) *feedbackDemo {
	return ui.Local(c.Root(), demoKey("feedback"), func() feedbackDemo {
		return feedbackDemo{progress: 0.3, hp: 0.5, level: 4}
	})
}

func demoProgress(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		feedback.Progress(c, "Download", s.progress, feedback.ProgressOptions{Percent: true})
		feedback.Progress(c, "Searching", -1, feedback.ProgressOptions{})
		ui.Row(c).Gap(k.D(4)).Children(func() {
			if input.Button(c, "+10%", input.ButtonOptions{Size: input.Small}).Clicked {
				s.progress = min(1, s.progress+0.1)
				log(fmt.Sprintf("%.0f%%", s.progress*100))
			}
			if input.Button(c, "RESET", input.ButtonOptions{Size: input.Small}).Clicked {
				s.progress = 0
				log("reset")
			}
		})
	})
}

func demoSpinner(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(12)).AlignItems(ui.Center).Children(func() {
		for sc := 1; sc <= 3; sc++ {
			feedback.Spinner(c, fmt.Sprintf("Loading ×%d", sc), feedback.SpinnerOptions{Scale: sc})
		}
	})
}

func demoMeter(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		feedback.Meter(c, "Heat", s.hp, feedback.MeterOptions{Warning: 0.6, Critical: 0.85})
		if input.Slider(c, "Heat level", &s.hp, input.SliderOptions{}).Changed {
			log(fmt.Sprintf("heat %.2f", s.hp))
		}
	})
}

func demoSkeleton(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(10)).Children(func() {
		feedback.Skeleton(c, feedback.SkeletonOptions{Avatar: true})
		feedback.Skeleton(c, feedback.SkeletonOptions{Lines: 2, Width: 120})
	})
}

func demoAlert(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		feedback.Alert(c, "Did you know?", "Pressing P paints the world anew.", feedback.AlertOptions{Kind: feedback.Info})
		feedback.Alert(c, "Quest complete", "The princess was in this castle.", feedback.AlertOptions{Kind: feedback.Success})
		feedback.Alert(c, "Low on potions", "Two left. Shop before the boss.", feedback.AlertOptions{Kind: feedback.Warning})
		if !s.closed && feedback.Alert(c, "Save corrupted", "Just kidding. Close me.", feedback.AlertOptions{Kind: feedback.Danger, Closable: true}).Closed {
			s.closed = true
			log("closed")
		}
	})
}

func demoBanner(c *ui.Context, log func(string)) {
	r := feedback.Banner(c, "A new patch is out: slimes now bounce higher.", feedback.BannerOptions{Kind: feedback.Info, Actions: []string{"UPDATE", "LATER"}, Dismissable: true})
	if r.Action != "" {
		log(r.Action)
	}
	if r.Dismissed {
		log("dismissed")
	}
}

func demoEmptyState(c *ui.Context, log func(string)) {
	if feedback.EmptyState(c, "Your inventory is empty.", feedback.EmptyStateOptions{Action: "VISIT SHOP"}).Acted {
		log("off to the shop")
	}
}

func demoErrorState(c *ui.Context, log func(string)) {
	s := feedbackState(c)
	// A new key starts a fresh countdown.
	ui.Box(c).Key(s.tries).Children(func() {
		if feedback.ErrorState(c, "The slime won this round.", feedback.ErrorStateOptions{}).Continue {
			s.tries++
			log(fmt.Sprintf("continue #%d", s.tries))
		}
	})
}

func demoSuccessState(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	ui.Column(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
		ui.Box(c).Key(s.replays).Children(func() {
			feedback.SuccessState(c, "All 99 coins found.", feedback.SuccessStateOptions{})
		})
		if input.Button(c, "REPLAY", input.ButtonOptions{Size: input.Small}).Clicked {
			s.replays++
			log("replay")
		}
	})
}

func demoLoadingScreen(c *ui.Context, log func(string)) {
	feedback.LoadingScreen(c, feedback.LoadingScreenOptions{Tips: []string{
		"Tip: slimes fear the color red.", "Tip: press P to change the palette.", "Tip: save often. The dragon does not wait.",
	}})
}

func demoDotsLoader(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		feedback.DotsLoader(c, "Thinking", feedback.DotsLoaderOptions{})
		pixel.Text(c, "The wizard is thinking", pixel.TextOptions{Slot: pixel.Faint})
	})
}

func demoBarsLoader(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		feedback.BarsLoader(c, "Tuning", feedback.BarsLoaderOptions{})
		pixel.Text(c, "Tuning the chiptune", pixel.TextOptions{Slot: pixel.Faint})
	})
}

func demoHourglass(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		feedback.Hourglass(c, "Waiting", feedback.HourglassOptions{})
		pixel.Text(c, "The turn passes", pixel.TextOptions{Slot: pixel.Faint})
	})
}

func demoToast(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	ui.Row(c).Gap(k.D(4)).Children(func() {
		if input.Button(c, "SAVE GAME", input.ButtonOptions{Size: input.Small}).Clicked {
			s.toasts = append(s.toasts, feedback.Toast{Message: "Game saved."})
			log("toast: saved")
		}
		if input.Button(c, "DROP SWORD", input.ButtonOptions{Size: input.Small}).Clicked {
			s.toasts = append(s.toasts, feedback.Toast{Message: "You dropped the sword.", Action: "UNDO", For: 5 * time.Second})
			log("toast: dropped")
		}
	})
	if r := feedback.Toasts(c, &s.toasts, feedback.ToastsOptions{}); r.Acted != nil {
		log(r.Acted.Action + ": " + r.Acted.Message)
	}
}

func demoInlineMessage(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	ui.Column(c).Gap(k.D(3)).Children(func() {
		if input.TextInput(c, "Hero name", &s.hero, input.TextInputOptions{Placeholder: "Your name"}).Changed {
			log("name " + s.hero)
		}
		if s.hero == "" {
			feedback.InlineMessage(c, "A hero needs a name.", feedback.InlineMessageOptions{Kind: feedback.Danger})
		} else {
			feedback.InlineMessage(c, "That name is free!", feedback.InlineMessageOptions{Kind: feedback.Success})
		}
	})
}

func demoCountUp(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		feedback.CountUp(c, s.score, feedback.CountUpOptions{Scale: 2, Slot: pixel.Accent})
		for _, add := range []int{100, 1000} {
			if input.Button(c, fmt.Sprintf("+%d", add), input.ButtonOptions{Size: input.Small}).Clicked {
				s.score += add
				log(fmt.Sprintf("score %d", s.score))
			}
		}
	})
}

func demoPulse(c *ui.Context, log func(string)) {
	feedback.Pulse(c, feedback.PulseOptions{}, func() {
		if input.IconButton(c, icon.Mail, "New mail", input.IconButtonOptions{}).Clicked {
			log("mail opened")
		}
	})
}

func demoLevelUp(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		if input.Button(c, "GAIN XP", input.ButtonOptions{Size: input.Small}).Clicked {
			s.level++
			log(fmt.Sprintf("level %d", s.level))
		}
		feedback.LevelUp(c, s.level, feedback.LevelUpOptions{})
	})
}

func demoAchievement(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	titles := []string{"Pressed Start", "Slime Tamer", "Touched Grass", "Read the Manual"}
	ui.Column(c).Gap(k.D(4)).Children(func() {
		if input.Button(c, "UNLOCK", input.ButtonOptions{Size: input.Small}).Clicked {
			s.cue++
			log("unlocked " + titles[s.cue%len(titles)])
		}
		feedback.Achievement(c, titles[s.cue%len(titles)], feedback.AchievementOptions{Cue: s.cue})
	})
}

func demoComboCounter(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		ui.Row(c).Gap(k.D(4)).Children(func() {
			if input.Button(c, "HIT", input.ButtonOptions{Size: input.Small, Variant: input.Primary}).Clicked {
				s.combo++
				log(fmt.Sprintf("combo %d", s.combo))
			}
			if input.Button(c, "MISS", input.ButtonOptions{Size: input.Small}).Clicked {
				s.combo = 0
				log("combo broken")
			}
		})
		feedback.ComboCounter(c, s.combo, feedback.ComboCounterOptions{})
	})
}

func demoContinuePrompt(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	ui.Column(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
		ui.Box(c).Key(s.prompts).Children(func() {
			r := feedback.ContinuePrompt(c, feedback.ContinuePromptOptions{})
			switch {
			case r.Yes:
				log("yes: one more go")
			case r.No:
				log("no: thanks for playing")
			case r.TimedOut:
				log("time up")
			}
		})
		if input.Button(c, "AGAIN", input.ButtonOptions{Size: input.Small}).Clicked {
			s.prompts++
			log("again")
		}
	})
}

func demoSaveIndicator(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	saving := k.Now().Before(s.saving)
	if saving {
		k.Tick(c, 4)
	}
	ui.Row(c).Gap(k.D(8)).AlignItems(ui.Center).Children(func() {
		if input.Button(c, "SAVE", input.ButtonOptions{Size: input.Small, Disabled: saving}).Clicked {
			s.saving = k.Now().Add(1500 * time.Millisecond)
			log("saving")
		}
		feedback.SaveIndicator(c, saving, feedback.SaveIndicatorOptions{})
	})
}

func demoBusyOverlay(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	busy := k.Now().Before(s.busy)
	if busy {
		k.Tick(c, 4)
	}
	feedback.BusyOverlay(c, busy, feedback.BusyOverlayOptions{}, func() {
		ui.Column(c).Gap(k.D(4)).Padding(k.D(4)).Children(func() {
			if input.TextInput(c, "Guild motto", &s.motto, input.TextInputOptions{Placeholder: "We came, we saved"}).Changed {
				log("motto " + s.motto)
			}
			if input.Button(c, "SEND", input.ButtonOptions{Size: input.Small}).Clicked {
				s.busy = k.Now().Add(2 * time.Second)
				log("sending…")
			}
		})
	})
}

func demoReadyGo(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	ui.Column(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
		if input.Button(c, "START", input.ButtonOptions{Size: input.Small, Variant: input.Primary}).Clicked {
			s.start++
			log("ready")
		}
		if feedback.ReadyGo(c, feedback.ReadyGoOptions{Cue: s.start}).Done {
			log("GO!")
		}
	})
}

func demoShake(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := feedbackState(c)
	ui.Column(c).Gap(k.D(4)).Children(func() {
		feedback.Shake(c, feedback.ShakeOptions{Cue: s.shake}, func() {
			input.PasswordInput(c, "Password", &s.secret, input.PasswordInputOptions{Placeholder: "Speak, friend"})
		})
		if input.Button(c, "ENTER", input.ButtonOptions{Size: input.Small}).Clicked {
			s.shake++
			log("wrong password")
		}
	})
}
