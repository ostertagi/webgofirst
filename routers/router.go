package routers

import (
	"webgofirst/controllers"
	"webgofirst/filters"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/employees", &controllers.FirstController{}, "get:GetEmployees")
	beego.Router("/dashboard", &controllers.FirstController{}, "get:Dashbaord")
	beego.Router("/home", &controllers.SessionController{}, "get:Home")
	beego.Router("/auth/login", &controllers.SessionController{}, "get:Login")
	beego.Router("/auth/logout", &controllers.SessionController{}, "get:Logout")

	beego.InsertFilter("/*", beego.BeforeRouter, filters.LogManager)
	beego.InsertFilter("/auth/*", beego.BeforeExec, filters.MyFilterFunc)

	//http://localhost:8080/employee?id=4
	beego.Router("/employee", &controllers.FirstController{}, "get:GetEmployee")

	beego.Router("/getFromCache", &controllers.CacheController{}, "get:GetFromCache")
}
