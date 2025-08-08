package controllers

import "github.com/beego/beego/v2/server/web"

type SessionController struct {
	web.Controller
}

func (ctrl *SessionController) Login() {
	ctrl.SetSession("authenticated", true)
	ctrl.Ctx.ResponseWriter.WriteHeader(200)
	ctrl.Ctx.WriteString("You have successfully logged in.")
}

func (ctrl *SessionController) Logout() {
	ctrl.SetSession("authenticated", false)
	ctrl.Ctx.ResponseWriter.WriteHeader(200)
	ctrl.Ctx.WriteString("You have successfully logged out.")
}

func (ctrl *SessionController) Home() {
	isAuthenticated := ctrl.GetSession("authenticated")
	if isAuthenticated == nil || isAuthenticated == false {
		ctrl.Ctx.WriteString("You are unauthorized to view the page.")
		return
	}
	ctrl.Ctx.ResponseWriter.WriteHeader(200)
	ctrl.Ctx.WriteString("Home Page")
}
