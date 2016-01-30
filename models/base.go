package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"go-web-game/conf"
	"log"
	"os"
	"strings"
	"time"
)

var Xorm *xorm.Engine

func init() {
	//获取服务器数据库相关信息
	dbhost, _ := conf.Cfg.GetValue("db", "db_host")
	dbport, _ := conf.Cfg.GetValue("db", "db_port")
	dbuser, _ := conf.Cfg.GetValue("db", "db_user")
	dbpass, _ := conf.Cfg.GetValue("db", "db_pass")
	dbname, _ := conf.Cfg.GetValue("db", "db_name")
	dbIdle, _ := conf.Cfg.Int("db", "db_max_idle_conn")
	dbConn, _ := conf.Cfg.Int("db", "db_max_open_conn")

	dburl := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"

	engine, err := xorm.NewEngine("mysql", dburl)
	if err != nil {
		log.Fatalf("数据库连接失败:%s", err)
	}

	//连接池空闲数和最大连接数
	engine.SetMaxIdleConns(dbIdle)
	engine.SetMaxOpenConns(dbConn)

	//engine.Ping测试数据库连接
	if err = engine.Ping(); err != nil {
		log.Fatalf("数据库连接失败:%s", err)
		os.Exit(1)
	}
	//同步数据结构
	engine.Sync2(new(User))

	runmode, _ := conf.Cfg.GetValue("", "runmode")
	if !strings.EqualFold(runmode, "dev") {
		// 输出sql到log
		var filelog string
		filelog = "log/xorm-sql-" + strings.Split(time.Now().String(), " ")[0] + ".log"
		f, err := os.OpenFile(filelog, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalf("日志文件建立失败:%s", err)
			os.Exit(1)
		}
		engine.Logger = xorm.NewSimpleLogger(f)
	}

	//读取配置文件，设置xorm日志状态
	showSQL, _ := conf.Cfg.GetValue("xorm", "showSql")
	showDebug, _ := conf.Cfg.GetValue("xorm", "showDebug")
	showError, _ := conf.Cfg.GetValue("xorm", "showError")
	showWarn, _ := conf.Cfg.GetValue("xorm", "showWarn")

	engine.ShowSQL = strings.EqualFold(showSQL, "true")
	engine.ShowDebug = strings.EqualFold(showDebug, "true")
	engine.ShowErr = strings.EqualFold(showError, "true")
	engine.ShowWarn = strings.EqualFold(showWarn, "true")

	Xorm = engine
}
