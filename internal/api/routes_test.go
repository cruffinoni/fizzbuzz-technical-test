package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cruffinoni/fizzbuzz/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRoutes struct {
	mock.Mock
	db database.Database
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) AddRequest(request *database.FizzBuzzRequest) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockDB) DeleteTask(taskID int64) error {
	args := m.Called(taskID)
	return args.Error(0)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestRoutes_PlayFizzBuzz(t *testing.T) {
	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedResult string
	}{
		{
			name:           "valid_input",
			payload:        `{"number1": 3, "number2": 5, "replace1": "fizz", "replace2": "buzz", "limit": 15}`,
			expectedStatus: http.StatusOK,
			expectedResult: "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz",
		},
		{
			name:           "empty_replaces",
			payload:        `{"number1": 3, "number2": 5, "replace1": "", "replace2": "", "limit": 15}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid_json_input",
			payload:        `{"invalid":}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid_limit",
			payload:        `{"number1": 3, "number2": 5, "replace1": "", "replace2": "", "limit": 0}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid_numbers",
			payload:        `{"number1": -1, "number2": 0, "replace1": "fizz", "replace2": "buzz", "limit": 15}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			var requestBody PlayFizzBuzzBody
			_ = json.Unmarshal([]byte(tt.payload), &requestBody)
			//mockDB.On("AddTask", requestBody.Task).Return(tt.dbFunc(requestBody.Task))
			routes := Routes{db: mockDB}
			router := gin.Default()
			router.POST("/play", routes.PlayFizzBuzz)

			req, _ := http.NewRequest(http.MethodPost, "/play", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.expectedStatus, resp.Code)
			if resp.Code == http.StatusOK {
				var respBody PlayFizzBuzzResponse
				err := json.Unmarshal(resp.Body.Bytes(), &respBody)
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedResult, respBody.Result)
			}
		})
	}
}
