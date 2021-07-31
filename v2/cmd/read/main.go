package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/toxygene/periphio-gpio-rotary-encoder/v2/device"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
	"syscall"
	"time"
)

func main() {
	pinAName := flag.String("pina", "", "pin name for a channel of rotary encoder")
	pinBName := flag.String("pinb", "", "pin name for b channel of rotary encoder")
	timeout := flag.Int("timeout", 2, "timeout (in seconds) for reading a pin")
	help := flag.Bool("h", false, "print help page")
	logging := flag.String("logging", "", "logging level")

	flag.Parse()

	if *help || *pinAName == "" || *pinBName == "" || *timeout == 0 {
		flag.Usage()
		os.Exit(0)
	}

	log := logrus.New()

	if *logging != "" {
		logLevel, err := logrus.ParseLevel(*logging)
		if err != nil {
			panic(err)
		}

		log.SetLevel(logLevel)
	}

	logger := logrus.NewEntry(log)

	if _, err := host.Init(); err != nil {
		panic(err)
	}

	aPin := gpioreg.ByName(*pinAName)
	if aPin == nil {
		logger.WithField("pin_a", *pinAName).Error("no gpio pin found for pin a")

		panic("could not find pin a")
	}

	if err := aPin.In(gpio.PullNoChange, gpio.BothEdges); err != nil {
		panic(err)
	}

	bPin := gpioreg.ByName(*pinBName)
	if bPin == nil {
		logger.WithField("pin_b", *pinBName).Error("no gpio pin found for pin b")

		panic("could not find pin b")
	}

	if err := bPin.In(gpio.PullNoChange, gpio.BothEdges); err != nil {
		panic(err)
	}

	re := device.NewRotaryEncoder(aPin, bPin, (time.Duration(*timeout))*time.Second, logger)

	g := errgroup.Group{}

	ctx, cancel := context.WithCancel(context.Background())
	actions := make(chan device.Action)

	// Run the rotary encoder
	g.Go(func() error {
		defer close(actions)

		logger.Trace("starting rotary encoder")

		err := re.Run(ctx, actions)

		logger.Trace("shutting down rotary encoder")

		return err
	})

	// Print the actions until cancellation
	g.Go(func() error {
		logger.Trace("starting action printer")

		for action := range actions {
			logger.WithField("action", action).Trace("action received")
			fmt.Println(action)
		}

		return nil
	})

	g.Go(func() error {
		logger.Trace("starting interrupt signal handler")

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		logger.Trace("waiting for sigterm")

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			<-c

			logger.Trace("sigterm caught, cancelling context")

			cancel()
			return nil
		}
	})

	logger.Trace("starting application run group")

	if err := g.Wait(); err != nil {
		panic(err)
	}

	logger.Trace("application run group complete")
}
