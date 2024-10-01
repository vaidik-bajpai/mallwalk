package main

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	client *mongo.Client
}

const (
	DBName   = "orders"
	CollName = "orders"
)

func NewStore(client *mongo.Client) *store {
	return &store{client: client}
}
