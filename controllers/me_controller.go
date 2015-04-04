package controllers

import (
	"github.com/vincent3i/beego-blog/engine"
)

type MeController struct {
	AdminController
}

func (this *MeController) Default() {
	this.Layout = "layout/admin.html"
	this.TplNames = "me/default.html"
}

func (this *MeController) ShutDown() {
	engine.CloseSearcher()
	this.Ctx.WriteString("已关闭索引项,可以shutdown server")
}
