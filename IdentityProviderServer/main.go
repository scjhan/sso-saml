package main

import (
	_ "chenjunhan/sso-saml/IdentityProviderServer/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
