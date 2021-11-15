package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"hennge/yerassyl/twitterclone/internal/db"
	"io"
	"net/http"
	"strings"
	"time"
)

var Oauth2Config oauth2.Config
var Provider *oidc.Provider
var Verifier *oidc.IDTokenVerifier

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	state, err := randString(16)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	setCallbackCookie(w, r, "state", state)

	http.Redirect(w, r, Oauth2Config.AuthCodeURL(state), http.StatusFound)
}

func HandleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	state, err := r.Cookie("state")
	if err != nil {
		http.Error(w, "state not found", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	oauth2Token, err := Oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := Provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		http.Error(w, "Failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !db.UserExists(userInfo.Email) {
		err = db.CreateUser(userInfo.Email)
		if err != nil {
			http.Error(w, "Failed to create user"+err.Error(), http.StatusInternalServerError)
		}
	}

	resp := struct {
		OAuth2Token *oauth2.Token
		UserInfo    *oidc.UserInfo
	}{oauth2Token, userInfo}
	data, err := json.MarshalIndent(resp, "", "    ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			idToken, err := Verifier.Verify(ctx, authHeader[1])
			if err != nil {
				HandleRedirect(w, r)
				return
			}
			var claims struct {
				Email    string `json:"email"`
				Verified bool   `json:"email_verified"`
			}
			if err := idToken.Claims(&claims); err != nil {
				http.Error(w, "Couldn't extract user", http.StatusInternalServerError)
				return
			}
			cntxt := context.WithValue(r.Context(), "user", claims.Email)
			next.ServeHTTP(w, r.WithContext(cntxt))
		}
	})
}
