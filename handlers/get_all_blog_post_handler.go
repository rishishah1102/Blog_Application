package handlers

import (
	"blog-application/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// GetAllBlogPosts fetches all the blog posts

// @Summary Get all blog posts
// @Description Get all blog posts
// @Tags blog
// @Accept json
// @Produce json
// @Success 200 {object} models.GetAllBlogPostsSuccessResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/blog-post [get]
func (a *API) GetAllBlogPosts(c *gin.Context) {
	var blogs = make([]models.BlogPost, 0)

	cursor, err := a.db.Collection("blog_posts").Find(context.Background(), bson.M{})
	if err != nil {
		a.logger.Error("failed to get blog posts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get blog posts"})
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &blogs); err != nil {
		a.logger.Error("failed to decode blog posts", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode blog posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successfully fetched all the blogs",
		"blogs":   blogs,
	})
}
