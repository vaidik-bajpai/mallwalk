package main

import (
	"context"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedProductServiceServer

	service ProductsService
}

func NewGRPCHandler(grpcServer *grpc.Server, service ProductsService) {
	handler := &grpcHandler{
		service: service,
	}
	pb.RegisterProductServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateProduct(ctx context.Context, cp *pb.CreateProductRequest) (*pb.Product, error) {
	//validation

	return h.service.CreateProduct(ctx, cp)
}

func (h *grpcHandler) GetProduct(ctx context.Context, gp *pb.GetProductRequest) (*pb.Product, error) {
	return h.service.GetProduct(ctx, gp)
}

func (h *grpcHandler) ListProduct(ctx context.Context, lp *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	return h.service.ListProduct(ctx, lp)
}

func (h *grpcHandler) UpdateProduct(ctx context.Context, up *pb.UpdateProductRequest) (*pb.Product, error) {
	return h.service.UpdateProduct(ctx, up)
}

func (h *grpcHandler) DeleteProduct(ctx context.Context, dp *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	return h.service.DeleteProduct(ctx, dp)
}
