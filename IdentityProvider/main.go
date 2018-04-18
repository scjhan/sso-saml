package main

import (
	_ "chenjunhan/sso-saml/IdentityProvider/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}
