package gapi

import (
	"context"

	"github.com/hoangphuc3064/MyBank/db/sqlc"
	"github.com/hoangphuc3064/MyBank/pb"
	"github.com/hoangphuc3064/MyBank/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (response *pb.CreateUserResponse, err error) {
	hashPassword, err := util.HashPassword(request.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot hash password: %s", err)
	}

	arg := sqlc.CreateUserParams{
		Username: request.GetUsername(),
		Password: hashPassword,
		Fullname: request.GetFullname(),
		Email: request.GetEmail(),
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user already exists: %s", arg.Username)
			default:
				return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
			}
		}

		return nil, status.Errorf(codes.Internal, "cannot create user: %v", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	
	return rsp, nil
}