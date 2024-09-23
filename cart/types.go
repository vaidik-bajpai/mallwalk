package main

import (
	"context"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
)

type CartService interface {
	AddToCart(ctx context.Context, it *pb.AddToCartRequest) (*pb.CartResponse, error)
	RemoveFromCart(ctx context.Context, ri *pb.RemoveItemRequest) (*pb.CartResponse, error)
	ViewCart(ctx context.Context, vc *pb.ViewCartRequest) (*pb.Cart, error)
	ValidateItem(context.Context, *pb.AddToCartRequest) error
}

type CartsStore interface {
	Add(ctx context.Context, cartID string, it *Item) error
	Remove(ctx context.Context, cartID, productID string) error
	View(ctx context.Context, cID string) (*pb.Cart, error)
}
