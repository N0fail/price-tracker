// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: api/api.proto

package api

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

// AdminClient is the client API for Admin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdminClient interface {
	// Создает продукт с переданным кодом и именем
	// может вернуть ошибки: ErrNameTooShortError, ErrProductExists
	ProductCreate(ctx context.Context, in *ProductCreateRequest, opts ...grpc.CallOption) (*ProductCreateResponse, error)
	// Возвращает список всех продуктов (код, имя, последняя цена(если есть))
	ProductList(ctx context.Context, in *ProductListRequest, opts ...grpc.CallOption) (*ProductListResponse, error)
	// Удаляет продукт с переданным кодом
	// может вернуть ошибки: ErrProductNotExist
	ProductDelete(ctx context.Context, in *ProductDeleteRequest, opts ...grpc.CallOption) (*ProductDeleteResponse, error)
	// Добавляет цену для продукта с переданным кодом, дата передается в Unix формате
	// может вернуть ошибки: ErrProductNotExist, ErrNegativePrice
	PriceTimeStampAdd(ctx context.Context, in *PriceTimeStampAddRequest, opts ...grpc.CallOption) (*PriceTimeStampAddResponse, error)
	// Возвращает массив всех цен для продукта в хронологичеком порядке (принимает код продукта)
	// может вернуть ошибки: ErrProductNotExist
	PriceHistory(ctx context.Context, in *PriceHistoryRequest, opts ...grpc.CallOption) (*PriceHistoryResponse, error)
}

type adminClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminClient(cc grpc.ClientConnInterface) AdminClient {
	return &adminClient{cc}
}

func (c *adminClient) ProductCreate(ctx context.Context, in *ProductCreateRequest, opts ...grpc.CallOption) (*ProductCreateResponse, error) {
	out := new(ProductCreateResponse)
	err := c.cc.Invoke(ctx, "/ozon.dev.price_tracker.api.Admin/ProductCreate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) ProductList(ctx context.Context, in *ProductListRequest, opts ...grpc.CallOption) (*ProductListResponse, error) {
	out := new(ProductListResponse)
	err := c.cc.Invoke(ctx, "/ozon.dev.price_tracker.api.Admin/ProductList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) ProductDelete(ctx context.Context, in *ProductDeleteRequest, opts ...grpc.CallOption) (*ProductDeleteResponse, error) {
	out := new(ProductDeleteResponse)
	err := c.cc.Invoke(ctx, "/ozon.dev.price_tracker.api.Admin/ProductDelete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) PriceTimeStampAdd(ctx context.Context, in *PriceTimeStampAddRequest, opts ...grpc.CallOption) (*PriceTimeStampAddResponse, error) {
	out := new(PriceTimeStampAddResponse)
	err := c.cc.Invoke(ctx, "/ozon.dev.price_tracker.api.Admin/PriceTimeStampAdd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminClient) PriceHistory(ctx context.Context, in *PriceHistoryRequest, opts ...grpc.CallOption) (*PriceHistoryResponse, error) {
	out := new(PriceHistoryResponse)
	err := c.cc.Invoke(ctx, "/ozon.dev.price_tracker.api.Admin/PriceHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminServer is the server API for Admin service.
// All implementations must embed UnimplementedAdminServer
// for forward compatibility
type AdminServer interface {
	// Создает продукт с переданным кодом и именем
	// может вернуть ошибки: ErrNameTooShortError, ErrProductExists
	ProductCreate(context.Context, *ProductCreateRequest) (*ProductCreateResponse, error)
	// Возвращает список всех продуктов (код, имя, последняя цена(если есть))
	ProductList(context.Context, *ProductListRequest) (*ProductListResponse, error)
	// Удаляет продукт с переданным кодом
	// может вернуть ошибки: ErrProductNotExist
	ProductDelete(context.Context, *ProductDeleteRequest) (*ProductDeleteResponse, error)
	// Добавляет цену для продукта с переданным кодом, дата передается в Unix формате
	// может вернуть ошибки: ErrProductNotExist, ErrNegativePrice
	PriceTimeStampAdd(context.Context, *PriceTimeStampAddRequest) (*PriceTimeStampAddResponse, error)
	// Возвращает массив всех цен для продукта в хронологичеком порядке (принимает код продукта)
	// может вернуть ошибки: ErrProductNotExist
	PriceHistory(context.Context, *PriceHistoryRequest) (*PriceHistoryResponse, error)
	mustEmbedUnimplementedAdminServer()
}

// UnimplementedAdminServer must be embedded to have forward compatible implementations.
type UnimplementedAdminServer struct {
}

func (UnimplementedAdminServer) ProductCreate(context.Context, *ProductCreateRequest) (*ProductCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProductCreate not implemented")
}
func (UnimplementedAdminServer) ProductList(context.Context, *ProductListRequest) (*ProductListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProductList not implemented")
}
func (UnimplementedAdminServer) ProductDelete(context.Context, *ProductDeleteRequest) (*ProductDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProductDelete not implemented")
}
func (UnimplementedAdminServer) PriceTimeStampAdd(context.Context, *PriceTimeStampAddRequest) (*PriceTimeStampAddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PriceTimeStampAdd not implemented")
}
func (UnimplementedAdminServer) PriceHistory(context.Context, *PriceHistoryRequest) (*PriceHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PriceHistory not implemented")
}
func (UnimplementedAdminServer) mustEmbedUnimplementedAdminServer() {}

// UnsafeAdminServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdminServer will
// result in compilation errors.
type UnsafeAdminServer interface {
	mustEmbedUnimplementedAdminServer()
}

func RegisterAdminServer(s grpc.ServiceRegistrar, srv AdminServer) {
	s.RegisterService(&Admin_ServiceDesc, srv)
}

func _Admin_ProductCreate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProductCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).ProductCreate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ozon.dev.price_tracker.api.Admin/ProductCreate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).ProductCreate(ctx, req.(*ProductCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_ProductList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProductListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).ProductList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ozon.dev.price_tracker.api.Admin/ProductList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).ProductList(ctx, req.(*ProductListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_ProductDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProductDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).ProductDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ozon.dev.price_tracker.api.Admin/ProductDelete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).ProductDelete(ctx, req.(*ProductDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_PriceTimeStampAdd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PriceTimeStampAddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).PriceTimeStampAdd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ozon.dev.price_tracker.api.Admin/PriceTimeStampAdd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).PriceTimeStampAdd(ctx, req.(*PriceTimeStampAddRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Admin_PriceHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PriceHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminServer).PriceHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ozon.dev.price_tracker.api.Admin/PriceHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminServer).PriceHistory(ctx, req.(*PriceHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Admin_ServiceDesc is the grpc.ServiceDesc for Admin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Admin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ozon.dev.price_tracker.api.Admin",
	HandlerType: (*AdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProductCreate",
			Handler:    _Admin_ProductCreate_Handler,
		},
		{
			MethodName: "ProductList",
			Handler:    _Admin_ProductList_Handler,
		},
		{
			MethodName: "ProductDelete",
			Handler:    _Admin_ProductDelete_Handler,
		},
		{
			MethodName: "PriceTimeStampAdd",
			Handler:    _Admin_PriceTimeStampAdd_Handler,
		},
		{
			MethodName: "PriceHistory",
			Handler:    _Admin_PriceHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/api.proto",
}
