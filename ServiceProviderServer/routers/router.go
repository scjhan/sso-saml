package routers

import (
	"chenjunhan/sso-saml/ServiceProviderServer/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
