package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/health", HealthCheckHandler)
	r.POST("/shorten", ShortenHandler)
	r.GET("/:code", RedirectHandler)
	r.GET("/urls", ListAllHandler)
	r.DELETE("/:id", DeleteHandler)
	return r
}

func TestHealthCheckHandler(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		expectedStatus int
		expectedBody string
	}{
		{
			name:         "Health Check",
			path:         "/health",
			expectedStatus: http.StatusOK,
			expectedBody: "Welcome to URL Shortener Service",
		},
		// {
		// 	name:         "Unknown",
		// 	path:         "/unknown",
		// 	expectedStatus: http.StatusNotFound,
		// 	expectedBody: "404 page not found",
		// },
	}

	r := setupTestRouter()

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
    name       string
    body       string
    expectedStatus int
    expectedBody   string
	}{
		{
			name:       "Empty Body",
			body:       "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
		{
			name:       "Malformed JSON",
			body:       `{"url":`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
		{
			name:       "Missing URL Field",
			body:       `{}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
		{
			name:       "Invalid URL",
			body:       `{"url":"hello"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid URL",
		},
	}

	r := setupTestRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body := strings.NewReader(tt.body)

			req, err := http.NewRequest(
				http.MethodPost,
				"/shorten",
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
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Errorf("test %q: body %q, want %q", tt.name, recorder.Body.String(), tt.expectedBody)
			}
		})
	}
}	

// func TestRedirectHandler(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		path         string
// 		expectedStatus int
// 		expectedBody string
// 	}{
// 		{
// 			name:         "Redirect",
// 			path:         "/abc",
// 			expectedStatus: http.StatusSeeOther,
// 			expectedBody: "https://google.com",
// 		},
// 		{
// 			name:         "Not Found",
// 			path:         "/xyz",
// 			expectedStatus: http.StatusNotFound,
// 			expectedBody: "URL not found",
// 		},
// 	}

// 	r := setupTestRouter()

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			recorder := httptest.NewRecorder()

// 			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
// 			if err != nil {
// 				t.Fatalf("failed to create request: %v", err)
// 			}
// 			r.ServeHTTP(recorder, req)

// 			if recorder.Code != tt.expectedStatus {
// 				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
// 			}
// 			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
// 				t.Fatalf("response body = %q, want it to contain %q",
// 					recorder.Body.String(),
// 					tt.expectedBody,
// 				)
// 			}
// 		})
// 	}
// }


// func TestDeleteHandler(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		path         string
// 		expectedStatus int
// 		expectedBody string
// 	}{
// 		{
// 			name:         "Delete",
// 			path:         "/abc",
// 			expectedStatus: http.StatusSeeOther,
// 			expectedBody: "https://google.com",
// 		},
// 		{
// 			name:         "Not Found",
// 			path:         "/xyz",
// 			expectedStatus: http.StatusNotFound,
// 			expectedBody: "URL not found",
// 		},
// 	}

// 	r := setupTestRouter()

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			recorder := httptest.NewRecorder()

// 			req, err := http.NewRequest(http.MethodDelete, tt.path, nil)
// 			if err != nil {
// 				t.Fatalf("failed to create request: %v", err)
// 			}
// 			r.ServeHTTP(recorder, req)

// 			if recorder.Code != tt.expectedStatus {
// 				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
// 			}
// 			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
// 				t.Fatalf("response body = %q, want it to contain %q",
// 					recorder.Body.String(),
// 					tt.expectedBody,
// 				)
// 			}
// 		})
// 	}
// }


// func TestListAllHandler(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		path         string
// 		expectedStatus int
// 		expectedBody string
// 	}{
// 		{
// 			name:         "List All",
// 			path:         "/urls",
// 			expectedStatus: http.StatusOK,
// 			expectedBody: "https://google.com",
// 		},
// 		{
// 			name: "Not Found",
// 			path: "/xyz",
// 			expectedStatus: http.StatusNotFound,
// 			expectedBody: "URL not found",
// 		},
// 	}

// 	r := setupTestRouter()

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			recorder := httptest.NewRecorder()

// 			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
// 			if err != nil {
// 				t.Fatalf("failed to create request: %v", err)
// 			}
// 			r.ServeHTTP(recorder, req)

// 			if recorder.Code != tt.expectedStatus {
// 				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
// 			}
// 			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
// 				t.Fatalf("response body = %q, want it to contain %q",
// 					recorder.Body.String(),
// 					tt.expectedBody,
// 				)
// 			}
// 		})
// 	}
// }
	
