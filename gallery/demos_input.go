package gallery

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/8bitMyGoUI/pixel"
)

var inputEntries = []Entry{
	{66, "Button", "按钮", "input", "Primary, normal, ghost and danger in three sizes. Pressing sinks the face a dot; loading spins and waits.", demoButton},
	{67, "IconButton", "图标按钮", "input", "A square button with an icon and a name for screen readers.", demoIconButton},
	{68, "ButtonGroup", "按钮组", "input", "Buttons joined edge to edge: one Tab stop, arrows move between them.", demoButtonGroup},
	{69, "Toggle", "开关按钮", "input", "A button that stays sunken while it is on.", demoToggle},
	{70, "ToggleGroup", "开关组", "input", "Joined toggles, for many choices or just one.", demoToggleGroup},
	{71, "Checkbox", "复选框", "input", "An 8×8 box with a hand-placed check. Space toggles it.", demoCheckbox},
	{72, "CheckboxGroup", "复选组", "input", "A parent box over children: mixed while some are checked, all or none on a click.", demoCheckboxGroup},
	{73, "Radio", "单选", "input", "A pixel radio dot that picks its value.", demoRadio},
	{74, "RadioGroup", "单选组", "input", "Radios as one Tab stop; the arrows choose.", demoRadioGroup},
	{75, "Switch", "拨杆", "input", "A lever on an OFF/ON plate. Space flips it.", demoSwitch},
	{76, "Slider", "滑块", "input", "A Well track, an Accent fill and a beveled knob: drag, arrows, Home, End.", demoSlider},
	{77, "StepSlider", "刻度滑块", "input", "A slider that snaps to steps, with tick dots and its value.", demoStepSlider},
	{78, "RangeSlider", "区间滑块", "input", "Two knobs for low and high that never cross.", demoRangeSlider},
	{79, "Fader", "推子", "input", "A mixing-desk fader with a cap and ticks.", demoFader},
	{80, "Stepper", "步进器", "input", "◀ value ▶ within bounds. Hold an arrow to repeat; keys step it too.", demoStepper},
	{81, "NumberInput", "数字框", "input", "Type a number or step it with ▲▼ and the arrow keys; bad text turns the edge red.", demoNumberInput},
	{82, "TextInput", "文本框", "input", "A one-line pixel field: caret, selection, input methods, clipboard. Enter submits.", demoTextInput},
	{83, "PasswordInput", "密码框", "input", "Masked with ★; the eye shows the text.", demoPasswordInput},
	{84, "TextArea", "多行文本", "input", "Many lines in its own scroll, following the caret, with a character count.", demoTextArea},
	{85, "SearchField", "搜索框", "input", "A field with a magnifier. Escape clears it, Enter searches.", demoSearchField},
}

type inputDemo struct {
	toggle, music, sound, check, lever bool
	toggles, mode, party               []bool
	radio                              string
	choice, level, count               int
	volume, steps, low, high, fader    float64
	price                              float64
	name, secret, notes, query         string
}

func inputState(c *ui.Context) *inputDemo {
	return ui.Local(c.Root(), demoKey("input"), func() inputDemo {
		return inputDemo{music: true, toggles: []bool{true, false, false}, mode: []bool{false, true, false}, party: []bool{true, false, true},
			radio: "knight", choice: 1, level: 3, count: 5, volume: 0.6, steps: 40, low: 20, high: 70, fader: 0.3, price: 9.5}
	})
}

func demoButton(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
			for i, v := range []input.Variant{input.Primary, input.Normal, input.Ghost, input.Danger} {
				name := []string{"PRIMARY", "NORMAL", "GHOST", "DANGER"}[i]
				if input.Button(c, name, input.ButtonOptions{Variant: v}).Clicked {
					log(name)
				}
			}
		})
		ui.Row(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
			for _, b := range []struct {
				name string
				o    input.ButtonOptions
			}{{"Small", input.ButtonOptions{Size: input.Small}}, {"Save", input.ButtonOptions{Icon: icon.Save}}, {"BIG", input.ButtonOptions{Size: input.Large, Variant: input.Primary}}} {
				if input.Button(c, b.name, b.o).Clicked {
					log(b.name)
				}
			}
			input.Button(c, "Saving", input.ButtonOptions{Loading: true})
			input.Button(c, "Off", input.ButtonOptions{Disabled: true})
		})
	})
}

func demoIconButton(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(4)).Children(func() {
		for _, n := range []icon.Named{{Name: "Play", Sprite: icon.Play}, {Name: "Pause", Sprite: icon.Pause}, {Name: "Stop", Sprite: icon.Stop}} {
			if input.IconButton(c, n.Sprite, n.Name, input.IconButtonOptions{}).Clicked {
				log(n.Name)
			}
		}
		if input.IconButton(c, icon.Trash, "Delete", input.IconButtonOptions{Variant: input.Danger}).Clicked {
			log("Delete")
		}
	})
}

func demoButtonGroup(c *ui.Context, log func(string)) {
	labels := []string{"FIGHT", "MAGIC", "ITEM", "RUN"}
	if r := input.ButtonGroup(c, labels, input.ButtonGroupOptions{Label: "Command"}); r.Clicked >= 0 {
		log(labels[r.Clicked])
	}
}

func demoToggle(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := inputState(c)
	ui.Row(c).Gap(k.D(4)).Children(func() {
		if input.Toggle(c, "MUSIC", &s.music, input.ToggleOptions{Icon: icon.Music}).Changed {
			log(fmt.Sprintf("music %v", s.music))
		}
		if input.Toggle(c, "SOUND", &s.sound, input.ToggleOptions{Icon: icon.Volume}).Changed {
			log(fmt.Sprintf("sound %v", s.sound))
		}
	})
}

