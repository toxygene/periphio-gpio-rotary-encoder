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
			time.Millisecond,
		)

		go func() {
			aPin.EdgesChan <- gpio.Low
			time.Sleep(time.Millisecond)
			bPin.EdgesChan <- gpio.Low
			time.Sleep(time.Millisecond)
			aPin.EdgesChan <- gpio.High
			time.Sleep(time.Millisecond)
			bPin.EdgesChan <- gpio.High
			time.Sleep(time.Millisecond)
			aPin.EdgesChan <- gpio.Low
			time.Sleep(time.Millisecond)
			bPin.EdgesChan <- gpio.Low
		}()

		assert.Equal(t, device.CW, rotaryEncoder.Read())
	})
}
