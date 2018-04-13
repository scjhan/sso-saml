package main

import (
	"fmt"
	"net/url"

	"github.com/astaxie/beego"
)

func main() {
	//u, err := url.Parse("https://google.com/search?q=golang&return_to=https%3A%2F%2Fwww.baidu.com")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(u.Query().Get("return_to"))
	// u.Scheme = "https"
	// u.Host = "google.com"
	// q := u.Query()
	// q.Set("q", "golang")
	// q.Set("return_to", "https://www.baidu.com")
	// u.RawQuery = q.Encode()
	// fmt.Println(u)

	token := "token"
	q := url.Values{}
	q.Set("token", token)
	u := url.URL{
		Scheme:   "http",
		Host:     beego.AppConfig.String("idp::host"),
		Path:     "sso/verify_token",
		RawQuery: q.Encode(),
	}

	fmt.Println(u.String())

	us2 := "https%3A%2F%2Fwww.baidu.com"
	u2, _ := url.Parse(us2)
	fmt.Println(u2.Path)

	us3 := "/search?q=golang&return_to=https%3A%2F%2Fwww.baidu.com"
	u3, _ := url.Parse(us3)
	fmt.Println(u3.Path)
}
