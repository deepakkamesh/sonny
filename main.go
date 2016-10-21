package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/rpc"
	pb "github.com/deepakkamesh/sonny/sonny"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/chip"
	"github.com/kidoman/embd/sensor/lsm303"
	"google.golang.org/grpc"
)

func main() {

	// Setup Flags.
	var (
		tty       = flag.String("tty", "/dev/ttyS0", "tty port")
		pirPin    = flag.String("pir_pin", "gpio0", "PIR gpio pin")
		magBus    = flag.Int("mag_bus", 2, "I2C bus for Compass")
		enCompass = flag.Bool("en_compass", false, "Enable Compass")
		enPic     = flag.Bool("en_pic", false, "Enable PIC")
		enPir     = flag.Bool("en_pir", true, "Enable PIR")
	)
	flag.Parse()
	// Initialize PIC Controller.
	var ctrl *devices.Controller
	var err error
	if *enPic {
		ctrl, err = devices.NewController(*tty, 115200)
		if err != nil {
			log.Fatalf("Error creating new controller %v", err)
		}
		ctrl.Start()
	}
	// Initialize magnetometer.
	var mag *lsm303.LSM303
	if *enCompass {
		mag = lsm303.New(embd.NewI2CBus(byte(*magBus)))
		if err := mag.Run(); err != nil {
			log.Fatalf("Failed to start magnetometer %v", err)
		}
	}

	// Initialize PIR sensor.
	if *enPir {
		if err := embd.InitGPIO(); err != nil {
			log.Fatalf("Failed to initialize GPIO %v", err)
		}
		embd.SetDirection(*pirPin, embd.In)
	}

	// Build device list.
	dev := &rpc.Devices{
		Ctrl: ctrl,
		Mag:  mag,
		Pir:  *pirPin,
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
