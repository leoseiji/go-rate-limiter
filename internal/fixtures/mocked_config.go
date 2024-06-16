package fixtures

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type MockedConfig struct {
	WebServerPort          string
	EnableRateLimitByIp    string
	EnableRateLimitByToken string
	MaxRequestsByIP        string
	BlockDurationIP        string
	BlockDurationToken     string
	StorageType            string
	TokenLimitList         string
	tokenLimitMap          map[string]int64
}

func (c *MockedConfig) GetWebServerPort() string {
	return c.WebServerPort
}

func (c *MockedConfig) LoadTokenLimitList() error {
	err := json.Unmarshal([]byte(c.TokenLimitList), &c.tokenLimitMap)
	if err != nil {
		return fmt.Errorf("error unmarshalling token limit list: %w", err)
	}
	return nil
}

func (c *MockedConfig) GetTokenLimit(token string) (int64, bool) {
	limit, exists := c.tokenLimitMap[token]
	return limit, exists
}

func (c *MockedConfig) IsRateLimitByIPEnabled() bool {
	return c.EnableRateLimitByIp == "true"
}

func (c *MockedConfig) IsRateLimitByTokenEnabled() bool {
	return c.EnableRateLimitByToken == "true"
}

func (c *MockedConfig) GetMaxRequestsByIP() int64 {
	maxRequest, err := strconv.ParseInt(c.MaxRequestsByIP, 10, 64)
	if err != nil {
		return 100
	}
	return maxRequest
}

func (c *MockedConfig) GetBlockDurationIP() int64 {
	blockDuration, err := strconv.ParseInt(c.BlockDurationIP, 10, 64)
	if err != nil {
		return 60
	}
	return blockDuration
}

func (c *MockedConfig) GetBlockDurationToken() int64 {
	blockDuration, err := strconv.ParseInt(c.BlockDurationToken, 10, 64)
	if err != nil {
		return 1
	}
	return blockDuration
}

func (c *MockedConfig) GetStorageType() string {
	return c.StorageType
}

func (c *MockedConfig) GetTokenLimitList() map[string]int64 {
	return c.tokenLimitMap
}
