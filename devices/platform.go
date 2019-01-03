package devices

import (
	"github.com/deepakkamesh/ydlidar"
	"github.com/saljam/mjpeg"
)

// Platform defines the interface for the platform.
type Platform interface {
	Sensors(byte) ([]byte, error)
	TiltHeading() (float64, error)
	CalibrateCompass(bool) error
	DirectDrive(int16, int16) error
	Gyro() (int16, int16, int16, error)
	Accelerometer() (int16, int16, int16, error)
	GetAuxPowerState() int
	AuxPower(enable bool) error
	GetI2CBusState() int
	BattState() (float32, error)
	LDR() (uint16, error)
	DHT11() (uint8, uint8, error)
	ServoRotate(byte, int) error
	I2CBusEnable(bool) error
	Ping() error
	LEDOn(bool) error
	LEDBlink(uint16, byte) error
	GetPIRState() int
	GetVideoStream() *mjpeg.Stream
	/* Roomba Functions */
	GetRoombaTelemetry() (map[byte]int16, error)
	Drive(int16, int16) error
	Full() error
	Reset() error
	Safe() error
	Power() error
	Passive() error
	SeekDock() error
	GetRoombaMode() int
	SetRoombaMode(byte) error
	RoombaInitialized() bool
	StartRoomba(bool) error
	/* Lidar Functions */
	LidarPower(bool) error
	LidarData() ydlidar.Packet
	StartScan()
	/* Utility functions */
	Shutdown()
	Startup()
	ControllerInitialized() bool
	MagnetometerInitialized() bool
	LidarInitialized() bool
	MoveForward(int, int) (float64, error)
	Turn(float64) (float64, error)
}
