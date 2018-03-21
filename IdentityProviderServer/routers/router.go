package routers

import (
	"chenjunhan/sso-saml/IdentityProviderServer/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/sso/check_login", &controllers.MainController{}, "get:CheckLogin")
	beego.Router("/sso/login", &controllers.MainController{}, "post:Login")
}
