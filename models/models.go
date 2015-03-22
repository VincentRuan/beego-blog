package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/qiniu/api/rs"
	"time"
)

type Catalog struct {
	Id           int64
	Ident        string `orm:"unique"`
	Name         string
	Resume       string
	DisplayOrder int
	ImgUrl       string
}

type Blog struct {
	Id                    int64
	Ident                 string `orm:"unique"`
	Title                 string
	Keywords              string       `orm:"null"`
	CatalogId             int64        `orm:"index"`
	Content               *BlogContent `orm:"-"`
	BlogContentId         int64        `orm:"unique"`
	BlogContentLastUpdate int64
	Type                  int8 /*0:original, 1:translate, 2:reprint*/
	Status                int8 /*0:draft, 1:release*/
	Views                 int64
	Created               time.Time `orm:"auto_now_add;type(datetime)"`
}

type RssFeeder struct {
	Id            int64     `json:"id" form:"id"`
	RSSDesc       string    `orm:"column(rss_desc)" json:"rss_desc" valid:"Required" form:"rss_desc"`
	RSSUrl        string    `orm:"column(rss_url)" json:"rss_url" valid:"Required" form:"rss_url"`
	CreateTime    time.Time `orm:"column(create_time)"json:"create_time" form:"-"`
	UpdateTime    time.Time `orm:"column(update_time)"json:"update_time" form:"-"`
	SubscribeTime time.Time `orm:"column(subscribe_time)"json:"subscribe_time" form:"-"`
}

type BlogContent struct {
	Id      int64
	Content string `orm:"type(text)"`
}

func (*Catalog) TableEngine() string {
	return engine()
}

func (*Blog) TableEngine() string {
	return engine()
}

func (*BlogContent) TableEngine() string {
	return engine()
}

func engine() string {
	return "INNODB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci"
}

func init() {
	orm.RegisterModelWithPrefix("bb_", new(Catalog), new(Blog), new(BlogContent))
	orm.RegisterModel(new(RssFeeder))
}

func (this *RssFeeder) TableName() string {
	return "rss_feeder_t"
}

// func main() {
// 	orm.RegisterDataBase("default", "mysql", "root:@/beego_blog?charset=utf8&loc=Asia%2FShanghai", 30, 200)
// 	orm.RunCommand()
// }

func QiniuDownloadUrl(domain, key string) string {
	//g.Log.Debug("Download domain is --->>> %s", domain)
	//g.Log.Debug("Download key is --->>> %s", key)
	baseUrl := rs.MakeBaseUrl(domain, key)
	policy := rs.GetPolicy{}
	return policy.MakeRequest(baseUrl, nil)
}
