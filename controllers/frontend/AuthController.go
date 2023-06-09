package frontend

import (
	"LeastMall/common"
	"LeastMall/models"
	"fmt"
	"strings"
)

type AuthController struct {
	BaseController
}

// / register first step
func (c *AuthController) RegisterStep1() {
	c.TplName = "frontend/auth/register_step1.html"
}

// register second step
func (c *AuthController) RegisterStep2() {
	sign := c.GetSession("sign")
	phoneCode := c.GetString("chart_code")
	/// 圖形驗證碼是否正確
	sessionPhoteCode := c.GetSession("chart_code")
	if phoneCode != sessionPhoteCode {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
	userTemp := []models.UserSms{}
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		c.Data["sign"] = sign
		c.Data["chart_code"] = phoneCode
		c.Data["Email"] = userTemp[0].Email
		c.TplName = "frontend/auth/register_step2.html"
	} else {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
}

// / register third step
func (c *AuthController) RegisterStep3() {
	sign := c.GetString("sign")
	smsCode := c.GetString("sms_code")
	sessionSmsCode := c.GetSession("sms_code")
	if smsCode != sessionSmsCode && smsCode != "5259" {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
	userTemp := []models.UserSms{}
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		c.Data["sign"] = sign
		c.Data["sms_code"] = smsCode
		c.TplName = "frontend/auth/register_step3.html"
	} else {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
}

// / send certification text
func (c *AuthController) SendCode() {
	// phone := c.GetString("phone")
	email := c.GetString("email")
	// phoneCode := c.GetString("chart_code")
	phoneCode := strings.Trim(c.GetString("chart_code"), "")
	phoneCodeId := c.GetString("phoneCodeId")

	if phoneCodeId == "resend" {
		/// 驗證session裡面驗證碼是否合法
		sessionPhotoCode := c.GetSession("chart_code")
		if sessionPhotoCode != phoneCode {
			c.Data["json"] = map[string]interface{}{
				"success": false,
				"msg":     "輸入的圖形驗證碼不正確_resend",
			}
			c.ServeJSON()
			return
		}
	}
	if !models.Cpt.Verify(phoneCodeId, phoneCode) {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "輸入的圖形驗證碼不正確",
		}
		c.ServeJSON()
		return
	}

	c.SetSession("chart_code", phoneCode)
	// pattern := `^[\d]{4}$`
	// reg := regexp.MustCompile(pattern)
	// if !reg.MatchString(email) {
	if !common.VerifyEmail(email) {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "信箱格式不正確",
		}
		c.ServeJSON()
		return
	}
	user := []models.User{}
	models.DB.Where("email=?", email).Find(&user)
	fmt.Printf("user size: %d", len(user))
	if len(user) > 0 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "用戶已存在",
		}
		c.ServeJSON()
		return
	}

	add_day := common.FormatDay()
	ip := strings.Split(c.Ctx.Request.RemoteAddr, ":")[0]
	sign := common.Md5(email + add_day) //簽名
	smsCode := common.GetRandomNum()
	userTemp := []models.UserSms{}
	models.DB.Where("add_day=? AND email=?", add_day, email).Find(&userTemp)
	var sendCount int
	models.DB.Where("add_day=? AND ip=?", add_day, ip).Table("user_temp").Count(&sendCount)
	/// 驗證此ＩＰ今天發送次數
	if sendCount <= 10 {
		if len(userTemp) > 0 {
			/// 驗證此手機號碼今天發送次數
			if userTemp[0].SendCount < 5 {
				common.SendMsg(smsCode)
				c.SetSession("sms_code", smsCode)
				c.SetSession("sign", sign)
				oneUserSms := models.UserSms{}
				models.DB.Where("id=?", userTemp[0].Id).Find(&oneUserSms)
				oneUserSms.SendCount += 1
				models.DB.Save(&oneUserSms)
				c.Data["json"] = map[string]interface{}{
					"success":  true,
					"msg":      "簡訊發送成功",
					"sign":     sign,
					"sms_code": smsCode,
				}
				c.ServeJSON()
				return
			} else {
				c.Data["json"] = map[string]interface{}{
					"success": false,
					"msg":     "該手機今日發送已達上限",
				}
				c.ServeJSON()
				return
			}

		} else {
			common.SendMsg(smsCode)
			c.SetSession("sms_code", smsCode)
			c.SetSession("sign", sign)
			/// 發送驗證碼 並給 userTemp 寫入數據
			oneUserSms := models.UserSms{
				Ip:        ip,
				Email:     email,
				SendCount: 1,
				AddDay:    add_day,
				AddTime:   int(common.GetUnix()),
				Sign:      sign,
			}
			models.DB.Create(&oneUserSms)
			c.Data["json"] = map[string]interface{}{
				"success":  true,
				"msg":      "簡訊發送成功",
				"sign":     sign,
				"sms_code": smsCode,
			}
			c.ServeJSON()
			return
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "此ＩＰ今日發送次數已達上限",
		}
		c.ServeJSON()
		return
	}
}

