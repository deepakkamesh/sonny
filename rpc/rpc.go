package rpc

import (
	"errors"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/sonny/devices"
	pb "github.com/deepakkamesh/sonny/sonny"
	google_pb "github.com/golang/protobuf/ptypes/empty"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/chip"
	"golang.org/x/net/context"
)

type Devices struct {
	Ctrl   *devices.Controller
	Lidar  *i2c.LIDARLiteDriver
	Mag    *i2c.HMC6352Driver
	Pir    *int
	Roomba *roomba.Roomba
	Chip   *chip.Adaptor
}

type Server struct {
	ctrl   *devices.Controller
	lidar  *i2c.LIDARLiteDriver
	mag    *i2c.HMC6352Driver
	pir    *int
	roomba *roomba.Roomba
}

func New(d *Devices) *Server {
	return &Server{
		ctrl:   d.Ctrl,
		mag:    d.Mag,
		pir:    d.Pir,
		roomba: d.Roomba,
	}
}

// Ping returns nil if the pic controller is up and responsive.
func (m *Server) Ping(ctx context.Context, in *google_pb.Empty) (*google_pb.Empty, error) {
	if m.ctrl == nil {
		return &google_pb.Empty{}, errors.New("controller not enabled")
	}
	return &google_pb.Empty{}, m.ctrl.Ping()
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

// Distance returns the forward clearance in cm using the ultrasonic sensor.
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
