package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	client *mongo.Client
}

const (
	DBName   = "users"
	CollName = "users"
)

func NewStore(client *mongo.Client) *store {
	return &store{client: client}
}

func (s *store) Create(ctx context.Context, user *User) error {
	fmt.Println("Create store")
	col := s.client.Database(DBName).Collection(CollName)

	res, err := col.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	id := res.InsertedID.(primitive.ObjectID)
	fmt.Println("inserted ID:", id)

	return nil
}

func (s *store) GetByEmail(ctx context.Context, email string) (*User, error) {
	col := s.client.Database(DBName).Collection(CollName)

	var user User
	err := col.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
