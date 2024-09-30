package main

import (
	"context"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
)

type StockService interface {
	CheckIfItemIsInStock(context.Context, *pb.CheckIfItemIsInStockRequest) (*pb.CheckIfItemIsInStockResponse, error)
	UpdateStock(context.Context, *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error)
}

type StockStore interface {
	GetItemByID(context.Context, string) (*ItemStock, error)
	Update(ctx context.Context, pID string, delta int64) error
}

type ItemStock struct {
	ID    string `bson:"_id"`
	Stock int64  `bson:"stock"`
}
