package rpc

import (
	"errors"
	"fmt"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/sonny/devices"
	pb "github.com/deepakkamesh/sonny/sonny"
	"github.com/golang/glog"
	google_pb "github.com/golang/protobuf/ptypes/empty"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"golang.org/x/net/context"
)

type Devices struct {
	Ctrl   *devices.Controller   // PIC controller.
	Lidar  *i2c.LIDARLiteDriver  // Lidar Lite.
	Mag    *i2c.HMC6352Driver    // Magnetometer HMC5663.
	Pir    *int                  // PIR.
	Roomba *roomba.Roomba        // Roomba controller.
	I2CEn  *gpio.DirectPinDriver // GPIO port control for I2C Bus.
}

type Server struct {
	ctrl   *devices.Controller
	lidar  *i2c.LIDARLiteDriver
	mag    *i2c.HMC6352Driver
	pir    *int
	roomba *roomba.Roomba
	i2cEn  *gpio.DirectPinDriver
}

func New(d *Devices) *Server {
	return &Server{
		ctrl:   d.Ctrl,
		mag:    d.Mag,
		pir:    d.Pir,
		roomba: d.Roomba,
		i2cEn:  d.I2CEn,
	}
}

// Ping returns nil if the pic controller is up and responsive.
func (m *Server) Ping(ctx context.Context, in *google_pb.Empty) (*google_pb.Empty, error) {
	if m.ctrl == nil {
		return &google_pb.Empty{}, errors.New("controller not enabled")
	}
	return &google_pb.Empty{}, m.ctrl.Ping()
}

// RoombaModeReq sets the roomba mode.
func (m *Server) SetRoombaMode(ctx context.Context, in *pb.RoombaModeReq) (*google_pb.Empty, error) {
	if m.roomba == nil {
		return &google_pb.Empty{}, errors.New("Roomba not enabled")
	}

	return &google_pb.Empty{}, devices.SetRoombaMode(m.roomba, byte(in.Mode))
}

