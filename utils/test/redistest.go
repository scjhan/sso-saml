package main

import (
	"chenjunhan/sso-saml/utils/redis"
	"fmt"
)

func init() {
	redis.InitRedis("tcp", "127.0.0.1:6379")
}

func main() {
	value, err := redis.GetString("hello")
	if err != nil {
		fmt.Printf("redis error, error = %s", err.Error())
	} else {
		fmt.Println("value = ", value)
	}

	err = redis.SetString("good", "very good", 20)
	if err != nil {
		fmt.Printf("redis error, error = %s", err.Error())
	}

	ret, err := redis.SMembers("a")
	for _, v := range ret {
		fmt.Print(v, ";")
	}
}
