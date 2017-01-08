package claw

import (
	"pitronics/claw/motor"

	"github.com/stianeikeland/go-rpio"
)

// State represents the motor activation state of the Claw.
// Each element corresponds to a Claw motor.
type State struct {
	Motors []motor.State
	LED bool
}

type Claw struct {
	Motors []*motor.Motor
	LED rpio.Pin
}

// Close stops all motors in the Claw.
func (c Claw) Close() {
	for _, m := range c.Motors {
		m.SetState(motor.State_STOPPED)
	}
	c.LED.High()
}

// SetState sets the desired state for the Claw motors.
func (c Claw) SetState(state State) {
	for i, motorState := range state.Motors {
		motor := c.Motors[i]
		motor.SetState(motorState)
	}

	if state.LED {
		c.LED.Low()
	} else {
		c.LED.High()
	}
}
