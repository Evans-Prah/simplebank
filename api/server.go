package api

import (
	"fmt"

	db "github.com/Evans-Prah/simplebank/db/sqlc"
	"github.com/Evans-Prah/simplebank/db/util"
	"github.com/Evans-Prah/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for simple bank service
type Server struct {
	store db.Store
	router *gin.Engine
	tokenMaker token.Maker
	config util.Config
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store: store,
		tokenMaker: tokenMaker,
		config: config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
		v.RegisterValidation("customUsername", validateCustomUsername)

	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.POST("/api/users", server.createUser)
	router.POST("/api/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/api/accounts", server.createAccount)
	authRoutes.GET("/api/accounts/:id", server.getAccount)
	authRoutes.GET("/api/accounts", server.getAccounts)
	authRoutes.PATCH("/api/accounts/:id", server.updateAccount)
	authRoutes.DELETE("/api/accounts/:id", server.deleteAccount)

	authRoutes.POST("/api/transfers", server.createTransfer)

	server.router = router

}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}