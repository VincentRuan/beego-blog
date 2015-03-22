package handler

import (
	"github.com/astaxie/beego"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/vincent3i/beego-blog/models"
	"github.com/vincent3i/beego-blog/utilities/mongo"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gopkg.in/mgo.v2/bson"
	"io"
)

type BlogRssHandler struct {
	RSSFeeder models.RssFeeder
}

type RSSData struct {
	ID_       bson.ObjectId `bson:"_id,omitempty"`
	RSSFeeder models.RssFeeder
	RSSItems  []*rss.Item
}

func (this *BlogRssHandler) ProcessChannels(feed *rss.Feed, newchannels []*rss.Channel) {
	beego.BeeLogger.Debug("%d new channel(s) in %s", len(newchannels), feed.Url)
}

func (this *BlogRssHandler) ProcessItems(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	beego.BeeLogger.Debug("%d new item(s) in %s", len(newitems), feed.Url)
	rf := RSSData{RSSFeeder: this.RSSFeeder, RSSItems: newitems}
	SaveOrUpdate(&rf)
}

func BlogRssFeed(uri string, timeout int, rssFeeder models.RssFeeder) {
	blogRssHandler := BlogRssHandler{RSSFeeder: rssFeeder}

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

func SaveOrUpdate(rf *RSSData) {
	session, err := mongo.CopyMonotonicSession()
	if err != nil {
		beego.Error(err)
	}
	defer mongo.CloseSession(session)

	c := session.DB("beego_blog").C("rssdata")

	rd := RSSData{}
	err = c.Find(bson.M{"rssfeeder.id": rf.RSSFeeder.Id}).One(&rd)
	if err != nil {
		beego.Error(err)

		err = c.Insert(rf)
		if err != nil {
			beego.Error(err)
		}

		return
	}

	oid := rd.ID_

	//开启事务
	//	txn.SetDebug(true)
	//	txn.SetLogger(log.New(os.Stderr, "", log.LstdFlags))
	//	runner := txn.NewRunner(c)
	//	ops := []txn.Op{{
	//		C:      "rssdata",
	//		Id:     oid,
	//		Remove: true,
	//	}, {
	//		C:      "rssdata",
	//		Id:     bson.NewObjectId(),
	//		Insert: rf,
	//	}}

	//	err = runner.Run(ops, "", nil)
	//	if err != nil {
	//		beego.Error(err)
	//	}

	err = c.Remove(bson.M{"_id": oid})
	if err != nil {
		beego.Error(err)
		return
	}

	err = c.Insert(rf)
	if err != nil {
		beego.Error(err)
	}
}
