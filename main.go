package main

import (
	"fmt"
	"github.com/warthog618/gpio"
	"sync"
)

func handler(pin *gpio.Pin) {
	res := pin.Read()
	fmt.Println(res)
}

func main() {
	err := gpio.Open()
	if err != nil {
		panic(err)
	}

	pin := gpio.NewPin(15)
	pin.Input()
	pin.PullDown()
	pin.Watch(gpio.EdgeBoth, handler)

	handler(pin)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
