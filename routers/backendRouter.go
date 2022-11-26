package routers

import (
	"LeastMall/common"
	"LeastMall/controllers/backend"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/"+beego.AppConfig.String("adminPath"),
		beego.NSBefore(common.BackendAuth),
		beego.NSRouter("/", &backend.MainController{}),

		beego.NSRouter("/welcome", &backend.MainController{}, "get:Welcome"),

		/// login
		beego.NSRouter("/login", &backend.LoginController{}),
		beego.NSRouter("/login/gologin", &backend.LoginController{}, "post:GoLogin"),
		/// administrator
		beego.NSRouter("/administrator", &backend.AdministratorController{}),
		beego.NSRouter("/administrator/add", &backend.AdministratorController{}, "get:Add"),
		beego.NSRouter("/administrator/edit", &backend.AdministratorController{}, "get:Edit"),
		beego.NSRouter("/administrator/goadd", &backend.AdministratorController{}, "post:GoAdd"),
		beego.NSRouter("/administrator/goedit", &backend.AdministratorController{}, "post:GoEdit"),
		beego.NSRouter("/administrator/delete", &backend.AdministratorController{}, "get:Delete"),
	)
	beego.AddNamespace(ns)
}
