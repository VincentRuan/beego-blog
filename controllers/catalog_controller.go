package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/ulricqin/goutils/filetool"
	"github.com/vincent3i/beego-blog/g"
	"github.com/vincent3i/beego-blog/models"
	"github.com/vincent3i/beego-blog/models/catalog"
	"strings"
	"time"
)

type CatalogController struct {
	AdminController
}

func (this *CatalogController) Add() {
	this.Data["IsAddCatalog"] = true
	this.Layout = "layout/admin.html"
	this.TplNames = "catalog/add.html"
}

func (this *CatalogController) Edit() {
	id, err := this.GetInt64("id")
	if err != nil {
		this.Ctx.WriteString("param id should be digit")
		return
	}

	c := catalog.OneById(id)
	if c == nil {
		this.Ctx.WriteString(fmt.Sprintf("no such catalog_id:%d", id))
		return
	}

	this.Data["Catalog"] = c
	this.Layout = "layout/admin.html"
	this.TplNames = "catalog/edit.html"
}

func (this *CatalogController) Del() {
	id, err := this.GetInt64("id")
	if err != nil {
		this.Ctx.WriteString("param id should be digit")
		return
	}

	c := catalog.OneById(id)
	if c == nil {
		this.Ctx.WriteString(fmt.Sprintf("no such catalog_id:%d", id))
		return
	}

	err = catalog.Del(c)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}

	this.Ctx.WriteString("del success")
	return
}

func (this *CatalogController) extractCatalog(imgMust bool) (*models.Catalog, error) {
	catalog := &models.Catalog{}
	catalog.Name = this.GetString("name")
	catalog.Ident = this.GetString("ident")
	catalog.Resume = this.GetString("resume")
	catalog.DisplayOrder = this.GetIntWithDefault("display_order", 0)

	if catalog.Name == "" {
		return nil, fmt.Errorf("name is blank")
	}

	if catalog.Ident == "" {
		return nil, fmt.Errorf("ident is blank")
	}

	_, header, err := this.GetFile("img")
	if err != nil && imgMust {
		return nil, err
	}

	if err == nil {
		ext := filetool.Ext(header.Filename)

		//beego.Debug("Current OS is %s", runtime.GOOS)
		imgPath := fmt.Sprintf("%s/%s_%d%s", g.LocalCatalogUploadPath, catalog.Ident, time.Now().Unix(), ext)
		filetool.InsureDir(g.LocalCatalogUploadPath)
		err = this.SaveToFile("img", imgPath)
		beego.BeeLogger.Debug("Saved file as %s", imgPath)
		if err != nil && imgMust {
			return nil, err
		}

		if err == nil {
			catalog.ImgUrl = imgPath

			//存储在服务器的文件名
			qiniuFileName := g.QiniuCatalogUploadPath + imgPath[strings.LastIndex(imgPath, "/"):]

			if g.UseQiniu {
				if addr, er := g.UploadFile(imgPath, qiniuFileName); er != nil {
					if imgMust {
						beego.BeeLogger.Error("Upload file [%s] to Qiniu cloud store error %s", imgPath, er.Error())
						return nil, er
					}
				} else {
					catalog.ImgUrl = addr
					filetool.Unlink(imgPath)
					beego.BeeLogger.Debug("Uploaded file [%s] success, remove file from [%s].", qiniuFileName, imgPath)
				}
			}
		}
	}

	return catalog, nil
}

func (this *CatalogController) DoEdit() {
	cid, err := this.GetInt64("catalog_id")
	if err != nil {
		this.Ctx.WriteString("catalog_id is illegal")
		return
	}

	old := catalog.OneById(cid)
	if old == nil {
		this.Ctx.WriteString(fmt.Sprintf("no such catalog_id: %d", cid))
		return
	}

	var o *models.Catalog
	o, err = this.extractCatalog(false)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}

	old.Ident = o.Ident
	old.Name = o.Name
	old.Resume = o.Resume
	old.DisplayOrder = o.DisplayOrder
	if o.ImgUrl != "" {
		old.ImgUrl = o.ImgUrl
	}

	if err = catalog.Update(old); err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}

	this.Redirect("/", 302)

}

func (this *CatalogController) DoAdd() {
	o, err := this.extractCatalog(true)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}

	_, err = catalog.Save(o)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}

	this.Redirect("/", 302)
}
