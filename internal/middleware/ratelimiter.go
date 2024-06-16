package middleware

import (
	"net"
	"net/http"

	"github.com/leoseiji/go-ratelimiter/internal/ratelimiter"
)

const TOO_MANY_REQUESTS_TEXT = "you have reached the maximum number of requests or actions allowed within a certain time frame"

func Limit(next http.HandlerFunc, ratelimit ratelimiter.RateLimiter) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if ratelimit.IsRateLimit(r.Header.Get("API_TOKEN"), ip) {
			http.Error(w, TOO_MANY_REQUESTS_TEXT, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
