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
	CW  Action = "clockwise"
	CCW Action = "counterClockwise"
)

type RotaryEncoder struct {
	aPin         gpio.PinIO
	bPin         gpio.PinIO
	encoderState uint8
	logger       *logrus.Entry
	timeout      time.Duration
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
		return t.waitForEdgeOnPin(ctx, &mu, actions, t.aPin)
	})

	g.Go(func() error {
		return t.waitForEdgeOnPin(ctx, &mu, actions, t.bPin)
	})

	t.logger.Trace("starting rotary encoder run group")

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error: %w", err)
	}

	t.logger.Trace("rotary encoder run group finished")

	return nil
}

func (t *RotaryEncoder) waitForEdgeOnPin(ctx context.Context, mu *sync.Mutex, c chan<- Action, pin gpio.PinIO) error {
	logger := t.logger.WithField("pin", pin)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			logger.WithField("timeout", t.timeout).Trace("waiting for edge")

			if pin.WaitForEdge(t.timeout) == false {
				logger.Trace("no edge detected within timeout")

				continue
			}

			logger.Trace("edge detected, checking encoder state")

			mu.Lock()

			select {
			case <-ctx.Done():
				mu.Unlock()
				return ctx.Err()
			default:
				encoderState := t.getEncoderState()

				encoderStateLogger := logger.WithField("encoder_state", fmt.Sprintf("%#08b", encoderState))

				if encoderState == 0x4b || encoderState == 0x2d || encoderState == 0xb4 || encoderState == 0xd2 {
					encoderStateLogger.Trace("clockwise rotation detected")

					t.encoderState = 0
					c <- CW
				} else if encoderState == 0x87 || encoderState == 0x1e || encoderState == 0x78 || encoderState == 0xe1 {
					encoderStateLogger.Trace("counter clockwise rotation detected")

					t.encoderState = 0
					c <- CCW
				}

				encoderStateLogger.Trace("no rotation detected")

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
