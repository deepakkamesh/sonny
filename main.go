package main

import (
	"fmt"
	"time"

	"github.com/deepakkamesh/sonny/devices"
)

func main() {

	// Inialize Pi
	chip := devices.NewChip()
	chip.BlinkGPIO8()
	/*
		ctrl, err := devices.NewController("/dev/ttyAMA0", 115200)
		if err != nil {
			log.Fatalf("Error creating new controller %v", err)
		}
		ctrl.Start()

		ctrl := &devices.Controller{}
		lis, err := net.Listen("tcp", ":2233")
		if err != nil {
			log.Fatalf("Failed to listen:%v", err)
		}
		s := grpc.NewServer()
		pb.RegisterDevicesRPCServer(s, rpc.New(ctrl))
		s.Serve(lis)
		/*
			if err := ctrl.Ping(); err != nil {
				fmt.Println("Error", err)
			}

			if err := ctrl.LedOn(true); err != nil {
				fmt.Println("Error", err)
			}
			time.Sleep(2 * time.Second)
			if err := ctrl.LedOn(false); err != nil {
				fmt.Println("Error", err)
			}
	*/
	fmt.Println("completed")

	for {
		time.Sleep(time.Millisecond * 20)
	}
}
