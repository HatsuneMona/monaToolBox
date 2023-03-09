package main

import (
	"monaToolBox/bootstrap"
	"monaToolBox/global"
)

func main() {

	bootstrap.InitConfig()

	global.Log = bootstrap.InitLog()
	global.Log.Info("Log system start ok.")
	// 初始化 gorm 数据库
	global.DB = bootstrap.InitDatabase()
	global.Log.Info("gorm database system init ok.")
	defer func() {
		if global.DB != nil {
			db, _ := global.DB.DB()
			db.Close()
		}
	}()

	// 初始化Redis
	global.Redis = bootstrap.InitializeRedis()
	global.Log.Info("redis database system init ok.")

	// 初始化memCached
	global.MemCached = bootstrap.InitializeMemcached()
	global.Log.Info("memCached system init ok.")

	// 初始化验证器
	bootstrap.InitializeValidator()
	global.Log.Info("Initialize Validator system init ok.")

	// 启动 go-gin web服务器
	global.Log.Info("runServer.")
	bootstrap.RunServer()

}
