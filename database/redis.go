package database

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
)

const (
	UserChannelPrefix = "user-"
)

type RedisDB struct {
	Client *redis.Client
}

func MakeRedisClient() *RedisDB {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis_address"),
		Username: viper.GetString("redis_username"),
		Password: viper.GetString("redis_password"),
		DB:       viper.GetInt("redis_db"),
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Println("error connecting to redis")
		log.Println(err)
	}
	redisDB := &RedisDB{
		Client: client,
	}
	return redisDB
}

func (r RedisDB) PublishGroupJoin(userID string, groupID string) {
	message := RedisMessage{
		RedisMessageType: RedisGroupsChanged,
		Body:             map[string]string{"group_joined": groupID},
	}
	rawMessage, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	r.Client.Publish(context.Background(), UserChannelPrefix+userID, rawMessage)
}

func (r RedisDB) PublishGroupLeft(userID string, groupID string) {
	message := RedisMessage{
		RedisMessageType: RedisGroupsChanged,
		Body:             map[string]string{"group_left": groupID},
	}
	rawMessage, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	r.Client.Publish(context.Background(), UserChannelPrefix+userID, rawMessage)
}
