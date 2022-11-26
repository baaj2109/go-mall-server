package common

import (
	"encoding/json"
	"strconv"

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
	i, nil := strconv.Atoi(level)
	logConf["level"] = i
	ii, nil := strconv.Atoi(maxlines)
	logConf["maxlines"] = ii
	confStr, err := json.Marshal(logConf)
	if err != nil {
		return
	}
	beego.SetLogger(logs.AdapterFile, string(confStr))
	beego.SetLogFuncCall(true)
}
