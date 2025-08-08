package main

import (
	"webgofirst/controllers"
	_ "webgofirst/routers"

	"github.com/beego/beego/v2/server/web"
)

func main() {
	// Enable and configure sessions
	// web.BConfig.WebConfig.Session.SessionOn = true
	// web.BConfig.WebConfig.Session.SessionProvider = "file" // memory, redis, mysql or other providers
	// web.BConfig.WebConfig.Session.SessionProviderConfig = "./sessions"//directory path where session files will be stored
	// web.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600
	web.ErrorController(&controllers.ErrorController{})
	web.Run()
}
