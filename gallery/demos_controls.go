package gallery

import (
	"fmt"
	"slices"
	"time"

	"github.com/egoist/mygo/ui"

	"github.com/ELY-Open-Source/8bitMyGoUI/input"
	"github.com/ELY-Open-Source/8bitMyGoUI/pixel"
)

var controlEntries = []Entry{
	{86, "Select", "下拉选择", "input", "The choice and ▼; the list opens in a window with a hand pointing. Arrows and Enter pick.", demoSelect},
	{87, "MultiSelect", "多选", "input", "A popup of checkboxes; what you pick rides in the trigger as tags.", demoMultiSelect},
	{88, "Combobox", "组合框", "input", "Type to narrow the list below; Enter or a click takes one.", demoCombobox},
	{89, "Autocomplete", "自动补全", "input", "The rest of the word waits in grey. Tab or Enter takes it.", demoAutocomplete},
	{90, "TokenField", "标记输入", "input", "Tags before the caret. Enter or a comma adds; Backspace on empty takes the last back.", demoTokenField},
	{91, "Segmented", "分段控件", "input", "Joined segments, one lit. Arrows move the light.", demoSegmented},
	{92, "Rating", "评分", "input", "Stars or hearts. Click or arrows to rate; click the same mark to clear.", demoRating},
	{93, "Knob", "转钮", "input", "A rotary knob that clicks through 16 angles. Drag up or down, or use the arrows.", demoKnob},
	{94, "PinInput", "验证码", "input", "One box per digit: typing moves on, Backspace steps back, paste fills them all.", demoPinInput},
	{95, "NameEntry", "姓名登记", "input", "HIGH SCORE! Three initials from the letter grid, then END.", demoNameEntry},
	{96, "VirtualKeyboard", "屏幕键盘", "input", "A keyboard on screen, for pads and touch: letters, digits, space, erase, enter.", demoVirtualKeyboard},
	{97, "DPad", "方向键", "input", "A cross pad. Press an arm or use the arrow keys.", demoDPad},
	{98, "HoldButton", "长按按钮", "input", "Hold it down while the bar fills in eight steps. Let go and it starts over.", demoHoldButton},
	{99, "EditableText", "可编辑文本", "input", "Double click or press Enter to rename. Enter keeps it; Escape takes it back.", demoEditableText},
	{100, "Field", "表单项", "input", "A label above, a hint below, and a red ! when something is wrong.", demoField},
	{101, "Fieldset", "字段组", "input", "A framed group of fields under a legend.", demoFieldset},
	{102, "ColorSlot", "色槽选择", "input", "Pick one of the sixteen palette slots. It follows the palette.", demoColorSlot},
	{103, "Joystick", "操作杆", "input", "Drag the stick inside its ring; it springs home. Arrows nudge it.", demoJoystick},
	{104, "CheatCode", "秘籍", "input", "Focus it and enter ↑↑↓↓←→←→BA. You know the one.", demoCheatCode},
	{105, "SlotPicker", "转轮选择", "input", "Reels that spin by drag, wheel or arrows. Line up three for the jackpot.", demoSlotPicker},
}

type controlDemo struct {
	class, item, spell, hero, call, guild string
	party, tags                           []string
	mode, banner, stars, lives            int
	reel1, reel2, reel3                   int
	volume, bass                          float64
	pin, initials, message                string
	cape                                  pixel.Slot
	dir                                   input.Direction
}

func controlState(c *ui.Context) *controlDemo {
	return ui.Local(c.Root(), demoKey("controls"), func() controlDemo {
		return controlDemo{class: "Knight", hero: "Pixel Knight", tags: []string{"retro", "pixel"}, mode: 1, stars: 4, lives: 3,
			volume: 0.7, bass: 0.3, cape: pixel.Accent, reel2: 1, reel3: 2}
	})
}

func demoSelect(c *ui.Context, log func(string)) {
	s := controlState(c)
	if input.Select(c, "Class", &s.class, []string{"Knight", "Mage", "Thief", "Bard"}, input.SelectOptions{}).Changed {
		log("class " + s.class)
	}
}

