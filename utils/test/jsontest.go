package main

import (
	"encoding/json"
	"fmt"
)

type Msg struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	msg := Msg{
		Id:   1,
		Name: "one",
	}

	if b, err := json.Marshal(msg); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(b))
	}
}
