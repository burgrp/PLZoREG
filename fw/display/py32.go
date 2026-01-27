//go:build py32

package display

import (
	"device/py32"
	"runtime/volatile"
	"unsafe"
)

func configurePeripherals(digits, intensity int) {
	configureGpio()

	py32.RCC.SetAPBENR2_LEDEN(1)

	py32.LED.SetPR(100)
	SetIntensity(intensity)

	py32.LED.SetCR_EHS(1)
	py32.LED.SetCR_LED_COM_SEL(uint32(digits) - 1)
	py32.LED.SetCR_LEDON(1)
}

func SetIntensity(i int) {
	py32.LED.SetTR_T1(uint32(i))
	py32.LED.SetTR_T2(255)
}

func GlyphAt(glyph, position int) {
	volatile.StoreUint32((*uint32)(unsafe.Add(unsafe.Pointer(&py32.LED.DR0.Reg), position*4)), uint32(glyph))
}
