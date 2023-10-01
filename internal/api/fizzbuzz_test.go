package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cruffinoni/fizzbuzz/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) AddRequest(request *database.FizzBuzzRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *MockDB) GetMostUsedRequest() (*database.MostUsedRequest, error) {
	args := m.Called()
	return args.Get(0).(*database.MostUsedRequest), args.Error(1)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestRoutes_FormatFizzBuzzPerformance(t *testing.T) {
	input := &PlayFizzBuzzBody{
		Number1:  3,
		Number2:  5,
		Replace1: "fizz",
		Replace2: "buzz",
		Limit:    99999, // this value impose a remainder for the soft limit
	}

	var (
		disabledPerfElapsed = time.Duration(0)
		disabledPerfResult  *PlayFizzBuzzResponse
	)
	t.Run("disabled_perf", func(t *testing.T) {
		start := time.Now()
		disabledPerfResult = basicFormatFizzBuzzFromBody(input)
		disabledPerfElapsed = time.Since(start)
		assert.True(t, disabledPerfElapsed.Seconds() < 30)
	})
	var (
		enabledPerfElapsed = time.Duration(0)
		enabledPerfResult  *PlayFizzBuzzResponse
	)
	t.Run("enabled_perf", func(t *testing.T) {
		start := time.Now()
		enabledPerfResult = formatFizzBuzzFromBodyWithPerformance(input)
		enabledPerfElapsed = time.Since(start)
		assert.True(t, enabledPerfElapsed.Seconds() < 10)
	})
	assert.True(t, enabledPerfElapsed < disabledPerfElapsed)
	assert.Equal(t, enabledPerfResult.Result, disabledPerfResult.Result)
}

func TestRoutes_PlayFizzBuzz(t *testing.T) {
	tests := []struct {
		name             string
		payload          string
		expectedStatus   int
		expectedResult   string
		addRequestMockDb func(request *database.FizzBuzzRequest) error
	}{
		{
			name:           "valid_input",
			payload:        `{"number1": 3, "number2": 5, "replace1": "fizz", "replace2": "buzz", "limit": 15}`,
			expectedStatus: http.StatusOK,
			expectedResult: "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz",
			addRequestMockDb: func(request *database.FizzBuzzRequest) error {
				return nil
			},
		},
		{
			name:           "empty_replaces",
			payload:        `{"number1": 3, "number2": 5, "replace1": "", "replace2": "", "limit": 15}`,
			expectedStatus: http.StatusBadRequest,
			addRequestMockDb: func(request *database.FizzBuzzRequest) error {
				return nil
			},
		},
		{
			name:           "invalid_json_input",
			payload:        `{"invalid":}`,
			expectedStatus: http.StatusBadRequest,
			addRequestMockDb: func(request *database.FizzBuzzRequest) error {
				return nil
			},
		},
		{
			name:           "invalid_limit",
			payload:        `{"number1": 3, "number2": 5, "replace1": "", "replace2": "", "limit": 0}`,
			expectedStatus: http.StatusBadRequest,
			addRequestMockDb: func(request *database.FizzBuzzRequest) error {
				return nil
			},
		},
		{
			name:           "invalid_numbers",
			payload:        `{"number1": -1, "number2": 0, "replace1": "fizz", "replace2": "buzz", "limit": 15}`,
			expectedStatus: http.StatusBadRequest,
			addRequestMockDb: func(request *database.FizzBuzzRequest) error {
				return nil
			},
		},
		{
			name:           "error_db",
			payload:        `{"number1": 3, "number2": 5, "replace1": "fizz", "replace2": "buzz", "limit": 15}`,
			expectedStatus: http.StatusInternalServerError,
			addRequestMockDb: func(request *database.FizzBuzzRequest) error {
				return assert.AnError
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			var requestBody PlayFizzBuzzBody
			_ = json.Unmarshal([]byte(tt.payload), &requestBody)
			reqDb := &database.FizzBuzzRequest{
				Int1:  requestBody.Number1,
				Int2:  requestBody.Number2,
				Limit: requestBody.Limit,
				Str1:  requestBody.Replace1,
				Str2:  requestBody.Replace2,
			}
			mockDB.On("AddRequest", reqDb).Return(tt.addRequestMockDb(reqDb))
			routes := NewRoutes(mockDB)
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

func TestRoutes_GetMostUsedRequest(t *testing.T) {
	tests := []struct {
		name             string
		expectedStatus   int
		expectedResult   *database.MostUsedRequest
		getRequestMockDb func() (*database.MostUsedRequest, error)
	}{
		{
			name:           "valid_input",
			expectedStatus: http.StatusOK,
			expectedResult: &database.MostUsedRequest{
				Int1:  3,
				Int2:  5,
				Hints: 10,
			},
			getRequestMockDb: func() (*database.MostUsedRequest, error) {
				return &database.MostUsedRequest{
					Int1:  3,
					Int2:  5,
					Hints: 10,
				}, nil
			},
		},
		{
			name:           "generic_error_db",
			expectedStatus: http.StatusInternalServerError,
			getRequestMockDb: func() (*database.MostUsedRequest, error) {
				return nil, assert.AnError
			},
		},
		{
			name:           "no_data_db",
			expectedStatus: http.StatusNotFound,
			getRequestMockDb: func() (*database.MostUsedRequest, error) {
				return nil, database.ErrNoRequest
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			mockDB.On("GetMostUsedRequest").Return(tt.getRequestMockDb())
			routes := NewRoutes(mockDB)
			router := gin.Default()
			router.GET("/most-used", routes.GetMostUsedRequest)

			req, _ := http.NewRequest(http.MethodGet, "/most-used", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.expectedStatus, resp.Code)
			if resp.Code == http.StatusOK {
				var respBody database.MostUsedRequest
				err := json.Unmarshal(resp.Body.Bytes(), &respBody)
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedResult, &respBody)
			}
		})
	}
}
