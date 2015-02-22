package main

import (
	"github.com/vincent3i/beego-blog/g"
	_ "github.com/vincent3i/beego-blog/routers"
	"github.com/astaxie/beego"
)

func main() {
	g.InitEnv()
	beego.Run()
}
