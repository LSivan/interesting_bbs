package common

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/astaxie/beego"
)

var Redis *redis.Client

func SetupRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     beego.AppConfig.String("redis.host")+":"+beego.AppConfig.String("redis.port"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := Redis.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
}
