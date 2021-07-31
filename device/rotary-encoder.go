package device

import (
	"context"
	"sync"
	"time"

	"periph.io/x/periph/conn/gpio"
)

type Action string

const (
	CW  Action = "clockwise"
	CCW Action = "counterClockwise"
)

type RotaryEncoder struct {
	aPin         gpio.PinIO
	bPin         gpio.PinIO
	encoderState uint8
	timeout      time.Duration
}

func NewRotaryEncoder(aPin gpio.PinIO, bPin gpio.PinIO) *RotaryEncoder {
	return &RotaryEncoder{
		aPin:    aPin,
		bPin:    bPin,
		timeout: 1 * time.Second,
	}
}

func NewRotaryEncoderWithCustomTimeout(aPin gpio.PinIO, bPin gpio.PinIO, timeout time.Duration) *RotaryEncoder {
	return &RotaryEncoder{
		aPin:    aPin,
		bPin:    bPin,
		timeout: timeout,
	}
}

func (t *RotaryEncoder) Read() Action {
	c := make(chan Action, 1)
	defer close(c)

	t.aPin.In(gpio.PullNoChange, gpio.BothEdges)
	t.bPin.In(gpio.PullNoChange, gpio.BothEdges)

	mu := &sync.Mutex{}

	ctx, cancel := context.WithCancel(context.Background())

	go t.waitForEdgeOnPin(ctx, cancel, mu, c, t.aPin)
	go t.waitForEdgeOnPin(ctx, cancel, mu, c, t.bPin)

	return <-c
}

func (t *RotaryEncoder) waitForEdgeOnPin(ctx context.Context, cancel context.CancelFunc, mu *sync.Mutex, c chan<- Action, pin gpio.PinIO) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if pin.WaitForEdge(t.timeout) == false {
				continue
			}

			mu.Lock()

			select {
			case <-ctx.Done():
				mu.Unlock()
				return
			default:
				encoderState := t.getEncoderState()

				if encoderState == 0x4b || encoderState == 0x2d || encoderState == 0xb4 || encoderState == 0xd2 {
					t.encoderState = 0
					cancel()
					c <- CW
					mu.Unlock()
					return
				} else if encoderState == 0x87 || encoderState == 0x1e || encoderState == 0x78 || encoderState == 0xe1 {
					t.encoderState = 0
					cancel()
					c <- CCW
					mu.Unlock()
					return
				}

				mu.Unlock()
			}
		}
	}
}

func (t *RotaryEncoder) getEncoderState() uint8 {
	encoderState := t.readEncoderState()

	if encoderState == (t.encoderState & 3) {
		return t.encoderState
	}

	t.encoderState = (t.encoderState << 2) | encoderState
	return t.encoderState
}

func (t *RotaryEncoder) readEncoderState() uint8 {
	x := uint8(0)

	if t.aPin.Read() == gpio.High {
		x |= 2
	}

	if t.bPin.Read() == gpio.High {
		x |= 1
	}

	return x
}
