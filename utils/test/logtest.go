package main

import (
	"chenjunhan/sso-saml/utils/log"
)

func init() {
	log.InitLogService("logtest", ".")
}

func main() {
	log.Debug("hello world")
	log.Error("hello world")
	log.Warning("hello world")
	log.Info("hello world")
}
