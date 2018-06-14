package devices

import "github.com/saljam/mjpeg"

// Platform defines the interface for the platform.
type Platform interface {
	Sensors(byte) ([]byte, error)
	TiltHeading() (float64, error)
	CalibrateCompass(bool) error
	ForwardSweep(int, int, int) ([]int32, error)
	DirectDrive(int16, int16) error
	Gyro() (int16, int16, int16, error)
	Accelerometer() (int16, int16, int16, error)
	GetAuxPowerState() int
	AuxPower(enable bool) error
	GetRoombaTelemetry() (map[byte]int16, error)
	GetI2CBusState() int
	BattState() (float32, error)
	LDR() (uint16, error)
	DHT11() (uint8, uint8, error)
	ServoRotate(byte, int) error
	I2CBusEnable(bool) error
	Ping() error
	Drive(int16, int16) error
	LEDOn(bool) error
	LEDBlink(uint16, byte) error
	Full() error
	Reset() error
	Safe() error
	Power() error
	Passive() error
	SeekDock() error
	GetRoombaMode() int
	SetRoombaMode(byte) error
	GetPIRState() int
	RoombaInitialized() bool
	StartRoomba(bool) error
	GetVideoStream() *mjpeg.Stream
	LidarPower(bool) error
	Distance() (int, error)
	ControllerInitialized() bool
	MagnetometerInitialized() bool
	LidarInitialized() bool
}
