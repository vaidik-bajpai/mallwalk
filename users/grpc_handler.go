package main

import (
	"context"

	pb "github.com/vaidik-bajpai/mallwalk/common/api"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedUserServiceServer

	service UsersService
}

func NewGRPCHandler(grpcServer *grpc.Server, service UsersService) {
	handler := &grpcHandler{
		service: service,
	}
	pb.RegisterUserServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateUser(ctx context.Context, userDetails *pb.CreateUserRequest) (*pb.User, error) {
	_, err := h.service.CreateUser(ctx, userDetails)
	if err != nil {
		return nil, err
	}

	return &pb.User{}, nil
}

func (h *grpcHandler) UserLogin(ctx context.Context, userDetails *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	return h.service.UserLogin(ctx, userDetails)
}

func HashPassword(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func MatchPassword(plaintextPassword, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	if err != nil {
		return false, err
	}

	return true, nil
}
