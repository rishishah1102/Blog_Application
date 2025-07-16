package handlers

import (
	"blog-application/logger"
	"blog-application/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCreateBlogPost(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	log := logger.NewLogger()

	testCases := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(mt *mtest.T)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful blog post creation",
			requestBody: models.BlogPost{
				Title:       "Test Blog",
				Description: "This is a test blog post",
				Image:       "https://test.jpg",
			},
			setupMocks: func(mt *mtest.T) {
				blogID := primitive.NewObjectID()
				mt.AddMockResponses(mtest.CreateSuccessResponse(
					bson.D{
						{Key: "ok", Value: 1},
						{Key: "n", Value: 1},
						{Key: "insertedID", Value: blogID},
					}...),
				)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "missing requireds field",
			requestBody: models.BlogPost{
				Description: "Blog without title",
				Image:       "https://test.jpg",
			},
			setupMocks:     func(mt *mtest.T) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Key: 'BlogPost.Title' Error:Field validation for 'Title' failed on the 'required' tag",
		},
		{
			name: "database insertion failure",
			requestBody: models.BlogPost{
				Title:       "Database Error",
				Description: "This should fail",
				Image:       "https://error.jpg",
			},
			setupMocks: func(mt *mtest.T) {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 0},
					{Key: "errorMessage", Value: "insert failed"},
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "failed to create blog post",
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
			router.POST("/api/blog-post", api.CreateBlogPost)

			reqBody, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)
			req, err := http.NewRequest("POST", "/api/blog-post", bytes.NewBuffer(reqBody))
			assert.NoError(mt, err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(mt, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(mt, err)

			if tt.expectedError != "" {
				assert.Equal(mt, tt.expectedError, response["error"])
			} else {
				assert.Equal(mt, "successfully created the blog", response["message"])
				assert.NotNil(mt, response["blog"])
			}
		})
	}
}
