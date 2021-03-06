// kbcontroller is a claw server client that listens for keyboard events
// and sends them to the claw server.
// Run it with: ./kbcontroller -logtostderr -server :8081
//
// Use the arrow keys to control the claw state.

package main

import (
	"flag"

	"github.com/gdamore/tcell"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"pitronics/claw/motor"
	"pitronics/clawserver"
)

// parseKBState converts the given keyboard controller State into
// the desired motor state for the Claw.
func parseKBState(event *tcell.EventKey) *clawserver.SetClawStateRequest {
	motorStates := make([]motor.State, 5)
	for i := range motorStates {
		motorStates[i] = motor.State_STOPPED
	}
	led := false

	switch event.Key() {
	case tcell.KeyUp:
		motorStates[1] = motor.State_FORWARD
	case tcell.KeyDown:
		motorStates[1] = motor.State_BACKWARD
	case tcell.KeyRight:
		motorStates[0] = motor.State_BACKWARD
	case tcell.KeyLeft:
		motorStates[0] = motor.State_FORWARD
	case tcell.KeyRune:
		switch event.Rune() {
		case 'a':
			motorStates[2] = motor.State_FORWARD
		case 'd':
			motorStates[2] = motor.State_BACKWARD
		case 'w':
			motorStates[3] = motor.State_FORWARD
		case 's':
			motorStates[3] = motor.State_BACKWARD
		case 'q':
			motorStates[4] = motor.State_FORWARD
		case 'e':
			motorStates[4] = motor.State_BACKWARD
		case 'f':
			led = true
		case 'r':
			led = false
		}
	}

	return &clawserver.SetClawStateRequest{
		MotorStates: motorStates,
		Led: led,
	}
}

func main() {
	clawServerConnStr := flag.String("server", "", "ClawServer connection string")
	flag.Parse()

	if *clawServerConnStr == "" {
		glog.Error("You must provide the connection string with -server")
		return
	}

	glog.Infof("Connecting to ClawServer at %v", *clawServerConnStr)
	// TODO(palpant): Security would be good.
	conn, err := grpc.Dial(*clawServerConnStr, grpc.WithInsecure())
	if err != nil {
		glog.Fatal(err)
	}
	defer conn.Close()
	client := clawserver.NewClawServiceClient(conn)

	glog.Info("Initializing tcell screen")
	screen, err := tcell.NewScreen()
	if err != nil {
		glog.Fatal(err)
	}
	if err := screen.Init(); err != nil {
		glog.Fatal(err)
	}
	defer screen.Fini()

	glog.Info("Opening streaming RPC to claw server")
	stream, err := client.SetClawState(context.Background())
	if err != nil {
		glog.Fatal(err)
	}

	glog.Info("Listening for keyboard input")
Loop:
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				break Loop
			default:
				glog.V(1).Infof("Key pressed: %v", ev)
				clawStateReq := parseKBState(ev)
				if err := stream.Send(clawStateReq); err != nil {
					glog.Fatal(err)
				}
			}
		}
	}

	glog.Info("Closing stream to claw server")
	if _, err := stream.CloseAndRecv(); err != nil {
		glog.Fatal(err)
	}
}
