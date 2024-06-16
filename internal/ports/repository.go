package ports

import (
	"time"
)

type Repository interface {
	Get(key string) (int64, bool, error)
	Set(key string, value int64, expiration time.Duration) error
	Incr(key string) (int64, error)
}
