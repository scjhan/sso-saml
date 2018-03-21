package controllers

import (
	"encoding/json"
	"fmt"

	"chenjunhan/sso-saml/utils/mysql"

	"github.com/astaxie/beego"
)

func init() {
	mysql.InitMySQL("root", "123456", "sso")
}

type MainController struct {
	beego.Controller
}

// LoginArg LoginArg
type LoginArg struct {
	UserName    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	RedirectURL string `json:"redirect_url,omitempty"`
}

// CheckLogin check if the user has login in sso
func (c *MainController) CheckLogin() {
	if false /* has logined */ {
		// redirect with global session to allow login the subsys
	} else {
		// not login,r eturn the login page

		c.Data["Website"] = c.GetString("extra")
		c.TplName = "login_page.tpl"
	}
}

// Login login and create global session and redirect to subsystem
func (c *MainController) Login() {
	// check argemenents
	arg := LoginArg{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &arg); err != nil {
		// handler json unmarshal error
	}

	// check redirect url error
	if len(arg.RedirectURL) == 0 {

	}

	// check name && passwd
	o, err := mysql.NewMySQL()
	if err != nil {
		// handler mysql error
	}
	defer o.Close()

	query := fmt.Sprintf("select passwd from idp_user_info where name=%q and passwd=%q",
		arg.UserName, arg.Password)

	fmt.Println("query = ", query)

	result, num := o.Query(query)
	if num == 0 {
		// user not exiests or passwd error
	} else if arg.Password != result["passwd"][0].ToString() {
		// user passwd error
	} else {
		// ok, redirect
		c.Redirect(arg.RedirectURL, 302)
	}
}
