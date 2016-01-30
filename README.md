网站

将项目代码放到gopath目录下

项目中用到的第三方包:
xorm	
安装命令:go get github.com/go-xorm/xorm
官网:http://xorm.io/

goconfig
安装命令:go get github.com/Unknwon/goconfig

mysql
安装命令:go get github.com/go-sql-driver/mysql

conf/base.go
读取配置文件

models/base.go
加载数据库

log
开发者模式，日志不会写入log文件夹中，而是直接在控制台打印出来