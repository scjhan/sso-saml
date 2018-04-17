package routers

import (
	"chenjunhan/sso-saml/ServiceProviderServer/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.MainController{}, "get:Login")
	beego.Router("/index", &controllers.MainController{}, "get:Index")
	beego.Router("/check_login_ret", &controllers.MainController{}, "get:CheckLoginRet")
	//beego.Router("/idp_notify", &controllers.MainController{}, "get:IdpNotify")
	beego.Router("/push", &controllers.MainController{}, "post:Push")
	beego.Router("/logout", &controllers.MainController{}, "get:Logout")
	beego.Router("/test", &controllers.MainController{}, "get:Test")
}
