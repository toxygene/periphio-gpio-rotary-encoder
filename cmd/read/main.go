package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/toxygene/periphio-gpio-rotary-encoder/device"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

func main() {
	pinAName := flag.String("pina", "", "pin name for a channel of rotary encoder")
	pinBName := flag.String("pinb", "", "pin name for a channel of rotary encoder")
	help := flag.Bool("h", false, "print help page")

	flag.Parse()

	if *help || *pinAName == "" || *pinBName == "" {
		flag.Usage()
		os.Exit(0)
	}

	if _, err := host.Init(); err != nil {
		panic(err)
	}

	aPin := gpioreg.ByName(*pinAName)
	bPin := gpioreg.ByName(*pinBName)

	re := device.NewRotaryEncoder(aPin, bPin)

	fmt.Println("reading...")

	for {
		spew.Dump(re.Read())
	}
}
