package claw

import (
	"pitronics/claw/motor"
)

// State represents the motor activation state of the Claw.
// Each element corresponds to a Claw motor.
type State []motor.State

type Claw []*motor.Motor

// Close stops all motors in the Claw.
func (c Claw) Close() {
	for _, m := range c {
		m.SetState(motor.Stopped)
	}
}

// SetState sets the desired state for the Claw motors.
func (c Claw) SetState(state State) {
	for i, motorState := range state {
		motor := c[i]
		motor.SetState(motorState)
	}
}
