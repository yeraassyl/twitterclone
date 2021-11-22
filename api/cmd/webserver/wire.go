//go:build wireinject
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

func InitializerTweetService() *webservice.TweetService {
	wire.Build(NewPool, db.NewUserRepository, db.NewTweetRepository, webservice.NewTweetService)
	return &webservice.TweetService{}
}