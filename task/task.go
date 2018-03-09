package task

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/robfig/cron"
	"github.com/vincentruan/beego-blog/g"
	"github.com/vincentruan/beego-blog/handler"
	"github.com/vincentruan/beego-blog/models"
	"github.com/vincentruan/beego-blog/models/blog"
	"github.com/vincentruan/beego-blog/models/rss"
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

func InitTasks() error {
	filePath := fmt.Sprintf("%s/conf/%s", beego.AppPath, "task.json")
	beego.BeeLogger.Debug("Init schedulers with configuration file [%s]...", filePath)

	b, err := ioutil.ReadFile(filePath)
	if nil != err {
		beego.BeeLogger.Error("Unable to open/read file ---- [%s]", filePath)
		return err
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

	c = cron.New()
	c.AddFunc("0 0/1 * * * *", persistentBlogView)
	c.Start()

	return nil
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

func persistentBlogView() {
	var blogs []models.Blog
	blog.Blogs().Limit(-1).All(&blogs)

	for _, b := range blogs {
		if bg := g.BlogCacheGet(fmt.Sprintf("%d", b.Id)); bg != nil {
			if b1, ok := bg.(models.Blog); ok && b1.Views > b.Views {
				beego.BeeLogger.Debug("Push cache [%d] into DB[ - blog id [%d], view [%d] - ]", b1.Views, b.Id, b.Views)
				b.Views = b1.Views
				blog.UpdateView(&b)
			}
		}
	}
}
