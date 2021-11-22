package db

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

type Tweet struct {
	timePosted string
	content    string
	likes      int
}

type CreateTweet struct {
	content string
}

type TweetRepository interface {
	ListTweets() ([]Tweet, error)
	CreateTweet(userId int, tweet CreateTweet) error
	GetTweet(id string) (*Tweet, error)
	Like(id string) error
}

type TweetRedisRepository struct {
	pool *redis.Pool
}

func NewTweetRepository(pool *redis.Pool) TweetRepository {
	return &TweetRedisRepository{pool: pool}
}

func (r *TweetRedisRepository) ListTweets() ([]Tweet, error) {
	return nil, nil
}

func (r *TweetRedisRepository) CreateTweet(userId int, tweet CreateTweet) error {
	conn := r.pool.Get()
	currentTime := time.Now().Unix()

	nextPostId, err := redis.Int(conn.Do("INCR", "next_post_id"))

	if err != nil {
		return err
	}

	_, err = conn.Do("HMSET", "post:"+strconv.Itoa(nextPostId), "time", currentTime, "content", tweet.content, "likes", 0)

	if err != nil {
		return err
	}

	_, err = conn.Do("ZADD", "post_of:"+strconv.Itoa(userId), currentTime, nextPostId)

	if err != nil {
		return err
	}

	return nil
}

func (r *TweetRedisRepository) GetTweet(id string) (*Tweet, error) {
	conn := r.pool.Get()

	tweet, err := redis.StringMap(conn.Do("HGETALL", "post:"+id))

	if err != nil {
		return nil, err
	}

	i, err := strconv.ParseInt(tweet["time"], 10, 64)

	timePosted := time.Unix(i, 0)

	if err != nil {
		return nil, err
	}

	likes, err := strconv.Atoi(tweet["likes"])

	if err != nil {
		return nil, err
	}

	return &Tweet{
		timePosted: timePosted.String(),
		content:    tweet["content"],
		likes:      likes,
	}, nil
}

func (r *TweetRedisRepository) Like(id string) error {
	conn := r.pool.Get()

	_, err := conn.Do("HINCRBY", "post:"+id, "likes")

	if err != nil {
		return err
	}
	return nil
}

//
//func (r *TweetRedisRepository) DeleteTweet()  {
//
//}
