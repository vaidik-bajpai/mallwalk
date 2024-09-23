package main

import (
	"context"
	"time"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductsService interface {
	CreateProduct(context.Context, *pb.CreateProductRequest) (*pb.Product, error)
	GetProduct(context.Context, *pb.GetProductRequest) (*pb.Product, error)
	ListProduct(context.Context, *pb.ListProductsRequest) (*pb.ListProductsResponse, error)
	UpdateProduct(context.Context, *pb.UpdateProductRequest) (*pb.Product, error)
	DeleteProduct(context.Context, *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error)
}

type ProductsStore interface {
	Create(ctx context.Context, cp *Product) (*Product, error)
	Get(ctx context.Context, id primitive.ObjectID) (*Product, error)
	List(ctx context.Context, lp *pb.ListProductsRequest) ([]*pb.ProductSummary, error)
	Update(ctx context.Context, up *pb.UpdateProductRequest) (*Product, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type Product struct {
	ID          primitive.ObjectID `json:"-" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Price       uint32             `json:"price" bson:"price"`
	Stock       uint32             `json:"stock" bson:"stock"`
	Description string             `json:"description" bson:"description"`
	Category    string             `json:"category" bson:"category"`
	Rating      float32            `json:"rating" bson:"rating"`
	Image       string             `json:"image" bson:"image"`
	CreatedAt   time.Time          `json:"-" bson:"created_at"`
	UpdatedAt   time.Time          `json:"-" bson:"updated_at"`
}

type ProductSummary struct {
	ID     primitive.ObjectID `bson:"_id"`
	Name   string             `bson:"name"`
	Price  uint32             `bson:"price"`
	Stock  uint32             `bson:"stock"`
	Rating float32            `bson:"rating"`
	Image  string             `bson:"image"`
}
