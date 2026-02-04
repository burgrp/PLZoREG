package main

import (
	"device/py32"
	"machine"
	"plzoreg/util"
	"unsafe"
)

const (
	pinVSense     = machine.PB0
	pinTSense     = machine.PB1
	channelVSense = 8
	channelTSense = 9
	channelTMCU   = 11
)

var adcBuffer [64]struct {
	vSense uint16
	tSense uint16
	tMcu   uint16
}

func initAdc() {
	pinVSense.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	pinTSense.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})

	py32.RCC.SetAPBRSTR2_SYSCFGRST(1) // Reset SYSCFG
	py32.RCC.SetAPBENR2_SYSCFGEN(1)   // Enable SYSCFG clock - needed for DMA mapping
	py32.RCC.SetAPBRSTR2_SYSCFGRST(0) // Release SYSCFG from reset

	py32.RCC.SetAPBRSTR2_ADCRST(1) // Reset ADC
	py32.RCC.SetAPBENR2_ADCEN(1)   // Enable ADC clock
	py32.RCC.SetAPBRSTR2_ADCRST(0) // Release ADC from reset

	py32.RCC.SetAHBRSTR_DMARST(1) // Reset DMA
	py32.RCC.SetAHBENR_DMAEN(1)   // Enable DMA clock
	py32.RCC.SetAHBRSTR_DMARST(0) // Release DMA from reset

	dma := &py32.DMA.CH[0]
	dma.SetPAR(uint32(uintptr(unsafe.Pointer(&py32.ADC.DR.Reg))))
	dma.SetMAR(uint32(uintptr(unsafe.Pointer(&adcBuffer[0]))))
	dma.SetNDTR_NDT(uint32(len(adcBuffer) * 3)) // three words per entry
	dma.SetCR_DIR(0)                            // 0: peripheral to memory
	dma.SetCR_PSIZE(1)                          // 01: 16-bit
	dma.SetCR_MSIZE(1)                          // 01: 16-bit
	dma.SetCR_MINC(1)                           // memory increment
	dma.SetCR_PINC(0)                           // peripheral no increment
	dma.SetCR_CIRC(1)                           // circular buffer
	dma.SetCR_EN(1)                             // enable DMA channel

	py32.ADC.SetCFGR2_CKMODE(py32.ADC_CFGR2_CKMODE_HSI_Div64)
	py32.ADC.SetSMPR_SMP(py32.ADC_SMPR_SMP_Cycles239_5)                                   // Sampling time 239.5 ADC clock cycles
	py32.ADC.SetCFGR1_CONT(1)                                                             // continuous conversion mode
	py32.ADC.SetCFGR1_OVRMOD(1)                                                           // overwrite mode
	py32.ADC.SetCFGR1_DMAEN(1)                                                            // enable DMA
	py32.ADC.SetCFGR1_DMACFG(1)                                                           // DMA in circular mode
	py32.ADC.CHSELR.Set((1 << channelVSense) | (1 << channelTSense) | (1 << channelTMCU)) // select channels 8, 9 and 11
	py32.ADC.SetCCR_TSEN(1)
	py32.ADC.SetCCR_VREFEN(1)
	py32.ADC.SetCR_ADEN(1)    // enable ADC
	py32.ADC.SetCR_ADSTART(1) // start ADC

}

var vSenseLut = [...]util.LutPoint[uint16, uint16]{
	{3859, 50},
	{3771, 60},
	{3671, 70},
	{3562, 80},
	{3445, 90},
	{3319, 100},
	{3184, 110},
	{3041, 120},
	{2895, 130},
	{2757, 140},
	{2602, 150},
	{2440, 160},
	{2272, 170},
	{2099, 180},
	{1937, 190},
	{1768, 200},
	{1597, 210},
	{1420, 220},
	{1260, 230},
	{1098, 240},
	{918, 250},
}

var TSCAL1 = (*int32)(unsafe.Pointer(uintptr(0x1FFF0F14)))
var TSCAL2 = (*int32)(unsafe.Pointer(uintptr(0x1FFF0F18)))

func readAdcValues() (vSense uint32, tSense, tMcu int32) {

	for _, e := range adcBuffer {
		vSense += uint32(e.vSense)
		tSense += int32(e.tSense)
		tMcu += int32(e.tMcu)
	}
	vSense /= uint32(len(adcBuffer))
	tSense /= int32(len(adcBuffer))
	tMcu /= int32(len(adcBuffer))

	tMcu = 30 + (tMcu-*TSCAL1)*(85-30)/(*TSCAL2-*TSCAL1)
	tSense = int32(util.Interpolate(uint16(tSense), util.Divider_47k_ntc3950[:]) / 10)

	vSense = uint32(util.Interpolate(uint16(vSense), vSenseLut[:]))
	vSense = uint32(util.Compensate(100, 75, 250, 193, tSense, int32(vSense)))

	return vSense, tSense, tMcu
}
