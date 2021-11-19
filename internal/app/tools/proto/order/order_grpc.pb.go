// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package orders

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

// OrderServiceClient is the client API for OrderService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderServiceClient interface {
	PutOrder(ctx context.Context, in *Order, opts ...grpc.CallOption) (*OrderCost, error)
	GetAllOrders(ctx context.Context, in *UserIdOrder, opts ...grpc.CallOption) (*ArrayOrders, error)
}

type orderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderServiceClient(cc grpc.ClientConnInterface) OrderServiceClient {
	return &orderServiceClient{cc}
}

func (c *orderServiceClient) PutOrder(ctx context.Context, in *Order, opts ...grpc.CallOption) (*OrderCost, error) {
	out := new(OrderCost)
	err := c.cc.Invoke(ctx, "/OrderService/PutOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) GetAllOrders(ctx context.Context, in *UserIdOrder, opts ...grpc.CallOption) (*ArrayOrders, error) {
	out := new(ArrayOrders)
	err := c.cc.Invoke(ctx, "/OrderService/GetAllOrders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServiceServer is the server API for OrderService service.
// All implementations must embed UnimplementedOrderServiceServer
// for forward compatibility
type OrderServiceServer interface {
	PutOrder(context.Context, *Order) (*OrderCost, error)
	GetAllOrders(context.Context, *UserIdOrder) (*ArrayOrders, error)
	mustEmbedUnimplementedOrderServiceServer()
}

// UnimplementedOrderServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOrderServiceServer struct {
}

func (UnimplementedOrderServiceServer) PutOrder(context.Context, *Order) (*OrderCost, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutOrder not implemented")
}
func (UnimplementedOrderServiceServer) GetAllOrders(context.Context, *UserIdOrder) (*ArrayOrders, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllOrders not implemented")
}
func (UnimplementedOrderServiceServer) mustEmbedUnimplementedOrderServiceServer() {}

// UnsafeOrderServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServiceServer will
// result in compilation errors.
type UnsafeOrderServiceServer interface {
	mustEmbedUnimplementedOrderServiceServer()
}

func RegisterOrderServiceServer(s grpc.ServiceRegistrar, srv OrderServiceServer) {
	s.RegisterService(&OrderService_ServiceDesc, srv)
}

func _OrderService_PutOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Order)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).PutOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OrderService/PutOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).PutOrder(ctx, req.(*Order))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_GetAllOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserIdOrder)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).GetAllOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OrderService/GetAllOrders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).GetAllOrders(ctx, req.(*UserIdOrder))
	}
	return interceptor(ctx, in, info, handler)
}

// OrderService_ServiceDesc is the grpc.ServiceDesc for OrderService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "OrderService",
	HandlerType: (*OrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PutOrder",
			Handler:    _OrderService_PutOrder_Handler,
		},
		{
			MethodName: "GetAllOrders",
			Handler:    _OrderService_GetAllOrders_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order.proto",
}