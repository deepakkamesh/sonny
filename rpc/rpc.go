package rpc

import (
	"github.com/deepakkamesh/sonny/devices"
	pb "github.com/deepakkamesh/sonny/sonny"
	google_pb "github.com/golang/protobuf/ptypes/empty"
	_ "github.com/kidoman/embd/host/chip"
	"github.com/kidoman/embd/sensor/lsm303"
	"golang.org/x/net/context"
)

type Devices struct {
	Ctrl *devices.Controller
	Mag  *lsm303.LSM303
}

type Server struct {
	ctrl *devices.Controller
	mag  *lsm303.LSM303
}

func New(d *Devices) *Server {
	return &Server{
		ctrl: d.Ctrl,
		mag:  d.Mag,
	}
}

// Ping returns nil if the pic controller is up and responsive.
func (m *Server) Ping(ctx context.Context, in *google_pb.Empty) (*google_pb.Empty, error) {
	return &google_pb.Empty{}, m.ctrl.Ping()
}

// LEDOn turns on/off the LED indicator.
func (m *Server) LEDOn(ctx context.Context, in *pb.LEDOnReq) (*google_pb.Empty, error) {
	return &google_pb.Empty{}, m.ctrl.LEDOn(in.On)
}

// LEDBlink blinks the LED.
func (m *Server) LEDBlink(ctx context.Context, in *pb.LEDBlinkReq) (*google_pb.Empty, error) {
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
