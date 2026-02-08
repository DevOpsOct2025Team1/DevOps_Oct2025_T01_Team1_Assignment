package server

import (
	"github.com/gin-contrib/cors"
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
	fileClient handlers.FileServiceClient
}

func New(authClient handlers.AuthServiceClient, userClient handlers.UserServiceClient, fileClient handlers.FileServiceClient, cfg *config.Config) *Server {
	router := gin.Default()
	router.Use(otelgin.Middleware("api-gateway"))

	allowOrigins := []string{"http://localhost:5173"}
	if cfg.FrontendURL != "" {
		allowOrigins = append(allowOrigins, cfg.FrontendURL)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	s := &Server{
		Router:     router,
		authClient: authClient,
		userClient: userClient,
		fileClient: fileClient,
	}

	s.setupRoutes(cfg)
	return s
}

func (s *Server) setupRoutes(cfg *config.Config) {
	authHandler := handlers.NewAuthHandler(s.authClient)
	userHandler := handlers.NewUserHandler(s.userClient)
	fileHandler := handlers.NewFileHandler(s.fileClient)

	s.Router.POST("/api/login", authHandler.Login)

	s.Router.POST("/api/admin/create_user", middleware.ValidateRole(s.authClient, []userv1.Role{userv1.Role_ROLE_ADMIN}), authHandler.SignUp)
	s.Router.DELETE("/api/admin/delete_user", middleware.ValidateRole(s.authClient, []userv1.Role{userv1.Role_ROLE_ADMIN}), userHandler.DeleteUser)
	s.Router.GET("/api/admin/list_users", middleware.ValidateRole(s.authClient, []userv1.Role{userv1.Role_ROLE_ADMIN}), userHandler.ListUsers)

	files := s.Router.Group("/api/files")
	files.Use(middleware.ValidateRole(s.authClient, []userv1.Role{userv1.Role_ROLE_USER, userv1.Role_ROLE_ADMIN}))
	{
		files.GET("", fileHandler.ListFiles)
		files.POST("", fileHandler.UploadFile)
		files.GET("/:id", fileHandler.GetFile)
		files.GET("/:id/download", fileHandler.DownloadFile)
		files.DELETE("/:id", fileHandler.DeleteFile)
	}

	s.Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	if cfg != nil && cfg.Environment == "development" {
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

	if s.fileClient != nil {
		if e := s.fileClient.Close(); e != nil && err == nil {
			err = e
		}
	}

	return err
}
