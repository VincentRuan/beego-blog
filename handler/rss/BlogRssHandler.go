// BlogRssHandler
package rsshandler

import (
	"github.com/astaxie/beego"
	rss "github.com/jteeuwen/go-pkg-rss"
	"github.com/vincent3i/beego-blog/utilities/mongo"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
)

type BlogRssHandler struct {
}

func (this *BlogRssHandler) ProcessChannels(feed *rss.Feed, newchannels []*rss.Channel) {
	beego.BeeLogger.Debug("%d new channel(s) in %s\n", len(newchannels), feed.Url)
}

func (this *BlogRssHandler) ProcessItems(feed *rss.Feed, ch *rss.Channel, newitems []*rss.Item) {
	beego.BeeLogger.Debug("%d new item(s) in %s\n", len(newitems), feed.Url)
}
