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

func TestGetBlogPostByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	log := logger.NewLogger()

	testCases := []struct {
		name           string
		blogID         string
		setupMocks     func(mt *mtest.T)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successfully get blog post by ID",
			blogID: primitive.NewObjectID().Hex(),
			setupMocks: func(mt *mtest.T) {
				objectID, _ := primitive.ObjectIDFromHex(primitive.NewObjectID().Hex())
				now := primitive.NewDateTimeFromTime(time.Now())
				mt.AddMockResponses(mtest.CreateCursorResponse(1, "db.blog_posts", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: objectID},
					{Key: "title", Value: "Test Blog"},
					{Key: "description", Value: "Test Content"},
					{Key: "image", Value: "https://test.jpg"},
					{Key: "createdAt", Value: now},
					{Key: "updatedAt", Value: now},
				}))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID format",
			blogID:         "4543664",
			setupMocks:     func(mt *mtest.T) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid ID format",
		},
		{
			name:   "blog post not found",
			blogID: primitive.NewObjectID().Hex(),
			setupMocks: func(mt *mtest.T) {
				mt.AddMockResponses(mtest.CreateCursorResponse(0, "db.blog_posts", mtest.FirstBatch))
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "blog post not found",
		},
		{
			name:   "database error in finding a blog",
			blogID: primitive.NewObjectID().Hex(),
			setupMocks: func(mt *mtest.T) {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 0},
					{Key: "errmsg", Value: "connection failed"},
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "failed to get blog post",
		},
		{
			name:   "malformed document data",
			blogID: primitive.NewObjectID().Hex(),
			setupMocks: func(mt *mtest.T) {
				now := primitive.NewDateTimeFromTime(time.Now())
				mt.AddMockResponses(mtest.CreateCursorResponse(1, "db.blog_posts", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: "4534331"},
					{Key: "title", Value: "Test Blog"},
					{Key: "description", Value: "Test Content"},
					{Key: "image", Value: "https://test.jpg"},
					{Key: "createdAt", Value: now},
					{Key: "updatedAt", Value: now},
				}))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "failed to get blog post",
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
			router.GET("/api/blog-post/:id", api.GetBlogPostByID)

			req, err := http.NewRequest("GET", "/api/blog-post/"+tt.blogID, nil)
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
				assert.Equal(mt, "successfully fetched a blog", response["message"])
				assert.Equal(mt, tt.blogID, response["blogID"])
				assert.NotNil(mt, response["blog"])
			}
		})
	}
}
