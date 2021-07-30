package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/toxygene/periphio-gpio-rotary-encoder/device"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpiotest"
	"testing"
	"time"
)

func TestRotaryEncoder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		aPin := &gpiotest.Pin{EdgesChan: make(chan gpio.Level)}
		bPin := &gpiotest.Pin{EdgesChan: make(chan gpio.Level)}

		rotaryEncoder := device.NewRotaryEncoderWithCustomTimeout(
			aPin,
			bPin,
			10*time.Second,
		)

		go func() {
			aPin.EdgesChan <- gpio.Low
			bPin.EdgesChan <- gpio.Low
			aPin.EdgesChan <- gpio.High
			bPin.EdgesChan <- gpio.High
		}()

		a := rotaryEncoder.Read()

		assert.Equal(t, device.CW, a)
	})
}
