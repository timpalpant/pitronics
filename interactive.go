package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/timpalpant/joystick"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"pitronics/claw/motor"
	"pitronics/clawserver"
)

// parseJSState converts the given joystick controller State into
// the desired motor state for the Claw.
func parseJSState(jsState joystick.State) *clawserver.SetClawStateRequest {
	m1 := clawserver.STOPPED
	leftRightAxis := jsState.AxisData[0]
	if leftRightAxis < 0 {
		m1 = clawserver.FORWARD
	} else if leftRightAxis > 0 {
		m1 = clawserver.BACKWARD
	}

	m2 := clawserver.BACKWARD
	upDownAxis := jsState.AxisData[1]
	if upDownAxis > 0 {
		m2 = clawserver.FORWARD
	} else if upDownAxis < 0 {
		m2 = clawserver.BACKWARD
	}

	return &clawserver.SetClawStateRequest{
		MotorStates: []motor.State{m1, m2},
	}
}

func main() {
	clawServerConnStr := flag.String("server", "", "ClawServer connection string")
	jsid := flag.Int("jsid", 0, "Joystick id (default: 0)")
	flag.Parse()

	glog.Infof("Connecting to ClawServer at %v", *clawServerConnStr)
	// TODO(palpant): Security would be good.
	conn, err := grpc.Dial(*clawServerConnStr, grpc.WithInsecure())
	if err != nil {
		glog.Fatal(err)
	}
	defer conn.Close()
	client := clawserver.NewClawServiceClient(conn)

	glog.Infof("Initializing joystick %d", *jsid)
	js, err := joystick.Open(*jsid)
	if err != nil {
		panic(err)
	}

	jsEvents := js.Events()
	// Open streaming RPC to ClawServer.
	stream, err := client.SetClawState(context.Background())
	if err != nil {
		glog.Fatal(err)
	}

	// On each joystick event, get the joystick state and send it to the
	// ClawServer until we receive SIGTERM.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
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
			clawStateReq := parseJSState(jsState)
			if err := stream.Send(clawStateReq); err != nil {
				glog.Fatal(err)
			}
		}
	}

	_, err := stream.CloseAndRecv()
	if err != nil {
		glog.Fatal(err)
	}
}
