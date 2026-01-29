package main

import (
	"machine"
	"template/display"
	"template/keyboard"
	"time"
)

type Page int

const (
	PageVSense Page = iota
	PageTSense
	PageVTarget
)

type SystemState struct {
	Page         Page
	VSense       int
	TSense       int
	VTarget      int
	SettingBlink bool
}

var State = SystemState{
	Page:    PageVSense,
	VSense:  234,
	TSense:  52,
	VTarget: 230,
}

func initDisplay() {
	display.Init(3, 200)
	display.GlyphAt(0xFF, 0)
	display.GlyphAt(0xFF, 1)
	display.GlyphAt(0xFF, 2)
	time.Sleep(1 * time.Second)
	display.GlyphAt(0, 0)
	display.GlyphAt(0, 1)
	display.GlyphAt(0, 2)
	time.Sleep(100 * time.Millisecond)
}

func updateDisplay() {
	switch State.Page {
	case PageVSense:
		display.NumberAt(State.VSense, false, -1, 2)
	case PageTSense:
		display.NumberAt(State.TSense, false, -1, 1)
		display.GlyphAt(0x63, 2)
	case PageVTarget:
		dp := -1
		if State.SettingBlink {
			dp = 2
		}
		display.NumberAt(State.VTarget, false, dp, 2)
	}
}

func handleKeyArrow(keyID int) {
	switch State.Page {
	case PageVSense:
		State.Page = PageTSense
	case PageTSense:
		State.Page = PageVSense
	case PageVTarget:
		v := State.VTarget + keyID
		if v < 100 {
			v = 100
		}
		if v > 240 {
			v = 240
		}
		State.VTarget = v
	}
	updateDisplay()
}

func handleKeyDoublePress() {
	switch State.Page {
	case PageVSense:
		State.SettingBlink = true
		State.Page = PageVTarget
	case PageVTarget:
		saveSettings()
		State.Page = PageVSense
	}
	updateDisplay()
}

func initKeyboard() {
	kbd := keyboard.New(keyboard.Keys{
		{Pin: machine.PA11, ID: 1},
		{Pin: machine.PA12, ID: -1},
	})

	var bothDownMs int64

	go kbd.Run(func(e keyboard.KeyEvent) {

		if e.Kind == keyboard.KeyPress {
			bothDown := true
			for _, key := range *kbd.Keys {
				if key.DownMs == 0 {
					bothDown = false
					break
				}
			}
			if bothDown {
				if bothDownMs == 0 {
					bothDownMs = time.Now().UnixMilli()
				} else if time.Now().UnixMilli()-bothDownMs > 1000 {
					handleKeyDoublePress()
					bothDownMs = 0
					updateDisplay()
				}
			} else {
				bothDownMs = 0
				handleKeyArrow(e.ID)
				updateDisplay()
			}
		}

	})

}

func saveSettings() {
}

func loadSettings() {
}

func main() {

	loadSettings()

	initDisplay()

	initKeyboard()

	updateDisplay()

	for {
		// not to only sleep here, we refresh the display periodically just for a case it needs to blink
		time.Sleep(500 * time.Millisecond)
		State.SettingBlink = !State.SettingBlink
		updateDisplay()
	}
}
