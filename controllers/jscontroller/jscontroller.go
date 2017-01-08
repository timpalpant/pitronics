// jscontroller is a claw server client that listens for joystick events
// and sends them to the claw server.
// Run it with: `sudo ./jscontroller -server :8081`.
// 
// Use the left joystick axis to control the claw motors.

package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/golang/glog"
	"github.com/timpalpant/joystick"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"pitronics/claw/motor"
	"pitronics/clawserver"
)

const axisThreshold = 10

// parseJSState converts the given joystick controller State into
// the desired motor state for the Claw.
func parseJSState(jsState joystick.State) *clawserver.SetClawStateRequest {
	m1 := motor.State_STOPPED
	leftRightAxis := jsState.AxisData[0]
	if leftRightAxis < -axisThreshold {
		m1 = motor.State_FORWARD
	} else if leftRightAxis > axisThreshold {
		m1 = motor.State_BACKWARD
	}

	m2 := motor.State_STOPPED
	upDownAxis := jsState.AxisData[1]
	if upDownAxis > axisThreshold {
		m2 = motor.State_FORWARD
	} else if upDownAxis < -axisThreshold {
		m2 = motor.State_BACKWARD
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
	defer func() {
		if _, err := stream.CloseAndRecv(); err != nil {
			glog.Fatal(err)
		}
	}()

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
				glog.Fatal("error reading JS state: %v\n", err)
				return
			}
			clawStateReq := parseJSState(jsState)
			if err := stream.Send(clawStateReq); err != nil {
				glog.Fatal(err)
			}
		}
	}
}
