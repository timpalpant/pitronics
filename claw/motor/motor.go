package motor

import (
	"github.com/stianeikeland/go-rpio"
)

type State int

const (
	Forward State = iota
	Backward
	Stopped
)

// Motor represents a motor that can be driven forward or backward
// by controlling two relays. The direction relay determines the
// direction the motor is driven, and the onOff relay determines
// whether the motor is running or stopped.
type Motor struct {
	direction rpio.Pin
	onOff rpio.Pin
}

func NewMotor(direction, onOff rpio.Pin) *Motor {
	direction.Mode(rpio.Output)
	onOff.Mode(rpio.Output)

	return &Motor{
		direction: direction,
		onOff: onOff,
	}
}

func (m *Motor) SetState(s State) {
	switch s {
	case Forward:
		m.direction.High()
		m.onOff.Low()
	case Backward:
		m.direction.Low()
		m.onOff.Low()
	case Stopped:
		m.onOff.High()
	}
}
