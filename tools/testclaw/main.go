package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const (
	pinOnOffNumber = 4
	pinDirectionNumber = 17
)

func main() {
	err := rpio.Open()
	if err != nil {
		fmt.Printf("Error opening gpio! %s\n", err)
		return
	}

	pinOnOff := rpio.Pin(pinOnOffNumber)
	pinOnOff.Output()
	pinDirection := rpio.Pin(pinDirectionNumber)
	pinDirection.Output()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		pinOnOff.High()
		rpio.Close()
		os.Exit(0)
	}()

	for {
		pinOnOff.Low()
		time.Sleep(time.Second)
		pinOnOff.High()
		pinDirection.Toggle()
		time.Sleep(2 * time.Second)
	}
}
