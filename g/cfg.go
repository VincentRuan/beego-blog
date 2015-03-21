package g

import (
	"github.com/qiniu/api/conf"
	"regexp"
	"strings"
)

var (
	BlogTitle              string
	BlogResume             string
	BlogLogo               string
	QiniuAccessKey         string
	QiniuSecretKey         string
	QiniuScope             string
	UseQiniu               bool
	QiniuHttpDomain        string
	IsQiniuPublicAccess    bool = true
	QiniuCatalogUploadPath string
	LocalCatalogUploadPath string
)

func initCfg() {
	BlogTitle = Cfg.String("blog_title")
	BlogResume = Cfg.String("blog_resume")
	BlogLogo = Cfg.String("blog_logo")
	QiniuAccessKey = Cfg.String("qiniu_access_key")
	QiniuSecretKey = Cfg.String("qiniu_secret_key")
	QiniuScope = Cfg.String("qiniu_scope")
	UseQiniu = QiniuAccessKey != "" && QiniuSecretKey != "" && QiniuScope != ""
	QiniuHttpDomain = strings.Trim(Cfg.String("qiniu_http_domain"), " ")
	conf.ACCESS_KEY = QiniuAccessKey
	conf.SECRET_KEY = QiniuSecretKey

	switch strings.ToUpper(Cfg.String("qiniu_access_control")) {
	case "PRIVATE":
		IsQiniuPublicAccess = false
	case "PUBLIC":
		IsQiniuPublicAccess = true
	default:
		IsQiniuPublicAccess = true
	}

	QiniuCatalogUploadPath = Cfg.String("qiniu_catalog_upload_path")
	//上传路径将\强制转换成/
	re, _ := regexp.Compile(`\\+`)
	LocalCatalogUploadPath = re.ReplaceAllString(Cfg.String("local_catalog_upload_path"), "/")

}
