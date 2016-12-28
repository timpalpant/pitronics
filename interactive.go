package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/stianeikeland/go-rpio"
	"github.com/timpalpant/joystick"

	"pitronics/claw"
	"pitronics/claw/motor"
)

const (
	jsid = 0
	leftRightPinNumber = 3
	leftRightOnOffPinNumber = 2
	upDownPinNumber = 17
	upDownOnOffPinNumber = 4
)

// parseJSState converts the given joystick controller State into
// the desired motor state for the Claw.
func parseJSState(jsState joystick.State) claw.State {
	m1 := motor.Stopped
	leftRightAxis := jsState.AxisData[0]
	if leftRightAxis < 0 {
		m1 = motor.Forward
	} else if leftRightAxis > 0 {
		m1 = motor.Backward
	}

	m2 := motor.Stopped
	upDownAxis := jsState.AxisData[1]
	if upDownAxis > 0 {
		m2 = motor.Forward
	} else if upDownAxis < 0 {
		m2 = motor.Backward
	}

	return claw.State{m1, m2}
}

func main() {
	js, err := joystick.Open(jsid)
	if err != nil {
		panic(err)
	}

	err = rpio.Open()
	if err != nil {
		panic(err)
	}
	defer rpio.Close()

	leftRightPin := rpio.Pin(leftRightPinNumber)
	leftRightOnOffPin := rpio.Pin(leftRightOnOffPinNumber)
	upDownPin := rpio.Pin(upDownPinNumber)
	upDownOnOffPin := rpio.Pin(upDownOnOffPinNumber)

	robot := &claw.Claw{
		motor.NewMotor(leftRightPin, leftRightOnOffPin),
		motor.NewMotor(upDownPin, upDownOnOffPin),
	}
	defer robot.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	jsEvents := js.Events()
	for {
		select {
		case <-c:
			return
		case <-jsEvents:
			jsState, err := js.Read()
			if err != nil {
				fmt.Printf("error reading JS state: %v\n", err)
				return
			}
			clawState := parseJSState(jsState)
			robot.SetState(clawState)
		}
	}
}
