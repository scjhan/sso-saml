package controllers

import (
	"chenjunhan/sso-saml/TBInc/models"
	"chenjunhan/sso-saml/proto"
	"chenjunhan/sso-saml/utils/log"
	"chenjunhan/sso-saml/utils/redis"
	"chenjunhan/sso-saml/utils/util"
	"encoding/json"
	"io/ioutil"

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

	log.InitLogService(beego.AppConfig.String("appname"), "/c/Users/Han/Desktop/logs")
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
			key1 := models.CreateRedisKey(session.SessionID, models.SessionKey)
			key2 := models.CreateRedisKey(session.UID, models.UID2SessionKey)
			redis.SetString(key1, string(cache), 60*60)
			redis.SetString(key2, session.SessionID, 60*60) // make sure we can find a sesionid by uid
			log.Debug("check_login_ret, key = " + key1 + ";" + key2)
		}

		// set cookie
		c.Ctx.SetCookie("sessionid", session.SessionID)
		log.Debug("I have set cookie: " + session.SessionID)
	}

	// redirect to return_to
	c.Redirect(returnTo, 302)
}

func (c *MainController) Login() {
	return2 := c.GetString("return_to")

	ru := models.CreateIdpCheckLoginURL(return2)

	c.Redirect(ru, 302)
}

func (c *MainController) Index() {
	sessionid := c.Ctx.GetCookie("sessionid")
	log.Debug("Index get cookie: " + sessionid)

	unknown := false

	if len(sessionid) == 0 {
		unknown = true
	}

	cache, err := redis.GetString(models.CreateRedisKey(sessionid, models.SessionKey))
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
		c.Data["IdpLoginHref"] = "/login?return_to=http%3A%2F%2Ftb.com%3A6060%2Findex"
		c.TplName = "index_unknown.tpl"
	} else {
		c.Data["UserName"] = session.Name
		c.TplName = "index.tpl"
	}
}

func (c *MainController) Logout() {
	sessionid := c.Ctx.GetCookie("sessionid")
	if len(sessionid) != 0 {
		c.Ctx.SetCookie("sessionid", "")
		sessionStr, _ := redis.GetString(models.CreateRedisKey(sessionid, models.SessionKey))
		redis.Delete(models.CreateRedisKey(sessionid, models.SessionKey))

		session := proto.Session{}
		if err := json.Unmarshal([]byte(sessionStr), &session); err == nil && len(session.UID) != 0 {
			redis.Delete(models.CreateRedisKey(session.UID, models.UID2SessionKey))
		}

		models.IdpLogout(session.UID)
	}
	c.Redirect("/index", 302)
}

func (c *MainController) Test() {
	c.Ctx.WriteString(c.Ctx.Request.Host)
}

func (c *MainController) Push() {
	body, _ := ioutil.ReadAll(c.Ctx.Request.Body)
	msg := proto.ToPushMsg(body)
	resp := proto.PushMsg{
		Type: proto.Ok,
	}

	for i := 0; i < 1; i++ {
		if msg.Type == proto.Error {
			resp.Type = proto.Error
			break
		}

		if msg.Type == proto.IdpLogout {
			sessionid, _ := redis.GetString(models.CreateRedisKey(msg.Content, models.UID2SessionKey))
			redis.Delete(models.CreateRedisKey(msg.Content, models.UID2SessionKey))
			if len(sessionid) != 0 {
				redis.Delete(models.CreateRedisKey(sessionid, models.SessionKey))
			}
			log.Debug("idp notify logout done")
		}
	}

	c.Ctx.WriteString(resp.String())

}
