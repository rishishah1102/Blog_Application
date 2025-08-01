package handlers

import (
	"blog-application/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// CreateBlogPost creates a blog post

// @Summary Create a new blog post
// @Description Create a new blog post
// @Tags blog
// @Accept json
// @Produce json
// @Param blog body models.BlogPostRequest true "Blog Post"
// @Success 201 {object} models.CreateBlogPostSuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/blog-post [post]
func (a *API) CreateBlogPost(c *gin.Context) {
	var blog models.BlogPost
	if err := c.ShouldBindJSON(&blog); err != nil {
		a.logger.Error("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()

	result, err := a.db.Collection("blog_posts").InsertOne(context.Background(), blog)
	if err != nil {
		a.logger.Error("failed to create blog post", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create blog post"})
		return
	}

	blog.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "successfully created the blog",
		"blog":    blog,
	})
}
