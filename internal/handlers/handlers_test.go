package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/auth"
	"github.com/muhammedshamil8/url-shortener/internal/cache"
	"github.com/muhammedshamil8/url-shortener/internal/config"
	"github.com/muhammedshamil8/url-shortener/internal/models"
)

func setupTestRouter(repo Repository, cache cache.Cache) *gin.Engine {
	r := gin.New()
	h := New(repo, cache, config.Config{})

	r.GET("/api/v1/live", h.LiveHandler)
	r.GET("/api/v1/ready", h.ReadyHandler)
	r.POST("/api/v1/shorten", h.ShortenHandler)
	r.GET("/:code", h.RedirectHandler)
	r.GET("/api/v1/admin/urls", h.ListAllHandler)
	r.DELETE("/api/v1/admin/urls/:id", h.DeleteHandler)
	r.GET("/api/v1/admin/users", h.AdminListUsers)
	r.DELETE("/api/v1/admin/users/:id", h.AdminDeleteUser)
	r.POST("/api/v1/auth/register", h.RegisterHandler)
	r.POST("/api/v1/auth/login", h.LoginHandler)
	r.POST("/api/v1/auth/refresh", h.RefreshHandler)

	authGroup := r.Group("/api/v1", func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Set("email", "shamil@example.com")
		c.Next()
	})
	{
		authGroup.GET("/me", h.GetProfileHandler)
		authGroup.GET("/my/urls", h.ListUserURLs)
		authGroup.DELETE("/my/urls/:id", h.DeleteURL)
	}

	return r
}

func TestLiveHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Health Check",
			path:           "/api/v1/live",
			expectedStatus: http.StatusOK,
			expectedBody:   "alive",
		},
	}

	repo := &FakeRepository{}
	r := setupTestRouter(repo, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestShortenHandler(t *testing.T) {
	tests := []struct {
		name             string
		body             string
		expectedStatus   int
		expectedContains string
	}{
		{
			name:             "Empty Body",
			body:             "",
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid request body",
		},
		{
			name:             "Malformed JSON",
			body:             `{"url":`,
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid request body",
		},
		{
			name:             "Missing URL Field",
			body:             `{}`,
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid request body",
		},
		{
			name:             "Invalid URL",
			body:             `{"url":"hello"}`,
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid URL",
		},
	}
	repo := &FakeRepository{}
	r := setupTestRouter(repo, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body := strings.NewReader(tt.body)

			req, err := http.NewRequest(
				http.MethodPost,
				"/api/v1/shorten",
				body,
			)
			if err != nil {
				t.Fatalf("test %q: failed to create request: %v", tt.name, err)
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("test %q: status code %d, want %d", tt.name, recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedContains) {
				t.Errorf("test %q: body %q, want %q", tt.name, recorder.Body.String(), tt.expectedContains)
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "Redirect",
			path:           "/abc",
			expectedStatus: http.StatusSeeOther,
			expectedBody:   "https://google.com",
			repo: &FakeRepository{
				GetURLByCodeFunc: func(code string) (string, error) {
					return "https://google.com", nil
				},
			},
		},
		{
			name:           "db down",
			path:           "/xyz",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Database error",
			repo: &FakeRepository{
				GetURLByCodeFunc: func(code string) (string, error) {
					return "", errors.New("database error")
				},
			},
		},
		{
			name:           "Not Found",
			path:           "/abcd",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "URL not found",
			repo: &FakeRepository{
				GetURLByCodeFunc: func(code string) (string, error) {
					return "", sql.ErrNoRows
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := setupTestRouter(tt.repo, nil)

			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "Delete",
			path:           "/api/v1/admin/urls/1",
			expectedStatus: http.StatusOK,
			expectedBody:   "URL deleted successfully",
			repo: &FakeRepository{
				DeleteURLFunc: func(id int) error {
					return nil
				},
			},
		},
		{
			name:           "Invalid ID",
			path:           "/api/v1/admin/urls/xyz",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
			repo:           &FakeRepository{},
		},
		{
			name:           "URL not found",
			path:           "/api/v1/admin/urls/2",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "URL not found",
			repo: &FakeRepository{
				DeleteURLFunc: func(id int) error {
					return sql.ErrNoRows
				},
			},
		},
		{
			name:           "db down",
			path:           "/api/v1/admin/urls/2",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to delete url",
			repo: &FakeRepository{
				DeleteURLFunc: func(id int) error {
					return errors.New("database error")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo, nil)

			req, err := http.NewRequest(http.MethodDelete, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestListAllHandler(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		expectedStatus   int
		expectedContains []string
		repo             *FakeRepository
	}{
		{
			name:           "List All",
			path:           "/api/v1/admin/urls",
			expectedStatus: http.StatusOK,
			expectedContains: []string{
				"https://google.com",
				"abc",
			},
			repo: &FakeRepository{
				GetAllURLsFunc: func(opts models.ListOptions) ([]models.URL, error) {
					return []models.URL{
						{
							ID:          1,
							OriginalURL: "https://google.com",
							ShortCode:   "abc",
							CreatedAt:   time.Now(),
							ClickCount:  0,
						},
					}, nil
				},
			},
		},
		{
			name:             "Repository Error",
			path:             "/api/v1/admin/urls",
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: []string{"Failed to get all urls"},
			repo: &FakeRepository{
				GetAllURLsFunc: func(opts models.ListOptions) ([]models.URL, error) {
					return nil, errors.New("database error")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo, nil)

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			for _, s := range tt.expectedContains {
				if !strings.Contains(recorder.Body.String(), s) {
					t.Fatalf(
						"response body = %q, want it to contain %q",
						recorder.Body.String(),
						s,
					)
				}
			}
		})
	}
}

func TestReadyHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "Ready",
			path:           "/api/v1/ready",
			expectedStatus: http.StatusOK,
			expectedBody:   "ready",
			repo: &FakeRepository{
				HealthFunc: func() error {
					return nil
				},
			},
		},
		{
			name:           "Not Ready",
			path:           "/api/v1/ready",
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   "database is unavailable",
			repo: &FakeRepository{
				HealthFunc: func() error {
					return errors.New("database unavailable")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo, nil)

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		body           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "Register",
			path:           "/api/v1/auth/register",
			body:           `{"username":"shamil","email":"shamil@example.com","password":"password123"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
			repo: &FakeRepository{
				CreateUserFunc: func(username, email, password string) (int64, error) {
					return 1, nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo, nil)

			req, err := http.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		body           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "Login",
			path:           "/api/v1/auth/login",
			body:           `{"email":"shamil@example.com","password":"password123"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
			repo: &FakeRepository{
				GetUserByEmailFunc: func(email string) (*models.User, error) {
					hash, _ := auth.HashPassword("password123")
					return &models.User{
						ID:           1,
						Username:     "shamil",
						Email:        "shamil@example.com",
						PasswordHash: hash,
					}, nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo, nil)

			req, err := http.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestDeleteUserURLHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		body           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "Delete User URL",
			path:           "/api/v1/my/urls/1",
			body:           `{"id":1}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
			repo: &FakeRepository{
				DeleteUserURLFunc: func(id int) error {
					return nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo, nil)

			req, err := http.NewRequest(http.MethodDelete, tt.path, strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestGetUserURLsHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		body           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "Get User URLs",
			path:           "/api/v1/my/urls",
			body:           `{"email":"shamil@example.com"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
			repo: &FakeRepository{
				GetAllURLsByUserEmailFunc: func(email string) ([]models.URL, error) {
					return []models.URL{
						{
							ID:          1,
							ShortCode:   "abc",
							OriginalURL: "https://google.com",
							CreatedAt:   time.Now(),
							ClickCount:  0,
						},
					}, nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo, nil)

			req, err := http.NewRequest(http.MethodGet, tt.path, strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestGetProfileHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "Get Profile Success",
			path:           "/api/v1/me",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
			repo: &FakeRepository{
				GetUserByEmailFunc: func(email string) (*models.User, error) {
					return &models.User{
						ID:       1,
						Username: "shamil",
						Email:    "shamil@example.com",
					}, nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo, nil)

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestRefreshHandler(t *testing.T) {
	token, err := auth.GenerateToken(1, "shamil@example.com", "test-refresh-secret", "15m")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	tests := []struct {
		name           string
		path           string
		body           string
		config         config.Config
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Refresh Success",
			path: "/api/v1/auth/refresh",
			body: `{"refresh_token":"` + token + `"}`,
			config: config.Config{
				JWT: config.JWTConfig{
					AccessTokenSecret:  "test-access-secret",
					AccessTokenExpiry:  "15m",
					RefreshTokenSecret: "test-refresh-secret",
					RefreshTokenExpiry: "7d",
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "access_token",
		},
		{
			name: "Invalid Token",
			path: "/api/v1/auth/refresh",
			body: `{"refresh_token":"invalid.token.signature"}`,
			config: config.Config{
				JWT: config.JWTConfig{
					RefreshTokenSecret: "test-refresh-secret",
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid or expired refresh token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := gin.New()
			h := New(&FakeRepository{}, nil, tt.config)
			r.POST("/api/v1/auth/refresh", h.RefreshHandler)

			req, err := http.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d (body: %s)", recorder.Code, tt.expectedStatus, recorder.Body.String())
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestAdminListUsers(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "List Users Success",
			path:           "/api/v1/admin/users",
			expectedStatus: http.StatusOK,
			expectedBody:   "admin-user",
			repo: &FakeRepository{
				GetAllUsersFunc: func() ([]models.User, error) {
					return []models.User{
						{
							ID:       1,
							Username: "admin-user",
							Email:    "admin@example.com",
						},
					}, nil
				},
			},
		},
		{
			name:           "List Users Repository Error",
			path:           "/api/v1/admin/users",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to retrieve users",
			repo: &FakeRepository{
				GetAllUsersFunc: func() ([]models.User, error) {
					return nil, errors.New("db error")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo, nil)

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestAdminDeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "Delete User Success",
			path:           "/api/v1/admin/users/1",
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
			repo: &FakeRepository{
				DeleteUserFunc: func(id int) error {
					return nil
				},
			},
		},
		{
			name:           "Delete User Invalid ID",
			path:           "/api/v1/admin/users/xyz",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid user ID",
			repo:           &FakeRepository{},
		},
		{
			name:           "Delete User Not Found",
			path:           "/api/v1/admin/users/2",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "User not found",
			repo: &FakeRepository{
				DeleteUserFunc: func(id int) error {
					return sql.ErrNoRows
				},
			},
		},
		{
			name:           "Delete User Repository Error",
			path:           "/api/v1/admin/users/2",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to delete user",
			repo: &FakeRepository{
				DeleteUserFunc: func(id int) error {
					return errors.New("db error")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo, nil)

			req, err := http.NewRequest(http.MethodDelete, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}
