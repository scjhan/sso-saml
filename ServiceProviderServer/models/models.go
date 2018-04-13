package models

import (
	"chenjunhan/sso-saml/proto"
	"chenjunhan/sso-saml/utils/util"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/astaxie/beego"
)

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

func CreateIdpCheckLoginURL(rawhost string, rawreturn string) string {
	rawq := url.Values{}
	rawq.Set("return_to", rawreturn)
	rawu := url.URL{
		Scheme:   "http",
		Host:     rawhost,
		Path:     "check_login_ret",
		RawQuery: rawq.Encode(),
	}

	q := url.Values{}
	q.Set("return_to", rawu.String())
	u := url.URL{
		Scheme:   "http",
		Host:     beego.AppConfig.String("idp::host"),
		Path:     "sso/check_login",
		RawQuery: q.Encode(),
	}

	return u.String()
}

func VerifyToken(token string) proto.TokenVerifyData {
	retval := proto.TokenVerifyData{}

	if len(token) == 0 {
		return retval
	}

	client := http.Client{}
	url := CreateVerifyTokenURL(token)

	if resp, err := client.Get(url); err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &retval)
	} else {
		util.Debug("VerifyToken client.Get error: " + err.Error())
	}

	return retval
}
