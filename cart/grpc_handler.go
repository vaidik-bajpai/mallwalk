package main

import (
	"context"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedCartServiceServer

	service CartService
}

func NewGRPCHandler(grpcServer *grpc.Server, service CartService) {
	handler := &grpcHandler{
		service: service,
	}

	pb.RegisterCartServiceServer(grpcServer, handler)
}

func (h *grpcHandler) AddToCart(ctx context.Context, it *pb.AddToCartRequest) (*pb.CartResponse, error) {
	err := h.service.ValidateItem(ctx, it)
	if err != nil {
		return nil, err
	}

	return h.service.AddToCart(ctx, it)
}

func (h *grpcHandler) RemoveFromCart(ctx context.Context, ri *pb.RemoveItemRequest) (*pb.CartResponse, error) {
	return h.service.RemoveFromCart(ctx, ri)
}

func (h *grpcHandler) ViewCart(ctx context.Context, vc *pb.ViewCartRequest) (*pb.Cart, error) {
	return h.service.ViewCart(ctx, vc)
}
