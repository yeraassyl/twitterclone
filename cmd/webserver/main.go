package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"hennge/yerassyl/twitterclone/internal/webservice"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	wb := webservice.New()

	r.Route("/user", func(r chi.Router) {
		r.Get("/", wb.ListUsers)
		r.Post("/", wb.CreateUser)
	})

	http.ListenAndServe(":8080", r)
}

