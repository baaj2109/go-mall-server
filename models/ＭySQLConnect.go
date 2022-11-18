package models

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB
var err error

func init() {
	mysqladmin := beego.AppConfig.String("mysqladmin")
	mysqlpwd := beego.AppConfig.String("mysqlpwd")
	mysqldb := beego.AppConfig.String("mysqldb")
	DB, err = gorm.Open("mysql", mysqladmin+":"+mysqlpwd+"@/"+mysqldb+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		beego.Error(err)
		beego.Error("連接ＭｙＳＱＬ資料庫失敗")
	} else {
		beego.Info("連接MySQL資料庫成功")
	}
}
