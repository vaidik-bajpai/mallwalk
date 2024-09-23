package main

import (
	"context"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	store ProductsStore
}

func NewService(store ProductsStore) *service {
	return &service{
		store: store,
	}
}

func (s *service) CreateProduct(ctx context.Context, cp *pb.CreateProductRequest) (*pb.Product, error) {
	p, err := s.store.Create(ctx, &Product{
		Name:        cp.Name,
		Price:       cp.Price,
		Description: cp.Description,
		Category:    cp.Category,
		Stock:       cp.Stock,
		Image:       cp.Image,
		Rating:      cp.Rating,
	})
	if err != nil {
		return &pb.Product{}, err
	}

	return toPBProduct(p), nil
}

func (s *service) GetProduct(ctx context.Context, gp *pb.GetProductRequest) (*pb.Product, error) {
	id, err := primitive.ObjectIDFromHex(gp.ProductID)
	if err != nil {
		return &pb.Product{}, err
	}

	p, err := s.store.Get(ctx, id)
	if err != nil {
		return &pb.Product{}, err
	}

	return toPBProduct(p), err
}

func (s *service) ListProduct(ctx context.Context, lp *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	p, err := s.store.List(ctx, lp)
	if err != nil {
		return &pb.ListProductsResponse{}, err
	}

	return &pb.ListProductsResponse{
		Products:      p,
		TotalProducts: uint32(len(p)),
	}, nil
}

func (s *service) UpdateProduct(ctx context.Context, up *pb.UpdateProductRequest) (*pb.Product, error) {
	p, err := s.store.Update(ctx, up)
	if err != nil {
		return &pb.Product{}, err
	}

	return toPBProduct(p), nil
}

func (s *service) DeleteProduct(ctx context.Context, dp *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	id, err := primitive.ObjectIDFromHex(dp.ProductID)
	if err != nil {
		return &pb.DeleteProductResponse{
			Success: false,
		}, err
	}

	err = s.store.Delete(ctx, id)
	if err != nil {
		return &pb.DeleteProductResponse{
			Success: false,
		}, err
	}

	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}

func toPBProduct(p *Product) *pb.Product {
	return &pb.Product{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Category:    p.Category,
		Rating:      p.Rating,
		Image:       p.Image,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   timestamppb.New(p.CreatedAt),
	}
}
