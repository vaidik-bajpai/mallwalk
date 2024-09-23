package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type store struct {
	client *mongo.Client
}

const (
	DBName   = "carts"
	CollName = "carts"
)

func NewStore(client *mongo.Client) *store {
	return &store{
		client: client,
	}
}

type Item struct {
	ProductID string `bson:"product_id"`
	Name      string `bson:"name"`
	Price     uint32 `bson:"price"`
	Image     string `bson:"image"`
	Quantity  uint32 `bson:"quantity"`
}

func (s *store) Add(ctx context.Context, cartID string, it *Item) error {
	col := s.client.Database(DBName).Collection(CollName)

	filter := bson.M{"cart_id": cartID}
	update := bson.M{
		"$push": bson.M{"items": it},
	}
	opts := options.Update().SetUpsert(true)

	_, err := col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to add item to cart: %v", err)
	}

	return nil
}

func (s *store) Remove(ctx context.Context, cartID, productID string) error {
	col := s.client.Database(DBName).Collection(CollName)

	filter := bson.M{"cart_id": cartID}

	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{"product_id": productID},
		},
	}

	_, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

type Cart struct {
	CartID string `bson:"cart_id"`
	Items  []Item `bson:"items"`
}

func (s *store) View(ctx context.Context, cID string) (*pb.Cart, error) {
	col := s.client.Database(DBName).Collection(CollName)

	// Fetch the cart document based on cart_id
	var cart Cart
	err := col.FindOne(ctx, bson.M{"cart_id": cID}).Decode(&cart)
	if err != nil {
		return nil, err
	}

	// Initialize variables for items and total price
	var items []*pb.Item
	var totalPrice uint32

	// Iterate over items in the cart and add them to the result
	for _, cartItem := range cart.Items {
		log.Println(cartItem)                            // Print the cart item for debugging
		items = append(items, toPBItem(&cartItem))       // Convert to pb.Item
		totalPrice += cartItem.Price * cartItem.Quantity // Calculate total price
	}

	// Return the pb.Cart with converted items and total price
	return &pb.Cart{
		CartID:     cart.CartID,
		Items:      items,
		TotalPrice: totalPrice,
	}, nil
}
