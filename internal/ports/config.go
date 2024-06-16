package ports

type Conf interface {
	GetWebServerPort() string
	LoadTokenLimitList() error
	GetTokenLimit(token string) (int64, bool)
	IsRateLimitByIPEnabled() bool
	IsRateLimitByTokenEnabled() bool
	GetMaxRequestsByIP() int64
	GetBlockDurationIP() int64
	GetBlockDurationToken() int64
	GetStorageType() string
}
