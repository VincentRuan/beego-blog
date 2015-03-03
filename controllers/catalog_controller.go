package controllers

import (
	"fmt"
	"runtime"
	"github.com/vincent3i/beego-blog/g"
	"github.com/vincent3i/beego-blog/models"
	"github.com/vincent3i/beego-blog/models/catalog"
	"github.com/ulricqin/goutils/filetool"
	"strings"
	"time"
)

const (
	CATALOG_IMG_DIR = "F:/dev/dm/beego-blog/static/uploads/catalogs"
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
	o := &models.Catalog{}
	o.Name = this.GetString("name")
	o.Ident = this.GetString("ident")
	o.Resume = this.GetString("resume")
	o.DisplayOrder = this.GetIntWithDefault("display_order", 0)

	if o.Name == "" {
		return nil, fmt.Errorf("name is blank")
	}

	if o.Ident == "" {
		return nil, fmt.Errorf("ident is blank")
	}

	_, header, err := this.GetFile("img")
	if err != nil && imgMust {
		return nil, err
	}

	if err == nil {
		ext := filetool.Ext(header.Filename)
		g.Log.Debug("Saving file into directory %s", CATALOG_IMG_DIR)
		g.Log.Debug("Current OS is %s", runtime.GOOS)
		
		imgPath := fmt.Sprintf("%s/%s_%d%s", CATALOG_IMG_DIR, o.Ident, time.Now().Unix(), ext)
		
		filetool.InsureDir(CATALOG_IMG_DIR)
		err = this.SaveToFile("img", imgPath)
		if err != nil && imgMust {
			return nil, err
		}

		if err == nil {
			var dest_qiniu_path string
			if strings.EqualFold("windows", runtime.GOOS) {
				dest_qiniu_path = string(imgPath[len("F:/dev/dm/beego-blog"):])
			} else {
				dest_qiniu_path = imgPath
			}
			o.ImgUrl = dest_qiniu_path
			
			if g.UseQiniu {
				if addr, er := g.UploadFile(imgPath, dest_qiniu_path); er != nil {
					if imgMust {
						g.Log.Error("Upload file to Qiniu error %v", er)
						return nil, er
					}
				} else {
					o.ImgUrl = addr
					filetool.Unlink(imgPath)
				}
			}
		}
	}

	return o, nil
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
