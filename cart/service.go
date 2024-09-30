package main

import (
	"context"
	"fmt"

	"github.com/vaidik-bajpai/mallwalk/cart/gateway"
	pb "github.com/vaidik-bajpai/mallwalk/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	store   CartsStore
	gateway gateway.StockGateway
}

func NewService(store CartsStore, gateway gateway.StockGateway) *service {
	return &service{
		store:   store,
		gateway: gateway,
	}
}

func (s *service) AddToCart(ctx context.Context, at *pb.AddToCartRequest) (*pb.CartResponse, error) {
	res, err := s.gateway.CheckIfItemIsInStock(ctx, &pb.CheckIfItemIsInStockRequest{
		ID: at.Item.ProductID,
	})
	if !res.InStock || err != nil {
		return nil, nil
	}

	err = s.store.Add(ctx, at.CartID, &Item{
		ProductID: at.Item.ProductID,
		Name:      at.Item.Name,
		Image:     at.Item.Image,
		Price:     at.Item.Price,
		Quantity:  at.Item.Quantity,
	})

	return &pb.CartResponse{}, err
}

func (s *service) RemoveFromCart(ctx context.Context, ri *pb.RemoveItemRequest) (*pb.CartResponse, error) {
	err := s.store.Remove(ctx, ri.CartID, ri.ProductID)
	return &pb.CartResponse{}, err
}

func (s *service) ViewCart(ctx context.Context, vc *pb.ViewCartRequest) (*pb.Cart, error) {
	return s.store.View(ctx, vc.CartID)
}

func (s *service) ValidateItem(ctx context.Context, at *pb.AddToCartRequest) error {
	if _, err := primitive.ObjectIDFromHex(at.CartID); err != nil {
		return err
	}

	if _, err := primitive.ObjectIDFromHex(at.Item.ProductID); err != nil {
		return err
	}

	if at.Item.Name == "" || at.Item.Image == "" || at.Item.Price <= 0 || at.Item.Quantity <= 0 {
		return fmt.Errorf("invalid item data")
	}

	return nil
}
