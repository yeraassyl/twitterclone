package auth

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"net/http"
	"strings"
)

var Verifier *oidc.IDTokenVerifier

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
				http.Error(w, "Invalid token", http.StatusUnauthorized)
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
			//if !db.UserExists(claims.Email) {
			//	err = db.CreateUser(claims.Email)
			//	if err != nil {
			//		http.Error(w, "Failed to create user"+err.Error(), http.StatusInternalServerError)
			//	}
			//}
			cntxt := context.WithValue(r.Context(), "user", claims.Email)
			next.ServeHTTP(w, r.WithContext(cntxt))
		}
	})
}
