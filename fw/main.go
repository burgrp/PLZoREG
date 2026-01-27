package main

import (
	"machine"
	"template/display"
	"template/keyboard"
	"time"
)

func main() {

	display.Init(3, 200)

	kbd := keyboard.New(keyboard.Keys{
		{Pin: machine.PA11, ID: 1},
		{Pin: machine.PA12, ID: -1},
	})

	go kbd.Run()

	n := 100

	var bothDownMs int64

	for {
		display.NumberAt(n, false, -1, 2)
		e := <-kbd.Events
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
				} else if time.Now().UnixMilli()-bothDownMs > 3000 {
					n = 100
					bothDownMs = 0
				}
			} else {
				bothDownMs = 0
				n += e.ID
			}
		}
		// if e.Kind == keyboard.KeyDown {
		// 	// detect both keys down for 1 second to reset
		// 	allDown := true
		// 	for _, key := range *kbd.Keys {
		// 		if key.DownUs == 0 {
		// 			allDown = false
		// 		}
		// 	}
		// }
	}
}
