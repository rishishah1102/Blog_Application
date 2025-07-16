package server

import (
	"blog-application/config"
	"blog-application/database"
	"blog-application/handlers"
	"blog-application/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server is the struct blog-application server
type Server struct {
	cfg    *config.Config
	logger *zap.Logger
	Router *gin.Engine
}

// NewServer creates a new server for the application
func NewServer(cfg *config.Config) (*Server, error) {
	// Initialize logger
	log := logger.NewLogger()

	// Initialize MongoDB
	mongoClient, err := database.NewMongoClient(cfg.MongoDB.MongoURI, cfg.MongoDB.Timeout)
	if err != nil {
		return nil, logger.WrapError(err, "failed to connect to MongoDB")
	}
	db := mongoClient.Database(cfg.MongoDB.DatabaseName)
	log.Info("successfully connected to db")

	// Create Gin router
	router := gin.Default()
	router.Use(gin.Recovery())

	// Initialize API
	api := handlers.NewAPI(log, db)
	api.RegisterRoutes(router)

	return &Server{
		cfg:    cfg,
		logger: log,
		Router: router,
	}, nil
}

// Start starts the server on default port
func (s *Server) Start() error {
	s.logger.Info("Starting server", zap.String("port", s.cfg.Server.Port))
	return s.Router.Run(":" + s.cfg.Server.Port)
}
