package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// DeleteBlogPost deletes a blog post for a given ID

// @Summary Delete a blog post
// @Description Delete a blog post by ID
// @Tags blog
// @Accept json
// @Produce json
// @Param id path string true "Blog Post ID"
// @Success 200 {object} models.DeleteBlogPostSuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/blog-post/{id} [delete]
func (a *API) DeleteBlogPost(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Error("invalid ID format", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	result, err := a.db.Collection("blog_posts").DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		a.logger.Error("failed to delete blog post", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete blog post"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "blog post deleted successfully",
		"blogID":  id,
	})
}
