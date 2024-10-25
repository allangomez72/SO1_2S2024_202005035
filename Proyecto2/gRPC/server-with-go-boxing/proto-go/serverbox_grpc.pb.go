// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: serverbox.proto

//nombre del package

package facultad

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	FacultadService_SendUserInfo_FullMethodName = "/facultad.FacultadService/SendUserInfo"
)

// FacultadServiceClient is the client API for FacultadService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Definimos el servicio
type FacultadServiceClient interface {
	// El servidor ricibirá el mensaje de User y devolverá una respuesta de UserResponse
	SendUserInfo(ctx context.Context, in *Student, opts ...grpc.CallOption) (*StudentResponse, error)
}

type facultadServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFacultadServiceClient(cc grpc.ClientConnInterface) FacultadServiceClient {
	return &facultadServiceClient{cc}
}

func (c *facultadServiceClient) SendUserInfo(ctx context.Context, in *Student, opts ...grpc.CallOption) (*StudentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StudentResponse)
	err := c.cc.Invoke(ctx, FacultadService_SendUserInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FacultadServiceServer is the server API for FacultadService service.
// All implementations must embed UnimplementedFacultadServiceServer
// for forward compatibility.
//
// Definimos el servicio
type FacultadServiceServer interface {
	// El servidor ricibirá el mensaje de User y devolverá una respuesta de UserResponse
	SendUserInfo(context.Context, *Student) (*StudentResponse, error)
	mustEmbedUnimplementedFacultadServiceServer()
}

// UnimplementedFacultadServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFacultadServiceServer struct{}

func (UnimplementedFacultadServiceServer) SendUserInfo(context.Context, *Student) (*StudentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendUserInfo not implemented")
}
func (UnimplementedFacultadServiceServer) mustEmbedUnimplementedFacultadServiceServer() {}
func (UnimplementedFacultadServiceServer) testEmbeddedByValue()                         {}

// UnsafeFacultadServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FacultadServiceServer will
// result in compilation errors.
type UnsafeFacultadServiceServer interface {
	mustEmbedUnimplementedFacultadServiceServer()
}

func RegisterFacultadServiceServer(s grpc.ServiceRegistrar, srv FacultadServiceServer) {
	// If the following call pancis, it indicates UnimplementedFacultadServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FacultadService_ServiceDesc, srv)
}

func _FacultadService_SendUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Student)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FacultadServiceServer).SendUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FacultadService_SendUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FacultadServiceServer).SendUserInfo(ctx, req.(*Student))
	}
	return interceptor(ctx, in, info, handler)
}

// FacultadService_ServiceDesc is the grpc.ServiceDesc for FacultadService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FacultadService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "facultad.FacultadService",
	HandlerType: (*FacultadServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendUserInfo",
			Handler:    _FacultadService_SendUserInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "serverbox.proto",
}