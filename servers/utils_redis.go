package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
)

var redisClient RedisClient

// RedisOptions init options
type RedisOptions struct {
	addr string
	password string
}

// RedisClient class
type RedisClient struct {
	client redis.Conn
	options RedisOptions

	KeyProxy string
}

// NewRedis create Redis class
func NewRedis(redisOptions RedisOptions) (RedisClient, error) {
	if redisOptions.addr == "" {
		redisOptions.addr = "127.0.0.1:6379"
	}
	redisClient := RedisClient{
		options: redisOptions,
		KeyProxy: "haipproxy:all",
	}
	
	c, err := redis.Dial("tcp", redisOptions.addr)
	if err != nil {
	}

	if redisOptions.password != "" {
		if _, err := c.Do("AUTH", redisOptions.password); err != nil {
			log.Error("connection error")
			c.Close()
			return redisClient, err
		}
	}

	redisClient.client = c

	fmt.Println(redisClient.options.addr)
	fmt.Println(redisClient.isConnected())
	return redisClient, nil
}

func (redisClient *RedisClient) get(key string) (string, error) {
	redisClient.client.Send("GET", key)
	redisClient.client.Flush()
	return redis.String(redisClient.client.Receive())
}

func (redisClient *RedisClient) set(key string, value string) (bool, error) {
	redisClient.client.Send("SET", key, value)
	redisClient.client.Flush()
	return true, nil
}

func (redisClient *RedisClient) srandmember(key string) (string, error) {
	redisClient.client.Send("SRANDMEMBER", key)
	redisClient.client.Flush()
	return redis.String(redisClient.client.Receive())
}

// judge whether or not connected
func (redisClient *RedisClient) isConnected() bool {
	if redisClient.options.addr == "" {
		return false
	}

	return true
}
