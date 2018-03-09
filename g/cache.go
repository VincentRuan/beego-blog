package g

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/vincentruan/beego-blog/cache"
)

const (
	blogPrefix        = "blog_"
	catalogPrefix     = "catalog_"
	blogViewPrefix    = "blog_view_"
	blogContentPrefix = "blog_content_"
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

func BlogContentCacheGet(key int64) string {
	var content string
	if err := cache.Unmarshal(MemcachedCache.Get(fmt.Sprintf("%s%d", blogContentPrefix, key)), &content); err != nil {
		beego.Error(err)
		return ""
	}
	return content
}

func BlogContentCachePut(key int64, val interface{}) error {
	return MemcachedCache.Put(fmt.Sprintf("%s%d", blogContentPrefix, key), val, blogCacheExpire)
}

func BlogContentCacheDel(key int64) error {
	return MemcachedCache.Delete(fmt.Sprintf("%s%d", blogContentPrefix, key))
}
