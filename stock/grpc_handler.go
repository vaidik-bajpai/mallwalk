package main

import (
	"context"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedStocksServiceServer

	service StockService
}

func NewGRPCHandler(grpcServer *grpc.Server, service StockService) {
	handler := &grpcHandler{
		service: service,
	}
	pb.RegisterStocksServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CheckIfItemIsInStock(ctx context.Context, p *pb.CheckIfItemIsInStockRequest) (*pb.CheckIfItemIsInStockResponse, error) {
	return h.service.CheckIfItemIsInStock(ctx, p)
}

func (h *grpcHandler) UpdateStock(ctx context.Context, us *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
	return h.service.UpdateStock(ctx, us)
}
