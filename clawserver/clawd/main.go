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
	leftRightPinNumber      = 3
	leftRightOnOffPinNumber = 2
	upDownPinNumber         = 17
	upDownOnOffPinNumber    = 4
)

func initClaw() *claw.Claw {
	err := rpio.Open()
	if err != nil {
		panic(err)
	}

	leftRightPin := rpio.Pin(leftRightPinNumber)
	leftRightOnOffPin := rpio.Pin(leftRightOnOffPinNumber)
	upDownPin := rpio.Pin(upDownPinNumber)
	upDownOnOffPin := rpio.Pin(upDownOnOffPinNumber)

	return &claw.Claw{
		motor.NewMotor(leftRightPin, leftRightOnOffPin),
		motor.NewMotor(upDownPin, upDownOnOffPin),
	}
}

func main() {
	port := flag.Int("port", 9123, "Port to run gRPC service on")
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
