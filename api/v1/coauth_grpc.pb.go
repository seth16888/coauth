// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.23.3
// source: v1/coauth.proto

package v1

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
	Coauth_Authorize_FullMethodName = "/api.coauth.v1.coauth/Authorize"
	Coauth_Token_FullMethodName     = "/api.coauth.v1.coauth/Token"
	Coauth_AddApp_FullMethodName    = "/api.coauth.v1.coauth/AddApp"
)

// CoauthClient is the client API for Coauth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CoauthClient interface {
	Authorize(ctx context.Context, in *AuthorizeRequest, opts ...grpc.CallOption) (*AuthorizeReply, error)
	Token(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*TokenReply, error)
	AddApp(ctx context.Context, in *AddAppRequest, opts ...grpc.CallOption) (*AddAppReply, error)
}

type coauthClient struct {
	cc grpc.ClientConnInterface
}

func NewCoauthClient(cc grpc.ClientConnInterface) CoauthClient {
	return &coauthClient{cc}
}

func (c *coauthClient) Authorize(ctx context.Context, in *AuthorizeRequest, opts ...grpc.CallOption) (*AuthorizeReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthorizeReply)
	err := c.cc.Invoke(ctx, Coauth_Authorize_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coauthClient) Token(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*TokenReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TokenReply)
	err := c.cc.Invoke(ctx, Coauth_Token_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coauthClient) AddApp(ctx context.Context, in *AddAppRequest, opts ...grpc.CallOption) (*AddAppReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddAppReply)
	err := c.cc.Invoke(ctx, Coauth_AddApp_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CoauthServer is the server API for Coauth service.
// All implementations must embed UnimplementedCoauthServer
// for forward compatibility.
type CoauthServer interface {
	Authorize(context.Context, *AuthorizeRequest) (*AuthorizeReply, error)
	Token(context.Context, *TokenRequest) (*TokenReply, error)
	AddApp(context.Context, *AddAppRequest) (*AddAppReply, error)
	mustEmbedUnimplementedCoauthServer()
}

// UnimplementedCoauthServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCoauthServer struct{}

func (UnimplementedCoauthServer) Authorize(context.Context, *AuthorizeRequest) (*AuthorizeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}
func (UnimplementedCoauthServer) Token(context.Context, *TokenRequest) (*TokenReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Token not implemented")
}
func (UnimplementedCoauthServer) AddApp(context.Context, *AddAppRequest) (*AddAppReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddApp not implemented")
}
func (UnimplementedCoauthServer) mustEmbedUnimplementedCoauthServer() {}
func (UnimplementedCoauthServer) testEmbeddedByValue()                {}

// UnsafeCoauthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CoauthServer will
// result in compilation errors.
type UnsafeCoauthServer interface {
	mustEmbedUnimplementedCoauthServer()
}

func RegisterCoauthServer(s grpc.ServiceRegistrar, srv CoauthServer) {
	// If the following call pancis, it indicates UnimplementedCoauthServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Coauth_ServiceDesc, srv)
}

func _Coauth_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthorizeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoauthServer).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coauth_Authorize_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoauthServer).Authorize(ctx, req.(*AuthorizeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coauth_Token_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoauthServer).Token(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coauth_Token_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoauthServer).Token(ctx, req.(*TokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Coauth_AddApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoauthServer).AddApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Coauth_AddApp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoauthServer).AddApp(ctx, req.(*AddAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Coauth_ServiceDesc is the grpc.ServiceDesc for Coauth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Coauth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.coauth.v1.coauth",
	HandlerType: (*CoauthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authorize",
			Handler:    _Coauth_Authorize_Handler,
		},
		{
			MethodName: "Token",
			Handler:    _Coauth_Token_Handler,
		},
		{
			MethodName: "AddApp",
			Handler:    _Coauth_AddApp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/coauth.proto",
}
