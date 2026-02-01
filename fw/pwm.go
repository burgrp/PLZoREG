package main

import (
	"device/py32"
	"machine"
	"runtime/interrupt"
	"time"
)

const (
	pinTriac = machine.PA8
	pinZCD   = machine.PA10
)

var zcdWidth uint32

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
	py32.RCC.SetAPBENR2_TIM1EN(1)     // Enable TIM1 clock
	py32.TIM1.SetPSC(23)              // 24MHz / (23+1) = 1MHz timer clock
	py32.TIM1.SetCR1_OPM(1)           // One pulse mode
	py32.TIM1.SetCCMR1_Output_OC1M(7) // PWM mode 2
	py32.TIM1.SetCCER_CC1E(1)         // Enable output on CH1
	py32.TIM1.SetBDTR_MOE(1)          // Main output enable

	py32.TIM1.SetARR(6000)
	py32.TIM1.SetCCR1(5000)

	go func() {
		for {
			println(zcdWidth)
			time.Sleep(1 * time.Second)
		}
	}()

}

func EXTI4_15_IRQHandler(i interrupt.Interrupt) {

	if py32.EXTI.GetPR_PR10() == 1 {

		py32.EXTI.SetPR_PR10(1) // Clear pending bit for line 10 (PA10)

		if pinZCD.Get() {
			// rising edge
			py32.TIM1.SetEGR_UG(1)  // Update generation - reset timer
			py32.TIM1.SetCR1_CEN(1) // Start timer

		} else {
			// falling edge
			zcdWidth = py32.TIM1.GetCNT()
			State.Error = ErrorNoSync
		}

	}
}
