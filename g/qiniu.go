package g

import (
	"github.com/astaxie/beego"
	"github.com/qiniu/api/io"
	"github.com/qiniu/api/rs"
)

func UploadFile(localFile string, destName string) (addr string, err error) {
	//policy := new(rs.PutPolicy)
	//policy.Scope = QiniuScope
	policy := rs.PutPolicy{Scope: QiniuScope}
	uptoken := policy.Token(nil)

	var ret io.PutRet
	var extra = new(io.PutExtra)

	err = io.PutFile(nil, &ret, uptoken, destName, localFile, extra)
	if err != nil {
		return
	}

	//addr = "http://" + QiniuScope + ".qiniudn.com" + destName

	//私有域名访问
	//http://developer.qiniu.com/docs/v6/api/reference/security/download-token.html
	//http://developer.qiniu.com/docs/v6/sdk/go-sdk.html#io-get-private
	addr = "http://" + QiniuHttpDomain + "/@" + destName
	beego.BeeLogger.Debug("Upload file address is --->>> %s", addr)
	return
}
