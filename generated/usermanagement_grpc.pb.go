// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: usermanagement.proto

package usermanagement

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

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	// Add the following RPC method for login
	LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error)
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
	StoreUser(ctx context.Context, in *StoreUserRequest, opts ...grpc.CallOption) (*StoreUserResponse, error)
	UploadImage(ctx context.Context, opts ...grpc.CallOption) (UserService_UploadImageClient, error)
	FetchImage(ctx context.Context, in *FetchImageRequest, opts ...grpc.CallOption) (UserService_FetchImageClient, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error) {
	out := new(LoginUserResponse)
	err := c.cc.Invoke(ctx, "/usermanagement.UserService/LoginUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error) {
	out := new(GetUserResponse)
	err := c.cc.Invoke(ctx, "/usermanagement.UserService/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) StoreUser(ctx context.Context, in *StoreUserRequest, opts ...grpc.CallOption) (*StoreUserResponse, error) {
	out := new(StoreUserResponse)
	err := c.cc.Invoke(ctx, "/usermanagement.UserService/StoreUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UploadImage(ctx context.Context, opts ...grpc.CallOption) (UserService_UploadImageClient, error) {
	stream, err := c.cc.NewStream(ctx, &UserService_ServiceDesc.Streams[0], "/usermanagement.UserService/UploadImage", opts...)
	if err != nil {
		return nil, err
	}
	x := &userServiceUploadImageClient{stream}
	return x, nil
}

type UserService_UploadImageClient interface {
	Send(*ImageData) error
	CloseAndRecv() (*UploadImageResponse, error)
	grpc.ClientStream
}

type userServiceUploadImageClient struct {
	grpc.ClientStream
}

func (x *userServiceUploadImageClient) Send(m *ImageData) error {
	return x.ClientStream.SendMsg(m)
}

func (x *userServiceUploadImageClient) CloseAndRecv() (*UploadImageResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadImageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *userServiceClient) FetchImage(ctx context.Context, in *FetchImageRequest, opts ...grpc.CallOption) (UserService_FetchImageClient, error) {
	stream, err := c.cc.NewStream(ctx, &UserService_ServiceDesc.Streams[1], "/usermanagement.UserService/FetchImage", opts...)
	if err != nil {
		return nil, err
	}
	x := &userServiceFetchImageClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type UserService_FetchImageClient interface {
	Recv() (*ImageData, error)
	grpc.ClientStream
}

type userServiceFetchImageClient struct {
	grpc.ClientStream
}

func (x *userServiceFetchImageClient) Recv() (*ImageData, error) {
	m := new(ImageData)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	// Add the following RPC method for login
	LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error)
	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
	StoreUser(context.Context, *StoreUserRequest) (*StoreUserResponse, error)
	UploadImage(UserService_UploadImageServer) error
	FetchImage(*FetchImageRequest, UserService_FetchImageServer) error
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedUserServiceServer) GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUserServiceServer) StoreUser(context.Context, *StoreUserRequest) (*StoreUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreUser not implemented")
}
func (UnimplementedUserServiceServer) UploadImage(UserService_UploadImageServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadImage not implemented")
}
func (UnimplementedUserServiceServer) FetchImage(*FetchImageRequest, UserService_FetchImageServer) error {
	return status.Errorf(codes.Unimplemented, "method FetchImage not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usermanagement.UserService/LoginUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).LoginUser(ctx, req.(*LoginUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usermanagement.UserService/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUser(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_StoreUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).StoreUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/usermanagement.UserService/StoreUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).StoreUser(ctx, req.(*StoreUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UploadImage_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UserServiceServer).UploadImage(&userServiceUploadImageServer{stream})
}

type UserService_UploadImageServer interface {
	SendAndClose(*UploadImageResponse) error
	Recv() (*ImageData, error)
	grpc.ServerStream
}

type userServiceUploadImageServer struct {
	grpc.ServerStream
}

func (x *userServiceUploadImageServer) SendAndClose(m *UploadImageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *userServiceUploadImageServer) Recv() (*ImageData, error) {
	m := new(ImageData)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _UserService_FetchImage_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FetchImageRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UserServiceServer).FetchImage(m, &userServiceFetchImageServer{stream})
}

type UserService_FetchImageServer interface {
	Send(*ImageData) error
	grpc.ServerStream
}

type userServiceFetchImageServer struct {
	grpc.ServerStream
}

func (x *userServiceFetchImageServer) Send(m *ImageData) error {
	return x.ServerStream.SendMsg(m)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "usermanagement.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LoginUser",
			Handler:    _UserService_LoginUser_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _UserService_GetUser_Handler,
		},
		{
			MethodName: "StoreUser",
			Handler:    _UserService_StoreUser_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadImage",
			Handler:       _UserService_UploadImage_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "FetchImage",
			Handler:       _UserService_FetchImage_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "usermanagement.proto",
}
