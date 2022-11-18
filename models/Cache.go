package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
)

type cacheDB struct{}

var redisClient cache.Cache
var enableRedis, _ = beego.AppConfig.Bool("enableRedis")
var redisTime, _ = beego.AppConfig.Int("redisTime")
var YzmClient cache.Cache
var CacheDB = &cacheDB{}

func init() {
	if enableRedis {
		config := map[string]string{
			"key":      beego.AppConfig.String("redisKey"),
			"conn":     beego.AppConfig.String("redisConn"),
			"dbNum":    beego.AppConfig.String("redisDBNum"),
			"password": beego.AppConfig.String("redisPwd"),
		}
		bytes, _ := json.Marshal(config)
		redisClient, err = cache.NewCache("redis", string(bytes))
		YzmClient, _ = cache.NewCache("redis", string(bytes))
		if err != nil {
			fmt.Println(err)

			beego.Error("連接Redis資料庫失敗")
		} else {
			beego.Info("連接Ｒｅｄｉｓ成功")
		}
	}
}

// / 寫入資料
func (c cacheDB) Set(key string, value interface{}) {
	if enableRedis {
		bytes, _ := json.Marshal(value)
		redisClient.Put(key, string(bytes), time.Second*time.Duration(redisTime))
	}
}

// / 接收資料
func (c cacheDB) Get(key string, obj interface{}) bool {
	if enableRedis {
		redisStr := redisClient.Get(key)
		if redisStr != nil {
			fmt.Println("load data inside Redis")
			redisReadData, ok := redisStr.([]uint8)
			if !ok {
				fmt.Println("failed to get Redis data")
				return false
			}

			json.Unmarshal([]byte(redisReadData), obj)
			return true
		}
	}
	return false
}