func demoToggleGroup(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := inputState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		if input.ToggleGroup(c, []string{"B", "I", "U"}, s.toggles, input.ToggleGroupOptions{Label: "Style"}).Changed {
			log(fmt.Sprintf("style %v", s.toggles))
		}
		if input.ToggleGroup(c, []string{"EASY", "NORMAL", "HARD"}, s.mode, input.ToggleGroupOptions{Single: true, Label: "Mode"}).Changed {
			log(fmt.Sprintf("mode %v", s.mode))
		}
	})
}

func demoCheckbox(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := inputState(c)
	ui.Column(c).Gap(k.D(2)).Children(func() {
		if input.Checkbox(c, "Remember me", &s.check, input.CheckboxOptions{}).Changed {
			log(fmt.Sprintf("remember %v", s.check))
		}
		locked := true
		input.Checkbox(c, "Hard mode (locked)", &locked, input.CheckboxOptions{Disabled: true})
	})
}

func demoCheckboxGroup(c *ui.Context, log func(string)) {
	s := inputState(c)
	if input.CheckboxGroup(c, "Party", []string{"Knight", "Mage", "Thief"}, s.party, input.CheckboxGroupOptions{}).Changed {
		log(fmt.Sprintf("party %v", s.party))
	}
}

func demoRadio(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := inputState(c)
	ui.Column(c).Gap(k.D(2)).Children(func() {
		for _, class := range []string{"knight", "mage", "thief"} {
			if input.Radio(c, class, class, &s.radio, input.RadioOptions{}).Changed {
				log("class " + class)
			}
		}
	})
}

func demoRadioGroup(c *ui.Context, log func(string)) {
	labels := []string{"Spring", "Summer", "Autumn", "Winter"}
	s := inputState(c)
	if input.RadioGroup(c, labels, &s.choice, input.RadioGroupOptions{Row: true, Label: "Season"}).Changed {
		log("season " + labels[s.choice])
	}
}

func demoSwitch(c *ui.Context, log func(string)) {
	s := inputState(c)
	if input.Switch(c, "Power", &s.lever, input.SwitchOptions{}).Changed {
		log(fmt.Sprintf("power %v", s.lever))
	}
}

func demoSlider(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := inputState(c)
	ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		if input.Slider(c, "Volume", &s.volume, input.SliderOptions{Width: 120}).Changed {
			log(fmt.Sprintf("volume %.2f", s.volume))
		}
		pixel.Text(c, fmt.Sprintf("%3.0f%%", s.volume*100), pixel.TextOptions{})
	})
}

func demoStepSlider(c *ui.Context, log func(string)) {
	s := inputState(c)
	if input.StepSlider(c, "Speed", &s.steps, input.StepSliderOptions{Min: 0, Max: 100, Step: 20, Width: 120}).Changed {
		log(fmt.Sprintf("speed %g", s.steps))
	}
}

func demoRangeSlider(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := inputState(c)
	ui.Row(c).Gap(k.D(6)).AlignItems(ui.Center).Children(func() {
		if input.RangeSlider(c, "Level range", &s.low, &s.high, input.RangeSliderOptions{Min: 1, Max: 99, Width: 140}).Changed {
			log(fmt.Sprintf("levels %.0f–%.0f", s.low, s.high))
		}
		pixel.Text(c, fmt.Sprintf("LV %.0f–%.0f", s.low, s.high), pixel.TextOptions{})
	})
}

func demoFader(c *ui.Context, log func(string)) {
	s := inputState(c)
	if input.Fader(c, "Master", &s.fader, input.FaderOptions{}).Changed {
		log(fmt.Sprintf("master %.2f", s.fader))
	}
}

func demoStepper(c *ui.Context, log func(string)) {
	s := inputState(c)
	if input.Stepper(c, "Lives", &s.count, input.StepperOptions{Min: 1, Max: 9}).Changed {
		log(fmt.Sprintf("lives %d", s.count))
	}
}

func demoNumberInput(c *ui.Context, log func(string)) {
	s := inputState(c)
	if input.NumberInput(c, "Price", &s.price, input.NumberInputOptions{Min: 0, Max: 999, Step: 0.5, Width: 80}).Changed {
		log(fmt.Sprintf("price %g", s.price))
	}
}

func demoTextInput(c *ui.Context, log func(string)) {
	s := inputState(c)
	r := input.TextInput(c, "Hero name", &s.name, input.TextInputOptions{Placeholder: "Your name 你的名字", MaxLength: 12})
	if r.Changed {
		log("name " + s.name)
	}
	if r.Submitted {
		log("submitted " + s.name)
	}
}

func demoPasswordInput(c *ui.Context, log func(string)) {
	s := inputState(c)
	r := input.PasswordInput(c, "Password", &s.secret, input.PasswordInputOptions{Placeholder: "Password"})
	if r.Changed {
		log(fmt.Sprintf("%d characters", len([]rune(s.secret))))
	}
	if r.Submitted {
		log("signed in")
	}
}

func demoTextArea(c *ui.Context, log func(string)) {
	s := inputState(c)
	if input.TextArea(c, "Notes", &s.notes, input.TextAreaOptions{Placeholder: "Dear diary...", MaxLength: 140, Width: 200, Lines: 4}).Changed {
		log(fmt.Sprintf("%d characters", len([]rune(s.notes))))
	}
}

func demoSearchField(c *ui.Context, log func(string)) {
	s := inputState(c)
	r := input.SearchField(c, "Find items", &s.query, input.SearchFieldOptions{Placeholder: "Potion, sword...", Width: 160})
	if r.Changed {
		log("query " + s.query)
	}
	if r.Submitted {
		log("search " + s.query)
	}
}
