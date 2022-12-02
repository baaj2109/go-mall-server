package routers

import (
	"LeastMall/controllers/frontend"

	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/", &frontend.IndexController{})

	/// google oauth
	beego.Router("/oauth/GoogleLogin", &frontend.GoogleOauthController{}, "get:HandleGoogleLogin")
	beego.Router("/oauth/GoogleCallBack", &frontend.GoogleOauthController{}, "get:HandleGoogleCallback")

	/// auth
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
	beego.Router("/cart/incCart", &frontend.CartController{}, "get:IncCart")
	beego.Router("/cart/decCart", &frontend.CartController{}, "get:DecCart")
	beego.Router("/cart/delCart", &frontend.CartController{}, "get:DelCart")
	beego.Router("/cart/changeOneCart", &frontend.CartController{}, "get:ChangeOneCart")
	beego.Router("/cart/changeAllCart", &frontend.CartController{}, "get:ChangeAllCart")

	/// user
	beego.Router("/user", &frontend.UserController{})
	beego.Router("/user/order", &frontend.UserController{}, "get:OrderList")
	beego.Router("/user/orderinfo", &frontend.UserController{}, "get:OrderInfo")

	// product
	beego.Router("/category_:id([0-9]+).html", &frontend.ProductController{}, "get:CategoryList")
	beego.Router("/item_:id([0-9]+).html", &frontend.ProductController{}, "get:ProductItem")
	beego.Router("/product/getImgList", &frontend.ProductController{}, "get:GetImgList")

	/// buy
	beego.Router("/buy/checkout", &frontend.CheckoutController{}, "get:Checkout")
	beego.Router("/buy/doOrder", &frontend.CheckoutController{}, "post:GoOrder")
	beego.Router("/buy/confirm", &frontend.CheckoutController{}, "get:Confirm")
	beego.Router("/buy/orderPayStatus", &frontend.CheckoutController{}, "get:OrderPayStatus")

	/// address
	beego.Router("/address/addAddress", &frontend.AddressController{}, "post:AddAddress")
	beego.Router("/address/getOneAddressList", &frontend.AddressController{}, "get:GetOneAddressList")
	beego.Router("/address/goEditAddressList", &frontend.AddressController{}, "post:GoEditAddressList")
	beego.Router("/address/changeDefaultAddress", &frontend.AddressController{}, "get:ChangeDefaultAddress")

}
