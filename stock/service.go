package main

import (
	"context"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
)

type service struct {
	store StockStore
}

func NewService(store StockStore) *service {
	return &service{
		store: store,
	}
}

func (s *service) CheckIfItemIsInStock(ctx context.Context, p *pb.CheckIfItemIsInStockRequest) (*pb.CheckIfItemIsInStockResponse, error) {
	item, err := s.store.GetItemByID(ctx, p.ID)
	if err != nil {
		return nil, err
	}
	if item.Stock >= p.Quantity {
		return &pb.CheckIfItemIsInStockResponse{
			InStock: true,
		}, nil
	}

	return &pb.CheckIfItemIsInStockResponse{
		InStock: false,
	}, nil
}

func (s *service) UpdateStock(ctx context.Context, us *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
	err := s.store.Update(ctx, us.ID, us.Delta)
	return &pb.UpdateStockResponse{}, err
}
