package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	fmt.Println(client)
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

}

func RedisSet(key, value string) (val string, err error) {
	err = client.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
		return "", err
	}

	val, err = client.Get(key).Result()
	if err != nil {
		panic(err)
		return "", err
	}
	fmt.Println(key, val)
	return val, nil
}
