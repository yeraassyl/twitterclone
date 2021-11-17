package db

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
)

var Pool *redis.Pool

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func ListUsers() []User {
	conn := Pool.Get()

	//TODO: Will fix this
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

func CreateUser(email string) error {
	conn := Pool.Get()

	nextUserId, err := redis.Int(conn.Do("INCR", "next_user_id"))
	if err != nil {
		return err
	}

	//TODO: make a generic create method
	_, err = conn.Do("HMSET", "user:"+strconv.Itoa(nextUserId), "email", email)

	if err != nil {
		return err
	}

	_, err = conn.Do("HSET", "users", email, nextUserId)

	if err != nil {
		return err
	}

	log.Println("CREATED USER " + email)

	return nil
}

func GetUser(id string) (*User, error) {
	conn := Pool.Get()

	user, err := redis.StringMap(conn.Do("HGETALL", "user:"+id))

	if err != nil {
		log.Fatal(err)
	}

	if len(user) == 0 {
		return nil, errors.New("no such user")
	}

	return &User{
		Id:    id,
		Email: user["email"],
	}, nil
}

func UserExists(email string) bool {
	conn := Pool.Get()

	userId, err := conn.Do("HGET", "users", email)

	if err != nil {
		log.Fatal(err)
	}

	if userId == nil {
		return false
	} else {
		return true
	}
}
