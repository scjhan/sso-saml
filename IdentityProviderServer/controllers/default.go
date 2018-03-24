package controllers

import (
	"encoding/json"
	"fmt"

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
	//RedirectURL string `json:"redirect_url,omitempty"`
}

// CheckLogin check if the user has login in sso
func (c *MainController) CheckLogin() {
	util.Debug("Idp CheckLogin, args = " + c.Input().Encode())

	user := c.GetString("uid")
	expire, err := redis.GetString(user)
	if err != nil {
		c.TplName = "500.tpl"
		c.Data["ErrorMsg"] = "CheckLogin redis get error:" + err.Error()
		return
	}

	returnTo := c.GetString("return_to")
	if len(returnTo) == 0 {
		c.TplName = "500.tpl"
		c.Data["ErrorMsg"] = "CheckLogin redirect but 'return_to' is null, raw url = " + c.Input().Encode()
		return
	}

	if len(expire) != 0 {
		c.Redirect(returnTo, 302)
	} else {
		redirectUrl := fmt.Sprintf("http://idp.com:9090/sso/login_page/?return_to=%s", returnTo)
		c.Redirect(redirectUrl, 302)
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

	util.Debug("Idp Login returnTo = " + returnTo)

	// check argemenents
	arg := LoginArg{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &arg); err != nil {
		c.TplName = "500.tpl"
		c.Data["ErrorMsg"] = "Login 'return_to' is null"
		return
	}

	// check name && passwd
	o, err := mysql.NewMySQL()
	if err != nil {
		c.TplName = "500.tpl"
		c.Data["ErrorMsg"] = "Login 'return_to' is null"
		return
	}
	defer o.Close()

	query := fmt.Sprintf("select passwd from idp_user_info where name=%q and passwd=%q",
		arg.UserName, arg.Password)

	fmt.Println("query = ", query)

	result, num := o.Query(query)
	if num == 0 {
		// user not exiests or passwd error
		util.Debug("Idp Login Sql query num = 0")
	} else if arg.Password != result["passwd"][0].ToString() {
		// user passwd error
		util.Debug("Idp Login Sql query password error")
	} else {
		// ok, redirect
		util.Debug("idp login durl:" + returnTo)
		c.Redirect(returnTo, 302)
	}
}
