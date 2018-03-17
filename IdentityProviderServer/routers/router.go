package routers

import (
	"chenjunhan/sso-saml/IdentityProviderServer/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
