package models

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type cookie struct{}

var Cookie = &cookie{}

// / 寫入數據
func (c cookie) Set(ctx *context.Context, key string, value interface{}) {
	bytes, _ := json.Marshal(value)
	ctx.SetSecureCookie(
		beego.AppConfig.String("secureCookie"),
		key,
		string(bytes),
		3600*24*30,
		"/",
		beego.AppConfig.String("domain"),
		nil,
		true,
	)
}

// / 刪除數據
func (c cookie) Remove(ctx *context.Context, key string, value interface{}) {
	bytes, _ := json.Marshal(value)
	ctx.SetSecureCookie(
		beego.AppConfig.String("secureCookie"),
		key,
		string(bytes),
		-1,
		"/",
		beego.AppConfig.String("domain"),
		nil,
		true,
	)
}

// / 獲取數據
func (c cookie) Get(ctx *context.Context, key string, obj interface{}) bool {
	tempData, ok := ctx.GetSecureCookie(beego.AppConfig.String("secureCookie"), key)
	if !ok {
		return false
	}
	json.Unmarshal([]byte(tempData), obj)
	return true
}
