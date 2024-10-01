package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type store struct {
	client *mongo.Client
}

const (
	DBName   = "stocks"
	CollName = "stocks"
)

func NewStore(client *mongo.Client) *store {
	return &store{client: client}
}

func (s *store) GetItemByID(ctx context.Context, pID string) (*ItemStock, error) {
	col := s.client.Database(DBName).Collection(CollName)

	var itemStock *ItemStock
	err := col.FindOne(ctx, bson.M{"product_id": pID}).Decode(&itemStock)
	if err != nil {
		return nil, err
	}
	return itemStock, nil
}

func (s *store) Update(ctx context.Context, pID string, delta int64) error {
	col := s.client.Database(DBName).Collection(CollName)

	filter := bson.M{"product_id": pID}

	update := bson.M{
		"$inc": bson.M{"stock": delta},
	}

	if delta < 0 {
		filter["stock"] = bson.M{"$gte": -delta}
	}

	opts := options.Update().SetUpsert(true)

	result, err := col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert stock for product %s: %v", pID, err)
	}

	if result.UpsertedCount > 0 {
		fmt.Printf("New product with ID %s created with initial stock %d\n", pID, delta)
	}

	return nil
}
