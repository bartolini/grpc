// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package miaow

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

// ChatServiceClient is the client API for ChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatServiceClient interface {
	SayHello(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
	StreamHello(ctx context.Context, opts ...grpc.CallOption) (ChatService_StreamHelloClient, error)
}

type chatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatServiceClient(cc grpc.ClientConnInterface) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) SayHello(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/ChatService/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) StreamHello(ctx context.Context, opts ...grpc.CallOption) (ChatService_StreamHelloClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChatService_ServiceDesc.Streams[0], "/ChatService/StreamHello", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceStreamHelloClient{stream}
	return x, nil
}

type ChatService_StreamHelloClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type chatServiceStreamHelloClient struct {
	grpc.ClientStream
}

func (x *chatServiceStreamHelloClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *chatServiceStreamHelloClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatServiceServer is the server API for ChatService service.
// All implementations must embed UnimplementedChatServiceServer
// for forward compatibility
type ChatServiceServer interface {
	SayHello(context.Context, *Message) (*Message, error)
	StreamHello(ChatService_StreamHelloServer) error
	mustEmbedUnimplementedChatServiceServer()
}

// UnimplementedChatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChatServiceServer struct {
}

func (UnimplementedChatServiceServer) SayHello(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedChatServiceServer) StreamHello(ChatService_StreamHelloServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamHello not implemented")
}
func (UnimplementedChatServiceServer) mustEmbedUnimplementedChatServiceServer() {}

// UnsafeChatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServiceServer will
// result in compilation errors.
type UnsafeChatServiceServer interface {
	mustEmbedUnimplementedChatServiceServer()
}

func RegisterChatServiceServer(s grpc.ServiceRegistrar, srv ChatServiceServer) {
	s.RegisterService(&ChatService_ServiceDesc, srv)
}

func _ChatService_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ChatService/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).SayHello(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_StreamHello_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ChatServiceServer).StreamHello(&chatServiceStreamHelloServer{stream})
}

type ChatService_StreamHelloServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type chatServiceStreamHelloServer struct {
	grpc.ServerStream
}

func (x *chatServiceStreamHelloServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *chatServiceStreamHelloServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ChatService_ServiceDesc is the grpc.ServiceDesc for ChatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _ChatService_SayHello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamHello",
			Handler:       _ChatService_StreamHello_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "miaow.proto",
}
