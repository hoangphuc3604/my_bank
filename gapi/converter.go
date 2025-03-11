package gapi

import (
	"github.com/hoangphuc3064/MyBank/db/sqlc"
	"github.com/hoangphuc3064/MyBank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user sqlc.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		Fullname: 		user.Fullname,
		Email: 		user.Email,
		CreatedAt: 	timestamppb.New(user.CreatedAt),
	}
}