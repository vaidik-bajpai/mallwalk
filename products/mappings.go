package main

import (
	pb "github.com/vaidik-bajpai/mallwalk/common/api"
)

func toPBUpdateProduct(pID string, p *UpdateProduct) *pb.Product {
	return &pb.Product{
		ID:          pID,
		Name:        p.Name,
		Description: p.Description,
		Category:    p.Category,
		Rating:      p.Rating,
		Image:       p.Image,
		Price:       p.Price,
	}
}

func toPBProduct(p *Product) *pb.Product {
	return &pb.Product{
		ID:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		Category:    p.Category,
		Rating:      p.Rating,
		Image:       p.Image,
		Price:       p.Price,
	}
}

func toPBProductSummary(p *ProductSummary) *pb.ProductSummary {
	return &pb.ProductSummary{
		ID:     p.ID.String(),
		Name:   p.Name,
		Rating: p.Rating,
		Image:  p.Image,
		Price:  p.Price,
	}
}
