package gateway

import (
	"context"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
)

type StockGateway interface {
	CheckIfItemIsInStock(context.Context, *pb.CheckIfItemIsInStockRequest) (*pb.CheckIfItemIsInStockResponse, error)
	UpdateStock(context.Context, *pb.CheckIfItemIsInStockRequest) (*pb.CartResponse, error)
}
