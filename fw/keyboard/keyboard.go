package keyboard

import (
	"machine"
	"time"
)

type Keys []struct {
	Pin    machine.Pin
	DownMs int64
	ID     int
}

type KeyEvent struct {
	ID   int
	Kind KeyEventKind
}

type KeyEventKind int

const (
	KeyUp KeyEventKind = iota
	KeyPress
	KeyDown
)

type Keyboard struct {
	Keys   *Keys
	Events chan KeyEvent
}

func New(keys Keys) *Keyboard {
	for _, key := range keys {
		key.Pin.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	return &Keyboard{Keys: &keys, Events: make(chan KeyEvent)}
}

func (kbd *Keyboard) Run() {
	for {
		for i := range *kbd.Keys {
			key := &(*kbd.Keys)[i]
			down := !key.Pin.Get()

			if down {
				if key.DownMs == 0 {
					key.DownMs = time.Now().UnixMilli()
					kbd.Events <- KeyEvent{ID: key.ID, Kind: KeyDown}
					kbd.Events <- KeyEvent{ID: key.ID, Kind: KeyPress}
				} else {
					now := time.Now().UnixMilli()
					if now-key.DownMs > 500 {
						kbd.Events <- KeyEvent{ID: key.ID, Kind: KeyPress}
					}
				}
			} else {
				if key.DownMs != 0 {
					key.DownMs = 0
					kbd.Events <- KeyEvent{ID: key.ID, Kind: KeyUp}
				}
			}

		}
		time.Sleep(100 * time.Millisecond)
	}
}
