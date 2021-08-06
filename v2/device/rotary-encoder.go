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

type State struct {
	Pin   gpio.PinIO
	Level gpio.Level
}

type Action string

const (
	CW  Action = "clockwise"
	CCW Action = "counterClockwise"
)

type RotaryEncoder struct {
	aPin    gpio.PinIO
	bPin    gpio.PinIO
	state   [4]State
	logger  *logrus.Entry
	timeout time.Duration

	cw  [4]State
	ccw [4]State
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
	t.cw = [4]State{
		{
			Pin:   t.aPin,
			Level: gpio.High,
		},
		{
			Pin:   t.bPin,
			Level: gpio.High,
		},
		{
			Pin:   t.aPin,
			Level: gpio.Low,
		},
		{
			Pin:   t.bPin,
			Level: gpio.Low,
		},
	}

	t.ccw = [4]State{
		{
			Pin:   t.bPin,
			Level: gpio.High,
		},
		{
			Pin:   t.aPin,
			Level: gpio.High,
		},
		{
			Pin:   t.bPin,
			Level: gpio.Low,
		},
		{
			Pin:   t.aPin,
			Level: gpio.Low,
		},
	}

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
		t.logger.
			WithError(err).
			Error("rotary encoder failed")

		return fmt.Errorf("rotary encoder: %w", err)
	}

	t.logger.Trace("rotary encoder run group finished")

	return nil
}

func (t *RotaryEncoder) waitForEdgeOnPin(ctx context.Context, mu *sync.Mutex, c chan<- Action, pin gpio.PinIO) error {
	logger := t.logger.WithField("pin", pin)

	for {
		select {
		case <-ctx.Done():
			logger.Trace("context cancelled, shutting down goroutine")

			return ctx.Err()
		default:
			logger.WithField("timeout", t.timeout).Trace("waiting for edge")

			if pin.WaitForEdge(t.timeout) == false {
				logger.Trace("no edge detected within timeout")

				continue
			}

			logger.Trace("edge detected, checking encoder states")

			mu.Lock()

			select {
			case <-ctx.Done():
				logger.Trace("context cancelled, shutting down goroutine")
				mu.Unlock()
				return ctx.Err()
			default:
				state := State{Level: pin.Read(), Pin: pin}

				stateLogger := logger.WithField("state", state)

				stateLogger.Trace("current state")

				// Push the current state
				t.state = [4]State{
					t.state[1],
					t.state[2],
					t.state[3],
					state,
				}

				logger.WithField("states", t.state).Trace("checking encoder states")

				if (t.state[0] == t.cw[0] && t.state[1] == t.cw[1] && t.state[2] == t.cw[2] && t.state[3] == t.cw[3]) || (t.state[0] == t.cw[1] && t.state[1] == t.cw[2] && t.state[2] == t.cw[3] && t.state[3] == t.cw[0]) || (t.state[0] == t.cw[2] && t.state[1] == t.cw[3] && t.state[2] == t.cw[0] && t.state[3] == t.cw[1]) || (t.state[0] == t.cw[3] && t.state[1] == t.cw[0] && t.state[2] == t.cw[1] && t.state[3] == t.cw[2]) {
					logger.Trace("clockwise rotation detected")

					t.state = [4]State{}
					c <- CW
				} else if (t.state[0] == t.ccw[0] && t.state[1] == t.ccw[1] && t.state[2] == t.ccw[2] && t.state[3] == t.ccw[3]) || (t.state[0] == t.ccw[1] && t.state[1] == t.ccw[2] && t.state[2] == t.ccw[3] && t.state[3] == t.ccw[0]) || (t.state[0] == t.ccw[2] && t.state[1] == t.ccw[3] && t.state[2] == t.ccw[0] && t.state[3] == t.ccw[1]) || (t.state[0] == t.ccw[3] && t.state[1] == t.ccw[0] && t.state[2] == t.ccw[1] && t.state[3] == t.ccw[2]) {
					logger.Trace("counter-clockwise rotation detected")
					t.state = [4]State{}
					c <- CCW
				}
			}

			mu.Unlock()
		}
	}
}
