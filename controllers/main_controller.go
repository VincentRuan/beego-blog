package controllers

import (
	"github.com/astaxie/beego"
	"github.com/vincent3i/beego-blog/g"
	"github.com/vincent3i/beego-blog/models"
	"github.com/vincent3i/beego-blog/models/blog"
	"github.com/vincent3i/beego-blog/models/catalog"
)

type MainController struct {
	BaseController
}

func (this *MainController) Get() {
	catalogs := catalog.All()
	this.Data["Catalogs"] = catalogs
	this.Data["PageTitle"] = "首页"
	this.Data["CatalogSize"] = len(catalogs)
	this.Layout = "layout/default.html"
	this.TplNames = "index.html"
}

func (this *MainController) Read() {
	ident := this.GetString(":ident")
	b := blog.OneByIdent(ident)
	if b == nil {
		this.Ctx.WriteString("no such article")
		return
	}

	beego.BeeLogger.Debug("User read blog [%d] [%s]", b.Id, b.Title)

	if vc := g.BlogViewCacheGet(b.Id); vc > 0 {
		b.Views = vc + 1
	} else {
		b.Views = b.Views + 1
	}
	g.BlogViewCachePut(b.Id, b.Views)

	//b.Views = b.Views + 1
	//blog.Update(b, "")

	this.Data["Blog"] = b
	this.Data["Content"] = g.RenderMarkdown(blog.ReadBlogContent(b).Content)
	this.Data["PageTitle"] = b.Title
	this.Data["Catalog"] = catalog.OneById(b.CatalogId)
	this.Layout = "layout/default.html"
	this.TplNames = "article/read.html"
}

func (this *MainController) ListByCatalog() {
	cata := this.Ctx.Input.Param(":ident")
	if cata == "" {
		this.Ctx.WriteString("catalog ident is blank")
		return
	}

	limit := this.GetIntWithDefault("limit", 10)

	c := catalog.OneByIdent(cata)
	if c == nil {
		this.Ctx.WriteString("catalog:" + cata + " not found")
		return
	}

	if c.IsAuth && !this.IsAdmin {
		beego.Debug("User can NOT access this page as verification is need!")
		this.SetPaginator(10, 0)
		this.Data["Catalog"] = c
		this.Data["Blogs"] = make([]*models.Blog, 0, 0)
		this.Data["PageTitle"] = c.Name
	} else {
		ids := blog.Ids(c.Id)
		pager := this.SetPaginator(limit, int64(len(ids)))
		blogs := blog.ByCatalog(c.Id, pager.Offset(), limit)

		this.Data["Catalog"] = c
		this.Data["Blogs"] = blogs
		this.Data["PageTitle"] = c.Name
	}

	this.Layout = "layout/default.html"
	this.TplNames = "article/by_catalog.html"
}
