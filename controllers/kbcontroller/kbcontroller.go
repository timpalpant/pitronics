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
func parseKBState(key tcell.Key) *clawserver.SetClawStateRequest {
	m1 := motor.State_STOPPED
	m2 := motor.State_STOPPED
	switch key {
	case tcell.KeyUp:
		m2 = motor.State_FORWARD
	case tcell.KeyDown:
		m2 = motor.State_BACKWARD
	case tcell.KeyRight:
		m1 = motor.State_BACKWARD
	case tcell.KeyLeft:
		m1 = motor.State_FORWARD
	}

	return &clawserver.SetClawStateRequest{
		MotorStates: []motor.State{m1, m2},
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
				glog.V(1).Infof("Key pressed: %v", ev.Key())
				clawStateReq := parseKBState(ev.Key())
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
