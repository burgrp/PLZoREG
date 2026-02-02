package main

import (
	"device/py32"
	"machine"
	"unsafe"
)

const (
	pinVSense = machine.PB0
	pinTSense = machine.PB1
)

var adcBuffer [128]struct {
	vSense uint16
	tSense uint16
}

func initAdc() {
	pinVSense.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	pinTSense.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})

	py32.RCC.SetAPBENR2_ADCEN(1)    // Enable ADC clock
	py32.RCC.SetAPBENR2_SYSCFGEN(1) // Enable SYSCFG clock - needed for DMA mapping
	py32.RCC.SetAHBENR_DMAEN(1)     // Enable DMA clock

	dma := &py32.DMA.CH[0]
	dma.SetPAR(uint32(uintptr(unsafe.Pointer(&py32.ADC.DR.Reg))))
	dma.SetMAR(uint32(uintptr(unsafe.Pointer(&adcBuffer[0]))))
	dma.SetNDTR_NDT(uint32(len(adcBuffer)))
	dma.SetCR_DIR(0)   // 0: peripheral to memory
	dma.SetCR_PSIZE(1) // 01: 16-bit
	dma.SetCR_MSIZE(1) // 01: 16-bit
	dma.SetCR_MINC(1)  // memory increment
	dma.SetCR_PINC(0)  // peripheral no increment
	dma.SetCR_CIRC(1)  // circular buffer
	dma.SetCR_EN(1)    // enable DMA channel

	py32.ADC.SetCFGR2_CKMODE(py32.ADC_CFGR2_CKMODE_HSI_Div64)
	py32.ADC.SetSMPR_SMP(py32.ADC_SMPR_SMP_Cycles239_5) // Sampling time 239.5 ADC clock cycles
	py32.ADC.SetCFGR1_CONT(1)                           // continuous conversion mode
	py32.ADC.SetCFGR1_OVRMOD(1)                         // overwrite mode
	py32.ADC.SetCFGR1_DMAEN(1)                          // enable DMA
	py32.ADC.SetCFGR1_DMACFG(1)                         // DMA in circular mode
	py32.ADC.CHSELR.Set((1 << 8) | (1 << 9))            // select channels 8 and 9
	py32.ADC.SetCR_ADEN(1)                              // enable ADC
	py32.ADC.SetCR_ADSTART(1)                           // start ADC

}

func readAdcValues() (vSense uint32, tSense uint32) {

	for _, e := range adcBuffer {
		vSense += uint32(e.vSense)
		tSense += uint32(e.tSense)
	}

	cnt := uint32(len(adcBuffer))
	return vSense / cnt / 8, tSense / cnt
}
