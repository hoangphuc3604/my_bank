package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/hoangphuc3064/MyBank/api"
	"github.com/hoangphuc3064/MyBank/db/sqlc"
	"github.com/hoangphuc3064/MyBank/gapi"
	"github.com/hoangphuc3064/MyBank/pb"
	"github.com/hoangphuc3064/MyBank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	store := sqlc.NewStore(testDB)
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store *sqlc.Store) {
	server, err := gapi.NewServer(config, *store)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterServiceMyBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Cannot create listener:", err)
	}
	
	log.Printf("Start GRPC server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Cannot start GRPC server:", err)
	}
}

func runGinServer(config util.Config, store *sqlc.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("Cannot start HTTP server:", err)
	}
}
