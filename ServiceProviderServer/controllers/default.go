package controllers

import (
	"chenjunhan/sso-saml/ServiceProviderServer/models"
	"chenjunhan/sso-saml/proto"
	"chenjunhan/sso-saml/utils/redis"
	"chenjunhan/sso-saml/utils/util"
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

type MainController struct {
	beego.Controller
}

func init() {
	redis.InitRedis("tcp", "127.0.0.1:6379")

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin, X-Requested-With, Content-Type, Accept"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["LoginUrl"] = "tb.com:7070/login"
	c.TplName = "index.tpl"
}

func (c *MainController) CheckLoginRet() {
	token := c.GetString("token", "")
	returnTo := c.GetString("return_to", "")

	vd := models.VerifyToken(token)
	if vd.Valid() {
		// create session
		session := proto.Session{
			SessionID: util.GetGUID(),
			UID:       vd.Uid,
			Name:      vd.Name,
		}
		if cache, err := json.Marshal(session); err == nil {
			redis.SetString(session.SessionID, string(cache), 60*60)
		}

		// set cookie
		c.Ctx.SetCookie("sessionid", session.SessionID)
		util.Debug("I have set cookie: " + session.SessionID)
	}

	// redirect to return_to
	c.Redirect(returnTo, 302)
}

func (c *MainController) Login() {
	return2 := c.GetString("return_to")

	ru := models.CreateIdpCheckLoginURL("tm.com:7070", return2)

	c.Redirect(ru, 302)
}

func (c *MainController) Index() {
	sessionid := c.Ctx.GetCookie("sessionid")
	util.Debug("Index get cookie: " + sessionid)

	unknown := false

	if len(sessionid) == 0 {
		unknown = true
	}

	cache, err := redis.GetString(sessionid)
	if err != nil {
		c.TplName = "500.tpl"
		return
	}
	if !unknown && len(cache) == 0 {
		unknown = true
	}

	session := proto.Session{}
	if json.Unmarshal([]byte(cache), &session) != nil {
		unknown = true
	}

	if unknown {
		c.Data["IdpLoginHref"] = "http://tm.com:7070/login?return_to=http%3A%2F%2Ftm.com%3A7070%2Findex"
		c.TplName = "index_unknown.tpl"
	} else {
		c.Data["UserName"] = session.Name
		c.TplName = "index.tpl"
	}
}

func (c *MainController) Test() {
	//http://tm.com/login?return_to=http%3A%2F%2Ftm.com%2Findex

}
