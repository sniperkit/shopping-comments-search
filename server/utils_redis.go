package main

import (
	"github.com/go-redis/redis"
	// log "github.com/sirupsen/logrus"
)

var redisClient RedisClient

// RedisOptions init options
type RedisOptions struct {
	addr     string
	password string
	db       int
}

// RedisClient class
type RedisClient struct {
	client  *redis.Client
	options RedisOptions

	KeyProxy string
}

// NewRedis create Redis class
func NewRedis(redisOptions RedisOptions) (RedisClient, error) {
	if redisOptions.addr == "" {
		redisOptions.addr = "127.0.0.1:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisOptions.addr,
		Password: redisOptions.password,
		DB:       redisOptions.db,
	})

	redisClient := RedisClient{
		client:   client,
		options:  redisOptions,
		KeyProxy: "haipproxy:speed:http",
	}

	return redisClient, nil
}

func (redisClient *RedisClient) get(key string) (string, error) {
	return redisClient.client.Get(key).Result()
}

func (redisClient *RedisClient) set(key string, value string) (bool, error) {
	redisClient.client.Set(key, value, 0)
	return true, nil
}

func (redisClient *RedisClient) del(key string) (bool, error) {
	redisClient.client.Del(key)
	return true, nil
}

func (redisClient *RedisClient) srandmember(key string) (string, error) {
	return redisClient.client.SRandMember(key).Result()
}

func (redisClient *RedisClient) zrange(key string, minScore int64, maxScore int64) ([]string, error) {
	return redisClient.client.ZRange(key, minScore, maxScore).Result()
}

// judge whether or not connected
func (redisClient *RedisClient) isConnected() bool {
	result, err := redisClient.client.Ping().Result()
	if err != nil {
		return false
	}
	return result == "PONG"
}
