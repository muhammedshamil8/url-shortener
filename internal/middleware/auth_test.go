package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/auth"
	"github.com/muhammedshamil8/url-shortener/internal/config"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			AccessTokenSecret: "test_secret_key",
			AccessTokenExpiry: "1h",
		},
	}

	tests := []struct {
		name           string
		setupHeader    func(req *http.Request)
		expectedStatus int
		expectedEmail  string
	}{
		{
			name: "Missing Authorization Header",
			setupHeader: func(req *http.Request) {
				// No header set
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid Header Format (No Bearer)",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "invalid_token_value")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Short Header Value (Panic Prevention)",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "Bear")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid Token Signature",
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalid.token.signature")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Valid Token",
			setupHeader: func(req *http.Request) {
				token, err := auth.GenerateToken(42, "user@example.com", cfg.JWT.AccessTokenSecret, cfg.JWT.AccessTokenExpiry)
				if err != nil {
					t.Fatalf("failed to generate token: %v", err)
				}
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
			expectedEmail:  "user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(AuthMiddleware(cfg))

			var actualEmail string
			var actualUserID any

			r.GET("/test", func(c *gin.Context) {
				actualEmail = c.GetString("email")
				actualUserID = c.Value("user_id")
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			tt.setupHeader(req)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				if actualEmail != tt.expectedEmail {
					t.Errorf("got email %q, want %q", actualEmail, tt.expectedEmail)
				}
				if actualUserID == nil {
					t.Errorf("expected user_id to be set in context")
				}
			}
		})
	}
}
