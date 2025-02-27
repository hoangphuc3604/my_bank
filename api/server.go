package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/hoangphuc3064/MyBank/db/sqlc"
)

type Server struct {
	store *sqlc.Store
	router *gin.Engine
}

func NewServer(store *sqlc.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	
	router.POST("/transfers", server.createTransfer)

	router.POST("/users", server.createUser)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}