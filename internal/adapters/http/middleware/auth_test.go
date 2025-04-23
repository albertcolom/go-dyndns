package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	data := []struct {
		name           string
		token          string
		header         string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Valid query token",
			token:          "valid_token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid header token",
			header:         "Bearer valid_token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid query token",
			token:          "invalid",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "Invalid header token",
			header:         "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "No token provided",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:           "Malformed header",
			header:         "InvalidHeader",
			expectedStatus: http.StatusUnauthorized,
			expectError:    true,
		},
	}

	for _, d := range data {
		gin.SetMode(gin.TestMode)

		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, router := gin.CreateTestContext(w)
			router.Use(AuthMiddleware("valid_token"))

			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			q := req.URL.Query()
			if d.token != "" {
				q.Add("token", d.token)
			}
			req.URL.RawQuery = q.Encode()

			if d.header != "" {
				req.Header.Set("Authorization", d.header)
			}

			ctx.Request = req
			router.ServeHTTP(w, req)

			assert.Equal(t, d.expectedStatus, w.Code)
			if d.expectError {
				assert.Contains(t, w.Body.String(), "Unauthorized")
			}
		})
	}
}
