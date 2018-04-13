package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"

	"chenjunhan/sso-saml/proto"
	"chenjunhan/sso-saml/utils/mysql"
	"chenjunhan/sso-saml/utils/redis"
	"chenjunhan/sso-saml/utils/util"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	mysql.InitMySQL("root", "123456", "sso")
	redis.InitRedis("tcp", "127.0.0.1:6379")

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin, X-Requested-With, Content-Type, Accept"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
}

type MainController struct {
	beego.Controller
}

// LoginArg LoginArg
type LoginArg struct {
	UserName string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginRet struct {
	Code int    `json:"code"`
	Href string `json:"href"`
}

// CheckLogin check if the user has login in sso
func (c *MainController) CheckLogin() {
	sessionid := c.Ctx.GetCookie("sessionid")

	if len(sessionid) == 0 {
		// not login
		c.Data["SsoLoginUrl"] = "http://idp.com:9090/sso/login?return_to=" + c.GetString("return_to")
		c.TplName = "login_page.tpl"
	} else {
		token, _ := redis.GetString(sessionid)

		ret2 := c.GetString("return_to")
		u2, _ := url.Parse(ret2)
		q2 := u2.Query()
		q2.Add("token", token)
		u2.RawQuery = q2.Encode()

		c.Redirect(u2.String(), 302)
	}
}

func (c *MainController) LoginPage() {
	returnTo := c.GetString("return_to")
	if len(returnTo) == 0 {
		c.TplName = "500.tpl"
		c.Data["ErrorMsg"] = "LoginPage 'return_to' is null"
		return
	}

	c.TplName = "login_page.tpl"
	c.Data["SsoLoginUrl"] = fmt.Sprintf("idp.com:9090/sso/login/?return_to=%s", returnTo)
}

// Login login and create global session and redirect to subsystem
func (c *MainController) Login() {
	returnTo := c.GetString("return_to")
	if len(returnTo) == 0 {
		c.TplName = "500.tpl"
		c.Data["ErrorMsg"] = "Login 'return_to' is null"
		return
	}

	// check argemenents
	arg := LoginArg{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &arg); err != nil {
		c.TplName = "500.tpl"
		c.Data["ErrorMsg"] = "RequestBody marshal error"
		return
	}

	// check name && passwd
	o, err := mysql.NewMySQL()
	if err != nil {
		c.TplName = "500.tpl"
		c.Data["ErrorMsg"] = "open mysql error"
		return
	}
	defer o.Close()

	query := fmt.Sprintf("select uid, passwd from idp_user_info where name=%q limit 1", arg.UserName)

	result, num := o.Query(query)
	if num == 0 {
		// user not exiests
		ret := LoginRet{Code: 1}
		retByte, _ := json.Marshal(ret)
		c.Ctx.WriteString(string(retByte))
	} else if arg.Password != result["passwd"][0].ToString() {
		// user passwd error
		ret := LoginRet{Code: 2}
		retByte, _ := json.Marshal(ret)
		c.Ctx.WriteString(string(retByte))
	} else {
		// create session
		sessionid := util.GetGUID()
		token := proto.TokenVerifyData{
			Token: util.GetGUID(),
			Uid:   result["uid"][0].ToString(),
			Name:  arg.UserName,
		}
		bytes, _ := json.Marshal(token)

		redis.SetString(sessionid, token.Token, 60*60)
		redis.SetString(token.Token, string(bytes), 60*60)

		// set cookie
		c.Ctx.SetCookie("sessionid", sessionid)

		u, _ := url.Parse(returnTo)
		q := u.Query()
		q.Add("token", token.Token)
		u.RawQuery = q.Encode()

		ret := LoginRet{
			Code: 0,
			Href: u.String(),
		}

		util.Debug("idp login href: " + ret.Href)

		retByte, _ := json.Marshal(ret)
		c.Ctx.WriteString(string(retByte))
	}
}

func (c *MainController) VerifyToken() {
	token := c.GetString("token", "")

	tokenJson, _ := redis.GetString(token)

	c.Ctx.WriteString(tokenJson)
}
