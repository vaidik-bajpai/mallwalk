package main

import pb "github.com/vaidik-bajpai/mallwalk/common/api"

func toPBItem(it *Item) *pb.Item {
	return &pb.Item{
		ProductID: it.ProductID,
		Name:      it.Name,
		Image:     it.Image,
		Price:     it.Price,
		Quantity:  it.Quantity,
	}
}
