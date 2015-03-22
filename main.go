package main

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/vincent3i/beego-blog/g"
	_ "github.com/vincent3i/beego-blog/routers"
	"github.com/vincent3i/beego-blog/task"
	"github.com/vincent3i/beego-blog/utilities/mongo"
)

func main() {
	g.InitEnv()
	beego.BeeLogger.Debug("Is session on [%t], session provide --->>> [%s]", beego.SessionOn, beego.SessionProvider)
	mongo.Startup()
	task.InitTasks()
	beego.Run()
}
