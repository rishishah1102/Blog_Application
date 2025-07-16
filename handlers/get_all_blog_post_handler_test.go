package handlers

import (
	"blog-application/logger"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetAllBlogPosts(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	log := logger.NewLogger()

	testCases := []struct {
		name           string
		setupMocks     func(mt *mtest.T)
		expectedStatus int
		expectedError  string
		expectedCount  int
	}{
		{
			name: "successfully get all blog posts",
			setupMocks: func(mt *mtest.T) {
				firstID := primitive.NewObjectID()
				secondID := primitive.NewObjectID()
				now := primitive.NewDateTimeFromTime(time.Now())

				firstBatch := mtest.CreateCursorResponse(1, "db.blog_posts", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: firstID},
					{Key: "title", Value: "First Blog"},
					{Key: "description", Value: "First Blog Content"},
					{Key: "image", Value: "https://test-1.jpg"},
					{Key: "created_at", Value: now},
					{Key: "updated_at", Value: now},
				})
				secondBatch := mtest.CreateCursorResponse(1, "db.blog_posts", mtest.NextBatch, bson.D{
					{Key: "_id", Value: secondID},
					{Key: "title", Value: "Second Blog"},
					{Key: "description", Value: "Second Blog Content"},
					{Key: "image", Value: "https://test-2.jpg"},
					{Key: "created_at", Value: now},
					{Key: "updated_at", Value: now},
				})

				// Final empty batch to indicate end of cursor
				killCursors := mtest.CreateCursorResponse(0, "db.blog_posts", mtest.NextBatch)

				mt.AddMockResponses(firstBatch, secondBatch, killCursors)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "empty blog posts list",
			setupMocks: func(mt *mtest.T) {
				firstBatch := mtest.CreateCursorResponse(0, "db.blog_posts", mtest.FirstBatch)
				killCursors := mtest.CreateCursorResponse(0, "db.blog_posts", mtest.NextBatch)
				mt.AddMockResponses(firstBatch, killCursors)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "database error during finding blogs",
			setupMocks: func(mt *mtest.T) {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 0},
					{Key: "errorMessage", Value: "find failed"},
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "failed to get blog posts",
		},
		{
			name: "database error during cursor decode blogs",
			setupMocks: func(mt *mtest.T) {
				now := primitive.NewDateTimeFromTime(time.Now())
				firstBatch := mtest.CreateCursorResponse(1, "db.blog_posts", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: "725886"},
					{Key: "title", Value: "First Blog"},
					{Key: "description", Value: "First Blog Content"},
					{Key: "image", Value: "https://test-1.jpg"},
					{Key: "createdAt", Value: now},
					{Key: "updatedAt", Value: now},
				})
				killCursors := mtest.CreateCursorResponse(0, "db.blog_posts", mtest.NextBatch)
				mt.AddMockResponses(firstBatch, killCursors)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "failed to decode blog posts",
		},
	}

	for _, tt := range testCases {
		mt.Run(tt.name, func(mt *mtest.T) {
			tt.setupMocks(mt)

			api := API{
				db:     mt.DB,
				logger: log,
			}

			gin.SetMode(gin.TestMode)
			router := gin.Default()
			router.GET("/api/blog-post", api.GetAllBlogPosts)

			req, err := http.NewRequest("GET", "/api/blog-post", nil)
			assert.NoError(mt, err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(mt, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(mt, err)

			if tt.expectedError != "" {
				assert.Contains(mt, response["error"], tt.expectedError)
			} else {
				assert.Equal(mt, "successfully fetched all the blogs", response["message"])
				blogs := response["blogs"].([]interface{})
				assert.Equal(mt, tt.expectedCount, len(blogs))
			}
		})
	}
}
