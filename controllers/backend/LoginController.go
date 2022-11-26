package backend

import (
	"LeastMall/common"
	"LeastMall/models"
	"strings"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Get() {
	c.TplName = "backend/login/login.html"
}

func (c *LoginController) GoLogin() {
	var flag = models.Cpt.VerifyReq(c.Ctx.Request)
	if flag {
		username := strings.Trim(c.GetString("username"), "")
		password := common.Md5(strings.Trim(c.GetString("password"), ""))

		// test_admin := []models.Administrator{}
		// models.DB.Where("username=? ", username).Find(&test_admin)

		administrator := []models.Administrator{}
		models.DB.Where("username=? AND password=? AND status=1", username, password).Find(&administrator)
		if len(administrator) == 1 {
			c.SetSession("userinfo", administrator[0])
			c.Success("登入成功", "/")
		} else {
			c.Error("無登入權限或用戶名密碼錯誤", "/login")
		}
	} else {
		c.Error("驗證碼錯誤", "/login")
	}
}

func (c *LoginController) LoginOut() {
	c.DelSession("userinfo")
	c.Success("退出登入成功,返回登入頁面！", "/login")
}
