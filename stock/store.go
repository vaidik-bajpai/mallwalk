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
	err := col.FindOne(ctx, bson.M{"_id": pID}).Decode(&itemStock)
	if err != nil {
		return nil, err
	}
	return itemStock, nil
}

func (s *store) Update(ctx context.Context, pID string, delta int64) error {
	// Reference the MongoDB collection
	col := s.client.Database(DBName).Collection(CollName)

	// Create a filter to find the document with the specified product ID
	filter := bson.M{"product_id": pID}

	// Define the update operation, incrementing the stock by delta
	update := bson.M{
		"$inc": bson.M{"stock": delta},
	}

	// Ensure that stock cannot fall below 0 (this condition applies only if stock exists)
	if delta < 0 {
		filter["stock"] = bson.M{"$gte": -delta} // Ensure stock is at least enough to handle the decrement
	}

	// Set the upsert option to true so it inserts a new document if no match is found
	opts := options.Update().SetUpsert(true)

	// Perform the upsert operation
	result, err := col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to upsert stock for product %s: %v", pID, err)
	}

	// If a new document was inserted
	if result.UpsertedCount > 0 {
		fmt.Printf("New product with ID %s created with initial stock %d\n", pID, delta)
	}

	return nil
}
