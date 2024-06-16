package ratelimiter

import (
	"log"
	"sync"
	"time"

	"github.com/leoseiji/go-ratelimiter/internal/ports"
)

type RateLimiter struct {
	mu         *sync.RWMutex
	Configs    ports.Conf
	repository ports.Repository
}

func NewRateLimiter(configs ports.Conf, repository ports.Repository) *RateLimiter {
	r := &RateLimiter{
		mu:         &sync.RWMutex{},
		Configs:    configs,
		repository: repository,
	}
	return r
}

func (r *RateLimiter) IsRateLimit(token, ip string) bool {
	if r.Configs.IsRateLimitByTokenEnabled() {
		if !r.IsTokenAllowed(token) {
			return true
		}
	}
	return !r.IsIPAllowed(ip)
}

func (r *RateLimiter) IsTokenAllowed(token string) bool {
	if token == "" {
		return true
	}
	tokenLimit, exists := r.Configs.GetTokenLimit(token)
	if !exists {
		return false
	}
	return r.IsAllowed(token, tokenLimit, r.Configs.GetBlockDurationToken())
}

func (r *RateLimiter) IsIPAllowed(ip string) bool {
	return r.IsAllowed(ip, r.Configs.GetMaxRequestsByIP(), r.Configs.GetBlockDurationIP())
}

func (r *RateLimiter) IsAllowed(key string, maxRequests, blockDuration int64) bool {
	userCounter, exists, err := r.repository.Get(key)
	if err != nil {
		log.Println("error getting user counter: ", err)
		return false
	}
	if !exists {
		err = r.repository.Set(key, 1, time.Duration(blockDuration)*time.Second)
		if err != nil {
			log.Println("error setting user counter: ", err)
			return false
		}
		return true
	}
	if userCounter >= maxRequests {
		return false
	}
	_, _ = r.repository.Incr(key)

	return true
}
