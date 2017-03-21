package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/httphandler"
	"github.com/deepakkamesh/sonny/rpc"
	pb "github.com/deepakkamesh/sonny/sonny"
	"github.com/golang/glog"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/chip"
	"github.com/kidoman/embd/sensor/lsm303"
	"github.com/kidoman/embd/sensor/us020"
	"google.golang.org/grpc"
)

func main() {

	// Setup Flags.
	var (
		tty       = flag.String("tty", "/dev/ttyS0", "tty port")
		pirPin    = flag.String("pir_pin", "gpio0", "PIR gpio pin")
		usTrigPin = flag.String("us_trig_pin", "gpio3", "Ultrasonic Trigger Pin")
		usEchoPin = flag.String("us_echo_pin", "gpio1", "Ultrasonic Echo Pin")
		magBus    = flag.Int("mag_bus", 2, "I2C bus for Compass")
		enCompass = flag.Bool("en_compass", false, "Enable Compass")
		enPic     = flag.Bool("en_pic", false, "Enable PIC")
		enPir     = flag.Bool("en_pir", true, "Enable PIR")
		enUS      = flag.Bool("en_us", true, "Enable UltraSonic Sensor")
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
			glog.Fatalf("Failed to initialize GPIO %v", err)
		}
		embd.SetDirection(*pirPin, embd.In)
	}

	// Initialize Ultrasonic sensor.
	var us *us020.US020
	if *enUS {
		echo, err := embd.NewDigitalPin(*usEchoPin)
		if err != nil {
			glog.Fatalf("Failed to init digital pin %v", err)
		}
		trig, err := embd.NewDigitalPin(*usTrigPin)
		if err != nil {
			glog.Fatalf("Failed to init digital pin %v", err)
		}
		us = us020.New(echo, trig, nil)
		//defer us.Close()
	}

	// Build device list.
	dev := &rpc.Devices{
		Ctrl: ctrl,
		Mag:  mag,
		Pir:  *pirPin,
		Us:   us,
	}

	// Startup RPC service.
	lis, err := net.Listen("tcp", ":2233")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDevicesRPCServer(s, rpc.New(dev))
	go s.Serve(lis)

	// Startup HTTP service.
	h := httphandler.New(dev, false)
	if err := h.Start(); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	for {
		time.Sleep(time.Millisecond * 20)
	}
}
