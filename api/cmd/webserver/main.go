package main

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gomodule/redigo/redis"
	"hennge/yerassyl/twitterclone/internal/auth"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	clientID = os.Getenv("clientId")
)

func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}

func main() {
	ctx := context.Background()
	userService := InitializeUserService()

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		log.Println(err)
	}

	auth.Verifier = provider.Verifier(&oidc.Config{ClientID: clientID})

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Logger)

	r.Route("/user", func(r chi.Router) {
		r.Use(auth.Authorization)
		r.Get("/", userService.ListUsers)
		r.Get("/{id}", userService.GetUser)
	})

	http.ListenAndServe(":8080", r)
}
