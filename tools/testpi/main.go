package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var pinNumbers = [...]int{2, 3, 4, 17}

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Printf("Error opening pin! %s\n", err)
		return
	}

	pins := make([]rpio.Pin, len(pinNumbers))
	for i, pinNum := range pinNumbers {
		pins[i] = rpio.Pin(pinNum)
		pins[i].Output()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		for _, pin := range pins {
			pin.High()
		}
		rpio.Close()
		os.Exit(0)
	}()

	for {
		for _, pin := range pins {
			pin.Toggle()
			time.Sleep(100 * time.Millisecond)
		}
		time.Sleep(time.Second)
	}
}
