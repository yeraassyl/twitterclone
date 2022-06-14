package db

import (
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

type Tweet struct {
	Id         int `json:"id"`
	TimePosted string `json:"time_posted"`
	Content    string `json:"content"`
	Likes      int `json:"likes"`
}

type CreateTweet struct {
	Content string `json:"content"`
}

type TweetRepository interface {
	CreateTweet(userId int, tweet CreateTweet) error
	GetTweet(id string) (*Tweet, error)
	Like(id string) error
	UserFeed(id int) ([]Tweet, error)
	ListTweets(id string) ([]Tweet, error)
}

type TweetRedisRepository struct {
	pool *redis.Pool
}

func NewTweetRepository(pool *redis.Pool) TweetRepository {
	return &TweetRedisRepository{pool: pool}
}

func (r *TweetRedisRepository) UserFeed(id int) ([]Tweet, error) {
	conn := r.pool.Get()

	followers, err := redis.Strings(conn.Do("SMEMBERS", "followers:"+strconv.Itoa(id)))

	if err != nil {
		return nil, err
	}

	var tweets []Tweet
	for _, follower := range followers {
		userTweets, err := r.ListTweets(follower)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, userTweets...)
	}

	return tweets, nil
}

func (r *TweetRedisRepository) ListTweets(id string) ([]Tweet, error) {
	conn := r.pool.Get()

	userTweets, err := redis.Strings(conn.Do("SMEMBERS", "post_of:"+id))

	if err != nil {
		return nil, err
	}

	var tweets []Tweet

	for _, tweetId := range userTweets {
		tweet, err := r.GetTweet(tweetId)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, Tweet{
			Id:         tweet.Id,
			TimePosted: tweet.TimePosted,
			Content:    tweet.Content,
			Likes:      tweet.Likes,
		})
	}

	return tweets, nil
}

func (r *TweetRedisRepository) CreateTweet(userId int, tweet CreateTweet) error {
	conn := r.pool.Get()
	currentTime := time.Now().Unix()

	nextPostId, err := redis.Int(conn.Do("INCR", "next_post_id"))

	if err != nil {
		return err
	}

	_, err = conn.Do("HMSET", "post:"+strconv.Itoa(nextPostId), "time", currentTime, "content", tweet.Content, "likes", 0)

	if err != nil {
		return err
	}

	_, err = conn.Do("SADD", "post_of:"+strconv.Itoa(userId), nextPostId)

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

	id_, err := strconv.Atoi(id)

	if err != nil {
		return nil, err
	}

	return &Tweet{
		Id:         id_,
		TimePosted: timePosted.String(),
		Content:    tweet["content"],
		Likes:      likes,
	}, nil
}

func (r *TweetRedisRepository) Like(id string) error {
	conn := r.pool.Get()

	_, err := conn.Do("HINCRBY", "post:"+id, "likes", 1)

	if err != nil {
		return err
	}
	return nil
}

//
//func (r *TweetRedisRepository) DeleteTweet()  {
//
//}
