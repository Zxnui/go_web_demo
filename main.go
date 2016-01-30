package main

import (
	"go-web-demo/conf"
	"go-web-demo/controllers"
	_ "go-web-demo/models"
	"net/http"
	"runtime"
)

var HttpPort string

func init() {
	port, _ := conf.Cfg.GetValue("", "httpport")
	HttpPort = ":" + port
}

func main() {
	//多核运行
	runtime.GOMAXPROCS(runtime.NumCPU())
	//静态文件
	http.Handle("/css/", http.FileServer(http.Dir("statics")))
	http.Handle("/js/", http.FileServer(http.Dir("statics")))
	http.Handle("/fonts/", http.FileServer(http.Dir("statics")))
	http.Handle("/img/", http.FileServer(http.Dir("statics")))

	//路由
	http.HandleFunc("/", controllers.Home)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/loginPost", controllers.LoginPost)
	http.HandleFunc("/loginInfo", controllers.LoginInfo)
	http.HandleFunc("/login/forget", controllers.LoginForget)
	http.HandleFunc("/login/changePasswd", controllers.LoginChangePasswd)

	http.Handle("/favicon.ico", http.FileServer(http.Dir(".")))

	http.ListenAndServe(HttpPort, nil)
}
