package main

import (
	"device/py32"
	"machine"
	"time"
)

func main() {

	machine.PB8.Configure(machine.PinConfig{Mode: machine.PinAlternate})
	machine.PA0.Configure(machine.PinConfig{Mode: machine.PinAlternate})
	machine.PA1.Configure(machine.PinConfig{Mode: machine.PinAlternate})
	machine.PA2.Configure(machine.PinConfig{Mode: machine.PinAlternate})
	machine.PA3.Configure(machine.PinConfig{Mode: machine.PinAlternate})
	machine.PA4.Configure(machine.PinConfig{Mode: machine.PinAlternate})
	machine.PA5.Configure(machine.PinConfig{Mode: machine.PinAlternate})
	machine.PA6.Configure(machine.PinConfig{Mode: machine.PinAlternate})
	machine.PB4.Configure(machine.PinConfig{Mode: machine.PinAlternate})
	machine.PB3.Configure(machine.PinConfig{Mode: machine.PinAlternate})
	machine.PA15.Configure(machine.PinConfig{Mode: machine.PinAlternate})

	machine.PB8.SetAltFunc(3)
	machine.PA0.SetAltFunc(3)
	machine.PA1.SetAltFunc(3)
	machine.PA2.SetAltFunc(3)
	machine.PA3.SetAltFunc(3)
	machine.PA4.SetAltFunc(3)
	machine.PA5.SetAltFunc(3)
	machine.PA6.SetAltFunc(3)
	machine.PB4.SetAltFunc(6)
	machine.PB3.SetAltFunc(6)
	machine.PA15.SetAltFunc(6)

	py32.RCC.SetAPBENR2_LEDEN(1)

	py32.LED.SetTR_T2(100)
	py32.LED.SetTR_T1(100)
	py32.LED.SetPR(7)

	py32.LED.SetCR_EHS(1)
	py32.LED.SetCR_LED_COM_SEL(2)
	py32.LED.SetCR_LEDON(1)

	py32.LED.SetDR2_DATA2_A(1)
	// py32.LED.SetDR2_DATA2_B(1)
	// py32.LED.SetDR3_DATA3_C(1)
	// py32.LED.SetDR0_DATA0_D(1)
	// py32.LED.SetDR2_DATA2_E(1)
	// py32.LED.SetDR2_DATA2_F(1)
	// py32.LED.SetDR2_DATA2_G(1)
	// py32.LED.SetDR2_DATA2_DP(1)

	for {
		time.Sleep(1 * time.Second)
	}
}