// Returns sensor information from Roomba.
func (m *Server) RoombaSensor(ctx context.Context, in *google_pb.Empty) (*pb.RoombaSensorRet, error) {
	if m.roomba == nil {
		return &pb.RoombaSensorRet{}, errors.New("Roomba not enabled")
	}

	data, err := devices.GetRoombaTelemetry(m.roomba)
	if err != nil {
		glog.Error("Failed to get sensor %v", err)
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
	if m.roomba == nil {
		return &google_pb.Empty{}, errors.New("roomba not enabled")
	}
	if in.On {
		return &google_pb.Empty{}, m.i2cEn.DigitalWrite(1)
	}
	return &google_pb.Empty{}, m.i2cEn.DigitalWrite(0)
}

// SecondaryPower turns on/off the secondary power supply for accessories.
func (m *Server) SecondaryPower(ctx context.Context, in *pb.SecPowerReq) (*google_pb.Empty, error) {
	if m.roomba == nil {
		return &google_pb.Empty{}, errors.New("roomba not enabled")
	}
	return &google_pb.Empty{}, m.roomba.MainBrush(in.On, true)
}

// LEDOn turns on/off the LED indicator.
func (m *Server) LEDOn(ctx context.Context, in *pb.LEDOnReq) (*google_pb.Empty, error) {
	if m.ctrl == nil {
		return &google_pb.Empty{}, errors.New("controller not enabled")
	}
	return &google_pb.Empty{}, m.ctrl.LEDOn(in.On)
}

// LEDBlink blinks the LED.
func (m *Server) LEDBlink(ctx context.Context, in *pb.LEDBlinkReq) (*google_pb.Empty, error) {
	if m.ctrl == nil {
		return &google_pb.Empty{}, errors.New("controller not enabled")
	}
	return &google_pb.Empty{}, m.ctrl.LEDBlink(uint16(in.Duration), byte(in.Times))
}

// Servo Rotate rotates the servo by angle.
func (m *Server) ServoRotate(ctx context.Context, in *pb.ServoReq) (*google_pb.Empty, error) {
	if m.ctrl == nil {
		return &google_pb.Empty{}, errors.New("controller not enabled")
	}
	return &google_pb.Empty{}, m.ctrl.ServoRotate(byte(in.Servo), int(in.Angle))
}

// Moves moves the motor.
func (m *Server) Move(ctx context.Context, in *pb.MoveReq) (*pb.MoveRet, error) {
	/*	if m.ctrl == nil {
			return nil, errors.New("controller not enabled")
		}
		m1, m2, err := m.ctrl.Move(int16(in.Turns), in.Fwd, byte(in.DutyPercent))
		if err != nil {
			return nil, err
		}
		return &pb.MoveRet{
			M1Turns: uint32(m1),
			M2Turns: uint32(m2),
		}, nil*/
	return nil, nil
}

// Turn rotates the motor.
func (m *Server) Turn(ctx context.Context, in *pb.TurnReq) (*pb.TurnRet, error) {
	/*	if m.ctrl == nil {
			return nil, errors.New("controller not enabled")
		}
		m1, m2, err := m.ctrl.Turn(int16(in.Turns), byte(in.RotateType), byte(in.DutyPercent))
		if err != nil {
			return nil, err
		}
		return &pb.TurnRet{
			M1Turns: uint32(m1),
			M2Turns: uint32(m2),
		}, nil*/
	return nil, nil
}

// Heading returns the magnetic heading.
func (m *Server) Heading(ctx context.Context, in *google_pb.Empty) (*pb.HeadingRet, error) {
	if m.mag == nil {
		return nil, errors.New("magnetometer not enabled")
	}
	// TODO: To be implemented.
	//heading, err := m.mag.Heading()
	heading := 0.0
	var err error
	if err != nil {
		return nil, err
	}
	return &pb.HeadingRet{Heading: heading}, nil
}

// Distance returns the forward clearance in cm using the LIDAR.
func (m *Server) Distance(ctx context.Context, in *google_pb.Empty) (*pb.USRet, error) {
	if m.ctrl == nil {
		return nil, errors.New("controller not enabled")
	}
	d, err := m.lidar.Distance()
	if err != nil {
		return nil, err
	}
	return &pb.USRet{Distance: int32(d)}, nil
}

// Accelerometer returns the dynamic and static acceleration from the accelerometer.
func (m *Server) Accelerometer(ctx context.Context, in *google_pb.Empty) (*pb.AccelRet, error) {
	if m.ctrl == nil {
		return nil, errors.New("controller not enabled")
	}

	x, y, z, err := m.ctrl.Accelerometer()
	if err != nil {
		return nil, err
	}
	return &pb.AccelRet{X: x, Y: y, Z: z}, nil
}

// ForwardSweep returns the distance to the nearest object sweeping angle degrees at a time.
func (m *Server) ForwardSweep(ctx context.Context, in *pb.SweepReq) (*pb.SweepRet, error) {
	v, err := devices.ForwardSweep(m.ctrl, int(in.Angle))

	if err != nil {
		return nil, err
	}
	return &pb.SweepRet{Distance: v}, nil
}

// PIRDetect retuns true if infrared signal is detected.
func (m *Server) PIRDetect(ctx context.Context, in *google_pb.Empty) (*pb.PIRRet, error) {

	if m.pir == nil {
		return nil, errors.New("PIR Sensor not initialized")
	}

	if *m.pir == 1 {
		return &pb.PIRRet{On: true}, nil
	}

	return &pb.PIRRet{On: false}, nil
}

// BattState returns the battery level from pic.
func (m *Server) BattState(ctx context.Context, in *google_pb.Empty) (*pb.BattRet, error) {
	if m.ctrl == nil {
		return nil, errors.New("controller not enabled")
	}

	v, err := m.ctrl.BattState()
	if err != nil {
		return nil, err
	}
	return &pb.BattRet{Volt: v}, nil
}

// LDR returns the light level from pic.
func (m *Server) LDR(ctx context.Context, in *google_pb.Empty) (*pb.LDRRet, error) {
	if m.ctrl == nil {
		return nil, errors.New("controller not enabled")
	}

	v, err := m.ctrl.LDR()
	if err != nil {
		return nil, err
	}
	return &pb.LDRRet{Adc: uint32(v)}, nil
}

// DHT11 returns the temp and humidity level from pic.
func (m *Server) DHT11(ctx context.Context, in *google_pb.Empty) (*pb.DHT11Ret, error) {
	if m.ctrl == nil {
		return nil, errors.New("controller not enabled")
	}

	temp, humidity, err := m.ctrl.DHT11()
	if err != nil {
		return nil, err
	}
	return &pb.DHT11Ret{Temp: uint32(temp), Humidity: uint32(humidity)}, nil
}
