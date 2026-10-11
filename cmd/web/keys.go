//go:build js && wasm

package main

import "github.com/egoist/mygo/ui"

// keys maps KeyboardEvent.code to MyGo keys.
var keys = func() map[string]ui.Key {
	m := map[string]ui.Key{
		"Enter": ui.KeyEnter, "NumpadEnter": ui.KeyEnter, "Escape": ui.KeyEscape, "Backspace": ui.KeyBackspace,
		"Tab": ui.KeyTab, "Space": ui.KeySpace, "Delete": ui.KeyDelete, "Insert": ui.KeyInsert,
		"Home": ui.KeyHome, "End": ui.KeyEnd, "PageUp": ui.KeyPageUp, "PageDown": ui.KeyPageDown,
		"ArrowLeft": ui.KeyLeft, "ArrowRight": ui.KeyRight, "ArrowUp": ui.KeyUp, "ArrowDown": ui.KeyDown,
		"Minus": ui.KeyMinus, "Equal": ui.KeyEqual, "Comma": ui.KeyComma, "Period": ui.KeyPeriod,
		"Slash": ui.KeySlash, "Semicolon": ui.KeySemicolon, "Quote": ui.KeyQuote,
		"BracketLeft": ui.KeyBracketLeft, "BracketRight": ui.KeyBracketRight, "Backslash": ui.KeyBackslash,
		"Backquote": ui.KeyBackquote, "ContextMenu": ui.KeyContextMenu,
		"BrowserBack": ui.KeyBack, "BrowserForward": ui.KeyForward,
		"F1": ui.KeyF1, "F2": ui.KeyF2, "F3": ui.KeyF3, "F4": ui.KeyF4, "F5": ui.KeyF5, "F6": ui.KeyF6,
		"F7": ui.KeyF7, "F8": ui.KeyF8, "F9": ui.KeyF9, "F10": ui.KeyF10, "F11": ui.KeyF11, "F12": ui.KeyF12,
	}
	for i := range 10 {
		m["Digit"+string(rune('0'+i))] = ui.Key0 + ui.Key(i)
	}
	for i := range 26 {
		m["Key"+string(rune('A'+i))] = ui.KeyA + ui.Key(i)
	}
	return m
}()
