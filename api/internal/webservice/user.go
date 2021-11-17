package webservice

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"hennge/yerassyl/twitterclone/internal/db"
	"net/http"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	js, err := json.Marshal(db.ListUsers())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(405), 405)
		return
	}
	userId := chi.URLParam(r, "id")

	user, err := db.GetUser(userId)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	response, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}
