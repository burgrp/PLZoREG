package main

import (
	"device/py32"
	"machine"
	"runtime/interrupt"
)

const (
	pinTriac     = machine.PA8
	pinZCD       = machine.PA10
	triggerWidth = 500 // triac trigger pulse width in microseconds
)

var (
	tMeasure = py32.TIM3
	tPwm     = py32.TIM1
	zcdWidth uint32
	period   uint32
)

func initPwm() {

	// pin setup
	pinTriac.SetAltFunc(2) // TIM1_CH1
	pinTriac.Configure(machine.PinConfig{Mode: machine.PinAlternate})
	pinZCD.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// exti setup
	interrupt.New(py32.IRQ_EXTI4_15, EXTI4_15_IRQHandler).Enable()
	py32.EXTI.SetRTSR_RT10(1) // Rising trigger for line 10 (PA10)
	py32.EXTI.SetFTSR_FT10(1) // Falling trigger for line 10 (PA10)
	py32.EXTI.SetIMR_IM10(1)  // Unmask interrupt for line 10 (PA10)

	// timer setup
	py32.RCC.SetAPBRSTR1_TIM3RST(1) // Reset TIM3
	py32.RCC.SetAPBENR1_TIM3EN(1)   // Enable TIM3 clock
	py32.RCC.SetAPBRSTR1_TIM3RST(0) // Release TIM3 from reset

	py32.RCC.SetAPBRSTR2_TIM1RST(1) // Reset TIM1
	py32.RCC.SetAPBENR2_TIM1EN(1)   // Enable TIM1 clock
	py32.RCC.SetAPBRSTR2_TIM1RST(0) // Release TIM1 from reset

	tMeasure.SetPSC(23)     // 24MHz / (23+1) = 1MHz timer clock
	tMeasure.SetCR1_OPM(1)  // One pulse mode
	tMeasure.SetARR(0xFFFF) // Reload means missing ZCD

	tPwm.SetPSC(23)              // the same as tMeasure
	tPwm.SetCR1_OPM(1)           // One pulse mode
	tPwm.SetCCMR1_Output_OC1M(7) // PWM mode 2
	tPwm.SetCCER_CC1E(1)         // Enable output on CH1
	tPwm.SetBDTR_MOE(1)          // Main output enable

	setPwm(0)

	// go func() {
	// 	for {
	// 		println("PWM period:", period, "ZCD width:", zcdWidth, "Duty:", State.Duty)
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

}

func setPwm(duty uint32) {
	if zcdWidth > 0 && period > 0 {
		start := (period*(100-duty))/100 + zcdWidth/2
		if start == 0 { // CC match doesn't work at 0
			start = 1
		}
		tPwm.SetCCR1(start)
		tPwm.SetARR(start + triggerWidth)
	}
}

func isSynchronized() bool {
	return tMeasure.GetCR1_CEN() == 1
}

func EXTI4_15_IRQHandler(i interrupt.Interrupt) {

	if py32.EXTI.GetPR_PR10() == 1 {

		py32.EXTI.SetPR_PR10(1) // Clear pending bit for line 10 (PA10)

		if pinZCD.Get() {

			// rising edge

			if tMeasure.GetCR1_CEN() == 1 {
				period = tMeasure.GetCNT()
			}

			tMeasure.SetEGR_UG(1)  // Update generation - reset timer
			tMeasure.SetCR1_CEN(1) // Start timer

			tPwm.SetEGR_UG(1)  // Update generation - reset timer
			tPwm.SetCR1_CEN(1) // Start timer

		} else {

			// falling edge

			zcdWidth = tMeasure.GetCNT()

		}

	}
}
