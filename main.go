package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/httphandler"
	"github.com/deepakkamesh/sonny/rpc"
	pb "github.com/deepakkamesh/sonny/sonny"
	"github.com/golang/glog"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/hcsr501"
	"github.com/kidoman/embd/sensor/hmc5883l"
	"google.golang.org/grpc"
)

var (
	buildtime string // Compiler Flags
	githash   string // Compiler Flags
)

func main() {

	var (
		brc     = flag.String("brc", "LCD-D22", "GPIO port for roomba BRC for keepalive")
		picAddr = flag.Int("pic_addr", 0x07, "I2C address of PIC controller")
		tty     = flag.String("tty", "/dev/ttyS0", "tty port")
		res     = flag.String("resources", "./resources", "resources directory")
		pirPin  = flag.String("pir_pin", "CSID0", "PIR gpio pin")
		I2CBus  = flag.Int("i2c_bus", 1, "I2C bus")

		enCompass = flag.Bool("en_compass", false, "Enable Compass")
		enRoomba  = flag.Bool("en_roomba", false, "Enable Roomba")
		enPic     = flag.Bool("en_pic", false, "Enable PIC")
		enIO      = flag.Bool("en_io", false, "Enable CHIP IO (GPIO, I2C). Enable before other pic, pir etc")
		enPir     = flag.Bool("en_pir", false, "Enable PIR")
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

	var i2c embd.I2CBus
	if *enIO {
		// Initialize GPIO.
		if err := embd.InitGPIO(); err != nil {
			glog.Fatalf("Failed to initialize GPIO %v", err)
		}

		// Initialize I2C.
		if err := embd.InitI2C(); err != nil {
			panic(err)
		}
		defer embd.CloseI2C()
		i2c = embd.NewI2CBus(byte(*I2CBus))
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
	}
	// Power up secondary battery on main brush.
	time.Sleep(100 * time.Millisecond) // Not sure why, but a little time is needed.
	if err := rb.MainBrush(true, true); err != nil {
		glog.Fatalf("Failed to turn on main brush: %v ")
	}

	// Initialize PIC Controller.
	var ctrl *devices.Controller
	if *enPic && *enIO {
		ctrl = devices.NewController(i2c, byte(*picAddr))
	}

	// Initialize magnetometer.
	var mag *hmc5883l.HMC5883L
	if *enCompass && *enIO {
		mag = hmc5883l.New(i2c)
		if err := mag.Run(); err != nil {
			glog.Fatalf("Failed to start magnetometer %v", err)
		}
	}

	// Initialize PIR sensor.
	var pir *hcsr501.HCSR501
	if *enPir && *enIO {
		gpio, err := embd.NewDigitalPin(*pirPin)
		if err != nil {
			glog.Fatalf("Unable to initialize pin %v error %v", *pirPin, err)
		}
		pir = hcsr501.New(gpio)
	}

	// Build device list.
	dev := &rpc.Devices{
		Ctrl:   ctrl,
		Mag:    mag,
		Pir:    pir,
		Roomba: rb,
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

}
