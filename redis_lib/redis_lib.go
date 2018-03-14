package redis_lib

import (
	"fmt"
	"time"

	"encoding/json"

	"github.com/go-redis/redis"

	"github.com/Terryhung/infohub_rest/news"
)

func NewClient() (*redis.Client, bool) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	_, err := client.Ping().Result()
	if err == nil {
		return client, true
	} else {
		return client, false
	}
}

func CheckExists(client *redis.Client, key string) []news.News {
	val, err := client.Get(key).Bytes()
	results := []news.News{}
	if err == nil {
		fmt.Print("Hit\n")
		json.Unmarshal(val, &results)
	} else {
		fmt.Print(err)
	}
	return results
}

func SetValue(client *redis.Client, key string, val interface{}, duration_time int) bool {
	data, err := json.Marshal(val)
	err = client.Set(key, data, time.Duration(duration_time)*time.Second).Err()
	if err != nil {
		return false
	} else {
		fmt.Print("Save Cache\n")
		return true
	}
}
