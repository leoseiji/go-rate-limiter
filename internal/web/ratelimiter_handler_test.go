package web

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/leoseiji/go-ratelimiter/internal/database"
	"github.com/leoseiji/go-ratelimiter/internal/fixtures"
	"github.com/leoseiji/go-ratelimiter/internal/middleware"
	"github.com/leoseiji/go-ratelimiter/internal/ratelimiter"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiterHandler(t *testing.T) {
	tests := []struct {
		name            string
		configs         fixtures.MockedConfig
		qtRequests      int
		expectedSuccess int
		expectedError   int
	}{
		{
			name: "Block after 10 requests from the same IP",
			configs: fixtures.MockedConfig{
				EnableRateLimitByIp: "true",
				MaxRequestsByIP:     "10",
				BlockDurationIP:     "1",
			},
			qtRequests:      200,
			expectedSuccess: 10,
			expectedError:   190,
		},
		{
			name: "Block after 10 requests from the same token",
			configs: fixtures.MockedConfig{
				EnableRateLimitByToken: "true",
				BlockDurationToken:     "1",
				TokenLimitList:         `{"ABC123": 10}`,
			},
			qtRequests:      200,
			expectedSuccess: 10,
			expectedError:   190,
		},
		{
			name: "Block after 10 requests from the same token by priority over IP",
			configs: fixtures.MockedConfig{
				EnableRateLimitByToken: "true",
				EnableRateLimitByIp:    "true",
				BlockDurationToken:     "1",
				BlockDurationIP:        "1",
				TokenLimitList:         `{"ABC123": 10}`,
			},
			qtRequests:      200,
			expectedSuccess: 10,
			expectedError:   190,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			_ = tt.configs.LoadTokenLimitList()
			rateLimiter := *ratelimiter.NewRateLimiter(&tt.configs, database.NewLocalDB())
			mux.HandleFunc("/ratelimiter", middleware.Limit(RateLimiterHandler, rateLimiter))

			var successes atomic.Int32
			var errors atomic.Int32
			for i := 0; i < tt.qtRequests; i++ {
				req, _ := http.NewRequest(http.MethodGet, "/ratelimiter", nil)
				req.RemoteAddr = "0.0.0.1:8000"
				if tt.configs.IsRateLimitByTokenEnabled() {
					req.Header.Set("API_TOKEN", "ABC123")
				}
				rr := httptest.NewRecorder()
				mux.ServeHTTP(rr, req)
				if rr.Code == http.StatusOK {
					successes.Add(1)
				} else {
					errors.Add(1)
				}
			}

			assert.Equal(t, tt.expectedSuccess, int(successes.Load()))
			assert.Equal(t, tt.expectedError, int(errors.Load()))
		})
	}
}
