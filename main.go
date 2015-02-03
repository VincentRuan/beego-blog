package main

import (
	"github.com/VincentRuan/beego-blog/g"
	_ "github.com/VincentRuan/beego-blog/routers"
	"github.com/astaxie/beego"
)

func main() {
	g.InitEnv()
	beego.Run()
}