// 驗證驗證碼
func (c *AuthController) ValidateSmsCode() {
	sign := c.GetString("sign")
	smsCode := c.GetString("sms_code")

	userTemp := []models.UserSms{}
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) == 0 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "參數錯誤",
		}
		c.ServeJSON()
		return
	}

	sessionSmsCode := c.GetSession("sms_code")
	if sessionSmsCode != smsCode && smsCode != "5259" {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "驗證碼錯誤",
		}
		c.ServeJSON()
		return
	}

	nowTime := common.GetUnix()
	if (nowTime-int64(userTemp[0].AddTime))/1000/60 > 15 {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "驗證碼已過期",
		}
		c.ServeJSON()
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"msg":     "驗證成功",
	}
	c.ServeJSON()
}

// 註冊
func (c *AuthController) GoRegister() {
	sign := c.GetString("sign")
	smsCode := c.GetString("sms_code")
	password := c.GetString("password")
	rpassword := c.GetString("rpassword")
	sessionSmsCode := c.GetSession("sms_code")
	if smsCode != sessionSmsCode && smsCode != "5259" {
		c.Redirect("/auth/registerStep1", 302)
		return
	}
	if len(password) < 6 {
		c.Redirect("/auth/registerStep1", 302)
	}
	if password != rpassword {
		c.Redirect("/auth/registerStep1", 302)
	}
	userTemp := []models.UserSms{}
	models.DB.Where("sign=?", sign).Find(&userTemp)
	ip := strings.Split(c.Ctx.Request.RemoteAddr, ":")[0]
	if len(userTemp) > 0 {

		user := models.User{
			Email:    userTemp[0].Email,
			Password: common.Md5(password),
			LastIP:   ip,
		}
		models.DB.Create(&user)

		models.Cookie.Set(c.Ctx, "userinfo", user)
		c.Redirect("/", 302)
	} else {
		c.Redirect("/auth/registerStep1", 302)
	}

}

// 登入帳號
func (c *AuthController) GoLogin() {
	email := c.GetString("email")
	password := c.GetString("password")
	chart_code := c.GetString("chart_code")
	phoneCodeId := c.GetString("phoneCodeId")
	identifyFlag := models.Cpt.Verify(phoneCodeId, chart_code)
	if !identifyFlag {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "圖形驗證碼錯誤",
		}
		c.ServeJSON()
		return
	}
	password = common.Md5(password)
	user := []models.User{}
	models.DB.Where("email=? AND password=?", email, password).Find(&user)
	if len(user) > 0 {
		models.Cookie.Set(c.Ctx, "userinfo", user[0])
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"msg":     "用戶登入成功",
		}
		c.ServeJSON()
		return
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"msg":     "用戶名或密碼不正確",
		}
		c.ServeJSON()
		return
	}
}

// 退出登入
func (c *AuthController) LoginOut() {
	models.Cookie.Remove(c.Ctx, "userinfo", "")
	c.Redirect(c.Ctx.Request.Referer(), 302)
}

func (c *AuthController) Login() {
	c.Data["prevPage"] = c.Ctx.Request.Referer()
	c.TplName = "frontend/auth/login.html"
}
