package main

import (
	"chenjunhan/sso-saml/utils/util"
	"fmt"
)

func main() {
	for i := 0; i < 10; i++ {
		fmt.Println("guid = ", util.GetGUID())
	}
}
