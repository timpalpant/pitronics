package motor

import (
	"fmt"

	"github.com/stianeikeland/go-rpio"
)

// Motor represents a motor that can be driven forward or backward
// by controlling two relays. The direction relay determines the
// direction the motor is driven, and the onOff relay determines
// whether the motor is running or stopped.
type Motor struct {
	direction rpio.Pin
	onOff     rpio.Pin
}

func NewMotor(direction, onOff rpio.Pin) *Motor {
	direction.Mode(rpio.Output)
	onOff.Mode(rpio.Output)

	return &Motor{
		direction: direction,
		onOff:     onOff,
	}
}

func (m *Motor) SetState(s State) error {
	switch s {
	case State_FORWARD:
		m.direction.High()
		m.onOff.Low()
	case State_BACKWARD:
		m.direction.Low()
		m.onOff.Low()
	case State_STOPPED:
		m.onOff.High()
	default:
		return fmt.Errorf("unknown motor state: %v", s)
	}

	return nil
}
