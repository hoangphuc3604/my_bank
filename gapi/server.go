package gapi

import (
	"fmt"

	"github.com/hoangphuc3064/MyBank/db/sqlc"
	"github.com/hoangphuc3064/MyBank/pb"
	"github.com/hoangphuc3064/MyBank/token"
	"github.com/hoangphuc3064/MyBank/util"
)

type Server struct {
	pb.UnimplementedServiceMyBankServer
	config     util.Config
	store      sqlc.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store sqlc.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}