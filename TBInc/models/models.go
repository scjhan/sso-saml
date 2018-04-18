package models

import (
	"bytes"
	"chenjunhan/sso-saml/proto"
	"chenjunhan/sso-saml/utils/log"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/astaxie/beego"
)

func init() {
	log.InitLogService(beego.AppConfig.String("appname"), "/c/Users/Han/Desktop/logs")
}

func GetHost() string {
	host := beego.AppConfig.String("host")
	if len(host) != 0 {
		return host
	}
	return "tb.com"
}

func GetPort() string {
	port := beego.AppConfig.String("httpport")
	if len(port) != 0 {
		return port
	}
	return "6060"
}

func GetHostPort() string {
	return GetHost() + ":" + GetPort()
}

func CreateVerifyTokenURL(token string) string {
	q := url.Values{}
	q.Set("token", token)
	u := url.URL{
		Scheme:   "http",
		Host:     beego.AppConfig.String("idp::host"),
		Path:     "sso/verify_token",
		RawQuery: q.Encode(),
	}

	return u.String()
}

func CreateIdpCheckLoginURL(rawreturn string) string {
	rawq := url.Values{}
	rawq.Set("return_to", rawreturn)
	rawu := url.URL{
		Scheme:   "http",
		Host:     GetHostPort(),
		Path:     "check_login_ret",
		RawQuery: rawq.Encode(),
	}

	q := url.Values{}
	q.Set("host", GetHost())
	q.Set("return_to", rawu.String())
	u := url.URL{
		Scheme:   "http",
		Host:     beego.AppConfig.String("idp::host"),
		Path:     "sso/check_login",
		RawQuery: q.Encode(),
	}

	return u.String()
}

func CreateIdpPushURL() string {
	u := url.URL{
		Scheme: "http",
		Host:   beego.AppConfig.String("idp::host"),
		Path:   "/sso/push",
	}

	return u.String()
}

func VerifyToken(token string) proto.TokenVerifyData {
	retval := proto.TokenVerifyData{}

	c := http.Client{}
	u := CreateIdpPushURL()
	msg := proto.PushMsg{
		Type:    proto.ClusterVerifyToken,
		Content: token,
	}

	if r, err := c.Post(u, "application/json", bytes.NewReader([]byte(msg.String()))); err == nil {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		resp := proto.PushMsg{}
		if err = json.Unmarshal(body, &resp); err == nil && resp.Type == proto.Ok {
			json.Unmarshal([]byte(resp.Content), &retval)
			log.Debug(resp.String())
		} else {
			log.Debug(err.Error())
		}
	} else {
		log.Debug(err.Error())
	}

	log.Debug(fmt.Sprintf("%s", retval))

	return retval
}

type KeyType string

const (
	SessionKey     = KeyType("S")
	UID2SessionKey = KeyType("U2S")
)

func CreateRedisKey(key string, kt KeyType) string {
	return beego.AppConfig.String("appname") + "_" + string(kt) + "_" + key
}

func IdpLogout(uid string) {
	c := http.Client{}
	u := CreateIdpPushURL()
	msg := proto.PushMsg{
		Type:    proto.ClusterLogout,
		Content: uid,
	}

	c.Post(u, "application/json", bytes.NewReader([]byte(msg.String())))
}
