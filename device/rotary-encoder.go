package device

import (
	"sync"

	"periph.io/x/periph/conn/gpio"
)

type Action string

const (
	none Action = "none"
	cw   Action = "clockwise"
	ccw  Action = "counterClockwise"
)

type RotaryEncoder struct {
	aPin                 gpio.PinIO
	bPin                 gpio.PinIO
	previousEncoderState uint8
}

func NewRotaryEncoder(aPin gpio.PinIO, bPin gpio.PinIO) *RotaryEncoder {
	aPin.In(gpio.PullUp, gpio.BothEdges)
	bPin.In(gpio.PullUp, gpio.BothEdges)

	return &RotaryEncoder{
		aPin: aPin,
		bPin: bPin,
	}
}

func (t *RotaryEncoder) Read() Action {
	c := make(chan Action, 1)
	defer close(c)

	mu := sync.Mutex{}
	stop := false

	go func() {
		for {
			if stop {
				return
			}
			if t.aPin.WaitForEdge(-1) == false {
				return
			}

			mu.Lock()
			a := t.handleEdge()

			if a == cw || a == ccw {
				if !stop {
					c <- a
				}
				stop = true
				mu.Unlock()
				return
			}

			mu.Unlock()
		}
	}()

	go func() {
		for {
			if stop {
				return
			}
			if t.bPin.WaitForEdge(-1) == false {
				return
			}

			mu.Lock()
			a := t.handleEdge()

			if a == cw || a == ccw {
				if !stop {
					c <- a
				}
				stop = true
				mu.Unlock()
				return
			}

			mu.Unlock()
		}
	}()

	return <-c
}

func (t *RotaryEncoder) handleEdge() Action {
	encoderValue := t.readCurrentEncoderValue()

	if encoderValue == (t.previousEncoderState & 3) {
		return none
	}

	encoderState := (t.previousEncoderState << 2) | encoderValue

	if encoderState == 0x1e || encoderState == 0xe1 || encoderState == 0x78 || encoderState == 0x87 {
		t.previousEncoderState = 0
		return cw
	} else if encoderState == 0xb4 || encoderState == 0x4b || encoderState == 0x2d || encoderState == 0xd2 {
		t.previousEncoderState = 0
		return ccw
	}

	t.previousEncoderState = encoderState

	return none
}

func (t *RotaryEncoder) readCurrentEncoderValue() uint8 {
	x := uint8(0)

	if t.bPin.Read() == gpio.High {
		x |= 2
	}

	if t.aPin.Read() == gpio.High {
		x |= 1
	}

	return x
}
