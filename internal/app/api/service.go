package service

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route/internal/app/models"
	"route/internal/app/module"
	order "route/pkg/api/proto/order/v1/order/v1"
)

type Service interface {
	AcceptOrder(context.Context, *order.OrderRequest) (*order.OrderResponse, error)
	ReturnOrder(context.Context, *order.OrderRequest) (*order.OrderResponse, error)
	IssueOrder(context.Context, *order.OrderRequest) (*order.OrderResponse, error)
	ListOrders(context.Context, *order.ListOrdersRequest) (*order.ListResponse, error)
	AcceptReturn(context.Context, *order.OrderRequest) (*order.OrderResponse, error)
	ListReturns(context.Context, *order.ListReturnsRequest) (*order.ListResponse, error)
}

type OrderService struct {
	mod module.Module
	order.UnimplementedOrderServiceServer
}

func New(mod module.Module) *OrderService {
	return &OrderService{mod: mod}
}

func (o *OrderService) AcceptOrder(_ context.Context, req *order.OrderRequest) (*order.OrderResponse, error) {
	or := orderToDomain(req)
	err := o.mod.AcceptOrder(&or, models.ToPackageType(req.GetPackagingType()))
	if err != nil {
		log.Printf("Error accepting order: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &order.OrderResponse{Status: "success"}, nil
}

func (o *OrderService) ReturnOrder(_ context.Context, req *order.OrderRequest) (*order.OrderResponse, error) {
	or := orderToDomain(req)
	err := o.mod.ReturnOrder(or.OrderID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &order.OrderResponse{Status: "success"}, nil
}

func (o *OrderService) IssueOrder(_ context.Context, req *order.OrderRequest) (*order.OrderResponse, error) {
	or := orderToDomain(req)
	err := o.mod.IssueOrder(or.OrderID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &order.OrderResponse{Status: "success"}, nil
}

func (o *OrderService) ListOrders(_ context.Context, req *order.ListOrdersRequest) (*order.ListResponse, error) {
	listOr, err := o.mod.ListOrders(int(req.GetUserId()), int(req.GetLastN()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	orders := make([]*order.OrderInfo, len(listOr))
	for i, or := range listOr {
		orders[i] = &order.OrderInfo{
			OrderId: int32(or.OrderID),
			UserId:  int32(or.UserID),
			Status:  "Success",
			Weight:  or.Weight,
		}
	}
	return &order.ListResponse{Orders: orders}, nil
}

func (o *OrderService) AcceptReturn(_ context.Context, req *order.OrderRequest) (*order.OrderResponse, error) {
	or := orderToDomain(req)
	err := o.mod.AcceptReturn(or.OrderID, or.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &order.OrderResponse{Status: "success"}, nil
}

func (o *OrderService) ListReturns(_ context.Context, req *order.ListReturnsRequest) (*order.ListResponse, error) {
	listRet, err := o.mod.ListReturns(int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	orders := make([]*order.OrderInfo, len(listRet))
	for i, or := range listRet {
		orders[i] = &order.OrderInfo{
			OrderId: int32(or.OrderID),
			UserId:  int32(or.UserID),
			Status:  "Success",
			Weight:  or.Weight,
		}
	}
	return &order.ListResponse{Orders: orders}, nil
}

func (o *OrderService) mustEmbedUnimplementedOrderServiceServer() {}

func orderToDomain(req *order.OrderRequest) models.Order {
	return models.Order{
		OrderID: int(req.GetOrderId()),
		UserID:  int(req.GetUserId()),
		Weight:  req.GetWeight(),
	}
}
