package controllers

import (
	"webgofirst/models"

	beego "github.com/beego/beego/v2/server/web"
)

type FirstController struct {
	beego.Controller
}

type Employees []models.Employee

var employees []models.Employee

func init() {
	employees = Employees{
		models.Employee{Id: 1, FirstName: "Foo", LastName: "Bar"},
		models.Employee{Id: 2, FirstName: "Abdelghani", LastName: "Sebbai"},
	}
}
func (ctrl *FirstController) GetEmployees() {
	ctrl.Ctx.ResponseWriter.WriteHeader(200)
	ctrl.Data["json"] = employees
	ctrl.ServeJSON()
}
func (ctrl *FirstController) Dashbaord() {
	ctrl.Data["employees"] = employees
	ctrl.TplName = "dashboard.tpl"
}

func (ctrl *FirstController) GetEmployee() {
	var id int
	ctrl.Ctx.Input.Bind(&id, "id")
	var isEmployeeExist bool
	var emps []models.Employee
	for _, employee := range employees {
		if employee.Id == id {
			emps = append(emps, models.Employee{
				Id:        employee.Id,
				FirstName: employee.FirstName, LastName: employee.LastName,
			})
			isEmployeeExist = true
			break
		}
	}
	if !isEmployeeExist {
		ctrl.Abort("Generic")
	} else {
		ctrl.Data["employees"] = emps
		ctrl.TplName = "dashboard.tpl"
	}
}
