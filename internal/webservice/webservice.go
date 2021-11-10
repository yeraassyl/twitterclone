package webservice

import "net/http"

type WebService struct {

}

func New() *WebService{
	return &WebService{}
}

func (wb *WebService) ListUsers(w http.ResponseWriter, r *http.Request)  {
	
}

func (wb *WebService) CreateUser(w http.ResponseWriter, r *http.Request) {

}