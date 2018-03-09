package routers

import (
	"github.com/astaxie/beego"
	"github.com/vincentruan/beego-blog/controllers"
)

func init() {

	beego.AutoRouter(&controllers.ApiController{})

	beego.Router("/", &controllers.MainController{})
	beego.Router("/article/:ident", &controllers.MainController{}, "get:Read")
	beego.Router("/catalog/:ident", &controllers.MainController{}, "get:ListByCatalog")

	beego.Router("/login", &controllers.LoginController{}, "get:Login;post:DoLogin")
	beego.Router("/logout", &controllers.LoginController{}, "get:Logout")

	beego.Router("/me", &controllers.MeController{}, "get:Default")
	beego.Router("/me/catalog/add", &controllers.CatalogController{}, "get:Add;post:DoAdd")
	beego.Router("/me/catalog/edit", &controllers.CatalogController{}, "get:Edit;post:DoEdit")
	beego.Router("/me/catalog/del", &controllers.CatalogController{}, "get:Del")
	beego.Router("/me/article/add", &controllers.ArticleController{}, "get:Add;post:DoAdd")
	beego.Router("/me/article/edit", &controllers.ArticleController{}, "get:Edit;post:DoEdit")
	beego.Router("/me/article/del", &controllers.ArticleController{}, "get:Del")
	beego.Router("/me/article/draft", &controllers.ArticleController{}, "get:Draft")

	beego.Router("/so", &controllers.MainController{}, "get:Query;post:DoQuery")
	beego.Router("/q", &controllers.MainController{}, "get:ElasticQuery;post:DoElasticQuery")
	beego.Router("/shutdown", &controllers.MeController{}, "get:ShutDown")

	beego.Router("/me/rss", &controllers.RSSController{}, "get:LoadPage")
	beego.Router("/me/rss/read", &controllers.RSSController{}, "get,post:Read")
	beego.Router("/me/rss/edit", &controllers.RSSController{}, "post:DoEdit")
	beego.Router("/me/rss/del", &controllers.RSSController{}, "post:DoDel")
	beego.Router("/me/rss/detail/:id([0-9]+)", &controllers.RSSController{}, "get:RSSData")
}
