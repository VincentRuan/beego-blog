package task

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
	"github.com/vincent3i/beego-blog/handler"
	"github.com/vincent3i/beego-blog/models/rss"
	"io/ioutil"
	"time"
)

type Task struct {
	Id          string
	Expressions []string
}

type TaskSlice struct {
	Tasks []Task
}

type Job interface {
	doInterval(a ...interface{}) bool
}

func InitTasks() {
	filePath := fmt.Sprintf("%s/conf/%s", beego.AppPath, "task.json")
	beego.BeeLogger.Debug("Init scheduler with configuration file %s...", filePath)

	b, err := ioutil.ReadFile(filePath)
	if nil != err {
		beego.BeeLogger.Error("Unable to open/read file ---- [%s]", filePath)
		return
	}

	var taskSlice TaskSlice
	json.Unmarshal(b, &taskSlice)

	var c *cron.Cron
	for _, task := range taskSlice.Tasks {
		fmt.Println(task)
		if len(task.Expressions) == 0 {
			continue
		}
		c = cron.New()
		for _, expression := range task.Expressions {
			c.AddFunc(expression, p)
		}
		c.Start()
	}

	c = cron.New()
	c.AddFunc("0 30 * * * *", feed)
	c.Start()
}

func p() {
	fmt.Println(fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05")))
}

func feed() {
	rfs := rss.AllRssFeeder()
	for i, rf := range rfs {
		beego.BeeLogger.Debug("Fetch rss [%d] by url [%s]", i, rf.RSSUrl)
		//抓取RSS内容
		handler.BlogRssFeed(rf.RSSUrl, 3600, rf)
	}
}
