package device

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"periph.io/x/periph/conn/gpio"
	"sync"
	"time"
)

type Action string

const (
	none Action = "none"
	CW   Action = "clockwise"
	CCW  Action = "counterClockwise"
)

type RotaryEncoder struct {
	aPin                 gpio.PinIO
	bPin                 gpio.PinIO
	logger               *logrus.Entry
	previousEncoderState uint8
	timeout              time.Duration
}

func NewRotaryEncoder(aPin gpio.PinIO, bPin gpio.PinIO, timeout time.Duration, logger *logrus.Entry) *RotaryEncoder {
	return &RotaryEncoder{
		aPin:    aPin,
		bPin:    bPin,
		logger:  logger,
		timeout: timeout,
	}
}

func (t *RotaryEncoder) Run(ctx context.Context, actions chan<- Action) error {
	mu := sync.Mutex{}

	g := new(errgroup.Group)

	g.Go(func() error {
		logger := t.logger.WithField("pin_a", t.aPin)

		t.logger.Trace("setting up pin a")

		if err := t.aPin.In(gpio.PullNoChange, gpio.BothEdges); err != nil {
			t.logger.WithError(err).Error("setup of pin a failed")
			return fmt.Errorf("setup of pin a failed: %w", err)
		}

		for {
			select {
			case <-ctx.Done():
				logger.Trace("aborting wait for edge on pin")
				return nil
			default:
				logger.Trace("waiting for edge on pin a")
				if t.aPin.WaitForEdge(t.timeout) == false {
					continue
				}

				logger.Trace("edge encountered on pin a, updating rotary encoder state")
				mu.Lock()
				a := t.handleEdge()
				mu.Unlock()

				logger.WithField("action", a).Trace("action calculated for pin a edge")
				if a == CW || a == CCW {
					actions <- a
				}
			}
		}
	})

	g.Go(func() error {
		logger := t.logger.WithField("pin_b", t.bPin)

		logger.Trace("setting up pin b")

		if err := t.bPin.In(gpio.PullNoChange, gpio.BothEdges); err != nil {
			t.logger.WithError(err).Error("setup of pin b failed")
			return fmt.Errorf("setup of pin b failed: %w", err)
		}

		for {
			select {
			case <-ctx.Done():
				logger.Trace("aborting wait for edge on pin")
				return nil
			default:
				logger.Trace("waiting for edge on pin b")
				if t.bPin.WaitForEdge(t.timeout) == false {
					continue
				}

				logger.Trace("edge encountered on pin b, updating rotary encoder state")
				mu.Lock()
				a := t.handleEdge()
				mu.Unlock()

				logger.WithField("action", a).Trace("action calculated for pin b edge")
				if a == CW || a == CCW {
					actions <- a
				}
			}
		}
	})

	t.logger.Trace("starting rotary encoder run group")

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error: %w", err)
	}

	t.logger.Trace("rotary encoder run group finished")

	return nil
}

func (t *RotaryEncoder) handleEdge() Action {
	encoderValue := t.readCurrentEncoderValue()

	if encoderValue == (t.previousEncoderState & 3) {
		return none
	}

	encoderState := (t.previousEncoderState << 2) | encoderValue

	if encoderState == 0x1e || encoderState == 0xe1 || encoderState == 0x78 || encoderState == 0x87 {
		t.previousEncoderState = 0
		return CW
	} else if encoderState == 0xb4 || encoderState == 0x4b || encoderState == 0x2d || encoderState == 0xd2 {
		t.previousEncoderState = 0
		return CCW
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
