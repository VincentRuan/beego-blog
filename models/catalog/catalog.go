package catalog

import (
	"fmt"
	"github.com/vincent3i/beego-blog/g"
	"net/url"
	//"strings"
	//"strconv"
	//.表示可以不用带包名访问里面的变量方法
	"github.com/astaxie/beego/orm"
	"github.com/qiniu/api/rs"
	. "github.com/vincent3i/beego-blog/models"
)

func OneById(id int64) *Catalog {
	if id == 0 {
		return nil
	}

	//or used strconv.Itoa to convert int o string
	//strconv.Itoa(i)
	key := fmt.Sprintf("id_catalog_%d", id)
	val := g.CatalogCacheGet(key)
	if val == nil {
		if cp := OneByIdInDB(id); cp != nil {
			g.CatalogCachePut(key, *cp)
			return cp
		} else {
			return nil
		}
	}
	ret := val.(Catalog)
	return &ret
}

func OneByIdInDB(id int64) *Catalog {
	if id == 0 {
		return nil
	}

	c := Catalog{Id: id}
	err := orm.NewOrm().Read(&c, "Id")
	if err != nil {
		return nil
	}
	return &c
}

func IdByIdent(ident string) int64 {
	if ident == "" {
		return 0
	}

	val := g.CatalogCacheGet(ident)
	if val == nil {
		if cp := OneByIdentInDB(ident); cp != nil {
			g.CatalogCachePut(ident, cp.Id)
			return cp.Id
		} else {
			return 0
		}
	}

	return val.(int64)
}

func IdentExists(ident string) bool {
	id := IdByIdent(ident)
	return id > 0
}

func OneByIdent(ident string) *Catalog {
	id := IdByIdent(ident)
	return OneById(id)
}

func OneByIdentInDB(ident string) *Catalog {
	if ident == "" {
		return nil
	}

	c := Catalog{Ident: ident}
	err := orm.NewOrm().Read(&c, "Ident")
	if err != nil {
		return nil
	}

	return &c
}

func AllIdsInDB() []int64 {
	var catalogs []Catalog
	Catalogs().OrderBy("-DisplayOrder").All(&catalogs, "Id")
	size := len(catalogs)
	if size == 0 {
		return []int64{}
	}

	ret := make([]int64, size)
	for i := 0; i < size; i++ {
		ret[i] = catalogs[i].Id
	}

	return ret
}

func AllCatalogsInDB() []*Catalog {
	var catalogs []*Catalog
	Catalogs().OrderBy("-DisplayOrder").All(&catalogs)

	return catalogs
}

func AllIds() []int64 {
	val := g.CatalogCacheGet("ids")
	if val == nil {
		if ids := AllIdsInDB(); len(ids) != 0 {
			g.CatalogCachePut("ids", ids)
			return ids
		} else {
			return []int64{}
		}
	}

	return val.([]int64)
}

func All() []*Catalog {
	catalogsInCache := g.CatalogCacheGet("catalogs")
	if nil == catalogsInCache {
		if catalogs := AllCatalogsInDB(); len(catalogs) != 0 {
			catalogs = tokenValidCatalogs(catalogs)
			g.CatalogCachePut("catalogs", catalogs)
			return catalogs
		}

		return []*Catalog{}
	}

	return tokenValidCatalogs(catalogsInCache.([]*Catalog))
}

//七牛对于私有空间需要token验证
//使用前请先申请自己的AK/CK、scope
func tokenValidCatalogs(catalogs []*Catalog) []*Catalog {
	//私有域名访问
	//http://developer.qiniu.com/docs/v6/api/reference/security/download-token.html
	//http://developer.qiniu.com/docs/v6/sdk/go-sdk.html#io-get-private
	if g.IsQiniuPublicAccess {
		return catalogs
	}

	var uri *url.URL
	var err error
	for _, catalog := range catalogs {
		uri, err = url.ParseRequestURI(catalog.ImgUrl)
		if nil == err {
			//g.Log.Debug(uri.Path)
			if len(uri.Path) > 1 {
				catalog.ImgUrl = QiniuDownloadUrl(uri.Host, uri.Path[1:])
			} else {
				catalog.ImgUrl = "/static/images/golang.jpg"
			}

		} else {
			catalog.ImgUrl = "/static/images/golang.jpg"
		}
	}

	return catalogs
}

func Save(this *Catalog) (int64, error) {
	if IdentExists(this.Ident) {
		return 0, fmt.Errorf("catalog english identity exists")
	}
	num, err := orm.NewOrm().Insert(this)
	if err == nil {
		g.CatalogCacheDel("ids")
		g.CatalogCacheDel("catalogs")
	}

	return num, err
}

func Del(c *Catalog) error {
	num, err := orm.NewOrm().Delete(c)
	if err != nil {
		return err
	}

	if num > 0 {
		g.CatalogCacheDel("ids")
		g.CatalogCacheDel(fmt.Sprintf("id_catalog_%d", c.Id))
		g.CatalogCacheDel("catalogs")
	}
	return nil
}

func Update(this *Catalog) error {
	if this.Id == 0 {
		return fmt.Errorf("primary key id not set")
	}
	_, err := orm.NewOrm().Update(this)
	if err == nil {
		g.CatalogCacheDel(fmt.Sprintf("id_catalog_%d", this.Id))
		g.CatalogCacheDel("catalogs")
	}
	return err
}

func Catalogs() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Catalog))
}

func downloadUrl(domain, key string) string {
	baseUrl := rs.MakeBaseUrl(domain, key)
	policy := rs.GetPolicy{}
	return policy.MakeRequest(baseUrl, nil)
}
