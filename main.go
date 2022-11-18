package main

import (
	"LeastMall/models"
	_ "LeastMall/routers"
	"encoding/gob"

	"LeastMall/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func main() {

	/// 增加 function 到 map, 供前端HTM使用
	beego.AddFuncMap("timestampToDate", common.TimestampToDate)
	models.DB.LogMode(true)
	beego.AddFuncMap("formatImage", common.FormatImage)
	beego.AddFuncMap("mul", common.Mul)
	beego.AddFuncMap("formatAttribute", common.FormatAttribute)
	beego.AddFuncMap("setting", models.GetSettingByColumn)

	/// 後台設定
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"127.0.0.1"},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Authorization", "Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers", "Content-Type",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Content-Type",
		},
		AllowCredentials: true, //是否允許 cookies
	}))
	/// 註冊模型
	gob.Register(models.Administrator{})
	/// 結束時關閉資料庫
	defer models.DB.Close()
	/// 設定 Redis用於儲存 session
	beego.BConfig.WebConfig.Session.SessionProvider = "redis"
	/// docker-compose 請設定為 redisServiceHost
	// beego.BConfig.WebConfig.Session.SessionProviderConfig="redisServiceHost:6379"

	/// local launch
	beego.BConfig.WebConfig.Session.SessionProviderConfig = "127.0.0.1:6379"
	beego.Run()
}
