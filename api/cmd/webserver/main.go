package main

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gomodule/redigo/redis"
	"hennge/yerassyl/twitterclone/internal/auth"
	"hennge/yerassyl/twitterclone/internal/db"
	"hennge/yerassyl/twitterclone/internal/webservice"
	"net/http"
	"time"
)

var (
	clientID     = "798066806591-sn722ltj9mus74s6985moee0mq9cnl0u.apps.googleusercontent.com"
)

func main() {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
	db.Pool = pool

	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		// handle error
	}

	auth.Verifier = provider.Verifier(&oidc.Config{ClientID: clientID})

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
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
		r.Get("/", webservice.ListUsers)
		r.Get("/{id}", webservice.GetUser)
	})

	http.ListenAndServe(":8080", r)
}
