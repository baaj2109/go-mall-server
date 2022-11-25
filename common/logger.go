package common

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func InitLogger() {
	maxlines := beego.AppConfig.String("maxlines")
	// maxsize := beego.AppConfig.String("maxsize")
	level := beego.AppConfig.String("log_level")
	logPath := beego.AppConfig.String("log_path")
	logConf := make(map[string]interface{})
	logConf["filename"] = logPath
	logConf["level"] = level
	logConf["maxlines"] = maxlines
	confStr, err := json.Marshal(logConf)
	if err != nil {
		return
	}
	beego.SetLogger(logs.AdapterFile, string(confStr))
	beego.SetLogFuncCall(true)
}
