package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/warthog618/gpio"
)

const (
	DOOR_OPEN = iota
	DOOR_SHUT = iota
)

var (
	addr      = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	doorState = DOOR_SHUT
	rate      = time.Second / 2
	throttle  = time.Tick(rate)
)

func readDoorState(w http.ResponseWriter, req *http.Request) {
	message := "open"
	if doorState == DOOR_SHUT {
		message = "shut"
	}

	fmt.Fprintf(w, message)
}

func handler(pin *gpio.Pin) {
	<-throttle //rate limit

	level := pin.Read()
	if level == gpio.Low {
		log.Println("low")
		doorState = DOOR_SHUT
	} else {
		log.Println("high")
		doorState = DOOR_OPEN
	}
}

func main() {
	flag.Parse()

	err := gpio.Open()
	if err != nil {
		panic(err)
	}

	outputPin := gpio.NewPin(2)
	outputPin.Output()

	go func() {
		for {
			outputPin.Toggle()
			time.Sleep(1 * time.Second)
		}
	}()

	inputPin := gpio.NewPin(6)
	inputPin.Input()
	inputPin.Watch(gpio.EdgeBoth, handler)

	http.HandleFunc("/state", readDoorState)
	http.ListenAndServe(*addr, nil)
}
