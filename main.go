package main

import (
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/leoseiji/go-ratelimiter/configs"
	"github.com/leoseiji/go-ratelimiter/internal/database"
	"github.com/leoseiji/go-ratelimiter/internal/middleware"
	"github.com/leoseiji/go-ratelimiter/internal/ports"
	"github.com/leoseiji/go-ratelimiter/internal/ratelimiter"
	"github.com/leoseiji/go-ratelimiter/internal/web"
)

func main() {
	configs, configErr := configs.LoadConfig(".")
	if configErr != nil {
		panic(configErr)
	}
	db := getDB(configs.GetStorageType())

	mux := http.NewServeMux()

	rateLimiter := *ratelimiter.NewRateLimiter(configs, db)
	mux.HandleFunc("/ratelimiter", middleware.Limit(web.RateLimiterHandler, rateLimiter))

	err := http.ListenAndServe(configs.GetWebServerPort(), mux)
	if err != nil {
		log.Fatalln("error starting server", err)
	}
}

func getDB(storageType string) ports.Repository {
	if storageType == "redis" {
		return database.NewRedisDB(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	} else {
		return database.NewLocalDB()
	}
}
