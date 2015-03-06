package g

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	//log "github.com/ulricqin/goutils/logtool"
	"github.com/astaxie/beego/logs"
	"os"
)

var Cache cache.Cache
var blogCacheExpire int64
var catalogCacheExpire int64
var RunMode string
var Cfg = beego.AppConfig
var Log *logs.BeeLogger

func InitEnv() {
	var err error

	// log
	logLevel := Cfg.String("log_level")
	//log.SetLevelWithDefault(logLevel, "info")

	Log = logs.NewLogger(1024)
	Log.SetLogger("console", "")
	Log.SetLevel(getLogLevel(logLevel))

	// cache
	Cache, err = cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		//log.Fetal("cache init fail :(")
		Log.Error("cache init fail :(")
		os.Exit(1)
	}
	blogCacheExpire, _ = Cfg.Int64("blog_cache_expire")
	catalogCacheExpire, _ = Cfg.Int64("catalog_cache_expire")

	// database
	dbUser := Cfg.String("db_user")
	dbPass := Cfg.String("db_pass")
	dbHost := Cfg.String("db_host")
	dbPort := Cfg.String("db_port")
	dbName := Cfg.String("db_name")
	maxIdleConn, _ := Cfg.Int("db_max_idle_conn")
	maxOpenConn, _ := Cfg.Int("db_max_open_conn")
	dbLink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Asia%%2FChongqing", dbUser, dbPass, dbHost, dbPort, dbName)
	Log.Debug("dbLink ---->>> %s .", dbLink)

	orm.RegisterDriver("mysql", orm.DR_MySQL)
	orm.RegisterDataBase("default", "mysql", dbLink, maxIdleConn, maxOpenConn)

	RunMode = Cfg.String("runmode")
	if RunMode == "dev" {
		orm.Debug = true
	}

	initCfg()
}

func getLogLevel(logLevel string) int {
	var logLvl int

	switch Log.Debug("Log level is %s", logLevel); strings.ToUpper(logLevel) {
	case "":
		logLvl = logs.LevelInformational
	case "EMERGENCY":
		logLvl = logs.LevelEmergency
	case "ALERT":
		logLvl = logs.LevelAlert
	case "CRITICAL":
		logLvl = logs.LevelCritical
	case "ERROR":
		logLvl = logs.LevelError
	case "WARNING":
		logLvl = logs.LevelWarning
	case "NOTICE":
		logLvl = logs.LevelNotice
	case "INFO":
		logLvl = logs.LevelInformational
	case "DEBUG":
		logLvl = logs.LevelDebug
	default:
		logLvl = logs.LevelInformational
	}
	return logLvl
}
