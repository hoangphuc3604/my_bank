package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hoangphuc3064/MyBank/db/sqlc"
)

type Server struct {
	store *sqlc.Store
	router *gin.Engine
}

func NewServer(store *sqlc.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}