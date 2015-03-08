package main

import (
	"github.com/astaxie/beego"
	"github.com/vincent3i/beego-blog/g"
	_ "github.com/vincent3i/beego-blog/routers"
)

func main() {
	g.InitEnv()
	g.Log.Debug("Is session on [%t], session provide --->>> [%s]", beego.SessionOn, beego.SessionProvider)
	beego.Run()
}
