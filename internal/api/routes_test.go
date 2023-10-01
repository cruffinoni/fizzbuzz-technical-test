package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRoutes_Ping(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid_input",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"pong"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routes := NewRoutes(nil)
			router := gin.Default()
			router.GET("/ping", routes.Ping)

			req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.Equal(t, tt.expectedBody, resp.Body.String())
		})
	}
}

func TestNewRoutes(t *testing.T) {
	tests := []struct {
		name string
		want *Routes
	}{
		{
			name: "valid_input",
			want: &Routes{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRoutes(nil)
			assert.NotNil(t, got)
		})
	}
}
