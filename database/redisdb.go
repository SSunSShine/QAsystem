package database

import (
	"context"
	"github.com/SSunSShine/QAsystem/conf"
	"github.com/go-redis/redis/v8"
	"log"
)

var RDB *redis.Client
var ctx = context.Background()

// 连接Redis
func init()  {

	RDB = redis.NewClient(&redis.Options{
		Addr: conf.Config().Redis.Addr,
		Password: conf.Config().Redis.Password,
		DB: conf.Config().Redis.Db,
		PoolSize: 100,
	})

	if dbsize, err := RDB.DBSize(ctx).Result(); err != nil {
		log.Println(err)
		panic("failed to connect redis !")
	} else {
		log.Println("redis size: dbsize", dbsize)
	}

}