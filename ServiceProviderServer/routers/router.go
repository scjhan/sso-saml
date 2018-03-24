package routers

import (
	"chenjunhan/sso-saml/ServiceProviderServer/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.MainController{}, "get:Login")
	beego.Router("/index", &controllers.MainController{}, "get:Index")
}
