package webservice

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"hennge/yerassyl/twitterclone/internal/db"
	"io/ioutil"
	"log"
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
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	//
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	var post db.CreateTweet
	err = json.Unmarshal(body, &post)

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
	if r.Method != http.MethodPost {
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

func (s *TweetService) UserFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	userEmail := r.Context().Value("user").(string)
	user, err := s.userRepository.GetUserSimple(userEmail)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	tweets, err := s.tweetRepository.UserFeed(user.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Println(tweets)

	response, err := json.Marshal(tweets)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Println(response)

	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}

func (s *TweetService) ListTweets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	userId := chi.URLParam(r, "id")

	tweets, err := s.tweetRepository.ListTweets(userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response, err := json.Marshal(tweets)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}
