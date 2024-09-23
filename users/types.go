package main

import (
	"context"
	"time"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
)

type UsersService interface {
	CreateUser(context.Context, *pb.CreateUserRequest) (*pb.User, error)
	UserLogin(context.Context, *pb.UserLoginRequest) (*pb.UserLoginResponse, error)
}

type UsersStore interface {
	Create(context.Context, *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type User struct {
	Username     string    `json:"username" bson:"username"`
	Email        string    `json:"email" bson:"email"`
	HashPassword string    `json:"password" bson:"password"`
	CreatedAt    time.Time `json:"-" bson:"created_at"`
}
