package conf

import (
	"github.com/Unknwon/goconfig"
	"log"
)

var Cfg *goconfig.ConfigFile
var HttpWeb string

//加载配置文件
func init() {
	cfg, err := goconfig.LoadConfigFile("conf/conf.ini")
	if err != nil {
		log.Println("无法加载配置文件:", err.Error())
	}

	Cfg = cfg
	httpport, _ := cfg.GetValue("", "httpport")
	HttpWeb, _ = cfg.GetValue("", "httpweb")
	HttpWeb = HttpWeb + ":" + httpport
}
