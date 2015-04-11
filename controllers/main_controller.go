package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/vincent3i/beego-blog/engine"
	"github.com/vincent3i/beego-blog/g"
	"github.com/vincent3i/beego-blog/models"
	"github.com/vincent3i/beego-blog/models/blog"
	"github.com/vincent3i/beego-blog/models/catalog"
	"time"
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

	if !this.IsAdmin && blog.IsAuth(b.Id) {
		this.Ctx.WriteString("You do not have permission!")
	}

	beego.BeeLogger.Debug("User read blog [%d] [%s]", b.Id, b.Title)

	//in here blog cache must be set
	if bg := g.BlogCacheGet(fmt.Sprintf("%d", b.Id)); bg != nil {
		b1 := bg.(models.Blog)
		b1.Views = b1.Views + 1

		g.BlogCachePut(fmt.Sprintf("%d", b.Id), b1)
	}

	//b.Views = b.Views + 1
	//blog.Update(b, "")

	this.Data["Blog"] = b

	beginT := time.Now()
	if content := g.BlogContentCacheGet(b.Id); content != "" {
		beego.Debug("Get blog content by cache in ", time.Since(beginT))
		this.Data["Content"] = content
	} else {
		beginT = time.Now()
		content = g.RenderMarkdown(blog.ReadBlogContent(b).Content)
		beego.Debug("Get blog content by render in ", time.Since(beginT))
		g.BlogContentCachePut(b.Id, content)
		this.Data["Content"] = content
	}
	//this.Data["Content"] = g.RenderMarkdown(blog.ReadBlogContent(b).Content)
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

func (this *MainController) Query() {
	this.Data["PageTitle"] = "搜索博客"
	this.Data["BlogTitle"] = "Vincent"
	this.Layout = "layout/default.html"
	this.TplNames = "so/so.html"
}

func (this *MainController) DoQuery() {
	query := this.GetString("q")
	beego.BeeLogger.Debug("Search [%s]", query)

	var result []models.Blog
	if query != "" {
		result = engine.SearchResult(query, this.IsAdmin)
	} else {
		result = make([]models.Blog, 0)
	}
	beego.BeeLogger.Debug("查询结果 %d 条", len(result))
	this.Data["json"] = result

	this.ServeJson()
}

func (this *MainController) ElasticQuery() {
	this.Data["PageTitle"] = "搜索博客"
	this.Data["BlogTitle"] = "Vincent"
	this.Layout = "layout/default.html"
	this.TplNames = "so/query.html"
}

func (this *MainController) DoElasticQuery() {
	query := this.GetString("q")
	beego.BeeLogger.Debug("Search [%s]", query)

	var result []models.Blog
	if query != "" {
		result = engine.ElasticSearch(query, this.IsAdmin)
	} else {
		result = make([]models.Blog, 0)
	}
	beego.BeeLogger.Debug("查询结果 %d 条", len(result))
	this.Data["json"] = result

	this.ServeJson()
}
