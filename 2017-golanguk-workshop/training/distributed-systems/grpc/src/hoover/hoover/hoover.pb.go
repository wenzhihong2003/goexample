// Code generated by protoc-gen-gogo.
// source: hoover.proto
// DO NOT EDIT!

/*
Package hoover is a generated protocol buffer package.

It is generated from these files:
	hoover.proto

It has these top-level messages:
	GetRequest
	GetReply
*/
package hoover

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/duration"

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
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type GetRequest struct {
	Url string `protobuf:"bytes,1,opt,name=Url,json=url,proto3" json:"Url,omitempty"`
}

func (m *GetRequest) Reset()                    { *m = GetRequest{} }
func (m *GetRequest) String() string            { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()               {}
func (*GetRequest) Descriptor() ([]byte, []int) { return fileDescriptorHoover, []int{0} }

func (m *GetRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type GetReply struct {
	ResponseCode int32                     `protobuf:"varint,1,opt,name=response_code,json=responseCode,proto3" json:"response_code,omitempty"`
	Body         string                    `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	Elapsed      *google_protobuf.Duration `protobuf:"bytes,3,opt,name=elapsed" json:"elapsed,omitempty"`
}

func (m *GetReply) Reset()                    { *m = GetReply{} }
func (m *GetReply) String() string            { return proto.CompactTextString(m) }
func (*GetReply) ProtoMessage()               {}
func (*GetReply) Descriptor() ([]byte, []int) { return fileDescriptorHoover, []int{1} }

func (m *GetReply) GetResponseCode() int32 {
	if m != nil {
		return m.ResponseCode
	}
	return 0
}

func (m *GetReply) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func (m *GetReply) GetElapsed() *google_protobuf.Duration {
	if m != nil {
		return m.Elapsed
	}
	return nil
}

func init() {
	proto.RegisterType((*GetRequest)(nil), "hoover.GetRequest")
	proto.RegisterType((*GetReply)(nil), "hoover.GetReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Service service

type ServiceClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error)
}

type serviceClient struct {
	cc *grpc.ClientConn
}

func NewServiceClient(cc *grpc.ClientConn) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetReply, error) {
	out := new(GetReply)
	err := grpc.Invoke(ctx, "/hoover.Service/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Service service

type ServiceServer interface {
	Get(context.Context, *GetRequest) (*GetReply, error)
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hoover.Service/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "hoover.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Service_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hoover.proto",
}

func init() { proto.RegisterFile("hoover.proto", fileDescriptorHoover) }

var fileDescriptorHoover = []byte{
	// 215 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8f, 0xcf, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x8d, 0xab, 0xad, 0x8e, 0x15, 0xca, 0x9c, 0x62, 0x0f, 0xa5, 0xc4, 0x4b, 0x2f, 0x6e,
	0xa1, 0xbd, 0x78, 0x57, 0xe8, 0x7d, 0xc5, 0xb3, 0x24, 0xd9, 0x31, 0x06, 0x96, 0xcc, 0xba, 0x3f,
	0x82, 0xf9, 0xef, 0x85, 0x4d, 0x82, 0x78, 0x9b, 0x79, 0xf3, 0xf1, 0xf8, 0x06, 0x56, 0x5f, 0xcc,
	0x3d, 0x39, 0x69, 0x1d, 0x07, 0xc6, 0xc5, 0xb8, 0x6d, 0xb6, 0x0d, 0x73, 0x63, 0xe8, 0x90, 0xd2,
	0x2a, 0x7e, 0x1e, 0x74, 0x74, 0x65, 0x68, 0xb9, 0x1b, 0xb9, 0x62, 0x0b, 0x70, 0xa6, 0xa0, 0xe8,
	0x3b, 0x92, 0x0f, 0xb8, 0x06, 0xf1, 0xee, 0x4c, 0x9e, 0xed, 0xb2, 0xfd, 0xad, 0x12, 0xd1, 0x99,
	0xe2, 0x07, 0x6e, 0xd2, 0xdd, 0x9a, 0x01, 0x1f, 0xe1, 0xde, 0x91, 0xb7, 0xdc, 0x79, 0xfa, 0xa8,
	0x59, 0x53, 0xe2, 0xae, 0xd5, 0x6a, 0x0e, 0x5f, 0x58, 0x13, 0x22, 0x5c, 0x55, 0xac, 0x87, 0xfc,
	0x32, 0x75, 0xa4, 0x19, 0x4f, 0xb0, 0x24, 0x53, 0x5a, 0x4f, 0x3a, 0x17, 0xbb, 0x6c, 0x7f, 0x77,
	0x7c, 0x90, 0xa3, 0x96, 0x9c, 0xb5, 0xe4, 0xeb, 0xa4, 0xa5, 0x66, 0xf2, 0xf8, 0x0c, 0xcb, 0x37,
	0x72, 0x7d, 0x5b, 0x13, 0x3e, 0x81, 0x38, 0x53, 0x40, 0x94, 0xd3, 0x8b, 0x7f, 0xc6, 0x9b, 0xf5,
	0xbf, 0xcc, 0x9a, 0xa1, 0xb8, 0xa8, 0x16, 0xa9, 0xf5, 0xf4, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x0f,
	0x48, 0x2d, 0xc4, 0x12, 0x01, 0x00, 0x00,
}