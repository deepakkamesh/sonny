package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/chip"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/httphandler"
	"github.com/deepakkamesh/sonny/rpc"
	pb "github.com/deepakkamesh/sonny/sonny"
	"github.com/golang/glog"
	"google.golang.org/grpc"
)

var (
	buildtime string // Compiler Flags
	githash   string // Compiler Flags
)

func main() {

	var (
		brc          = flag.String("brc", "LCD-D22", "GPIO port for roomba BRC for keepalive")
		picAddr      = flag.Int("pic_addr", 0x07, "I2C address of PIC controller")
		tty          = flag.String("tty", "/dev/ttyS0", "tty port")
		res          = flag.String("resources", "./resources", "resources directory")
		pirPin       = flag.String("pir_pin", "LCD-D21", "PIR gpio pin")
		lidarI2CBus  = flag.Int("lidar_i2c_bus", 1, "I2C bus Lidar")
		magI2CBus    = flag.Int("mag_i2c_bus", 2, "I2C bus magnetometer")
		picI2CBus    = flag.Int("pic_i2c_bus", 2, "I2C bus pic")
		rpcPort      = flag.String("rpc_port", ":2233", "host:port number for rpc")
		httpHostPort = flag.String("http_port", ":8080", "host:port number for http")

		enCompass = flag.Bool("en_compass", false, "Enable Compass")
		enRoomba  = flag.Bool("en_roomba", false, "Enable Roomba")
		enPic     = flag.Bool("en_pic", false, "Enable PIC")
		enPir     = flag.Bool("en_pir", false, "Enable PIR")
		enLidar   = flag.Bool("en_lidar", false, "Enable Lidar")
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
	// Log Flush.
	go func() {
		for {
			glog.Flush()
			time.Sleep(300 * time.Millisecond)
		}
	}()

	// Initialize CHIP Adaptor.
	ch := chip.NewAdaptor()
	if err := ch.Connect(); err != nil {
		glog.Fatalf("Failed to initialize CHIP:%v", err)
	}

	// Initialize Roomba.
	var rb *roomba.Roomba
	if *enRoomba {
		glog.V(1).Infof("Enabling Roomba...")
		var err error
		if rb, err = roomba.MakeRoomba(*tty, *brc); err != nil {
			glog.Fatalf("Failed to initialize roomba: %v", err)
		}
		if err = rb.Start(true); err != nil {
			glog.Fatalf("Failed to start roomba: %v", err)
		}
		rb.Safe()

		// Power up secondary battery on main brush.
		time.Sleep(100 * time.Millisecond) // Not sure why, but a little time is needed.
		if err := rb.MainBrush(true, true); err != nil {
			glog.Fatalf("Failed to turn on main brush: %v ")
		}
	}

	// Initialize PIC I2C Controller.
	var ctrl *devices.Controller
	if *enPic {
		ctrl = devices.NewController(ch,
			i2c.WithBus(*picI2CBus),
			i2c.WithAddress(*picAddr))
		if err := ctrl.Start(); err != nil {
			glog.Fatalf("Failed to initialize controller:%v")
		}
	}

	// Initialize magnetometer.
	// TODO: Need driver for hmc5883L in gobot.
	var mag *i2c.HMC6352Driver
	if *enCompass {
		_ = magI2CBus
	}

	// Initialize Lidar.
	var lidar *i2c.LIDARLiteDriver
	if *enLidar {
		lidar = i2c.NewLIDARLiteDriver(ch,
			i2c.WithBus(*lidarI2CBus))
		if err := lidar.Start(); err != nil {
			glog.Fatalf("Failed to initialize lidar: %v")
		}
	}

	// Initialize PIR sensor.
	var pirVal int
	if *enPir {
		pir := gpio.NewPIRMotionDriver(ch, *pirPin)
		if err := pir.Start(); err != nil {
			glog.Fatalf("Unable to initialize pin %v error %v", *pirPin, err)
		}
		pirCh := pir.Subscribe()
		go func() {
			for {
				evt := <-pirCh
				pirVal = evt.Data.(int)
				glog.V(3).Infof("Got pir data %v %v", evt.Name, evt.Data.(int))
			}
		}()
	}

	// Build device list.
	dev := &rpc.Devices{
		Ctrl:   ctrl,
		Mag:    mag,
		Pir:    &pirVal,
		Roomba: rb,
		Lidar:  lidar,
	}

	// Startup RPC service.
	lis, err := net.Listen("tcp", *rpcPort)
	if err != nil {
		glog.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDevicesRPCServer(s, rpc.New(dev))
	go s.Serve(lis)

	// Startup HTTP service.
	h := httphandler.New(dev, false, *res)
	if err := h.Start(*httpHostPort); err != nil {
		glog.Fatalf("Failed to listen: %v", err)
	}

}
