package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/vincent3i/beego-blog/handler"
	. "github.com/vincent3i/beego-blog/models"
	"github.com/vincent3i/beego-blog/models/rss"
	"github.com/vincent3i/beego-blog/utilities/mongo"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

type RSSController struct {
	AdminController
}

func (this *RSSController) LoadPage() {
	this.Data["IsRSS"] = true
	this.Layout = "layout/admin.html"
	this.TplNames = "rss/rss.html"
}

func (this *RSSController) Read() {
	limit, offset := this.GetPaginationParam()

	rfs := rss.SearchRssFeeders(this.GetString("search"), this.GetString("sort"), this.GetString("order"), limit, offset)

	result := make(map[string]interface{})
	result["total"] = len(rfs)
	result["rows"] = rfs
	this.Data["json"] = result

	this.ServeJson()
}

func (this *RSSController) DoEdit() {
	rf := RssFeeder{}
	this.ParseForm(&rf)

	result := make(map[string]interface{})
	valid := validation.Validation{}
	b, err := valid.Valid(&rf)
	if err != nil {
		result["status"] = "fail"
		result["result"] = err.Error()
	} else {
		if b {
			rf.UpdateTime = time.Now()
			err = rss.Update(&rf, "rss_desc", "rss_url", "update_time")
			if err != nil {
				result["status"] = "fail"
				result["result"] = err.Error()
			} else {
				beego.Debug("update RssFeeder[", rf.Id, "] success")
				result["status"] = "success"
			}
		} else {
			//验证不通过
			result["status"] = "fail"
			var errStr string
			for _, err := range valid.Errors {
				errStr += err.Key + err.Message
			}
			result["result"] = errStr
		}
	}
	this.Data["json"] = result
	this.ServeJson()
}

func (this *RSSController) DoDel() {
	id, err := this.GetInt64("id")
	result := make(map[string]interface{})
	if err != nil {
		result["status"] = "fail"
		result["result"] = err.Error()
	} else {
		rf := RssFeeder{Id: id}
		err = rss.Del(&rf)
		if err != nil {
			result["status"] = "fail"
			result["result"] = err.Error()
		} else {
			beego.Debug("remove RssFeeder[", id, "] success")
			result["status"] = "success"
		}
	}
	this.Data["json"] = result
	this.ServeJson()
}

func (this *RSSController) RSSData() {
	id, _ := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)

	session, err := mongo.CopyMonotonicSession()
	if err != nil {
		beego.Error(err)
	}
	defer mongo.CloseSession(session)
	c := session.DB("beego_blog").C("rssdata")

	rd := handler.RSSData{}
	err = c.Find(bson.M{"rssfeeder.id": id}).One(&rd)
	if err != nil {
		beego.Error(err)
	}
	this.Data["json"] = rd
	this.ServeJson()
}
