package controllers

import "github.com/beego/beego/v2/server/web"

type ErrorController struct {
	web.Controller
}

func (c *ErrorController) Error404() {
	c.Data["content"] = "Page Not Found"
	c.TplName = "404.tpl"
}
func (c *ErrorController) Error500() {
	c.Data["content"] = "Internal Server Error"
	c.TplName = "500.tpl"
}
func (c *ErrorController) ErrorGeneric() {
	c.Data["content"] = "Some Error Occurred"
	c.TplName = "genericerror.tpl"
}

func (c *ErrorController) ErrorDb() {
	c.Data["content"] = "Database access is impossible"
	c.TplName = "genericerror.tpl"
}
