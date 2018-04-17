package main

import (
	"fmt"
)

func motify(key *string) {
	if key != nil {
		*key = "hello"
	}
}

func main() {
	key := "world"
	motify(&key)
	fmt.Println(key)
}
