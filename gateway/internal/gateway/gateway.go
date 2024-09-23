package gateway

import (
	"context"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
)

type APIGateway interface {
	CreateUser(context.Context, *pb.CreateUserRequest) (*pb.User, error)
	UserLogin(context.Context, *pb.UserLoginRequest) (*pb.UserLoginResponse, error)

	CreateProduct(context.Context, *pb.CreateProductRequest) (*pb.Product, error)
	GetProduct(context.Context, *pb.GetProductRequest) (*pb.Product, error)
	ListProduct(context.Context, *pb.ListProductsRequest) (*pb.ListProductsResponse, error)
	UpdateProduct(context.Context, *pb.UpdateProductRequest) (*pb.Product, error)
	DeleteProduct(context.Context, *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error)

	AddToCart(context.Context, *pb.AddToCartRequest) (*pb.CartResponse, error)
	RemoveFromCart(context.Context, *pb.RemoveItemRequest) (*pb.CartResponse, error)
	ViewCart(context.Context, *pb.ViewCartRequest) (*pb.Cart, error)

	PlaceOrder(context.Context) error
	MakePayment(context.Context) error
}
