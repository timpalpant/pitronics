// clawd is the main entrypoint for the claw server.
// Run it with: `sudo ./clawd -port 8081 -logtostderr`.

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"pitronics/claw"
	"pitronics/claw/motor"
	"pitronics/clawserver"

	"github.com/golang/glog"
	"github.com/stianeikeland/go-rpio"
	"google.golang.org/grpc"
)

const (
	m1DirectionPinNumber = 1
	m1OnOffPinNumber = 2
	m2DirectionPinNumber = 3
	m2OnOffPinNumber = 4
	m3DirectionPinNumber = 5
	m3OnOffPinNumber = 6
	m4DirectionPinNumber = 7
	m4OnOffPinNumber = 8
	m5DirectionPinNumber = 9
	m5OnOffPinNumber = 10
	ledOnOffPinNumber = 11
)

func initClaw() *claw.Claw {
	err := rpio.Open()
	if err != nil {
		panic(err)
	}

	m1DirectionPin := rpio.Pin(m1DirectionPinNumber)
	m1OnOffPin := rpio.Pin(m1OnOffPinNumber)
	m2DirectionPin := rpio.Pin(m2DirectionPinNumber)
	m2OnOffPin := rpio.Pin(m2OnOffPinNumber)
	m3DirectionPin := rpio.Pin(m3DirectionPinNumber)
	m3OnOffPin := rpio.Pin(m3OnOffPinNumber)
	m4DirectionPin := rpio.Pin(m4DirectionPinNumber)
	m4OnOffPin := rpio.Pin(m4OnOffPinNumber)
	m5DirectionPin := rpio.Pin(m5DirectionPinNumber)
	m5OnOffPin := rpio.Pin(m5OnOffPinNumber)
	ledPin := rpio.Pin(ledOnOffPinNumber)

	return &claw.Claw{
		Motors: []motor.Motor{
			motor.NewMotor(leftRightPin, leftRightOnOffPin),
			motor.NewMotor(upDownPin, upDownOnOffPin),
		},
		LED: ledPin,
	}
}

func main() {
	port := flag.Int("port", 8081, "Port to run gRPC service on")
	flag.Parse()

	glog.Info("Initializing claw")
	robot := initClaw()
	defer rpio.Close()
	defer robot.Close()

	endpoint := fmt.Sprintf(":%d", *port)
	lis, err := net.Listen("tcp", endpoint)
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}

	glog.Info("Initializing RPC server")
	grpcServer := grpc.NewServer()
	clawServer := &clawserver.ClawServer{robot}
	clawserver.RegisterClawServiceServer(grpcServer, clawServer)
	go grpcServer.Serve(lis)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	glog.Info(<-ch)
	glog.Info("Shutting down")

	grpcServer.Stop()
}
