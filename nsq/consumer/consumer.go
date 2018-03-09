package consumer

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/bitly/go-nsq"
	"github.com/vincentruan/beego-blog/engine"
	"github.com/vincentruan/beego-blog/g"
	"github.com/vincentruan/beego-blog/models"
	"github.com/vincentruan/beego-blog/models/blog"
	"gopkg.in/vmihailenco/msgpack.v2"
	"strconv"
)

type Handler func(*nsq.Message)

type queue struct {
	callback Handler
	*nsq.Consumer
}

func (q *queue) HandleMessage(message *nsq.Message) error {
	q.callback(message)
	return nil
}

func InitNSQCunsumer() error {
	if g.NSQAddr == "" {
		return errors.New("Unable to read NSQ address from config file!")
	}

	err := initBlogConsumer()
	if err != nil {
		return err
	}

	return nil
}

func initBlogConsumer() error {
	c, err := nsq.NewConsumer("elastic-blog", "blog-chan", nsq.NewConfig())
	if err != nil {
		return err
	}
	c.SetLogger(beego.BeeLogger, nsq.LogLevelInfo)

	//add handler
	q := &queue{HandleElasticBlogs, c}
	c.AddHandler(q)

	err = c.ConnectToNSQD(g.NSQAddr)
	if err != nil {
		return err
	}

	return nil
}

func HandleElasticBlogs(msg *nsq.Message) {
	bb := models.Blog{}
	err := msgpack.Unmarshal(msg.Body, &bb)
	if err != nil {
		beego.Error(err)
		return
	}

	elasticBlog := engine.ElasticBlog{strconv.FormatInt(bb.Id, 10), bb.Title + blog.ReadBlogContent(&bb).Content}
	put, err := engine.ElasticClient.Index().
		Index(engine.Blog_Index_Name).
		Type("blog").
		Id(elasticBlog.Id).
		BodyJson(elasticBlog).
		Do()
	if err != nil {
		beego.Error(err)
		return
	}
	beego.BeeLogger.Debug("Indexed blog %s to index %s, type %s", put.Id, put.Index, put.Type)
	//eat message
	msg.Finish()
}
