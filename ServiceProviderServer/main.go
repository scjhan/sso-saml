package main

import (
	_ "chenjunhan/sso-saml/ServiceProviderServer/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

