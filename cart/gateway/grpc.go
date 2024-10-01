package gateway

import (
	"context"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
	"github.com/vaidik-bajpai/mallwalk/common/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func (g *Gateway) UpdateStock(context.Context, *pb.CheckIfItemIsInStockRequest) (*pb.CartResponse, error) {
	panic("unimplemented")
}

func NewStockGateway(registry discovery.Registry) *Gateway {
	return &Gateway{
		registry: registry,
	}
}

func (g *Gateway) CheckIfItemIsInStock(ctx context.Context, p *pb.CheckIfItemIsInStockRequest) (*pb.CheckIfItemIsInStockResponse, error) {
	conn, err := discovery.ServiceConnection(ctx, "stock", g.registry)
	if err != nil {
		return nil, nil
	}

	c := pb.NewStocksServiceClient(conn)

	return c.CheckIfItemIsInStock(ctx, p)
}