func demoMultiSelect(c *ui.Context, log func(string)) {
	s := controlState(c)
	if input.MultiSelect(c, "Party", &s.party, []string{"Warrior", "Healer", "Rogue", "Wizard"}, input.MultiSelectOptions{}).Changed {
		log(fmt.Sprintf("party %v", s.party))
	}
}

func demoCombobox(c *ui.Context, log func(string)) {
	s := controlState(c)
	items := []string{"Potion", "Hi-Potion", "Ether", "Elixir", "Antidote", "Phoenix Down"}
	r := input.Combobox(c, "Item", &s.item, items, input.ComboboxOptions{Placeholder: "Use an item…", Width: 140})
	if r.Changed {
		log("item " + s.item)
	}
	if r.Chosen {
		log("used " + s.item)
	}
}

func demoAutocomplete(c *ui.Context, log func(string)) {
	s := controlState(c)
	spells := []string{"Fire", "Fira", "Firaga", "Blizzard", "Thunder", "Cure", "Curaga"}
	r := input.Autocomplete(c, "Spell", &s.spell, spells, input.AutocompleteOptions{Placeholder: "Cast…"})
	if r.Changed {
		log("spell " + s.spell)
	}
	if r.Accepted {
		log("cast " + s.spell)
	}
	if r.Submitted {
		log("fizzle: " + s.spell)
	}
}

func demoTokenField(c *ui.Context, log func(string)) {
	s := controlState(c)
	if input.TokenField(c, "Tags", &s.tags, input.TokenFieldOptions{Placeholder: "add tag"}).Changed {
		log(fmt.Sprintf("tags %v", s.tags))
	}
}

func demoSegmented(c *ui.Context, log func(string)) {
	s := controlState(c)
	modes := []string{"EASY", "NORMAL", "HARD"}
	if input.Segmented(c, "Difficulty", modes, &s.mode, input.SegmentedOptions{}).Changed {
		log("difficulty " + modes[s.mode])
	}
}

func demoRating(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := controlState(c)
	ui.Column(c).Gap(k.D(6)).Children(func() {
		if input.Rating(c, "Dungeon", &s.stars, input.RatingOptions{}).Changed {
			log(fmt.Sprintf("%d stars", s.stars))
		}
		if input.Rating(c, "Lives", &s.lives, input.RatingOptions{Max: 3, Hearts: true}).Changed {
			log(fmt.Sprintf("%d lives", s.lives))
		}
	})
}

func demoKnob(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := controlState(c)
	ui.Row(c).Gap(k.D(12)).AlignItems(ui.Center).Children(func() {
		for _, kn := range []struct {
			name string
			v    *float64
		}{{"VOL", &s.volume}, {"BASS", &s.bass}} {
			ui.Column(c).Gap(k.D(2)).AlignItems(ui.Center).Children(func() {
				if input.Knob(c, kn.name, kn.v, input.KnobOptions{Step: 0.05}).Changed {
					log(fmt.Sprintf("%s %.2f", kn.name, *kn.v))
				}
				pixel.Text(c, fmt.Sprintf("%s %d", kn.name, int(*kn.v*100+0.5)), pixel.TextOptions{Slot: pixel.Faint})
			})
		}
	})
}

func demoPinInput(c *ui.Context, log func(string)) {
	s := controlState(c)
	r := input.PinInput(c, "Vault code", &s.pin, input.PinInputOptions{})
	if r.Changed {
		log(fmt.Sprintf("%d digits", len(s.pin)))
	}
	if r.Complete {
		if s.pin == "0451" {
			log("ACCESS GRANTED")
		} else {
			log("WRONG CODE")
		}
	}
}

func demoNameEntry(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := controlState(c)
	ui.Column(c).Gap(k.D(4)).AlignItems(ui.Center).Children(func() {
		pixel.Text(c, "NEW HIGH SCORE  99999", pixel.TextOptions{Slot: pixel.Accent})
		r := input.NameEntry(c, "Initials", &s.initials, input.NameEntryOptions{})
		if r.Changed {
			log("initials " + s.initials)
		}
		if r.Confirmed {
			log("1ST " + s.initials + " 99999")
		}
	})
}

func demoVirtualKeyboard(c *ui.Context, log func(string)) {
	s := controlState(c)
	r := input.VirtualKeyboard(c, "Message", &s.message, input.VirtualKeyboardOptions{})
	if r.Changed {
		log("message " + s.message)
	}
	if r.Submitted {
		log("sent " + s.message)
		s.message = ""
	}
}

