package controllers

import (
	"github.com/vincent3i/beego-blog/models/rss"
)

type RSSController struct {
	AdminController
}

func (this *RSSController) LoadPage() {
	this.Layout = "layout/admin.html"
	this.TplNames = "rss/rss.html"
}

func (this *RSSController) Read() {
	limit, _ := this.GetInt("limit")
	offset, _ := this.GetInt("offset")
	if limit == 0 {
		limit = -1
	}
	rfs := rss.AllRssFeeders(limit, offset)

	result := make(map[string]interface{})
	result["total"] = len(rfs)
	result["rows"] = rfs
	this.Data["json"] = result

	this.ServeJson()
}
