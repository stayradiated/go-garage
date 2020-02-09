package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

const (
	DOOR_OPEN = iota
	DOOR_SHUT = iota
)

var (
	addr      = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	doorState = DOOR_SHUT
)

func readDoorState(w http.ResponseWriter, req *http.Request) {
	message := "open"
	if doorState == DOOR_SHUT {
		message = "shut"
	}

	fmt.Fprintf(w, message)
}

func main() {
	flag.Parse()

	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	led := 0

	l, err := c.RequestLine(rpi.GPIO2, gpiod.AsOutput(led))
	if err != nil {
		panic(err)
	}
	defer func() {
		l.Reconfigure(gpiod.AsInput)
		l.Close()
	}()

	// capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				led ^= 1
				l.SetValue(led)
			case <-quit:
				return
			}
		}
	}()

	l2, err := c.RequestLine(rpi.GPIO6, gpiod.WithBothEdges(func(evt gpiod.LineEvent) {
		if evt.Type == gpiod.LineEventFallingEdge {
			doorState = DOOR_SHUT
			log.Println("DOOR_SHUT")
		} else {
			doorState = DOOR_OPEN
			log.Println("DOOR_OPEN")
		}
	}))
	if err != nil {
		panic(err)
	}
	defer l2.Close()

	http.HandleFunc("/state", readDoorState)
	http.ListenAndServe(*addr, nil)
}
