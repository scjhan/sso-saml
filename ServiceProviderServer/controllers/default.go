package controllers

import (
	"chenjunhan/sso-saml/utils/redis"
	"chenjunhan/sso-saml/utils/util"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

type MainController struct {
	beego.Controller
}

//var gSessionsMgr *session.Manager

func init() {
	redis.InitRedis("tcp", "127.0.0.1:6379")

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin, X-Requested-With, Content-Type, Accept"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	// beego.BConfig.WebConfig.Session.SessionOn = true
	// sessionCfg := &session.ManagerConfig{
	// 	CookieName:      "tbsession",
	// 	EnableSetCookie: true,
	// 	Gclifetime:      3600,
	// 	CookieLifeTime:  3600,
	// }
	//var err error
	//gSessionsMgr, err = session.NewManager("cookie", sessionCfg)
	//if err != nil {
	//	fmt.Println("NewManager error: ", err.Error())
	//}
	//go gSessionsMgr.GC()
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["LoginUrl"] = "tb.com:7070/login"
	c.TplName = "index.tpl"
}

func (c *MainController) Login() {
	uid := c.Ctx.GetCookie("uid")
	if len(uid) != 0 {
		cache, err := redis.GetString(uid)
		if err != nil {
			c.TplName = "500.tpl"
			return
		} else if len(cache) != 0 {
			redis.SetString(uid, uid, 3600)
			c.Data["UserName"] = uid
			c.TplName = "res.tpl"
			return
		}
	}

	url := fmt.Sprintf("http://idp.com:9090/sso/check_login/?uid=%s&return_to=%s",
		uid, "http://tb.com:7070/index")
	util.Debug("need check_login, url = " + url)
	c.Redirect(url, 302)
}

func (c *MainController) Index() {
	//c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	c.Data["UserName"] = c.GetString("username", "default")
	c.TplName = "res.tpl"
}
