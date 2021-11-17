package main

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gomodule/redigo/redis"
	"golang.org/x/oauth2"
	"hennge/yerassyl/twitterclone/internal/auth"
	"hennge/yerassyl/twitterclone/internal/db"
	"hennge/yerassyl/twitterclone/internal/webservice"
	"net/http"
	"os"
	"time"
)

var (
	clientID     = os.Getenv("clientId")
	clientSecret = os.Getenv("clientSecret")
	redirectURL  = os.Getenv("redirectUrl")
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

	auth.Provider = provider

	auth.Verifier = provider.Verifier(&oidc.Config{ClientID: clientID})

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	auth.Oauth2Config = oauth2Config

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", auth.HandleRedirect)
	r.Get("/auth/google/callback", auth.HandleOAuth2Callback)
	r.Route("/user", func(r chi.Router) {
		r.Use(auth.Authorization)
		r.Get("/", webservice.ListUsers)
		r.Get("/{id}", webservice.GetUser)
	})

	http.ListenAndServe(":8080", r)
}
