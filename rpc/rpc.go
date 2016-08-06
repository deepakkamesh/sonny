package rpc

import (
	"github.com/deepakkamesh/sonny/devices"
	google_pb "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

type Server struct {
	controller *devices.Controller
}

func New(c *devices.Controller) *Server {
	return &Server{
		controller: c,
	}
}

func (m *Server) Ping(ctx context.Context, in *google_pb.Empty) (*google_pb.Empty, error) {
	return &google_pb.Empty{}, nil
}
