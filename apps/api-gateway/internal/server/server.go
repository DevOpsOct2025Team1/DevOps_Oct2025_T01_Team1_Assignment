package server

import (
	"github.com/gin-gonic/gin"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/handlers"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/middleware"
)

type Server struct {
	Router     *gin.Engine
	authClient handlers.AuthServiceClient
}

func New(authClient handlers.AuthServiceClient) *Server {
	router := gin.Default()

	s := &Server{
		Router:     router,
		authClient: authClient,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	authHandler := handlers.NewAuthHandler(s.authClient)

	s.Router.POST("/api/login", authHandler.Login)
	s.Router.POST("/api/admin/create_user", middleware.ValidateRole(s.authClient, []string{"admin"}), authHandler.SignUp)

	s.Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}

func (s *Server) Close() error {
	if s.authClient != nil {
		return s.authClient.Close()
	}
	return nil
}
