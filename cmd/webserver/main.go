package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gomodule/redigo/redis"
	"hennge/yerassyl/twitterclone/internal/db"
	"hennge/yerassyl/twitterclone/internal/webservice"
	"net/http"
	"time"
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)


	r.Route("/user", func(r chi.Router) {
		r.Get("/", webservice.ListUsers)
		r.Post("/", webservice.CreateUser)
		r.Get("/{id}", webservice.GetUser)
	})

	http.ListenAndServe(":8080", r)
}

