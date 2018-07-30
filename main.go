package main

import (
	"encoding/gob"
	"port-forward/models"
	_ "port-forward/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	logs.SetLogger(logs.AdapterConsole, `{"level":7}`)
	logs.SetLogger(logs.AdapterFile, `{"filename":"app.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)

	//输出文件名和行号
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	//为了让日志输出不影响性能，开启异步日志
	logs.Async()

	//开启seesion支持，默认使用的存储引擎为：memory
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "sessionID"
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
	//beego.BConfig.WebConfig.Session.SessionProvider = "file"
	//beego.BConfig.WebConfig.Session.SessionProviderConfig = "./session"
	gob.Register(&models.LoginUser{})


	logs.Debug("================================================================")
	logs.Debug("               SDN网元数据保护系统v3.0 启动...")
	logs.Debug("    _____        _           _____           _            _")
	logs.Debug("   |  __ \\      | |         |  __ \\         | |          | |")
	logs.Debug("   | |  | | __ _| |_ __ _   | |__) | __ ___ | |_ ___  ___| |_")
	logs.Debug("   | |  | |/ _` | __/ _` |  |  ___/ '__/ _ \\| __/ _ \\/ __| __|")
	logs.Debug("   | |__| | (_| | || (_| |  | |   | | | (_) | ||  __/ (__| |_")
	logs.Debug("   |_____/ \\__,_|\\__\\__,_|  |_|   |_|  \\___/ \\__\\___|\\___|\\__|")
	logs.Debug("")
	logs.Debug("=================================================================")

	//默认static目录是可以直接访问的，其它目录需要单独指定
	beego.SetStaticPath("/theme", "theme")
	beego.SetStaticPath("/flowcheck", "flowcheck")

	//启动应用
	beego.Run()

}
