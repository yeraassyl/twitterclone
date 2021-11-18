package db

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
)

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type UserRepository interface {
	ListUsers() []User
	CreateUser(email string) error
	GetUser(id string) (*User, error)
	UserExists(email string) bool
}

type UserRedisRepository struct {
	pool *redis.Pool
}

func NewUserRepository(pool *redis.Pool) UserRepository{
	return &UserRedisRepository{pool: pool}
}

func (r *UserRedisRepository) ListUsers() []User {
	conn := r.pool.Get()

	emailToId, err := redis.StringMap(conn.Do("HGETALL", "users"))
	if err != nil {
		log.Fatal(err)
	}

	var users []User

	for email, id := range emailToId {
		//user, err := redis.StringMap(conn.Do("HGETALL", "user:"+id))
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, User{id, email})
	}
	return users
}

func (r *UserRedisRepository) CreateUser(email string) error {
	conn := r.pool.Get()

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

func (r *UserRedisRepository) GetUser(id string) (*User, error) {
	conn := r.pool.Get()

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

func (r *UserRedisRepository) UserExists(email string) bool {
	conn := r.pool.Get()

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

func (r *UserRedisRepository) follow(userId string, followerId string) error {
	conn := r.pool.Get()

	email, err := redis.String(conn.Do("HGET", "user:"+followerId, "email"))

	if err != nil {
		return err
	}

	_, err = conn.Do("ZADD", "followers:"+userId, followerId, email)

	if err != nil {
		return err
	}

	return nil
}