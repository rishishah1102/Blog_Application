package handlers

import (
	"blog-application/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// GetBlogPostByID fetches a blog for a given ID

// @Summary Get a blog post by ID
// @Description Get a single blog post by ID
// @Tags blog
// @Accept json
// @Produce json
// @Param id path string true "Blog Post ID"
// @Success 200 {object} models.GetBlogPostByIDSuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/blog-post/{id} [get]
func (a *API) GetBlogPostByID(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Error("invalid ID format", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var blog models.BlogPost
	err = a.db.Collection("blog_posts").FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&blog)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "blog post not found"})
			return
		}
		a.logger.Error("failed to get blog post", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get blog post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully fetched a blog",
		"blog":    blog,
		"blogID":  id,
	})
}
