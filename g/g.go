package g

import (
	"crypto/rsa"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/square/go-jose"
	_ "github.com/vincent3i/beego-blog/cache"
	"io/ioutil"
	"os"
	"strings"
)

var MemcachedCache cache.Cache
var Cache cache.Cache
var blogCacheExpire int64
var catalogCacheExpire int64
var RunMode string
var Cfg = beego.AppConfig
var RSAPrivateKey *rsa.PrivateKey
var RSAPublicKey *rsa.PublicKey

func InitEnv() {
	// log
	logLevel := Cfg.String("log_level")
	//log.SetLevelWithDefault(logLevel, "info")

	beego.SetLogger("file", fmt.Sprintf(`{"filename":"%s"}`, Cfg.String("log_file_path")))
	beego.SetLevel(getLogLevel(logLevel))
	beego.SetLogFuncCall(true)

	initCfg()

	LoadRSACipher()
	registerDB()

	initCache()

}

func getLogLevel(logLevel string) int {
	var logLvl int

	switch strings.ToUpper(logLevel) {
	case "":
		logLvl = beego.LevelInformational
	case "EMERGENCY":
		logLvl = beego.LevelEmergency
	case "ALERT":
		logLvl = beego.LevelAlert
	case "CRITICAL":
		logLvl = beego.LevelCritical
	case "ERROR":
		logLvl = beego.LevelError
	case "WARNING":
		logLvl = beego.LevelWarning
	case "NOTICE":
		logLvl = beego.LevelNotice
	case "INFO":
		logLvl = beego.LevelInformational
	case "DEBUG":
		logLvl = beego.LevelDebug
	default:
		logLvl = beego.LevelInformational
	}
	return logLvl
}

func initCache() {
	// cache
	var err error
	Cache, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		beego.Error("cache init fail :(")
		beego.BeeLogger.Flush()
		os.Exit(1)
	}

	cacheConfig := fmt.Sprintf(`{"conn":"%s"}`, Cfg.String("memcache_addresses"))
	//注意memcache只能存储字符串，如需存储对象，需要先序列化
	//序列化手段有encoding/gob,json ,bson, msgpack(性能最佳)
	MemcachedCache, err = cache.NewCache("packmemcache", cacheConfig)
	beego.BeeLogger.Debug("Loaded cache module with config %s.", cacheConfig)
	if err != nil {
		beego.Error("memcache init fail :(")
		beego.BeeLogger.Flush()
		os.Exit(1)
	}
	//设置缓存项的到期时间
	blogCacheExpire, _ = Cfg.Int64("blog_cache_expire")
	catalogCacheExpire, _ = Cfg.Int64("catalog_cache_expire")
}

//加载RSA密钥/公钥
func LoadRSACipher() {
	privateCipherPath := Cfg.String("private_cipher_path")
	publicCipherPath := Cfg.String("public_cipher_path")
	if privateCipherPath == "" || publicCipherPath == "" {
		beego.Error("Unale to found private_cipher_path or public_cipher_path from app.cnf :(")
		beego.BeeLogger.Flush()
		os.Exit(1)
	}

	//读取私钥
	fileBytes, err := ioutil.ReadFile(privateCipherPath)
	if err != nil {
		beego.BeeLogger.Error("Unale to read file [%s] :(", privateCipherPath)
		beego.BeeLogger.Flush()
		os.Exit(1)
	}
	privateKey, err := jose.LoadPrivateKey(fileBytes)
	switch privateKey.(type) {
	case *rsa.PrivateKey:
		beego.Debug("Found RSA private cipher!")
		RSAPrivateKey = privateKey.(*rsa.PrivateKey)
	default:
		beego.BeeLogger.Error("Key from [%s] is not a valid private key for RSA :(", privateCipherPath)
		beego.BeeLogger.Flush()
		os.Exit(1)
	}

	//读取公钥
	fileBytes, err = ioutil.ReadFile(publicCipherPath)
	if err != nil {
		beego.BeeLogger.Error("Unale to read file [%s] :(", publicCipherPath)
		beego.BeeLogger.Flush()
		os.Exit(1)
	}
	publicKey, err := jose.LoadPublicKey(fileBytes)
	switch publicKey.(type) {
	case *rsa.PublicKey:
		beego.Debug("Found RSA public cipher!")
		RSAPublicKey = publicKey.(*rsa.PublicKey)
	default:
		beego.BeeLogger.Error("Key from [%s] is not a valid public key for RSA :(", publicCipherPath)
		beego.BeeLogger.Flush()
		os.Exit(1)
	}

	beego.Debug("Load RSA private and public cipher from app.cnf completed!")
}

func Encrypt(plaintext string) string {
	// Instantiate an encrypter using RSA-OAEP with AES128-GCM. An error would
	// indicate that the selected algorithm(s) are not currently supported.
	encrypter, err := jose.NewEncrypter(jose.RSA_OAEP, jose.A128GCM, RSAPublicKey)
	if err != nil {
		beego.BeeLogger.Error("Unable to instantiate an encrypter using RSA-OAEP with AES128-GCM. %s", err)
		return plaintext
	}
	// Encrypt a sample plaintext. Calling the encrypter returns an encrypted
	// JWE object, which can then be serialized for output afterwards. An error
	// would indicate a problem in an underlying cryptographic primitive.
	object, err := encrypter.Encrypt([]byte(plaintext))
	if err != nil {
		beego.BeeLogger.Error("Unable to calling the encrypter returns an encrypted JWE object. %s", err)
		return plaintext
	}

	serialized, err := object.CompactSerialize()
	if err != nil {
		beego.BeeLogger.Error("Unable to compactSerialize serializes an object using the compact serialization format. %s", err)
		return plaintext
	}

	return serialized
}

func Decrypt(serialized string) string {
	// Parse the serialized, encrypted JWE object. An error would indicate that
	// the given input did not represent a valid message.
	object, err := jose.ParseEncrypted(serialized)
	if err != nil {
		beego.BeeLogger.Error("Unable to parse the serialized, encrypted JWE object. %s", err)
		return serialized
	}

	// Now we can decrypt and get back our original plaintext. An error here
	// would indicate the the message failed to decrypt, e.g. because the auth
	// tag was broken or the message was tampered with.
	decrypted, err := object.Decrypt(RSAPrivateKey)
	if err != nil {
		beego.BeeLogger.Error("Unable to decrypt and get back the original plaintext, encrypted JWE object. %s", err)
		return serialized
	}

	return string(decrypted)
}

func registerDB() {
	// database
	dbUser := Cfg.String("db_user")
	dbPass := Cfg.String("db_pass")
	dbHost := Cfg.String("db_host")
	dbPort := Cfg.String("db_port")
	dbName := Cfg.String("db_name")
	maxIdleConn, _ := Cfg.Int("db_max_idle_conn")
	maxOpenConn, _ := Cfg.Int("db_max_open_conn")
	//%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Asia%%2FChongqing&allowOldPasswords=true
	dbLink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Asia%%2FChongqing&allowOldPasswords=true", dbUser, Decrypt(dbPass), dbHost, dbPort, dbName)
	beego.BeeLogger.Debug("dbLink ---->>> %s .", dbLink)

	orm.RegisterDriver("mysql", orm.DR_MySQL)
	err := orm.RegisterDataBase("default", "mysql", dbLink, maxIdleConn, maxOpenConn)
	if nil != err {
		beego.Error(err.Error())
	}

	RunMode = Cfg.String("runmode")
	if RunMode == "dev" {
		orm.Debug = true
	}
}
