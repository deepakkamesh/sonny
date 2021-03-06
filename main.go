package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/httphandler"
	"github.com/deepakkamesh/sonny/ros"
	"github.com/deepakkamesh/sonny/rpc"
	pb "github.com/deepakkamesh/sonny/sonny"
	"github.com/deepakkamesh/ydlidar"
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
		picAddr      = flag.Int("pic_addr", 0x55, "I2C address of PIC controller")
		roombaTTY    = flag.String("tty_roomba", "/dev/ttyS0", "tty port of roomba")
		lidarTTY     = flag.String("tty_lidar", "/dev/ttyUSB0", "tty port of LIDAR")
		res          = flag.String("resources", "./resources", "resources directory")
		pirPin       = flag.String("pir_pin", "22", "PIR gpio pin")
		enI2CPin     = flag.String("i2c_en_pin", "11", "I2C enable pin (high to enable I2C chip)")
		magI2CBus    = flag.Int("mag_i2c_bus", 0, "I2C bus magnetometer")
		gyroI2CBus   = flag.Int("gyro_i2c_bus", 0, "I2C gyro bus")
		picI2CBus    = flag.Int("pic_i2c_bus", 0, "I2C bus pic")
		rpcPort      = flag.String("rpc_port", ":2233", "host:port number for rpc")
		httpHostPort = flag.String("http_port", ":8080", "host:port number for http")

		roombaMode = flag.Uint("roomba_mode", 1, "0=Off 1=Passive 2=Safe 3=Full")
		version    = flag.Bool("version", false, "display version")

		vidHeight = flag.Uint("vid_height", 120, "Video Height")
		vidWidth  = flag.Uint("vid_width", 160, "Video Width")

		// Enable feature flags.
		enCompass    = flag.Bool("en_compass", false, "Enable Compass")
		enGyro       = flag.Bool("en_gyro", false, "Enable Gyro")
		enRoomba     = flag.Bool("en_roomba", false, "Enable Roomba")
		enPic        = flag.Bool("en_pic", false, "Enable PIC")
		enPir        = flag.Bool("en_pir", false, "Enable PIR")
		enLidar      = flag.Bool("en_lidar", false, "Enable Lidar")
		enI2C        = flag.Bool("en_i2c", false, "Enable I2C Connect")
		enAuxPower   = flag.Bool("en_aux_power", false, "Enable Auxillary Power")
		enVid        = flag.Bool("en_video", false, "Enable video")
		enDataStream = flag.Bool("en_data_stream", false, "Enable data stream for http")
		enROS        = flag.Bool("en_ros", false, "Enable ROS connection")
	)
	flag.Parse()

	// Print version and exit.
	if *version {
		fmt.Printf("Version commit hash %s\n", githash)
		fmt.Printf("Build date %s\n", buildtime)
		os.Exit(0)
	}

	glog.Infof("Starting Sonny ver %s build on) %s", githash, buildtime)

	// Log flush Routine.
	go func() {
		for {
			glog.Flush()
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Initialize PI Adaptor.
	pi := raspi.NewAdaptor()
	if err := pi.Connect(); err != nil {
		glog.Fatalf("Failed to initialize Adapter:%v", err)
	}

	// I2C enable control.
	i2cEn := gpio.NewDirectPinDriver(pi, *enI2CPin)
	if err := i2cEn.Start(); err != nil {
		glog.Fatalf("Failed to initialize I2C en pin: %v", err)
	}

	// Initialize Roomba.
	var rb *roomba.Roomba
	if *enRoomba {
		glog.V(1).Infof("Initializing Roomba...")

		// Setup BRC pin for roomba keep-alive.
		brcPin := gpio.NewDirectPinDriver(pi, *brc)
		if err := brcPin.Start(); err != nil {
			glog.Fatalf("Failed to setup BRC pin: %v", err)
		}
		var err error
		if rb, err = roomba.MakeRoomba(*roombaTTY, brcPin); err != nil {
			glog.Fatalf("Failed to initialize roomba: %v", err)
		}
		if err = rb.Start(true); err != nil {
			glog.Fatalf("Failed to start roomba: %v", err)
		}
		// Verify we can talk to roomba by querying mode.
		v, e := rb.QueryList([]byte{35})
		if e != nil {
			glog.Fatalf("Failed to initialize Roomba:%v (got %v)", e, v)
		}
		if v[0][0] != 1 && v[0][0] != 2 && v[0][0] != 3 {
			glog.Fatalf("Bad roomba mode: %v ", v[0][0])
		}
		glog.V(1).Infof("Roomba ready in mode: %d", v[0][0])
	}

	// Initialize PIC I2C Controller.
	var ctrl *devices.Controller
	if *enPic {
		glog.V(1).Infof("Initializing PIC Controller...")
		ctrl = devices.NewController(pi,
			i2c.WithBus(*picI2CBus),
			i2c.WithAddress(*picAddr))
		if err := ctrl.Start(); err != nil {
			glog.Fatalf("Failed to initialize controller:%v", err)
		}
	}

	// Initialize magnetometer.
	var mag *i2c.QMC5883Driver
	if *enCompass {
		glog.V(1).Infof("Initializing Compass...")
		mag = i2c.NewQMC5883Driver(pi, i2c.WithBus(*magI2CBus))
		mag.SetConfig(i2c.QMC5883Continuous | i2c.QMC5883ODR200Hz | i2c.QMC5883RNG2G | i2c.QMC5883OSR512)
	}

	// Initialize MPU6050.
	var gyro *i2c.MPU6050Driver
	if *enGyro {
		glog.V(1).Infof("Initializing Gyro...")
		gyro = i2c.NewMPU6050Driver(pi, i2c.WithBus(*gyroI2CBus))
	}

	// Initialize Lidar and related systems.
	var (
		lidar *ydlidar.YDLidar
	)
	if *enLidar {
		glog.V(1).Infof("Initializing Lidar...")
		lidar = ydlidar.NewLidar()
		ser, err := ydlidar.GetSerialPort(*lidarTTY)
		if err != nil {
			glog.Fatalf("Failed to initialize LIDAR: %v", err)
		}
		lidar.SetSerial(ser)
		if err := lidar.SetDTR(false); err != nil {
			glog.Fatalf("Failed to set DTR during initialization: %v", err)
		}
	}

	// Initialize PIR sensor.
	var pir *gpio.PIRMotionDriver
	if *enPir {
		glog.V(1).Infof("Initializing PIR Sensor...")
		pir = gpio.NewPIRMotionDriver(pi, *pirPin)
		if err := pir.Start(); err != nil {
			glog.Fatalf("Unable to initialize pin %v error %v", *pirPin, err)
		}
	}

	// Initialize video device.
	var vid *devices.Video
	if *enVid {
		glog.V(1).Infof("Initializing Video...")
		vid = devices.NewVideo(devices.MJPEG, uint32(*vidWidth), uint32(*vidHeight), 10)
		vid.StartVideoStream()
	}

	// Build Devices.
	sonny := devices.NewSonny(ctrl, lidar, mag, gyro, rb, i2cEn, pir, vid)

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

	// Setup the any post init routines tied to Aux power.
	sonny.SetAuxPostInit(
		// Function is called whenever aux power is turned on.
		func() error {
			if sonny.GetI2CBusState() == 0 {
				err := fmt.Errorf("I2C bus not enabled.")
				glog.Warningf("%v", err)
				return err
			}
			if *enCompass {
				// Magnetometer needs to be (re)configured after every power up.
				if err := mag.Start(); err != nil {
					glog.Fatalf("Failed to initialize magnetometer:%v", err)
				}
			}
			if *enGyro {
				if err := gyro.Start(); err != nil {
					glog.Fatalf("Failed to initialize Gyro:%v", err)
				}
			}
			return nil
		},
		// Function is called whenever aux power is turned off.
		func() error {
			return nil
		})

	// Easier to set roomba mode once the sonny struct is ready since
	// sonny has a simpler function to set mode.
	if *enRoomba {
		if err := sonny.SetRoombaMode(byte(*roombaMode)); err != nil {
			glog.Fatalf("Failed to set roomba mode:%v", err)
		}

		if err := sonny.AuxPower(*enAuxPower); err != nil {
			glog.Fatalf("Failed to turn on Aux Power: %v ", err)
		}
	}

	// finally startup Sonny.
	sonny.Startup()

	// Power up sequence complete
	glog.Infof("Sonny device initialization complete")
	time.Sleep(500 * time.Millisecond) // Sleep to allow devices to power up.

	// Catch interrupts to exit clean.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case sig := <-c:
			glog.Infof("Got %s signal. Aborting...\n", sig)
			sonny.Shutdown()
			os.Exit(1)
		}
	}()

	/******************** Startup Major Services ********************/
	// Startup ROS connection.
	if *enROS {
		glog.V(1).Info("Starting up ROS...")
		ros := ros.NewRos(sonny)
		if err := ros.StartNode("rover"); err != nil {
			glog.Fatalf("Failed to start ROS:%v", err)
		}
		defer ros.Shutdown()
	}

	// Startup RPC service.
	lis, err := net.Listen("tcp", *rpcPort)
	if err != nil {
		glog.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDevicesRPCServer(s, rpc.New(sonny))
	go s.Serve(lis)

	// Startup HTTP service.
	h := httphandler.New(sonny, false, *res, *enVid, *enDataStream)
	if err := h.Start(*httpHostPort); err != nil {
		glog.Fatalf("Failed to listen: %v", err)
	}

}
