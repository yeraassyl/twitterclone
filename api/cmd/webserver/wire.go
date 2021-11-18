// +build wireinject

package main

import (
	"github.com/google/wire"
	"hennge/yerassyl/twitterclone/internal/db"
	"hennge/yerassyl/twitterclone/internal/webservice"
)

func InitializeUserService() *webservice.UserService {
	wire.Build(NewPool, db.NewUserRepository, webservice.NewUserService)
	return &webservice.UserService{}
}