func demoDPad(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := controlState(c)
	ui.Row(c).Gap(k.D(12)).AlignItems(ui.Center).Children(func() {
		if d := input.DPad(c, "Pad", input.DPadOptions{}).Pressed; d != input.None {
			s.dir = d
			log(d.String())
		}
		pixel.Text(c, "last: "+s.dir.String(), pixel.TextOptions{Slot: pixel.Faint})
	})
}

func demoHoldButton(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(6)).Children(func() {
		if input.HoldButton(c, "SKIP CUTSCENE", input.HoldButtonOptions{}).Confirmed {
			log("cutscene skipped")
		}
		if input.HoldButton(c, "DELETE SAVE", input.HoldButtonOptions{Duration: 2 * time.Second, Danger: true}).Confirmed {
			log("save deleted (not really)")
		}
	})
}

func demoEditableText(c *ui.Context, log func(string)) {
	s := controlState(c)
	if input.EditableText(c, "Hero name", &s.hero, input.EditableTextOptions{Placeholder: "Nameless"}).Changed {
		log("renamed to " + s.hero)
	}
}

func demoField(c *ui.Context, log func(string)) {
	s := controlState(c)
	o := input.FieldOptions{Description: "Shown on the leaderboard."}
	if s.call == "" {
		o.Error = "A hero needs a call sign!"
	}
	input.Field(c, "Call sign", o, func() {
		r := input.TextInput(c, "Call sign", &s.call, input.TextInputOptions{MaxLength: 10})
		if r.Changed {
			log("call sign " + s.call)
		}
		if r.Submitted {
			log("signed " + s.call)
		}
	})
}

func demoFieldset(c *ui.Context, log func(string)) {
	s := controlState(c)
	input.Fieldset(c, "Guild", input.FieldsetOptions{}, func() {
		input.Field(c, "Name", input.FieldOptions{}, func() {
			r := input.TextInput(c, "Guild name", &s.guild, input.TextInputOptions{Placeholder: "The Dot Matrix"})
			if r.Changed {
				log("guild " + s.guild)
			}
			if r.Submitted {
				log("founded " + s.guild)
			}
		})
		banners := []string{"RED", "BLUE", "GOLD"}
		if input.Segmented(c, "Banner", banners, &s.banner, input.SegmentedOptions{}).Changed {
			log("banner " + banners[s.banner])
		}
	})
}

func demoColorSlot(c *ui.Context, log func(string)) {
	s := controlState(c)
	if input.ColorSlot(c, "Cape", &s.cape, input.ColorSlotOptions{}).Changed {
		log("cape " + s.cape.String())
	}
}

func demoJoystick(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	ui.Row(c).Gap(k.D(12)).AlignItems(ui.Center).Children(func() {
		r := input.Joystick(c, "Stick", input.JoystickOptions{})
		if r.Moved {
			log(fmt.Sprintf("x %+.2f y %+.2f", r.X, r.Y))
		}
		pixel.Text(c, fmt.Sprintf("X %+.2f  Y %+.2f", r.X, r.Y), pixel.TextOptions{Slot: pixel.Faint})
	})
}

func demoCheatCode(c *ui.Context, log func(string)) {
	if input.CheatCode(c, "Cheat", input.CheatCodeOptions{}).Unlocked {
		log("30 LIVES!")
	}
}

func demoSlotPicker(c *ui.Context, log func(string)) {
	k := pixel.K(c)
	s := controlState(c)
	fruit := []string{"Cherry", "Bell", "Seven", "Melon", "BAR"}
	ui.Row(c).Gap(k.D(4)).Children(func() {
		changed := false
		for i, v := range []*int{&s.reel1, &s.reel2, &s.reel3} {
			if input.SlotPicker(c, fmt.Sprintf("Reel %d", i+1), v, fruit, input.SlotPickerOptions{}).Changed {
				changed = true
				log(fmt.Sprintf("reel %d %s", i+1, fruit[*v]))
			}
		}
		if all := []int{s.reel1, s.reel2, s.reel3}; changed && slices.Max(all) == slices.Min(all) {
			log("JACKPOT! " + fruit[s.reel1])
		}
	})
}
