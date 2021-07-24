package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/toxygene/periphio-gpio-rotary-encoder/device"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
	"syscall"
	"time"
)

func main() {
	pinAName := flag.String("pina", "", "pin name for a channel of rotary encoder")
	pinBName := flag.String("pinb", "", "pin name for b channel of rotary encoder")
	buttonPinName := flag.String("button", "", "pin name for the button channel of the rotary encoder")
	timeout := flag.Int("timeout", 2, "timeout (in seconds) for reading a pin")
	help := flag.Bool("h", false, "print help page")
	verbose := flag.Bool("verbose", false, "print debugging information")

	flag.Parse()

	if *help || *pinAName == "" || *pinBName == "" || *buttonPinName == "" || *timeout == 0 {
		flag.Usage()
		os.Exit(0)
	}

	log := logrus.New()

	if *verbose {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.ErrorLevel)
	}

	logger := logrus.NewEntry(log)

	if _, err := host.Init(); err != nil {
		panic(err)
	}

	aPin := gpioreg.ByName(*pinAName)
	if aPin == nil {
		logger.WithField("pin_a", *pinAName).Error("no gpio pin found for pin a")
		os.Exit(1)
	}

	bPin := gpioreg.ByName(*pinBName)
	if bPin == nil {
		logger.WithField("pin_b", *pinBName).Error("no gpio pin found for pin b")
		os.Exit(1)
	}

	buttonPin := gpioreg.ByName(*buttonPinName)
	if buttonPin == nil {
		logger.WithField("button", *buttonPinName).Error("no gpio bin found for button")
		os.Exit(1)
	}

	re := device.NewRotaryEncoder(aPin, bPin, buttonPin, (time.Duration(*timeout))*time.Second, logger)

	g := errgroup.Group{}

	ctx, cancel := context.WithCancel(context.Background())
	actions := make(chan device.Action)

	// Run the rotary encoder
	g.Go(func() error {
		logger.Info("starting rotary encoder")

		err := re.Run(ctx, actions)

		logger.WithField("channel", actions).Info("rotary encoder done, closing results channel")

		close(actions)

		logger.Info("shutting down rotary encoder")

		return err
	})

	// Print the actions until cancellation
	g.Go(func() error {
		logger.Info("starting action printer")

		for {
			logger.Info("waiting for action")

			a, ok := <-actions

			if !ok {
				logger.WithField("channel", actions).Info("channel closed, shutting down")
				return nil
			}

			logger.WithField("action", a).Info("action received")
			fmt.Println(a)
		}
	})

	g.Go(func() error {
		logger.Info("starting interrupt signal handler")

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		logger.Info("waiting for sigterm")

		<-c

		logger.Info("sigterm caught, cancelling context")

		cancel()

		return nil
	})

	logger.Info("starting application run group")

	if err := g.Wait(); err != nil {
		panic(err)
	}

	logger.Info("application run group complete")
}
