package main

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/vincent3i/beego-blog/engine"
	"github.com/vincent3i/beego-blog/g"
	"github.com/vincent3i/beego-blog/nsq/consumer"
	"github.com/vincent3i/beego-blog/nsq/producer"
	_ "github.com/vincent3i/beego-blog/routers"
	"github.com/vincent3i/beego-blog/task"
	"github.com/vincent3i/beego-blog/utilities/mongo"
)

func main() {
	g.InitEnv()
	beego.AddAPPStartHook(mongo.Startup)
	beego.AddAPPStartHook(task.InitTasks)
	beego.AddAPPStartHook(consumer.InitNSQCunsumer)
	beego.AddAPPStartHook(producer.InitNSQProducer)
	beego.AddAPPStartHook(engine.InitSearcher)
	beego.AddAPPStartHook(engine.InitElasticSearch)
	beego.Run()
}
