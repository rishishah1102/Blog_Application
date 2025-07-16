package handlers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	_ "blog-application/docs"
)

// API is the struct for all the handlers
type API struct {
	logger *zap.Logger
	db     *mongo.Database
}

// NewAPI creates a new API instance
func NewAPI(logger *zap.Logger, db *mongo.Database) *API {
	return &API{
		logger: logger,
		db:     db,
	}
}

// RegisterRoutes register all the handlers to gin router
func (a *API) RegisterRoutes(router *gin.Engine) {
	// Server Check
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome to the blog API",
		})
	})

	// Blog routes
	blogRoutes := router.Group("/api/blog-post")
	{
		blogRoutes.POST("/", a.CreateBlogPost)
		blogRoutes.GET("/", a.GetAllBlogPosts)
		blogRoutes.GET("/:id", a.GetBlogPostByID)
		blogRoutes.PATCH("/:id", a.UpdateBlogPost)
		blogRoutes.DELETE("/:id", a.DeleteBlogPost)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
