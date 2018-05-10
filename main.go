package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/httphandler"
	"github.com/deepakkamesh/sonny/navigator"
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
		brc          = flag.String("brc", "7", "GPIO port for roomba BRC for keepalive")
		picAddr      = flag.Int("pic_addr", 0x08, "I2C address of PIC controller")
		tty          = flag.String("tty", "/dev/ttyS0", "tty port")
		res          = flag.String("resources", "./resources", "resources directory")
		pirPin       = flag.String("pir_pin", "22", "PIR gpio pin")
		enI2CPin     = flag.String("i2c_en_pin", "11", "I2C enable pin (high to enable I2C chip)")
		enLidarPin   = flag.String("lidar_en_pin", "24", "LIDAR power enable pin (high to enable lidar)")
		lidarI2CBus  = flag.Int("lidar_i2c_bus", 1, "I2C bus Lidar")
		magI2CBus    = flag.Int("mag_i2c_bus", 3, "I2C bus magnetometer")
		picI2CBus    = flag.Int("pic_i2c_bus", 3, "I2C bus pic")
		rpcPort      = flag.String("rpc_port", ":2233", "host:port number for rpc")
		httpHostPort = flag.String("http_port", ":8080", "host:port number for http")

		roombaMode = flag.Uint("roomba_mode", 1, "0=Off 1=Passive 2=Safe 3=Full")
		version    = flag.Bool("version", false, "display version")

		vidHeight = flag.Uint("vid_height", 480, "Video Height")
		vidWidth  = flag.Uint("vid_width", 640, "Video Width")

		enCompass  = flag.Bool("en_compass", false, "Enable Compass")
		enRoomba   = flag.Bool("en_roomba", false, "Enable Roomba")
		enPic      = flag.Bool("en_pic", false, "Enable PIC")
		enPir      = flag.Bool("en_pir", false, "Enable PIR")
		enLidar    = flag.Bool("en_lidar", false, "Enable Lidar")
		enI2C      = flag.Bool("en_i2c", false, "Enable I2C Connect")
		enAuxPower = flag.Bool("en_aux_power", false, "Enable Auxillary Power")
		enVid      = flag.Bool("en_video", true, "Enable video")
	)
	flag.Parse()

	// Print version and exit.
	if *version {
		fmt.Printf("Version commit hash %s\n", githash)
		fmt.Printf("Build date %s\n", buildtime)
		os.Exit(0)
	}

	glog.Infof("Starting Sonny ver %s build on %s", githash, buildtime)

	// Log flush Routine.
	go func() {
		for {
			glog.Flush()
			time.Sleep(300 * time.Millisecond)
		}
	}()

	// Initialize PI Adaptor.
	pi := raspi.NewAdaptor()
	if err := pi.Connect(); err != nil {
		glog.Fatalf("Failed to initialize Adapter:%v", err)
	}

	// Initialize Roomba.
	var rb *roomba.Roomba
	if *enRoomba {
		glog.V(1).Infof("Enabling Roomba...")

		// Setup BRC pin for roomba keep-alive.
		brcPin := gpio.NewDirectPinDriver(pi, *brc)
		if err := brcPin.Start(); err != nil {
			glog.Fatalf("Failed to setup BRC pin: %v", err)
		}
		var err error
		if rb, err = roomba.MakeRoomba(*tty, brcPin); err != nil {
			glog.Fatalf("Failed to initialize roomba: %v", err)
		}
		if err = rb.Start(true); err != nil {
			glog.Fatalf("Failed to start roomba: %v", err)
		}
	}

	// I2C enable control.
	i2cEn := gpio.NewDirectPinDriver(pi, *enI2CPin)
	if err := i2cEn.Start(); err != nil {
		glog.Fatalf("Failed to initialize I2C en pin: %v", err)
	}

	// Initialize PIC I2C Controller.
	var ctrl *devices.Controller
	if *enPic {
		ctrl = devices.NewController(pi,
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

	// Initialize Lidar and related systems.
	var (
		lidar      *i2c.LIDARLiteDriver
		lidarEnPin *gpio.DirectPinDriver
	)
	if *enLidar {
		lidar = i2c.NewLIDARLiteDriver(pi,
			i2c.WithBus(*lidarI2CBus))
		if err := lidar.Start(); err != nil {
			glog.Fatalf("Failed to initialize lidar: %v", err)
		}
		// Lidar enable control. Needed so we bring devices online
		// in a phased manner as to not overload the power.
		// Pull high to disable.
		lidarEnPin = gpio.NewDirectPinDriver(pi, *enLidarPin)
		if err := lidarEnPin.Start(); err != nil {
			glog.Fatalf("Failed to initialize Lidar enable pin: %v", err)
		}
		lidarEnPin.DigitalWrite(0)
	}

	// Initialize PIR sensor.
	var pir *gpio.PIRMotionDriver
	if *enPir {
		pir = gpio.NewPIRMotionDriver(pi, *pirPin)
		if err := pir.Start(); err != nil {
			glog.Fatalf("Unable to initialize pin %v error %v", *pirPin, err)
		}
	}

	// Initialize video device.
	var vid *devices.Video
	if *enVid {
		vid = devices.NewVideo(devices.MJPEG, uint32(*vidWidth), uint32(*vidHeight), 2)
		vid.StartVideoStream()
	}

	// Build Devices.
	sonny := devices.NewSonny(ctrl, lidar, mag, rb, i2cEn, pir, lidarEnPin, vid)

	// Enable I2C Bus if flag is set.
	// Explicit disable is needed as the gpio may be high from prior run.
	if err := sonny.I2CBusEnable(false); err != nil {
		glog.Fatalf("Failed to disable I2C Bus")
	}
	if *enI2C {
		if err := sonny.I2CBusEnable(true); err != nil {
			glog.Fatalf("Failed to enable I2C Bus")
		}
	}

	sonny.PIREventLoop()

	// Easier to set roomba mode once the sonny struct is ready.
	if *enRoomba {
		if err := sonny.SetRoombaMode(byte(*roombaMode)); err != nil {
			glog.Fatalf("Failed to set roomba mode:%v", err)
		}

		// Power up auxillary battery on main brush.
		time.Sleep(100 * time.Millisecond) // Not sure why, but a little time is needed.
		if err := sonny.AuxPower(*enAuxPower); err != nil {
			glog.Fatalf("Failed to turn on Aux Power: %v ", err)
		}
	}

	// Power up sequence complete
	glog.Info("Sonny device initialization complete")

	// Start up navigation routines.
	navi := navigator.NewOgrid(sonny)

	// Startup RPC service.
	lis, err := net.Listen("tcp", *rpcPort)
	if err != nil {
		glog.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDevicesRPCServer(s, rpc.New(sonny))
	go s.Serve(lis)

	// Startup HTTP service.
	h := httphandler.New(sonny, navi, false, *res)
	if err := h.Start(*httpHostPort); err != nil {
		glog.Fatalf("Failed to listen: %v", err)
	}

}
