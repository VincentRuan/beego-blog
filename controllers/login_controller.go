package controllers

import (
	"github.com/astaxie/beego"
	"github.com/vincent3i/beego-blog/models/admin"
)

type LoginController struct {
	BaseController
}

func (this *LoginController) Login() {
	if this.IsAdmin {
		beego.Debug(this.GetSession("username").(string), "had been login, redirect to home page!")
		this.Redirect("/", 302)
	}
	this.TplNames = "login/login.html"
}

func (this *LoginController) DoLogin() {
	name := this.GetString("name")
	if name == "" {
		this.Ctx.WriteString("name is blank")
		return
	}
	password := this.GetString("password")
	if password == "" {
		this.Ctx.WriteString("password is blank")
		return
	}

	if ok := admin.CheckUser(name, password); ok {
		beego.BeeLogger.Debug("User [%s] login success.", name)
	}

	this.Ctx.SetCookie("bb_name", name, 2592000, "/")
	this.Ctx.ResponseWriter.Header().Add("Set-Cookie", "bb_password="+password+"; Max-Age=2592000; Path=/; httponly")

	beego.BeeLogger.Debug("Put user name [%s] into session", name)
	this.SetSession("username", name)
	this.Ctx.WriteString("")
}

func (this *LoginController) Logout() {
	this.Ctx.SetCookie("bb_name", "", 0, "/")
	this.Ctx.ResponseWriter.Header().Add("Set-Cookie", "bb_password=; Max-Age=0; Path=/; httponly")

	this.DelSession("username")

	//this.Ctx.WriteString("logout")
	this.Redirect("/", 302)
}
