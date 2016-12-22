// Code generated by protoc-gen-go.
// source: sonny.proto
// DO NOT EDIT!

/*
Package sonny is a generated protocol buffer package.

It is generated from these files:
	sonny.proto

It has these top-level messages:
	HeadingRet
	LEDBlinkReq
	LEDOnReq
	PIRRet
	ServoReq
	USRet
	SweepReq
	SweepRet
*/
package sonny

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type HeadingRet struct {
	Heading float64 `protobuf:"fixed64,1,opt,name=heading" json:"heading,omitempty"`
}

func (m *HeadingRet) Reset()                    { *m = HeadingRet{} }
func (m *HeadingRet) String() string            { return proto.CompactTextString(m) }
func (*HeadingRet) ProtoMessage()               {}
func (*HeadingRet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type LEDBlinkReq struct {
	Duration uint32 `protobuf:"varint,1,opt,name=duration" json:"duration,omitempty"`
	Times    uint32 `protobuf:"varint,2,opt,name=times" json:"times,omitempty"`
}

func (m *LEDBlinkReq) Reset()                    { *m = LEDBlinkReq{} }
func (m *LEDBlinkReq) String() string            { return proto.CompactTextString(m) }
func (*LEDBlinkReq) ProtoMessage()               {}
func (*LEDBlinkReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type LEDOnReq struct {
	On bool `protobuf:"varint,1,opt,name=On" json:"On,omitempty"`
}

func (m *LEDOnReq) Reset()                    { *m = LEDOnReq{} }
func (m *LEDOnReq) String() string            { return proto.CompactTextString(m) }
func (*LEDOnReq) ProtoMessage()               {}
func (*LEDOnReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type PIRRet struct {
	On bool `protobuf:"varint,1,opt,name=On" json:"On,omitempty"`
}

func (m *PIRRet) Reset()                    { *m = PIRRet{} }
func (m *PIRRet) String() string            { return proto.CompactTextString(m) }
func (*PIRRet) ProtoMessage()               {}
func (*PIRRet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type ServoReq struct {
	Servo uint32 `protobuf:"varint,1,opt,name=servo" json:"servo,omitempty"`
	Angle uint32 `protobuf:"varint,2,opt,name=angle" json:"angle,omitempty"`
}

func (m *ServoReq) Reset()                    { *m = ServoReq{} }
func (m *ServoReq) String() string            { return proto.CompactTextString(m) }
func (*ServoReq) ProtoMessage()               {}
func (*ServoReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type USRet struct {
	Distance int32 `protobuf:"varint,1,opt,name=distance" json:"distance,omitempty"`
}

func (m *USRet) Reset()                    { *m = USRet{} }
func (m *USRet) String() string            { return proto.CompactTextString(m) }
func (*USRet) ProtoMessage()               {}
func (*USRet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type SweepReq struct {
	Angle int32 `protobuf:"varint,1,opt,name=angle" json:"angle,omitempty"`
}

func (m *SweepReq) Reset()                    { *m = SweepReq{} }
func (m *SweepReq) String() string            { return proto.CompactTextString(m) }
func (*SweepReq) ProtoMessage()               {}
func (*SweepReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type SweepRet struct {
	Distance []int32 `protobuf:"varint,1,rep,packed,name=distance" json:"distance,omitempty"`
}

func (m *SweepRet) Reset()                    { *m = SweepRet{} }
func (m *SweepRet) String() string            { return proto.CompactTextString(m) }
func (*SweepRet) ProtoMessage()               {}
func (*SweepRet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func init() {
	proto.RegisterType((*HeadingRet)(nil), "sonny.HeadingRet")
	proto.RegisterType((*LEDBlinkReq)(nil), "sonny.LEDBlinkReq")
	proto.RegisterType((*LEDOnReq)(nil), "sonny.LEDOnReq")
	proto.RegisterType((*PIRRet)(nil), "sonny.PIRRet")
	proto.RegisterType((*ServoReq)(nil), "sonny.ServoReq")
	proto.RegisterType((*USRet)(nil), "sonny.USRet")
	proto.RegisterType((*SweepReq)(nil), "sonny.SweepReq")
	proto.RegisterType((*SweepRet)(nil), "sonny.SweepRet")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for DevicesRPC service

type DevicesRPCClient interface {
	Ping(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	LEDBlink(ctx context.Context, in *LEDBlinkReq, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	LEDOn(ctx context.Context, in *LEDOnReq, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	Heading(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*HeadingRet, error)
	PIRDetect(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*PIRRet, error)
	ServoRotate(ctx context.Context, in *ServoReq, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	Distance(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*USRet, error)
	ForwardSweep(ctx context.Context, in *SweepReq, opts ...grpc.CallOption) (*SweepRet, error)
}

type devicesRPCClient struct {
	cc *grpc.ClientConn
}

func NewDevicesRPCClient(cc *grpc.ClientConn) DevicesRPCClient {
	return &devicesRPCClient{cc}
}

func (c *devicesRPCClient) Ping(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/sonny.DevicesRPC/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *devicesRPCClient) LEDBlink(ctx context.Context, in *LEDBlinkReq, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/sonny.DevicesRPC/LEDBlink", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *devicesRPCClient) LEDOn(ctx context.Context, in *LEDOnReq, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/sonny.DevicesRPC/LEDOn", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *devicesRPCClient) Heading(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*HeadingRet, error) {
	out := new(HeadingRet)
	err := grpc.Invoke(ctx, "/sonny.DevicesRPC/Heading", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *devicesRPCClient) PIRDetect(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*PIRRet, error) {
	out := new(PIRRet)
	err := grpc.Invoke(ctx, "/sonny.DevicesRPC/PIRDetect", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *devicesRPCClient) ServoRotate(ctx context.Context, in *ServoReq, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/sonny.DevicesRPC/ServoRotate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *devicesRPCClient) Distance(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*USRet, error) {
	out := new(USRet)
	err := grpc.Invoke(ctx, "/sonny.DevicesRPC/Distance", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *devicesRPCClient) ForwardSweep(ctx context.Context, in *SweepReq, opts ...grpc.CallOption) (*SweepRet, error) {
	out := new(SweepRet)
	err := grpc.Invoke(ctx, "/sonny.DevicesRPC/ForwardSweep", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DevicesRPC service

type DevicesRPCServer interface {
	Ping(context.Context, *google_protobuf.Empty) (*google_protobuf.Empty, error)
	LEDBlink(context.Context, *LEDBlinkReq) (*google_protobuf.Empty, error)
	LEDOn(context.Context, *LEDOnReq) (*google_protobuf.Empty, error)
	Heading(context.Context, *google_protobuf.Empty) (*HeadingRet, error)
	PIRDetect(context.Context, *google_protobuf.Empty) (*PIRRet, error)
	ServoRotate(context.Context, *ServoReq) (*google_protobuf.Empty, error)
	Distance(context.Context, *google_protobuf.Empty) (*USRet, error)
	ForwardSweep(context.Context, *SweepReq) (*SweepRet, error)
}

func RegisterDevicesRPCServer(s *grpc.Server, srv DevicesRPCServer) {
	s.RegisterService(&_DevicesRPC_serviceDesc, srv)
}

func _DevicesRPC_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevicesRPCServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonny.DevicesRPC/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevicesRPCServer).Ping(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DevicesRPC_LEDBlink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LEDBlinkReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevicesRPCServer).LEDBlink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonny.DevicesRPC/LEDBlink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevicesRPCServer).LEDBlink(ctx, req.(*LEDBlinkReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DevicesRPC_LEDOn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LEDOnReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevicesRPCServer).LEDOn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonny.DevicesRPC/LEDOn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevicesRPCServer).LEDOn(ctx, req.(*LEDOnReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DevicesRPC_Heading_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevicesRPCServer).Heading(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonny.DevicesRPC/Heading",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevicesRPCServer).Heading(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DevicesRPC_PIRDetect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevicesRPCServer).PIRDetect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonny.DevicesRPC/PIRDetect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevicesRPCServer).PIRDetect(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DevicesRPC_ServoRotate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevicesRPCServer).ServoRotate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonny.DevicesRPC/ServoRotate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevicesRPCServer).ServoRotate(ctx, req.(*ServoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _DevicesRPC_Distance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevicesRPCServer).Distance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonny.DevicesRPC/Distance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevicesRPCServer).Distance(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DevicesRPC_ForwardSweep_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SweepReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DevicesRPCServer).ForwardSweep(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sonny.DevicesRPC/ForwardSweep",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DevicesRPCServer).ForwardSweep(ctx, req.(*SweepReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _DevicesRPC_serviceDesc = grpc.ServiceDesc{
	ServiceName: "sonny.DevicesRPC",
	HandlerType: (*DevicesRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _DevicesRPC_Ping_Handler,
		},
		{
			MethodName: "LEDBlink",
			Handler:    _DevicesRPC_LEDBlink_Handler,
		},
		{
			MethodName: "LEDOn",
			Handler:    _DevicesRPC_LEDOn_Handler,
		},
		{
			MethodName: "Heading",
			Handler:    _DevicesRPC_Heading_Handler,
		},
		{
			MethodName: "PIRDetect",
			Handler:    _DevicesRPC_PIRDetect_Handler,
		},
		{
			MethodName: "ServoRotate",
			Handler:    _DevicesRPC_ServoRotate_Handler,
		},
		{
			MethodName: "Distance",
			Handler:    _DevicesRPC_Distance_Handler,
		},
		{
			MethodName: "ForwardSweep",
			Handler:    _DevicesRPC_ForwardSweep_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("sonny.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 377 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x92, 0x4f, 0x6b, 0xea, 0x40,
	0x14, 0xc5, 0x9f, 0xfa, 0xa2, 0x79, 0x37, 0xfa, 0x4a, 0x87, 0x52, 0x42, 0x56, 0x92, 0x82, 0x74,
	0x95, 0x85, 0x2d, 0xd2, 0x5d, 0xa1, 0x8d, 0xa5, 0x85, 0x82, 0x32, 0xd2, 0x0f, 0x10, 0xf5, 0xd6,
	0x86, 0xea, 0x8c, 0x4d, 0x46, 0xa5, 0xdf, 0xaf, 0x1f, 0xac, 0xf3, 0x2f, 0x51, 0x84, 0xb8, 0xcb,
	0x99, 0xcc, 0xef, 0xdc, 0x9b, 0x73, 0x02, 0x5e, 0xce, 0x19, 0xfb, 0x8e, 0xd6, 0x19, 0x17, 0x9c,
	0x38, 0x5a, 0x04, 0x1e, 0xae, 0xd6, 0xc2, 0x9e, 0x85, 0x3d, 0x80, 0x67, 0x4c, 0xe6, 0x29, 0x5b,
	0x50, 0x14, 0xc4, 0x87, 0xd6, 0x87, 0x51, 0x7e, 0xad, 0x5b, 0xbb, 0xae, 0xd1, 0x42, 0x86, 0xf7,
	0xe0, 0xbd, 0x0e, 0xe3, 0x87, 0x65, 0xca, 0x3e, 0x29, 0x7e, 0x91, 0x00, 0xdc, 0xf9, 0x26, 0x4b,
	0x44, 0xca, 0x99, 0xbe, 0xd9, 0xa1, 0xa5, 0x26, 0x17, 0xe0, 0x88, 0x74, 0x85, 0xb9, 0x5f, 0xd7,
	0x2f, 0x8c, 0x08, 0x25, 0x21, 0x0d, 0x46, 0x4c, 0xd1, 0xff, 0xa1, 0x3e, 0x32, 0x9c, 0x4b, 0xe5,
	0x53, 0xe8, 0x43, 0x73, 0xfc, 0x42, 0xd5, 0x02, 0xc7, 0x6f, 0x06, 0xe0, 0x4e, 0x30, 0xdb, 0x72,
	0x45, 0x49, 0xdf, 0x5c, 0x3d, 0xdb, 0x81, 0x46, 0xa8, 0xd3, 0x84, 0x2d, 0x96, 0x58, 0x4c, 0xd3,
	0x22, 0xbc, 0x02, 0xe7, 0x6d, 0xa2, 0x0c, 0xd5, 0xa2, 0x69, 0x2e, 0x12, 0x36, 0x43, 0xcd, 0x39,
	0xb4, 0xd4, 0x61, 0x57, 0x9a, 0xef, 0x10, 0xd7, 0xd6, 0xdc, 0xd8, 0x98, 0x4b, 0xd6, 0xa6, 0x57,
	0xde, 0x38, 0x76, 0x6a, 0x1c, 0x3a, 0xf5, 0x7f, 0x1a, 0x00, 0x31, 0x6e, 0xd3, 0x19, 0xe6, 0x74,
	0xfc, 0x48, 0xee, 0xe0, 0xef, 0x58, 0x86, 0x46, 0x2e, 0xa3, 0x05, 0xe7, 0xd2, 0xc7, 0x64, 0x3d,
	0xdd, 0xbc, 0x47, 0x43, 0x15, 0x7d, 0x50, 0x71, 0x1e, 0xfe, 0x91, 0xa4, 0x5b, 0xc4, 0x4c, 0x48,
	0x64, 0xca, 0x3b, 0xc8, 0xfd, 0x04, 0xd9, 0x07, 0x47, 0xe7, 0x4b, 0xce, 0xf6, 0x98, 0x4e, 0xfb,
	0x04, 0x33, 0x80, 0x96, 0x2d, 0xbf, 0x72, 0xd5, 0x73, 0xeb, 0xb6, 0xff, 0x49, 0x24, 0x77, 0x0b,
	0xff, 0x64, 0x5f, 0x31, 0x0a, 0x9c, 0x89, 0x4a, 0xb2, 0x63, 0x49, 0xd3, 0xac, 0xfe, 0x36, 0xcf,
	0x74, 0xc9, 0x45, 0x22, 0xb0, 0xdc, 0xb3, 0xe8, 0xf7, 0xe4, 0xb7, 0xb9, 0xb1, 0x8d, 0xba, 0x72,
	0x5c, 0xdb, 0xda, 0xe9, 0xda, 0x35, 0xd3, 0x7e, 0xe2, 0xd9, 0x2e, 0xc9, 0xe6, 0xba, 0xc1, 0xfd,
	0x38, 0xdb, 0x78, 0x70, 0x74, 0x20, 0x99, 0x69, 0x53, 0x7b, 0xde, 0xfc, 0x06, 0x00, 0x00, 0xff,
	0xff, 0xdc, 0xef, 0x72, 0xf3, 0x36, 0x03, 0x00, 0x00,
}
