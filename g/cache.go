package g

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/vincent3i/beego-blog/cache"
)

const (
	blogPrefix     = "blog_"
	catalogPrefix  = "catalog_"
	blogViewPrefix = "blog_view_"
)

func BlogCachePut(key string, val interface{}) error {
	return Cache.Put(blogPrefix+key, val, blogCacheExpire)
}

func CatalogCachePut(key string, val interface{}) error {
	return Cache.Put(catalogPrefix+key, val, catalogCacheExpire)
}

func BlogCacheGet(key string) interface{} {
	return Cache.Get(blogPrefix + key)
}

func CatalogCacheGet(key string) interface{} {
	return Cache.Get(catalogPrefix + key)
}

func CatalogCacheDel(key string) error {
	return Cache.Delete(catalogPrefix + key)
}

func BlogCacheDel(key string) error {
	return Cache.Delete(blogPrefix + key)
}

func BlogViewCacheDel(key int64) error {
	return MemcachedCache.Delete(fmt.Sprintf("%s%d", blogViewPrefix, key))
}

func BlogViewCacheGet(key int64) int64 {
	var view int64
	if err := cache.Unmarshal(MemcachedCache.Get(fmt.Sprintf("%s%d", blogViewPrefix, key)), &view); err != nil {
		beego.Error(err)
		return -1
	}
	return view
}

func BlogViewCachePut(key int64, val interface{}) error {
	return MemcachedCache.Put(fmt.Sprintf("%s%d", blogViewPrefix, key), val, 0)
}
