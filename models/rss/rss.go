package rss

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	rss "github.com/jteeuwen/go-pkg-rss"
	. "github.com/vincent3i/beego-blog/models"
	"github.com/vincent3i/beego-blog/utilities/mongo"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
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

func Update(rf *RssFeeder) error {
	if rf.Id == 0 {
		return fmt.Errorf("primary key:id not set")
	}

	_, err := orm.NewOrm().Update(rf)
	if err != nil {
		beego.Error(err)
	}

	return err
}

func AllRssFeeders(limit, offset int) []RssFeeder {
	var rfs []RssFeeder
	_, err := RssFeeders().OrderBy("-SubscribeTime").Limit(limit, offset).All(&rfs)
	if err != nil {
		beego.Error(err)
		return make([]RssFeeder, 0)
	}

	return rfs
}

func RssFeeders() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(RssFeeder))
}

type BlogRssHandler struct {
}

func (this *BlogRssHandler) ProcessChannels(feed *rss.Feed, newchannels []*rss.Channel) {
	beego.BeeLogger.Debug("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func (this *BlogRssHandler) ProcessItems(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	beego.BeeLogger.Debug("%d new item(s) in %s\n", len(newitems), feed.Url)
	save(ch)
}

func BlogRssFeed(uri string, timeout int) {
	blogRssHandler := BlogRssHandler{}

	feed := rss.NewWithHandlers(timeout, true, &blogRssHandler, &blogRssHandler)
	if err := feed.Fetch(uri, rssFetchCharset); err != nil {
		beego.BeeLogger.Error("[e] %s: %s", uri, err)
		return
	}
}

func rssFetchCharset(charset string, input io.Reader) (io.Reader, error) {
	beego.Debug(charset)
	var r *transform.Reader
	switch charset {
	case "iso-8859-1":
		r = transform.NewReader(input, charmap.Windows1252.NewDecoder())
	case "gbk":
		r = transform.NewReader(input, simplifiedchinese.GBK.NewDecoder())
	case "gb2312":
		r = transform.NewReader(input, simplifiedchinese.GBK.NewDecoder())
	}

	return r, nil
}

func save(ch *rss.Channel) {
	session, err := mongo.CopyMonotonicSession()
	if err != nil {
		panic(err)
	}
	defer mongo.CloseSession(session)

	dbNames, _ := session.DatabaseNames()
	beego.BeeLogger.Debug("%v", dbNames)

	//c := session.DB("beego_blog").C("channel")

	var chs []rss.Channel
	f := func(collection *mgo.Collection) error {
		queryMap := bson.M{}
		return collection.Find(queryMap).All(&chs)
	}
	mongo.Execute(session, "", "channel", f)

	beego.BeeLogger.Debug("%d", len(chs))
	for _, ch := range chs {
		beego.BeeLogger.Debug("%s", ch)
	}
	//	err = c.Insert(ch)
	//	if err != nil {
	//		panic(err)
	//	}

	//	info, err := c.RemoveAll(bson.M{"name": "Ale"})
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	fmt.Println(info.Removed)

	//	result := rss.Channel{}
	//	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	//	if err != nil {
	//		panic(err)
	//		//fmt.Println(err)
	//	}

	//	fmt.Println("Phone:", result.Phone)
}
