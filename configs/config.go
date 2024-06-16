package configs

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/leoseiji/go-ratelimiter/internal/ports"
	"github.com/spf13/viper"
)

type Conf struct {
	WebServerPort          string `mapstructure:"WEB_SERVER_PORT"`
	EnableRateLimitByIp    string `mapstructure:"ENABLE_RATE_LIMIT_BY_IP"`
	EnableRateLimitByToken string `mapstructure:"ENABLE_RATE_LIMIT_BY_TOKEN"`
	MaxRequestsByIP        string `mapstructure:"MAX_REQUESTS_BY_IP"`
	BlockDurationIP        string `mapstructure:"BLOCK_DURATION_IP"`
	BlockDurationToken     string `mapstructure:"BLOCK_DURATION_TOKEN"`
	StorageType            string `mapstructure:"STORAGE_TYPE"`
	TokenLimitList         string `mapstructure:"TOKEN_LIMIT_LIST"`
	tokenLimitMap          map[string]int64
}

func (c *Conf) GetWebServerPort() string {
	return c.WebServerPort
}

func (c *Conf) LoadTokenLimitList() error {
	err := json.Unmarshal([]byte(c.TokenLimitList), &c.tokenLimitMap)
	if err != nil {
		return fmt.Errorf("error unmarshalling token limit list: %w", err)
	}
	return nil
}

func (c *Conf) GetTokenLimit(token string) (int64, bool) {
	limit, exists := c.tokenLimitMap[token]
	return limit, exists
}

func (c *Conf) IsRateLimitByIPEnabled() bool {
	return c.EnableRateLimitByIp == "true"
}

func (c *Conf) IsRateLimitByTokenEnabled() bool {
	return c.EnableRateLimitByToken == "true"
}

func (c *Conf) GetMaxRequestsByIP() int64 {
	maxRequest, err := strconv.ParseInt(c.MaxRequestsByIP, 10, 64)
	if err != nil {
		return 100
	}
	return maxRequest
}

func (c *Conf) GetBlockDurationIP() int64 {
	blockDuration, err := strconv.ParseInt(c.BlockDurationIP, 10, 64)
	if err != nil {
		return 60
	}
	return blockDuration
}

func (c *Conf) GetBlockDurationToken() int64 {
	blockDuration, err := strconv.ParseInt(c.BlockDurationToken, 10, 64)
	if err != nil {
		return 1
	}
	return blockDuration
}

func (c *Conf) GetStorageType() string {
	return c.StorageType
}

func LoadConfig(path string) (ports.Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	err = cfg.LoadTokenLimitList()

	return cfg, err
}
