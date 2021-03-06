package rpc

import (
	"errors"
	"fmt"

	"github.com/deepakkamesh/sonny/devices"
	pb "github.com/deepakkamesh/sonny/sonny"
	"github.com/golang/glog"
	google_pb "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

type Server struct {
	sonny devices.Platform
}

func New(s devices.Platform) *Server {
	return &Server{
		sonny: s,
	}
}

// Ping returns nil if the pic controller is up and responsive.
func (m *Server) Ping(ctx context.Context, in *google_pb.Empty) (*google_pb.Empty, error) {
	if !m.sonny.ControllerInitialized() {
		return &google_pb.Empty{}, errors.New("Controller not enabled")
	}

	return &google_pb.Empty{}, m.sonny.Ping()
}

// RoombaModeReq sets the roomba mode.
func (m *Server) SetRoombaMode(ctx context.Context, in *pb.RoombaModeReq) (*google_pb.Empty, error) {
	if !m.sonny.RoombaInitialized() {
		return &google_pb.Empty{}, errors.New("Roomba not enabled")
	}

	return &google_pb.Empty{}, m.sonny.SetRoombaMode(byte(in.Mode))
}

// Returns sensor information from Roomba.
func (m *Server) RoombaSensor(ctx context.Context, in *google_pb.Empty) (*pb.RoombaSensorRet, error) {
	if !m.sonny.RoombaInitialized() {
		return &pb.RoombaSensorRet{}, errors.New("Roomba not enabled")
	}

	data, err := m.sonny.GetRoombaTelemetry()
	if err != nil {
		glog.Errorf("Failed to get sensor %v", err)
		return &pb.RoombaSensorRet{}, fmt.Errorf("%v", err)
	}
	// Convert to the right data format for proto.
	dataOnWire := make(map[uint32]int32)
	for k, v := range data {
		dataOnWire[uint32(k)] = int32(v)
	}
	return &pb.RoombaSensorRet{Data: dataOnWire}, nil
}

// I2CBusEn enables/disables the I2C master bus.
func (m *Server) I2CBusEn(ctx context.Context, in *pb.I2CBusEnReq) (*google_pb.Empty, error) {
	return &google_pb.Empty{}, m.sonny.I2CBusEnable(in.On)
}

// LidarPower turns on/off the Lidar power.
func (m *Server) LidarPower(ctx context.Context, in *pb.LidarPowerReq) (*google_pb.Empty, error) {
	if !m.sonny.RoombaInitialized() {
		return &google_pb.Empty{}, errors.New("roomba not enabled")
	}
	return &google_pb.Empty{}, m.sonny.LidarPower(in.On)
}

// SecondaryPower turns on/off the secondary power supply for accessories.
func (m *Server) SecondaryPower(ctx context.Context, in *pb.SecPowerReq) (*google_pb.Empty, error) {
	if !m.sonny.RoombaInitialized() {
		return &google_pb.Empty{}, errors.New("roomba not enabled")
	}
	return &google_pb.Empty{}, m.sonny.AuxPower(in.On)
}

// LEDOn turns on/off the LED indicator.
func (m *Server) LEDOn(ctx context.Context, in *pb.LEDOnReq) (*google_pb.Empty, error) {
	if !m.sonny.ControllerInitialized() {
		return &google_pb.Empty{}, errors.New("Controller not enabled")
	}
	return &google_pb.Empty{}, m.sonny.LEDOn(in.On)
}

// LEDBlink blinks the LED.
func (m *Server) LEDBlink(ctx context.Context, in *pb.LEDBlinkReq) (*google_pb.Empty, error) {
	if !m.sonny.ControllerInitialized() {
		return &google_pb.Empty{}, errors.New("Controller not enabled")
	}
	return &google_pb.Empty{}, m.sonny.LEDBlink(uint16(in.Duration), byte(in.Times))
}

// Servo Rotate rotates the servo by angle.
func (m *Server) ServoRotate(ctx context.Context, in *pb.ServoReq) (*google_pb.Empty, error) {
	if !m.sonny.ControllerInitialized() {
		return &google_pb.Empty{}, errors.New("Controller not enabled")
	}
	return &google_pb.Empty{}, m.sonny.ServoRotate(byte(in.Servo), int(in.Angle))
}

// Moves moves the motor.
func (m *Server) Move(ctx context.Context, in *pb.MoveReq) (*pb.MoveRet, error) {
	r, err := m.sonny.MoveForward(int(in.Dist), int(in.Vel))
	if err != nil {
		return nil, err
	}
	return &pb.MoveRet{Dist: float32(r)}, nil
}

// Turn rotates the motor.
func (m *Server) Turn(ctx context.Context, in *pb.TurnReq) (*pb.TurnRet, error) {
	a, err := m.sonny.Turn(float64(in.Angle))
	if err != nil {
		return nil, err
	}
	return &pb.TurnRet{Delta: float32(a)}, nil

}

// Heading returns the magnetic heading.
func (m *Server) Heading(ctx context.Context, in *google_pb.Empty) (*pb.HeadingRet, error) {
	if !m.sonny.MagnetometerInitialized() {
		return nil, errors.New("magnetometer not enabled")
	}
	heading, err := m.sonny.TiltHeading()
	if err != nil {
		return nil, err
	}
	return &pb.HeadingRet{Heading: heading}, nil
}

// Distance returns the forward clearance in cm using the LIDAR.
func (m *Server) Distance(ctx context.Context, in *google_pb.Empty) (*pb.USRet, error) {
	return nil, errors.New("Not implemented")
}

// Accelerometer returns the dynamic and static acceleration from the accelerometer.
func (m *Server) Accelerometer(ctx context.Context, in *google_pb.Empty) (*pb.AccelRet, error) {

	x, y, z, err := m.sonny.Accelerometer()
	if err != nil {
		return nil, err
	}
	return &pb.AccelRet{X: int32(x), Y: int32(y), Z: int32(z)}, nil
}

// ForwardSweep returns the distance to the nearest object sweeping angle degrees at a time.
func (m *Server) ForwardSweep(ctx context.Context, in *pb.SweepReq) (*pb.SweepRet, error) {
	return nil, errors.New("Not implemented")
}

// PIRDetect retuns true if infrared signal is detected.
func (m *Server) PIRDetect(ctx context.Context, in *google_pb.Empty) (*pb.PIRRet, error) {
	if m.sonny.GetPIRState() == 1 {
		return &pb.PIRRet{On: true}, nil
	}
	return &pb.PIRRet{On: false}, nil
}

// BattState returns the battery level from pic.
func (m *Server) BattState(ctx context.Context, in *google_pb.Empty) (*pb.BattRet, error) {
	if !m.sonny.ControllerInitialized() {
		return &pb.BattRet{}, errors.New("Controller not enabled")
	}
	v, err := m.sonny.BattState()
	if err != nil {
		return nil, err
	}
	return &pb.BattRet{Volt: v}, nil
}

// LDR returns the light level from pic.
func (m *Server) LDR(ctx context.Context, in *google_pb.Empty) (*pb.LDRRet, error) {
	if !m.sonny.ControllerInitialized() {
		return &pb.LDRRet{}, errors.New("Controller not enabled")
	}
	v, err := m.sonny.LDR()
	if err != nil {
		return nil, err
	}
	return &pb.LDRRet{Adc: uint32(v)}, nil
}

// DHT11 returns the temp and humidity level from pic.
func (m *Server) DHT11(ctx context.Context, in *google_pb.Empty) (*pb.DHT11Ret, error) {
	if !m.sonny.ControllerInitialized() {
		return &pb.DHT11Ret{}, errors.New("Controller not enabled")
	}
	temp, humidity, err := m.sonny.DHT11()
	if err != nil {
		return nil, err
	}
	return &pb.DHT11Ret{Temp: uint32(temp), Humidity: uint32(humidity)}, nil
}
