package main

import (
	"device/py32"
	"machine"
	"runtime"
)

const (
	pinVSense = machine.PB0
	pinTSense = machine.PB1
)

func initAdc() {
	pinVSense.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	pinTSense.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})

	py32.RCC.SetAPBENR2_ADCEN(1) // Enable ADC clock

	py32.ADC.SetSMPR_SMP(py32.ADC_SMPR_SMP_Cycles239_5) // Sampling time 239.5 ADC clock cycles

}

var prevVSense uint32
var prevTSense uint32

func readAdcValues() (vSense uint32, tSense uint32) {

	for _ = range 128 {
		vSense = readAdcChannel(8)
		if prevVSense != 0 {
			vSense = (vSense + prevVSense*127) / 128
		}
		prevVSense = vSense
	}

	tSense = readAdcChannel(9)
	println("ADC read VSense:", vSense, "TSense:", tSense)
	return vSense / 16, tSense
}

func readAdcChannel(channel uint32) uint32 {

	py32.ADC.CHSELR.Set(1 << channel)
	py32.ADC.SetCR_ADEN(1)
	py32.ADC.SetCR_ADSTART(1)

	for py32.ADC.GetISR_EOC() == 0 {
		runtime.Gosched()
	}

	return py32.ADC.GetDR_DATA()
}
