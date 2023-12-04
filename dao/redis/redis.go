package redis

import (
	"fmt"
	"main/settings"

	"github.com/go-redis/redis"
)

//var rdb *redis.Client
//
//func Init(cfg *settings.RedisConfig) (err error) {
//	// 注意这里不再声明rdb变量了，不然之后对rdb操作会产生空指针错误
//	rdb = redis.NewClient(&redis.Options{
//		Addr: fmt.Sprintf("%s:%d",
//			cfg.Host,
//			cfg.Port,
//		),
//		Password: cfg.Password, // no password set
//		DB:       cfg.DB,       // use default DB
//		PoolSize: cfg.PoolSize,
//	})
//	_, err = rdb.Ping().Result()
//	return err
//
//}
//
//func Close() {
//	_ = rdb.Close()
//}

var (
	client *redis.Client
	Nil    = redis.Nil
)

// Init 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
		PoolSize: cfg.PoolSize,
	})

	_, err = client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = client.Close()
}
