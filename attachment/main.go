package main

import (
	"beeblog/models"
	"beeblog/routers"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	// 注册数据库
	models.RegisterDB()
}

func main() {
	// 开启 ORM 调试模式
	orm.Debug = true
	// 自动建表
	orm.RunSyncdb("default", false, true)

	os.Mkdir("attachment", os.ModePerm)

	// 注册 beego 路由
	routers.Init()

	// 启动 beego
	beego.Run()
}
