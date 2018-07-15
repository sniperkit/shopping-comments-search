package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error

	log.Info("reading configuration")
	configPath := "/Users/haofly/workspace/shopping-commont-search/servers/conf.json"
	config, err = NewConfig(configPath)
	if err != nil {log.Error(err)}

	log.Info("connecting mongodb...")
	mongoClient, err = NewMongo(MongoOptions{
		uri: config.get("mongo.uri"),
		db:  config.get("mongo.db"),
	})
	if err != nil {log.Error(err)}

	// log.Info("connecting redis...")
	// redisClient, err = NewRedis(RedisOptions{
	// 	addr: config.get("redis.addr"),
	// 	password: config.get("redis.password"),
	// })
	// if err != nil {log.Error(err)}

	itemID := "538232353890"
	sellerID := "1862759827"

	tmall := NewTmall()
	tmall.getComments(itemID, sellerID)
}
