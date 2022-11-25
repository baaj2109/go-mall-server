package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/gomarkdown/markdown"
)

// / 時間戳記改寫成日期格式
func TimestampToDate(t int) string {
	unix_t := time.Unix(int64(t), 0)
	return unix_t.Format("2006-01-02 15:04:05")
}

// / 格式化圖片
func FormatImage(picPath string) string {
	ossStatus, err := beego.AppConfig.Bool("ossStatus")
	if err != nil {
		// 判斷目錄前面是否有 "/"
		flag := strings.Contains(picPath, "/static")
		if flag {
			return picPath
		}
		return "/" + picPath
	}
	if ossStatus {
		return beego.AppConfig.String("ossDomain") + "/" + picPath
	} else {
		flag := strings.Contains(picPath, "/static")
		if flag {
			return picPath
		} else {
			return "/" + picPath
		}
	}
}

// / 計算乘法
func Mul(price float64, num int) float64 {
	return price * float64(num)
}

// / 格式化標題
func FormatAttribute(s string) string {
	md := []byte(s)
	html := markdown.ToHTML(md, nil, nil)
	return string(html)
}

// / 獲取日期
func FormatDay() string {
	template := "20060102"
	return time.Now().Format(template)
}

// / get random runber
func GetRandomNum() string {
	var str string
	for i := 0; i < 4; i++ {
		current := rand.Intn(10) //0-9   "math/rand"
		str += strconv.Itoa(current)
	}
	return str
}

// Md5加密
func Md5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return string(hex.EncodeToString(m.Sum(nil)))
}

// get current time
func GetUnix() int64 {
	fmt.Println(time.Now().Unix())
	return time.Now().Unix()
}

func SendMsg(str string) {
	// 短信验证码需要到相关网站申请
	// 目前先固定一个值
	ioutil.WriteFile("test_send.txt", []byte(str), 06666)
}

func GenerateOrderId() string {
	template := "200601021504"
	return time.Now().Format(template) + GetRandomNum()
}
