package webservice

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"hennge/yerassyl/twitterclone/internal/db"
	"net/http"
	"strconv"
)

type UserService struct {
	userRepository db.UserRepository
}

func NewUserService(repository db.UserRepository) *UserService {
	return &UserService{userRepository: repository}
}

func (s *UserService) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	users, err := s.userRepository.ListUsers()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response, err := json.Marshal(users)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (s *UserService) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	userId := chi.URLParam(r, "id")

	user, err := s.userRepository.GetUser(userId)

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

func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	userEmail, ok := r.Context().Value("user").(string)

	if !ok {
		http.Error(w, http.StatusText(401), http.StatusBadRequest)
	}

	if s.userRepository.UserExists(userEmail) {
		userSimple, err := s.userRepository.GetUserSimple(userEmail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		user, err := s.userRepository.GetUser(strconv.Itoa(userSimple.Id))

		response, err := json.Marshal(user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-type", "application/json")
		w.Write(response)
		return
	}

	user, err := s.userRepository.CreateUser(userEmail)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}

func (s *UserService) Follow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	userEmail := r.Context().Value("user").(string)
	userId, err := s.userRepository.GetUserSimple(userEmail)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	userToFollow, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = s.userRepository.Follow(userId.Id, userToFollow)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
