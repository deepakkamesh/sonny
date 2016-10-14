package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/deepakkamesh/sonny/rpc"
	pb "github.com/deepakkamesh/sonny/sonny"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/chip"
	"github.com/kidoman/embd/sensor/lsm303"
	"google.golang.org/grpc"
)

func main() {

	// Initialize PIC Controller.
	/*	ctrl, err := devices.NewController("/dev/ttyS0", 115200)
		if err != nil {
			log.Fatalf("Error creating new controller %v", err)
		}
		ctrl.Start()
	*/
	// Initialize magnetometer .
	bus := embd.NewI2CBus(2)
	mag := lsm303.New(bus)
	if err := mag.Run(); err != nil {
		log.Fatalf("Failed to start magnetometer %v", err)
	}
	h, e := mag.Heading()
	if e != nil {
		log.Fatalf("Got e %v", e)
	}
	fmt.Printf("Heading %v\n", h)

	// Build device list.
	dev := &rpc.Devices{
		//	Ctrl: ctrl,
		Mag: mag,
	}
	// Startup RPC service.
	lis, err := net.Listen("tcp", ":2233")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDevicesRPCServer(s, rpc.New(dev))
	s.Serve(lis)

	// Inialize Pi
	/*	for {
		bus := embd.NewI2CBus(2)
		mag := lsm303.New(bus)

		h, e := mag.Heading()
		if e != nil {
			fmt.Printf("Got e %v", e)
			continue
		}
		fmt.Printf("Heading %v\n", h)
		time.Sleep(time.Millisecond * 500)
	} */
	/*
		/*
					log.Printf("Ping to controller failed: %v", err)
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
