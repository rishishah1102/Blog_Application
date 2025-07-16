package handlers

import (
	"blog-application/logger"
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

func TestDeleteBlogPost(t *testing.T) {
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
			name:   "successful deletion",
			blogID: primitive.NewObjectID().Hex(),
			setupMocks: func(mt *mtest.T) {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 1},
					{Key: "n", Value: 1},
				})
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID format",
			blogID:         "346343454",
			setupMocks:     func(mt *mtest.T) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid ID format",
		},
		{
			name:   "blog post not found",
			blogID: primitive.NewObjectID().Hex(),
			setupMocks: func(mt *mtest.T) {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 1},
					{Key: "n", Value: 0},
				})
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "blog post not found",
		},
		{
			name:   "database deletion failure",
			blogID: primitive.NewObjectID().Hex(),
			setupMocks: func(mt *mtest.T) {
				mt.AddMockResponses(bson.D{
					{Key: "ok", Value: 0},
					{Key: "errorMessage", Value: "delete failed"},
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "failed to delete blog post",
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
			router.DELETE("/api/blog-post/:id", api.DeleteBlogPost)

			req, err := http.NewRequest("DELETE", "/api/blog-post/"+tt.blogID, nil)
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
				assert.Equal(mt, "blog post deleted successfully", response["message"])
				assert.Equal(mt, tt.blogID, response["blogID"])
			}
		})
	}
}
