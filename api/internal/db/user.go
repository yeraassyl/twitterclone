package db

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
)

type User struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	Followers []int  `json:"followers"`
	Following []int  `json:"following"`
}

type UserSimple struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type CreateUser struct {
	Email string `json:"email"`
}

type Follow struct {
	UserToFollow int `json:"user_to_follow"`
}

type UserRepository interface {
	ListUsers() ([]UserSimple, error)
	CreateUser(user CreateUser) error
	GetUser(id string) (*User, error)
	UserExists(email string) bool
	Follow(userId int, anotherUserId int) error
	GetUserSimple(email string) (*UserSimple, error)
}

type UserRedisRepository struct {
	pool *redis.Pool
}

func NewUserRepository(pool *redis.Pool) UserRepository {
	return &UserRedisRepository{pool: pool}
}

func (r *UserRedisRepository) ListUsers() ([]UserSimple, error) {
	conn := r.pool.Get()

	emailToId, err := redis.StringMap(conn.Do("HGETALL", "users"))
	if err != nil {
		return nil, err
	}

	var users []UserSimple

	for email, id := range emailToId {
		//user, err := redis.StringMap(conn.Do("HGETALL", "user:"+id))
		if err != nil {
			return nil, err
		}
		id_, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		users = append(users, UserSimple{id_, email})
	}
	return users, nil
}

func (r *UserRedisRepository) CreateUser(user CreateUser) error {
	conn := r.pool.Get()

	nextUserId, err := redis.Int(conn.Do("INCR", "next_user_id"))
	if err != nil {
		return err
	}

	_, err = conn.Do("HMSET", "user:"+strconv.Itoa(nextUserId), "email", user.Email)

	if err != nil {
		return err
	}

	_, err = conn.Do("HSET", "users", user.Email, nextUserId)

	if err != nil {
		return err
	}

	log.Println("CREATED USER " + user.Email)

	return nil
}

func (r *UserRedisRepository) GetUser(id string) (*User, error) {
	conn := r.pool.Get()

	email, err := redis.String(conn.Do("HGET", "user:"+id, "email"))

	if err != nil {
		return nil, err
	}

	followers, err := redis.Ints(conn.Do("HGETALL", "followers:"+id))

	if err != nil {
		return nil, err
	}

	following, err := redis.Ints(conn.Do("HGETALL", "following:"+id))

	if err != nil {
		return nil, err
	}

	return &User{
		Id:        id,
		Email:     email,
		Followers: followers,
		Following: following,
	}, nil
}

func (r *UserRedisRepository) GetUserSimple(email string) (*UserSimple, error) {
	conn := r.pool.Get()

	userId, err := redis.Int(conn.Do("HGET", "users", email))

	if err != nil {
		return nil, err
	}

	return &UserSimple{
		Id:    userId,
		Email: email}, nil
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

//func (r *UserRedisRepository) UserExistsById(id string)  {
//	conn := r.pool.Get()
//
//	email, err := conn.Do("HGET", "")
//}

func (r *UserRedisRepository) Follow(userId int, anotherUserId int) error {
	conn := r.pool.Get()

	_, err := conn.Do("SADD", "following:"+strconv.Itoa(userId), anotherUserId)
	if err != nil {
		return err
	}

	_, err = conn.Do("SADD", "followers:"+strconv.Itoa(anotherUserId), userId)
	if err != nil {
		return err
	}

	return nil
}
