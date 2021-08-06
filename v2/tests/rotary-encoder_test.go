package tests

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/toxygene/periphio-gpio-rotary-encoder/v2/device"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpiotest"
	"testing"
	"time"
)

func TestRotaryEncoder(t *testing.T) {
	t.Run("success clockwise 1", func(t *testing.T) {
		aPin := &gpiotest.Pin{N: "A", EdgesChan: make(chan gpio.Level)}
		bPin := &gpiotest.Pin{N: "B", EdgesChan: make(chan gpio.Level)}

		rotaryEncoder := device.NewRotaryEncoder(
			aPin,
			bPin,
			time.Millisecond,
			logrus.NewEntry(logrus.New()),
		)

		ctx, cancel := context.WithCancel(context.Background())
		actions := make(chan device.Action)

		go func() {
			err := rotaryEncoder.Run(ctx, actions)

			assert.Error(t, err, "context cancellation")
		}()

		aPin.EdgesChan <- gpio.Low
		time.Sleep(time.Millisecond)
		bPin.EdgesChan <- gpio.Low
		time.Sleep(time.Millisecond)
		aPin.EdgesChan <- gpio.High
		time.Sleep(time.Millisecond)
		bPin.EdgesChan <- gpio.High

		assert.Equal(t, device.CW, <-actions)

		cancel()
	})

	t.Run("success clockwise 2", func(t *testing.T) {
		aPin := &gpiotest.Pin{N: "A", EdgesChan: make(chan gpio.Level)}
		bPin := &gpiotest.Pin{N: "B", EdgesChan: make(chan gpio.Level)}

		rotaryEncoder := device.NewRotaryEncoder(
			aPin,
			bPin,
			time.Millisecond,
			logrus.NewEntry(logrus.New()),
		)

		ctx, cancel := context.WithCancel(context.Background())
		actions := make(chan device.Action)

		go func() {
			err := rotaryEncoder.Run(ctx, actions)

			assert.Error(t, err, "context cancellation")
		}()

		bPin.EdgesChan <- gpio.Low
		time.Sleep(time.Millisecond)
		aPin.EdgesChan <- gpio.High
		time.Sleep(time.Millisecond)
		bPin.EdgesChan <- gpio.High
		time.Sleep(time.Millisecond)
		aPin.EdgesChan <- gpio.Low

		assert.Equal(t, device.CW, <-actions)

		cancel()
	})

	t.Run("success clockwise 3", func(t *testing.T) {
		aPin := &gpiotest.Pin{N: "A", EdgesChan: make(chan gpio.Level)}
		bPin := &gpiotest.Pin{N: "B", EdgesChan: make(chan gpio.Level)}

		rotaryEncoder := device.NewRotaryEncoder(
			aPin,
			bPin,
			time.Millisecond,
			logrus.NewEntry(logrus.New()),
		)

		ctx, cancel := context.WithCancel(context.Background())
		actions := make(chan device.Action)

		go func() {
			err := rotaryEncoder.Run(ctx, actions)

			assert.Error(t, err, "context cancellation")
		}()

		aPin.EdgesChan <- gpio.High
		time.Sleep(time.Millisecond)
		bPin.EdgesChan <- gpio.High
		time.Sleep(time.Millisecond)
		aPin.EdgesChan <- gpio.Low
		time.Sleep(time.Millisecond)
		bPin.EdgesChan <- gpio.Low

		assert.Equal(t, device.CW, <-actions)

		cancel()
	})

	t.Run("success clockwise 4", func(t *testing.T) {
		aPin := &gpiotest.Pin{N: "A", EdgesChan: make(chan gpio.Level)}
		bPin := &gpiotest.Pin{N: "B", EdgesChan: make(chan gpio.Level)}

		rotaryEncoder := device.NewRotaryEncoder(
			aPin,
			bPin,
			time.Millisecond,
			logrus.NewEntry(logrus.New()),
		)

		ctx, cancel := context.WithCancel(context.Background())
		actions := make(chan device.Action)

		go func() {
			err := rotaryEncoder.Run(ctx, actions)

			assert.Error(t, err, "context cancellation")
		}()

		bPin.EdgesChan <- gpio.High
		time.Sleep(time.Millisecond)
		aPin.EdgesChan <- gpio.Low
		time.Sleep(time.Millisecond)
		bPin.EdgesChan <- gpio.Low
		time.Sleep(time.Millisecond)
		aPin.EdgesChan <- gpio.High

		assert.Equal(t, device.CW, <-actions)

		cancel()
	})

	t.Run("success counter clockwise 1", func(t *testing.T) {
		aPin := &gpiotest.Pin{EdgesChan: make(chan gpio.Level)}
		bPin := &gpiotest.Pin{EdgesChan: make(chan gpio.Level)}

		rotaryEncoder := device.NewRotaryEncoder(
			aPin,
			bPin,
			time.Millisecond,
			logrus.NewEntry(logrus.New()),
		)

		ctx, cancel := context.WithCancel(context.Background())
		actions := make(chan device.Action)

		go func() {
			err := rotaryEncoder.Run(ctx, actions)

			assert.Error(t, err, "context cancellation")
		}()

		bPin.EdgesChan <- gpio.Low
		time.Sleep(time.Millisecond)
		aPin.EdgesChan <- gpio.Low
		time.Sleep(time.Millisecond)
		bPin.EdgesChan <- gpio.High
		time.Sleep(time.Millisecond)
		aPin.EdgesChan <- gpio.High

		assert.Equal(t, device.CCW, <-actions)

		cancel()
	})

	t.Run("success counter clockwise 2", func(t *testing.T) {
		aPin := &gpiotest.Pin{EdgesChan: make(chan gpio.Level)}
		bPin := &gpiotest.Pin{EdgesChan: make(chan gpio.Level)}

		rotaryEncoder := device.NewRotaryEncoder(
			aPin,
			bPin,
			time.Millisecond,
			logrus.NewEntry(logrus.New()),
		)

		ctx, cancel := context.WithCancel(context.Background())
		actions := make(chan device.Action)

		go func() {
			err := rotaryEncoder.Run(ctx, actions)

			assert.Error(t, err, "context cancellation")
		}()

		aPin.EdgesChan <- gpio.Low
		time.Sleep(time.Millisecond)
		bPin.EdgesChan <- gpio.High
		time.Sleep(time.Millisecond)
		aPin.EdgesChan <- gpio.High
		time.Sleep(time.Millisecond)
		bPin.EdgesChan <- gpio.Low

		assert.Equal(t, device.CCW, <-actions)

		cancel()
	})

	t.Run("success counter clockwise 3", func(t *testing.T) {
		aPin := &gpiotest.Pin{EdgesChan: make(chan gpio.Level)}
		bPin := &gpiotest.Pin{EdgesChan: make(chan gpio.Level)}

		rotaryEncoder := device.NewRotaryEncoder(
			aPin,
			bPin,
			time.Millisecond,
			logrus.NewEntry(logrus.New()),
		)

		ctx, cancel := context.WithCancel(context.Background())
		actions := make(chan device.Action)

		go func() {
			err := rotaryEncoder.Run(ctx, actions)

			assert.Error(t, err, "context cancellation")
		}()

		bPin.EdgesChan <- gpio.High
		time.Sleep(time.Millisecond)
		aPin.EdgesChan <- gpio.High
		time.Sleep(time.Millisecond)
		bPin.EdgesChan <- gpio.Low
		time.Sleep(time.Millisecond)
		aPin.EdgesChan <- gpio.Low

		assert.Equal(t, device.CCW, <-actions)

		cancel()
	})

	t.Run("success counter clockwise 4", func(t *testing.T) {
		aPin := &gpiotest.Pin{EdgesChan: make(chan gpio.Level)}
		bPin := &gpiotest.Pin{EdgesChan: make(chan gpio.Level)}

		rotaryEncoder := device.NewRotaryEncoder(
			aPin,
			bPin,
			time.Millisecond,
			logrus.NewEntry(logrus.New()),
		)

		ctx, cancel := context.WithCancel(context.Background())
		actions := make(chan device.Action)

		go func() {
			err := rotaryEncoder.Run(ctx, actions)

			assert.Error(t, err, "context cancellation")
		}()

		aPin.EdgesChan <- gpio.High
		time.Sleep(time.Millisecond)
		bPin.EdgesChan <- gpio.Low
		time.Sleep(time.Millisecond)
		aPin.EdgesChan <- gpio.Low
		time.Sleep(time.Millisecond)
		bPin.EdgesChan <- gpio.High

		assert.Equal(t, device.CCW, <-actions)

		cancel()
	})
}
