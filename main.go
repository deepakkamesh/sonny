package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/httphandler"
	"github.com/deepakkamesh/sonny/rpc"
	pb "github.com/deepakkamesh/sonny/sonny"
	"github.com/golang/glog"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/chip"
	"github.com/kidoman/embd/sensor/hcsr501"
	"github.com/kidoman/embd/sensor/hmc5883l"
	"google.golang.org/grpc"
)

var (
	buildtime string // Compiler Flags
	githash   string // Compiler Flags
)

func main() {

	// Setup Flags.
	var (
		baud      = flag.Int("baud", 115200, "TTY Baud rate")
		tty       = flag.String("tty", "/dev/ttyS0", "tty port")
		res       = flag.String("resources", "./resources", "resources directory")
		pirPin    = flag.String("pir_pin", "132", "PIR gpio pin")
		magBus    = flag.Int("mag_bus", 2, "I2C bus for Compass")
		enCompass = flag.Bool("en_compass", false, "Enable Compass")
		enPic     = flag.Bool("en_pic", false, "Enable PIC")
		enPir     = flag.Bool("en_pir", true, "Enable PIR")
		version   = flag.Bool("version", false, "display version")
	)
	flag.Parse()

	// Print version and exit.
	if *version {
		fmt.Printf("Version commit hash %s\n", githash)
		fmt.Printf("Build date %s\n", buildtime)
		os.Exit(0)
	}

	glog.Infof("Starting Sonny ver %s build on %s", githash, buildtime)
	defer glog.Flush()

	// Initialize PIC Controller.
	var ctrl *devices.Controller
	var err error
	if *enPic {
		ctrl, err = devices.NewController(*tty, *baud)
		if err != nil {
			glog.Fatalf("Error creating new controller %v", err)

		}
		ctrl.Start()
	}

	// Initialize I2C.
	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	// Initialize GPIO.
	if err := embd.InitGPIO(); err != nil {
		glog.Fatalf("Failed to initialize GPIO %v", err)
	}

	// Initialize magnetometer.
	var mag *hmc5883l.HMC5883L
	if *enCompass {
		mag = hmc5883l.New(embd.NewI2CBus(byte(*magBus)))
		if err := mag.Run(); err != nil {
			glog.Fatalf("Failed to start magnetometer %v", err)
		}
	}

	// Initialize PIR sensor.
	var pir *hcsr501.HCSR501
	if *enPir {
		gpio, err := embd.NewDigitalPin(*pirPin)
		if err != nil {
			glog.Fatalf("Unable to initialize pin %v error %v", *pirPin, err)
		}
		pir = hcsr501.New(gpio)
	}

	// Build device list.
	dev := &rpc.Devices{
		Ctrl: ctrl,
		Mag:  mag,
		Pir:  pir,
	}

	// Startup RPC service.
	lis, err := net.Listen("tcp", ":2233")
	if err != nil {
		glog.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDevicesRPCServer(s, rpc.New(dev))
	go s.Serve(lis)

	// Startup HTTP service.
	h := httphandler.New(dev, false, *res)
	if err := h.Start(); err != nil {
		glog.Fatalf("Failed to listen: %v", err)
	}

	for {
		time.Sleep(time.Millisecond * 20)
	}
}
