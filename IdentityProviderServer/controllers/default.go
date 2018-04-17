package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"

	"chenjunhan/sso-saml/IdentityProviderServer/models"
	"chenjunhan/sso-saml/proto"
	"chenjunhan/sso-saml/utils/log"
	"chenjunhan/sso-saml/utils/mysql"
	"chenjunhan/sso-saml/utils/redis"
	"chenjunhan/sso-saml/utils/util"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
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

	//add interceptor
	beego.InsertFilter("/sso/*", beego.BeforeExec, func(ctx *context.Context) {
		util.Debug("Referer = " + ctx.Request.Referer())
		//host := strings.Split(ctx.Request.Host, ":")[0]
		// if !models.CheckHost(ctx.Input.Host()) {
		// 	util.Debug("beego.InsertFilter host = " + ctx.Input.IP())
		// 	ctx.Redirect(302, "/not_allowed")
		// }
	})
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

func (c *MainController) NotAllowed() {
	c.TplName = "403.tpl"
}

// CheckLogin check if the user has login in sso
func (c *MainController) CheckLogin() {
	sessionid := c.Ctx.GetCookie("sessionid")

	if len(sessionid) != 0 {
		token, _ := redis.GetString(models.CreateRedisKey(sessionid, models.SessionTokenKey))

		if len(token) != 0 {
			ret2 := c.GetString("return_to")
			u2, _ := url.Parse(ret2)
			q2 := u2.Query()
			q2.Add("token", token)
			u2.RawQuery = q2.Encode()

			c.Redirect(u2.String(), 302)
		} else {
			// means has logout but local session isn't deleted
			c.Ctx.SetCookie("sessionid", "")
		}
	}

	c.Data["SsoLoginUrl"] = "/sso/login?return_to=" + c.GetString("return_to")
	c.TplName = "login_page.tpl"
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

		redis.SetString(models.CreateRedisKey(sessionid, models.SessionTokenKey), token.Token, 60*60)
		redis.SetString(models.CreateRedisKey(token.Token, models.TokenValueKey), string(bytes), 60*60)
		redis.SetString(models.CreateRedisKey(token.Uid, models.UIDSessionKey), sessionid, 60*60)

		// set cookie
		c.Ctx.SetCookie("sessionid", sessionid)

		// cache host
		host := strings.Split(c.Ctx.Request.Host, ":")[0]
		redis.SetAdd(models.CreateRedisKey(token.Uid, models.HostSetKey), host)

		u, _ := url.Parse(returnTo)
		q := u.Query()
		q.Add("token", token.Token)
		u.RawQuery = q.Encode()

		ret := LoginRet{
			Code: 0,
			Href: u.String(),
		}

		retByte, _ := json.Marshal(ret)
		c.Ctx.WriteString(string(retByte))
	}
}

func (c *MainController) Push() {
	body, _ := ioutil.ReadAll(c.Ctx.Request.Body)
	msg := proto.ToPushMsg(body)
	resp := proto.PushMsg{
		Type: proto.Ok,
	}

	log.Debug(string(body))

	for i := 0; i < 1; i++ {
		if msg.Type == proto.Error {
			resp.Type = proto.Error
			break
		}

		if msg.Type == proto.ClusterLogout {
			uid := msg.Content

			models.DeleteUIDCache(uid)
			models.NotifyLogout(uid)
			break
		}

		if msg.Type == proto.ClusterVerifyToken {
			tokenStr, _ := redis.GetString(models.CreateRedisKey(msg.Content, models.TokenValueKey))
			resp.Content = tokenStr
			break
		}
	}

	c.Ctx.WriteString(resp.String())
}
