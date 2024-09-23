package gateway

import (
	"context"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
	discovery "github.com/vaidik-bajpai/mallwalk/common/discovery"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway {
	return &gateway{registry: registry}
}

func (g *gateway) CreateUser(ctx context.Context, userDetails *pb.CreateUserRequest) (*pb.User, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "user", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewUserServiceClient(conn)

	_, err = c.CreateUser(ctx, userDetails)

	return nil, err
}

func (g *gateway) UserLogin(ctx context.Context, userDetails *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "user", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewUserServiceClient(conn)

	return c.UserLogin(ctx, userDetails)
}

func (g *gateway) CreateProduct(ctx context.Context, cp *pb.CreateProductRequest) (*pb.Product, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "product", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewProductServiceClient(conn)

	return c.CreateProduct(ctx, cp)
}

func (g *gateway) GetProduct(ctx context.Context, gp *pb.GetProductRequest) (*pb.Product, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "product", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewProductServiceClient(conn)

	return c.GetProduct(ctx, gp)
}

func (g *gateway) ListProduct(ctx context.Context, lp *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "product", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewProductServiceClient(conn)

	return c.ListProduct(ctx, lp)
}

func (g *gateway) UpdateProduct(ctx context.Context, up *pb.UpdateProductRequest) (*pb.Product, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "product", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewProductServiceClient(conn)

	return c.UpdateProduct(ctx, up)
}

func (g *gateway) DeleteProduct(ctx context.Context, dp *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "product", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewProductServiceClient(conn)

	return c.DeleteProduct(ctx, dp)
}

func (g *gateway) AddToCart(ctx context.Context, it *pb.AddToCartRequest) (*pb.CartResponse, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "cart", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewCartServiceClient(conn)

	return c.AddToCart(ctx, it)
}

func (g *gateway) RemoveFromCart(ctx context.Context, ri *pb.RemoveItemRequest) (*pb.CartResponse, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "cart", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewCartServiceClient(conn)

	return c.RemoveFromCart(ctx, ri)
}

func (g *gateway) ViewCart(ctx context.Context, vc *pb.ViewCartRequest) (*pb.Cart, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "cart", g.registry)
	if err != nil {
		return nil, err
	}

	c := pb.NewCartServiceClient(conn)

	return c.ViewCart(ctx, vc)
}

func (g *gateway) PlaceOrder(ctx context.Context) error  { return nil }
func (g *gateway) MakePayment(ctx context.Context) error { return nil }
