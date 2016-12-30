package clawserver

import (
	"io"

	"pitronics/claw"

	"github.com/golang/glog"
)

// ClawServer implements ClawService and can be used to remotely
// manage the state of a Claw.
type ClawServer struct {
	Claw *claw.Claw
}

// SetClawState is a streaming RPC that receives a stream of SetClawStateRequests.
// The Claw state is updated for each request. When the stream is closed, the
// claw state is automatically reset to stopped.
func (s *ClawServer) SetClawState(stream ClawService_SetClawStateServer) error {
	glog.Info("Beginning SetClawState stream")
	defer s.Claw.Close()

	resp := &SetClawStateResponse{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			glog.Info("SetClawState stream closing")
			return stream.SendAndClose(resp)
		}
		if err != nil {
			glog.Errorf("SetClawState stream error: %v", err)
			return err
		}

		clawState := claw.State(req.MotorStates)
		glog.Infof("Setting claw state to: %v", clawState)
		s.Claw.SetState(clawState)
	}
}
