package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/config"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/internal/handlers"
	"github.com/provsalt/DOP_P01_Team1/api-gateway/middleware"
	userv1 "github.com/provsalt/DOP_P01_Team1/common/user/v1"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Server struct {
	Router     *gin.Engine
	authClient handlers.AuthServiceClient
	userClient handlers.UserServiceClient
}

func New(authClient handlers.AuthServiceClient, userClient handlers.UserServiceClient) *Server {
	router := gin.Default()
	router.Use(otelgin.Middleware("api-gateway"))

	s := &Server{
		Router:     router,
		authClient: authClient,
		userClient: userClient,
	}

	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	authHandler := handlers.NewAuthHandler(s.authClient)
	userHandler := handlers.NewUserHandler(s.userClient)

	s.Router.POST("/api/login", authHandler.Login)

	s.Router.POST("/api/admin/create_user", middleware.ValidateRole(s.authClient, []userv1.Role{userv1.Role_ROLE_ADMIN}), authHandler.SignUp)
	s.Router.DELETE("/api/admin/delete_user", middleware.ValidateRole(s.authClient, []userv1.Role{userv1.Role_ROLE_ADMIN}), userHandler.DeleteUser)

	s.Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Environment == "development" {
		s.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func (s *Server) Run(addr string) error {
	return s.Router.Run(addr)
}

func (s *Server) Close() error {
	var err error

	if s.authClient != nil {
		if e := s.authClient.Close(); e != nil && err == nil {
			err = e
		}
	}

	if s.userClient != nil {
		if e := s.userClient.Close(); e != nil && err == nil {
			err = e
		}
	}

	return err
}
