package main

import (
	"machine"
	"plzoreg/display"
	"plzoreg/flash"
	"plzoreg/keyboard"
	"time"
	"unsafe"
)

type Page int

const (
	pinKeyUp   = machine.PA11
	pinKeyDown = machine.PA12
)

const (
	PageVSense Page = iota
	PageVTarget
	PageTSense
	PageDuty
	PageCount
)

const FlashKeyV1 = 0xDEADBEEF

const (
	KeyUp   = 1
	KeyDown = -1
)

const (
	VTargetMin     = 100
	VTargetMax     = 240
	VTargetDefault = 200
)

type ErrorCode int

const (
	ErrorNone     ErrorCode = iota
	ErrorNoSync             // No ZeroCross sync detected
	ErrorOverTemp           // Overtemperature detected
)

type SystemState struct {
	Page         Page
	VSense       uint32
	VTarget      uint32
	TSense       uint32
	InSetting    bool
	DisplayBlink bool
	Error        ErrorCode
	Duty         uint32
}

var State = SystemState{
	Page:      PageTSense,
	InSetting: false,
	VSense:    0,
	TSense:    0,
	VTarget:   VTargetDefault,
	Error:     ErrorNone,
	Duty:      0,
}

func initDisplay() {
	display.Init(3, 200)
	display.GlyphAt(0xFF, 0)
	display.GlyphAt(0xFF, 1)
	display.GlyphAt(0xFF, 2)
	time.Sleep(500 * time.Millisecond)
	display.GlyphAt(0, 0)
	display.GlyphAt(0, 1)
	display.GlyphAt(0, 2)
	time.Sleep(50 * time.Millisecond)
}

func updateDisplay() {

	if State.Error != ErrorNone {
		display.GlyphAt(0x79, 0) // E
		display.NumberAt(uint32(State.Error), true, -1, 2, 2)
		return
	}

	switch State.Page {
	case PageVSense:
		display.NumberAt(State.VSense, false, -1, 2, 3)
	case PageVTarget:
		dp := -1
		if !State.InSetting || State.DisplayBlink {
			dp = 2
		}
		display.NumberAt(State.VTarget, false, dp, 2, 3)
	case PageTSense:
		display.NumberAt(State.TSense, false, -1, 1, 2)
		display.GlyphAt(0x63, 2)
	case PageDuty:
		dp := -1
		if !State.InSetting || State.DisplayBlink {
			dp = 0
		}
		display.NumberAt(State.Duty, true, dp, 2, 3)
	}
}

func handleKeyArrow(keyID int) {

	if State.Error != ErrorNone {
		return
	}

	switch {
	case State.Page == PageVTarget && State.InSetting:
		v := int(State.VTarget) + keyID
		if v < VTargetMin {
			v = VTargetMin
		}
		if v > VTargetMax {
			v = VTargetMax
		}
		State.VTarget = uint32(v)

	case State.Page == PageDuty && State.InSetting:
		v := int(State.Duty) + keyID*5
		if v < 0 {
			v = 0
		}
		if v > 100 {
			v = 100
		}
		State.Duty = uint32(v)

	default:
		p := Page(int(State.Page) - keyID)
		if p < 0 {
			p = PageCount - 1
		}
		if p >= PageCount {
			p = 0
		}
		State.Page = p
	}

	updateDisplay()
}

func handleKeyDoublePress() {

	if State.Error != ErrorNone {
		return
	}

	switch State.Page {
	case PageVTarget:
		if !State.InSetting {
			State.DisplayBlink = true
			State.InSetting = true
		} else {
			State.InSetting = false
			saveSettings()
		}
	case PageDuty:
		State.InSetting = !State.InSetting
	}
	updateDisplay()
}

func initKeyboard() {
	kbd := keyboard.New(keyboard.Keys{
		{Pin: pinKeyUp, ID: KeyUp},
		{Pin: pinKeyDown, ID: KeyDown},
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

func GetFlashDataAddr() uint32 {
	return flash.MainFlashBase + flash.MainFlashSize - flash.PageSize
}

func GetFlashDataWords() *[flash.WordsPerPage]uint32 {
	return (*[flash.WordsPerPage]uint32)(unsafe.Pointer(uintptr(GetFlashDataAddr())))
}

func saveSettings() {
	var data [flash.WordsPerPage]uint32
	data[0] = uint32(FlashKeyV1)
	data[1] = uint32(State.VTarget)

	FlashDataPageAddrAsUint32 := uint32(GetFlashDataAddr())

	flash.ErasePage(FlashDataPageAddrAsUint32)
	flash.ProgramPage(FlashDataPageAddrAsUint32, &data)
}

func loadSettings() {
	flashDataPageWords := GetFlashDataWords()
	if flashDataPageWords[0] != uint32(FlashKeyV1) {
		// no valid data, use defaults
		State.VTarget = VTargetDefault
		return
	}

	State.VTarget = flashDataPageWords[1]
}

func main() {

	loadSettings()
	initDisplay()
	initKeyboard()
	initAdc()
	initPwm()

	updateDisplay()

	for {
		// not to only sleep here, we refresh the display periodically just for a case it needs to blink
		time.Sleep(500 * time.Millisecond)
		State.DisplayBlink = !State.DisplayBlink

		if !isSynchronized() {
			State.Error = ErrorNoSync
		} else {
			if State.Error == ErrorNoSync {
				State.Error = ErrorNone
			}
		}

		v, t := readAdcValues()
		State.VSense = v
		State.TSense = t

		// if State.TSense > 90 {
		// 	State.Error = ErrorOverTemp
		// }

		if State.Error == ErrorNone {

			if !(State.Page == PageDuty && State.InSetting) {

				d := State.Duty

				if State.VSense < State.VTarget {
					if d > 0 {
						d--
					}

				} else {
					if d < 100 {
						d++
					}
				}

				State.Duty = d

			}

		} else {
			State.Duty = 0
		}

		setPwm(State.Duty)

		updateDisplay()
	}
}
