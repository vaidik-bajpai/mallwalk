package handler

import (
	pb "github.com/vaidik-bajpai/mallwalk/common/api"
)

func toPBCreateUser(user *UserRegister) *pb.CreateUserRequest {
	return &pb.CreateUserRequest{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}

func toPBLoginUser(user *UserLogin) *pb.UserLoginRequest {
	return &pb.UserLoginRequest{
		Email:    user.Email,
		Password: user.Password,
	}
}

func toPBCreateProductRequest(p Product) *pb.CreateProductRequest {
	return &pb.CreateProductRequest{
		Name:        p.Name,
		Description: p.Description,
		Category:    p.Category,
		Stock:       p.Stock,
		Price:       p.Price,
		Rating:      p.Rating,
		Image:       p.Image,
	}
}

func toPBGetProduct(productID string) *pb.GetProductRequest {
	return &pb.GetProductRequest{
		ProductID: productID,
	}
}

func toPBListProductRequest(pageNumber, pageSize uint32, category string, minRating float32) *pb.ListProductsRequest {
	return &pb.ListProductsRequest{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Category:   category,
		MinRating:  minRating,
	}
}
func toPBUpdateProduct(p Product) *pb.UpdateProductRequest {
	return &pb.UpdateProductRequest{
		ID:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
		Category:    p.Category,
		Price:       p.Price,
		Stock:       p.Stock,
		Rating:      p.Rating,
		Image:       p.Image,
	}
}

func toPBDeleteProductRequest(pID string) *pb.DeleteProductRequest {
	return &pb.DeleteProductRequest{
		ProductID: pID,
	}
}

func toPBAddToCart(cID string, item *CartProduct) *pb.AddToCartRequest {
	return &pb.AddToCartRequest{
		CartID: cID,
		Item: &pb.Item{
			ProductID: item.ProductID,
			Name:      item.Name,
			Image:     item.Image,
			Price:     uint32(item.Price),
			Quantity:  item.Quantity,
		},
	}
}

func toPBRemoveFromCart(pID, cID string) *pb.RemoveItemRequest {
	return &pb.RemoveItemRequest{
		ProductID: pID,
		CartID:    cID,
	}
}

func toPBViewCart(cID string) *pb.ViewCartRequest {
	return &pb.ViewCartRequest{
		CartID: cID,
	}
}
