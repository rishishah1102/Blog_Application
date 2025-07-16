package handlers

import (
	"blog-application/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// UpdateBlogPost modifies a blog post for a given ID

// @Summary Update a blog post
// @Description Update a blog post by ID
// @Tags blog
// @Accept json
// @Produce json
// @Param id path string true "Blog Post ID"
// @Param blog body models.BlogPostRequest true "Blog Post"
// @Success 200 {object} models.UpdateBlogPostSuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/blog-post/{id} [patch]
func (a *API) UpdateBlogPost(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Error("invalid ID format", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var blog models.BlogPost
	if err := c.ShouldBindJSON(&blog); err != nil {
		a.logger.Error("invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":       blog.Title,
			"description": blog.Description,
			"body":        blog.Image,
			"updated_at":  time.Now(),
		},
	}

	result, err := a.db.Collection("blog_posts").UpdateOne(context.Background(), bson.M{"_id": objectID}, update)
	if err != nil {
		a.logger.Error("failed to update blog post", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update blog post"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog post not found"})
		return
	}

	// Fetch updated document
	var updatedBlog models.BlogPost
	err = a.db.Collection("blog_posts").FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&updatedBlog)
	if err != nil {
		a.logger.Error("failed to fetch updated blog post", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated blog post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully updated a blog",
		"blog":    updatedBlog,
		"blogID":  id,
	})
}
