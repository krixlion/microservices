// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: pkg/grpc/pb/user.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
	// Get all Users with filter - A server-to-client streaming RPC.
	Index(ctx context.Context, in *UserFilter, opts ...grpc.CallOption) (User_IndexClient, error)
	// Create a new User - A simple RPC
	Create(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) Index(ctx context.Context, in *UserFilter, opts ...grpc.CallOption) (User_IndexClient, error) {
	stream, err := c.cc.NewStream(ctx, &User_ServiceDesc.Streams[0], "/user.User/Index", opts...)
	if err != nil {
		return nil, err
	}
	x := &userIndexClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type User_IndexClient interface {
	Recv() (*UserRequest, error)
	grpc.ClientStream
}

type userIndexClient struct {
	grpc.ClientStream
}

func (x *userIndexClient) Recv() (*UserRequest, error) {
	m := new(UserRequest)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userClient) Create(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/user.User/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
// All implementations must embed UnimplementedUserServer
// for forward compatibility
type UserServer interface {
	// Get all Users with filter - A server-to-client streaming RPC.
	Index(*UserFilter, User_IndexServer) error
	// Create a new User - A simple RPC
	Create(context.Context, *UserRequest) (*UserResponse, error)
	mustEmbedUnimplementedUserServer()
}

// UnimplementedUserServer must be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (UnimplementedUserServer) Index(*UserFilter, User_IndexServer) error {
	return status.Errorf(codes.Unimplemented, "method Index not implemented")
}
func (UnimplementedUserServer) Create(context.Context, *UserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedUserServer) mustEmbedUnimplementedUserServer() {}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_Index_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(UserFilter)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserServer).Index(m, &userIndexServer{stream})
}

type User_IndexServer interface {
	Send(*UserRequest) error
	grpc.ServerStream
}

type userIndexServer struct {
	grpc.ServerStream
}

func (x *userIndexServer) Send(m *UserRequest) error {
	return x.ServerStream.SendMsg(m)
}

func _User_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.User/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Create(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _User_Create_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Index",
			Handler:       _User_Index_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/grpc/pb/user.proto",
}
