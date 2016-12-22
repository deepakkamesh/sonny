package rpc

import (
	"errors"

	"github.com/deepakkamesh/sonny/devices"
	pb "github.com/deepakkamesh/sonny/sonny"
	google_pb "github.com/golang/protobuf/ptypes/empty"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/chip"
	"github.com/kidoman/embd/sensor/lsm303"
	"github.com/kidoman/embd/sensor/us020"
	"golang.org/x/net/context"
)

type Devices struct {
	Ctrl *devices.Controller
	Mag  *lsm303.LSM303
	Pir  string
	Us   *us020.US020
}

type Server struct {
	ctrl *devices.Controller
	mag  *lsm303.LSM303
	us   *us020.US020
	pir  string
}

func New(d *Devices) *Server {
	return &Server{
		ctrl: d.Ctrl,
		mag:  d.Mag,
		pir:  d.Pir,
		us:   d.Us,
	}
}

// Ping returns nil if the pic controller is up and responsive.
func (m *Server) Ping(ctx context.Context, in *google_pb.Empty) (*google_pb.Empty, error) {
	if m.ctrl == nil {
		return &google_pb.Empty{}, errors.New("controller not enabled")
	}
	return &google_pb.Empty{}, m.ctrl.Ping()
}

// Servo Rotate rotates the servo by angle.
func (m *Server) ServoRotate(ctx context.Context, in *pb.ServoReq) (*google_pb.Empty, error) {
	if m.ctrl == nil {
		return &google_pb.Empty{}, errors.New("controller not enabled")
	}
	return &google_pb.Empty{}, m.ctrl.ServoRotate(byte(in.Servo), byte(in.Angle))
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

// Heading returns the magnetic heading.
func (m *Server) Heading(ctx context.Context, in *google_pb.Empty) (*pb.HeadingRet, error) {
	heading, err := m.mag.Heading()
	if err != nil {
		return nil, err
	}
	return &pb.HeadingRet{Heading: heading}, nil
}

// Distance returns the forward clearance in cm using the ultrasonic sensor.
func (m *Server) Distance(ctx context.Context, in *google_pb.Empty) (*pb.USRet, error) {
	d, err := m.us.Distance()
	if err != nil {
		return nil, err
	}
	return &pb.USRet{Distance: int32(d)}, nil
}

// ForwardSweep returns the distance to the nearest object sweeping angle degrees at a time.
func (m *Server) ForwardSweep(ctx context.Context, in *pb.SweepReq) (*pb.SweepRet, error) {
	v, err := devices.ForwardSweep(m.ctrl, m.us, int(in.Angle))

	if err != nil {
		return nil, err
	}
	return &pb.SweepRet{Distance: v}, nil
}

// PIRDetect retuns true if Infrared signal is detected.
func (m *Server) PIRDetect(ctx context.Context, in *google_pb.Empty) (*pb.PIRRet, error) {
	v, err := embd.DigitalRead(m.pir)
	if err != nil {
		return nil, err
	}
	if v == embd.High {
		return &pb.PIRRet{On: true}, nil
	}
	return &pb.PIRRet{On: false}, nil
}
