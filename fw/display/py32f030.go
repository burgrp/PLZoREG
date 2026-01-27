//go:build py32f030

package display

import (
	"machine"
)

func configureGpio() {

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
}
