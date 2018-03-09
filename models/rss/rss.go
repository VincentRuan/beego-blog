package rss

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	. "github.com/vincentruan/beego-blog/models"
	"net/url"
	"strings"
	"time"
)

func Save(this *RssFeeder) (*RssFeeder, error) {
	if CheckUrl(this.RSSUrl) {
		return nil, fmt.Errorf("Url [%s] is illegal.", this.RSSUrl)
	}

	or := orm.NewOrm()
	this.CreateTime = time.Now()
	this.UpdateTime = time.Now()
	rfId, e := or.Insert(this)
	if e != nil {
		beego.Error(e)
		return nil, e
	}

	this.Id = rfId

	return this, nil
}

//检查输入的RSS地址是否合法
func CheckUrl(rssUrl string) bool {
	if rssUrl == "" {
		beego.Error("Parameter can not be empty!")
		return false
	}

	uri, err := url.ParseRequestURI(rssUrl)
	if err != nil {
		beego.Error("Unsupport URL format type, ", rssUrl, err)
		return false
	}

	var rfs []RssFeeder
	count, err := RssFeeders().All(&rfs, "rss_url")

	if err != nil {
		beego.Error(err)
		return false
	}

	if count > 0 {
		var uriInDB *url.URL
		for _, feeder := range rfs {
			uriInDB, err = url.ParseRequestURI(feeder.RSSUrl)

			if err != nil {
				beego.Error(err)
			}

			if strings.EqualFold(uri.Host, uriInDB.Host) && uri.Path == uriInDB.Path {
				beego.BeeLogger.Debug("Rss URL [%s] had existed in DB!", feeder.RSSUrl)
				return false
			}
		}
	}

	return true
}

func Del(rf *RssFeeder) error {
	_, err := orm.NewOrm().Delete(rf)
	if err != nil {
		beego.Error(err)
		return err
	}

	return nil
}

func Update(rf *RssFeeder, fields ...string) error {
	if rf.Id == 0 {
		return fmt.Errorf("primary key:id not set")
	}

	_, err := orm.NewOrm().Update(rf, fields...)
	if err != nil {
		beego.Error(err)
	}

	return err
}

func SearchRssFeeders(search, sort, order string, limit, offset int) []RssFeeder {
	var rfs []RssFeeder

	qSetter := RssFeeders().Limit(limit, offset)
	if search != "" {
		qSetter = qSetter.Filter("rss_desc__icontains", search)
	}
	if sort != "" {
		var ob string
		if strings.EqualFold(order, "DESC") {
			ob = "-"
		}
		qSetter = qSetter.OrderBy(ob + sort)
	} else {
		qSetter = qSetter.OrderBy("-subscribe_time")
	}
	_, err := qSetter.All(&rfs)
	if err != nil {
		beego.Error(err)
		return make([]RssFeeder, 0)
	}

	return rfs
}

func AllRssFeeder() []RssFeeder {
	var rfs []RssFeeder

	_, err := RssFeeders().OrderBy("-SubscribeTime").Limit(-1).All(&rfs)
	if err != nil {
		beego.Error(err)
		return make([]RssFeeder, 0)
	}

	return rfs
}

func RssFeeders() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(RssFeeder))
}
