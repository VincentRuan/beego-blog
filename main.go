package main

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/vincent3i/beego-blog/g"
	"github.com/vincent3i/beego-blog/models/rss"
	_ "github.com/vincent3i/beego-blog/routers"
	"github.com/vincent3i/beego-blog/task"
	"github.com/vincent3i/beego-blog/utilities/mongo"
)

func main() {
	g.InitEnv()
	beego.BeeLogger.Debug("Is session on [%t], session provide --->>> [%s]", beego.SessionOn, beego.SessionProvider)
	task.InitTasks()
	mongo.Startup()
	rss.BlogRssFeed("http://blog.case.edu/news/feed.atom", 5)
	beego.Run()
}
