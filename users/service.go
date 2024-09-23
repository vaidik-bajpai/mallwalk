package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pascaldekloe/jwt"
	pb "github.com/vaidik-bajpai/mallwalk/common/api"
)

type service struct {
	store     UsersStore
	jwtSecret string
}

func NewService(store UsersStore, jwtSecret string) *service {
	return &service{
		store:     store,
		jwtSecret: jwtSecret,
	}
}

func (s *service) CreateUser(ctx context.Context, userDetails *pb.CreateUserRequest) (*pb.User, error) {
	fmt.Println("Create User grpc handler")
	hashedPassword, err := HashPassword(userDetails.Password)
	if err != nil {
		return nil, err
	}

	err = s.store.Create(ctx, &User{
		Username:     userDetails.Username,
		Email:        userDetails.Email,
		HashPassword: string(hashedPassword),
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *service) UserLogin(ctx context.Context, userDetails *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	user, err := s.store.GetByEmail(ctx, userDetails.Email)
	if err != nil {
		return nil, err
	}

	fmt.Println("user:", user)

	//verify the email and password
	ok, err := MatchPassword(userDetails.Password, user.HashPassword)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &pb.UserLoginResponse{}, fmt.Errorf("invalid email or password")
	}

	var claims jwt.Claims
	claims.Subject = user.Email
	claims.Issuer = "vaidik bajpai"
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Audiences = []string{"mallwalk audience"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(s.jwtSecret))
	if err != nil {
		return &pb.UserLoginResponse{}, err
	}

	return &pb.UserLoginResponse{
		Token: string(jwtBytes),
	}, nil
}
