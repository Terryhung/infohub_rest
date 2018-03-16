package redis_lib

import (
	"fmt"
	"time"

	"encoding/json"

	"github.com/go-redis/redis"

	"github.com/Terryhung/infohub_rest/gifimage"
	"github.com/Terryhung/infohub_rest/news"
	"github.com/Terryhung/infohub_rest/video"
)

func NewClient() (*redis.Client, bool) {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	_, err := client.Ping().Result()
	if err == nil {
		fmt.Print("New Cache Connection!\n")
		return client, true
	} else {
		fmt.Print("No Redis!\n")
		return client, false
	}
}

func CheckExists(client *redis.Client, key string, result interface{}) {
	val, err := client.Get(key).Bytes()
	switch t := result.(type) {
	case *[]video.Video, *[]news.News, *[]gifimage.GifImage:
		if err == nil {
			fmt.Print("Hit\n")
			json.Unmarshal(val, &t)
		} else {
			fmt.Print(err)
		}

	}
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
