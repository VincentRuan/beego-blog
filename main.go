package main

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/vincentruan/beego-blog/engine"
	"github.com/vincentruan/beego-blog/g"
	"github.com/vincentruan/beego-blog/nsq/consumer"
	"github.com/vincentruan/beego-blog/nsq/producer"
	_ "github.com/vincentruan/beego-blog/routers"
	"github.com/vincentruan/beego-blog/task"
	"github.com/vincentruan/beego-blog/utilities/mongo"
)

func main() {
	beego.AddAPPStartHook(g.InitEnv)
	beego.AddAPPStartHook(mongo.Startup)
	beego.AddAPPStartHook(task.InitTasks)
	beego.AddAPPStartHook(consumer.InitNSQCunsumer)
	beego.AddAPPStartHook(producer.InitNSQProducer)
	beego.AddAPPStartHook(engine.InitElasticSearch)
	beego.AddAPPStartHook(engine.InitSearcher)
	beego.Run()
}
