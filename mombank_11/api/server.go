package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/hippo-an/tiny-go-challenges/mombank_11/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.listAccount)
	router.GET("/accounts/:id", server.getAccount)
	server.router = router
	return server
}
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
