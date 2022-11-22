package routers

import (
	"LeastMall/controllers/frontend"

	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/", &frontend.IndexController{})

	/// augh
	beego.Router("/auth/registerStep1", &frontend.AuthController{}, "get:RegisterStep1")
	beego.Router("/auth/registerStep2", &frontend.AuthController{}, "get:RegisterStep2")
	beego.Router("/auth/registerStep3", &frontend.AuthController{}, "get:RegisterStep3")
	beego.Router("/auth/login", &frontend.AuthController{}, "get:Login")
	beego.Router("/auth/goLogin", &frontend.AuthController{}, "post:GoLogin")
	beego.Router("/auth/loginOut", &frontend.AuthController{}, "get:LoginOut")

	beego.Router("/auth/sendCode", &frontend.AuthController{}, "get:SendCode")
	beego.Router("/auth/doRegister", &frontend.AuthController{}, "post:GoRegister")
	beego.Router("/auth/validateSmsCode", &frontend.AuthController{}, "get:ValidateSmsCode")

	/// cart
	beego.Router("/cart", &frontend.CartController{})
	beego.Router("/cart/addCart", &frontend.CartController{}, "get:AddCart")
}
