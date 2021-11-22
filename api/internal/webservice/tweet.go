package webservice

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"hennge/yerassyl/twitterclone/internal/db"
	"net/http"
)

type TweetService struct {
	tweetRepository db.TweetRepository
	userRepository  db.UserRepository
}

func NewTweetService(repository db.TweetRepository, userRepository db.UserRepository) *TweetService {
	return &TweetService{tweetRepository: repository, userRepository: userRepository}
}

func (s *TweetService) CreateTweet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var post db.CreateTweet
	err := decoder.Decode(&post)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	userEmail, ok := r.Context().Value("user").(string)

	if !ok {
		http.Error(w, "Couldn't get user email", http.StatusInternalServerError)
	}

	user, err := s.userRepository.GetUserSimple(userEmail)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = s.tweetRepository.CreateTweet(user.Id, post)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *TweetService) GetTweet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	postId := chi.URLParam(r, "id")

	post, err := s.tweetRepository.GetTweet(postId)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	response, err := json.Marshal(post)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}

func (s *TweetService) Like(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	postId := chi.URLParam(r, "id")

	err := s.tweetRepository.Like(postId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
