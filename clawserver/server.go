package clawserver

import (
	"io"

	"pitronics/claw"
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
	defer s.Claw.Close()

	resp := &SetClawStateResponse{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(resp)
		}
		if err != nil {
			return err
		}

		s.Claw.SetState(claw.State(req.MotorStates))
	}
}
