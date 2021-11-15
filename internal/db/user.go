package db

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
)

var Pool *redis.Pool

type User struct {
	Id string `json:"id"`
	Username string `json:"username"`
}

func ListUsers() []User {
	conn := Pool.Get()
	keys, err := redis.Strings(conn.Do("SCAN", "0", "MATCH", "user:*"))
	if err != nil {
		log.Fatal(err)
	}

	var users []User

	for _, key := range keys {
		user, err1 := redis.StringMap(conn.Do("HGETALL", key))
		if err1 != nil {
			log.Fatal(err1)
		}
		users = append(users, User{user["id"], user["username"]})
	}
	return users
}

func CreateUser (id int, username string) {
	conn := Pool.Get()

	nextUserId, err := redis.String(conn.Do("INCR", "next_user_id"))
	if err != nil {
		log.Fatal(err)
	}

	//TODO: make a generic create method
	_, err1 := conn.Do("HMSET", "user:" + nextUserId, "id", id, "username", username)

	if err1 != nil {
		log.Fatal(err)
	}

	fmt.Println("CREATED USER " + username)
}

func GetUser (id int) *User{
	conn := Pool.Get()

	user, err := redis.StringMap(conn.Do("HGETALL", "user:" + strconv.Itoa(id)))

	if err != nil {
		log.Fatal(err)
	}

	return &User{
		Id:       user["id"],
		Username: user["username"],
	}
}