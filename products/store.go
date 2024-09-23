package main

import (
	"context"
	"time"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type store struct {
	client *mongo.Client
}

const (
	DBName   = "products"
	CollName = "products"
)

func NewStore(client *mongo.Client) *store {
	return &store{client: client}
}

func (s *store) Create(ctx context.Context, cp *Product) (*Product, error) {
	col := s.client.Database(DBName).Collection(CollName)

	cp.CreatedAt = time.Now()
	cp.ID = primitive.NewObjectID()
	res, err := col.InsertOne(ctx, cp)
	if err != nil {
		return nil, err
	}

	id := res.InsertedID.(primitive.ObjectID)

	cp.ID = id

	return cp, nil
}

func (s *store) Get(ctx context.Context, id primitive.ObjectID) (*Product, error) {
	col := s.client.Database(DBName).Collection(CollName)

	var p Product
	err := col.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *store) List(ctx context.Context, lp *pb.ListProductsRequest) ([]*pb.ProductSummary, error) {
	col := s.client.Database(DBName).Collection(CollName)

	skip := (lp.PageNumber - 1) * lp.PageSize

	fOpts := options.Find()
	fOpts.SetLimit(int64(lp.PageSize))
	fOpts.SetSkip(int64(skip))

	filter := bson.M{}

	if lp.MinRating > 0 {
		filter["rating"] = lp.MinRating
	}

	if lp.Category != "" {
		filter["category"] = lp.Category
	}

	cursor, err := col.Find(ctx, filter, fOpts)
	if err != nil {
		return nil, err
	}

	var products []*pb.ProductSummary
	for cursor.Next(ctx) {
		var product ProductSummary
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, toPBProductSummary(&product))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *store) Update(ctx context.Context, up *pb.UpdateProductRequest) (*Product, error) {
	col := s.client.Database(DBName).Collection(CollName)

	p := toProduct(up)

	id, err := primitive.ObjectIDFromHex(up.ID)
	if err != nil {
		return nil, err
	}

	_, err = col.UpdateByID(ctx, bson.M{"_id": id}, bson.M{
		"$set": p,
	})
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *store) Delete(ctx context.Context, id primitive.ObjectID) error {
	col := s.client.Database(DBName).Collection(CollName)
	_, err := col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func toProduct(up *pb.UpdateProductRequest) *Product {
	return &Product{
		Name:        up.Name,
		Price:       up.Price,
		Description: up.Description,
		Rating:      up.Rating,
		Stock:       up.Stock,
		Category:    up.Category,
		Image:       up.Image,
		UpdatedAt:   time.Now(),
	}
}

func toPBProductSummary(p *ProductSummary) *pb.ProductSummary {
	return &pb.ProductSummary{
		ID:     p.ID.String(),
		Name:   p.Name,
		Rating: p.Rating,
		Image:  p.Image,
		Price:  p.Price,
		Stock:  p.Stock,
	}
}
