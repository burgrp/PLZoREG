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
	ErrorNone   ErrorCode = iota
	ErrorNoSync           // No ZeroCross sync detected
)

type SystemState struct {
	Page         Page
	VSense       int
	VTarget      int
	TSense       int
	InSetting    bool
	DisplayBlink bool
	Error        ErrorCode
	TestDuty     int
}

var State = SystemState{
	Page:      PageDuty,
	InSetting: true,
	VSense:    234,
	TSense:    52,
	VTarget:   230,
	Error:     ErrorNone,
	TestDuty:  50,
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
		display.NumberAt(int(State.Error), true, -1, 2, 2)
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
		display.NumberAt(State.TestDuty, true, dp, 2, 3)
	}
}

func handleKeyArrow(keyID int) {

	switch {
	case State.Page == PageVTarget && State.InSetting:
		v := State.VTarget + keyID
		if v < VTargetMin {
			v = VTargetMin
		}
		if v > VTargetMax {
			v = VTargetMax
		}
		State.VTarget = v

	case State.Page == PageDuty && State.InSetting:
		v := State.TestDuty + keyID*5
		if v < 0 {
			v = 0
		}
		if v > 100 {
			v = 100
		}
		State.TestDuty = v

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
		{Pin: machine.PA11, ID: KeyUp},
		{Pin: machine.PA12, ID: KeyDown},
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

	State.VTarget = int(flashDataPageWords[1])
}

func main() {

	loadSettings()

	initDisplay()

	initKeyboard()

	updateDisplay()

	for {
		// not to only sleep here, we refresh the display periodically just for a case it needs to blink
		time.Sleep(500 * time.Millisecond)
		State.DisplayBlink = !State.DisplayBlink
		updateDisplay()
	}
}
