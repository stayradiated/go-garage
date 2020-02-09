package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/warthog618/gpio"
)

var (
	addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
)

func handler(pin *gpio.Pin) {
	level := pin.Read()
	if level == gpio.Low {
		log.Println("low")
	} else {
		log.Println("high")
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

	http.ListenAndServe(*addr, nil)
}
