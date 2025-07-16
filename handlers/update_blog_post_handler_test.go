package handlers

import (
	"blog-application/logger"
	"blog-application/models"
	"bytes"
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

func TestUpdateBlogPost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	log := logger.NewLogger()

	testCases := []struct {
		name           string
		blogID         string
		requestBody    interface{}
		setupMocks     func(mt *mtest.T)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successfully update blog post",
			blogID: primitive.NewObjectID().Hex(),
			requestBody: models.BlogPost{
				Title:       "Updated Title",
				Description: "Updated Description",
				Image:       "https://updated.jpg",
			},
			setupMocks: func(mt *mtest.T) {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 1},
					{Key: "nModified", Value: 1},
					{Key: "n", Value: 1},
				})

				// Mock find response for updated document
				objectID, _ := primitive.ObjectIDFromHex(primitive.NewObjectID().Hex())
				mt.AddMockResponses(mtest.CreateCursorResponse(1, "db.blog_posts", mtest.FirstBatch, bson.D{
					{Key: "_id", Value: objectID},
					{Key: "title", Value: "Updated Title"},
					{Key: "description", Value: "Updated Description"},
					{Key: "image", Value: "https://updated.jpg"},
					{Key: "created_att", Value: primitive.NewDateTimeFromTime(time.Now())},
					{Key: "updated_at", Value: primitive.NewDateTimeFromTime(time.Now())},
				}))
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID format",
			blogID:         "554454533",
			requestBody:    models.BlogPost{},
			setupMocks:     func(mt *mtest.T) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid ID format",
		},
		{
			name:   "invalid request body",
			blogID: primitive.NewObjectID().Hex(),
			requestBody: models.BlogPost{
				Description: "Updated Description",
				Image:       "https://updated.jpg",
			},
			setupMocks:     func(mt *mtest.T) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Key: 'BlogPost.Title' Error:Field validation for 'Title' failed on the 'required' tag",
		},
		{
			name:   "blog post not found",
			blogID: primitive.NewObjectID().Hex(),
			requestBody: models.BlogPost{
				Title:       "Non-existent",
				Description: "Should fail",
				Image:       "http://test.jpg",
			},
			setupMocks: func(mt *mtest.T) {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 1},
					{Key: "nModified", Value: 0},
					{Key: "n", Value: 0},
				})
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Blog post not found",
		},
		{
			name:   "database update error",
			blogID: primitive.NewObjectID().Hex(),
			requestBody: models.BlogPost{
				Title:       "Error Case",
				Description: "Should fail",
				Image:       "error.jpg",
			},
			setupMocks: func(mt *mtest.T) {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 0},
					{Key: "errmsg", Value: "update failed"},
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "failed to update blog post",
		},
		{
			name:   "database find error after update",
			blogID: primitive.NewObjectID().Hex(),
			requestBody: models.BlogPost{
				Title:       "Find Error",
				Description: "Should fail on find",
				Image:       "https://finderror.jpg",
			},
			setupMocks: func(mt *mtest.T) {
				// Successful update
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 1},
					{Key: "nModified", Value: 1},
					{Key: "n", Value: 1},
				})

				// Failed find
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 0},
					{Key: "errmsg", Value: "find failed"},
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "failed to fetch updated blog post",
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
			router.PATCH("/api/blog-post/:id", api.UpdateBlogPost)

			reqBody, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)
			req, err := http.NewRequest("PATCH", "/api/blog-post/"+tt.blogID, bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError != "" {
				assert.Contains(t, response["error"], tt.expectedError)
			} else {
				assert.Equal(t, "successfully updated a blog", response["message"])
				assert.Equal(t, tt.blogID, response["blogID"])
				assert.NotNil(mt, response["blog"])
			}
		})
	}
}
