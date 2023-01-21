package main

import (
	// "fmt"
	"rmpParser/controller"
	"rmpParser/worker"
)

func main() {
	// create a controller
	c := controller.GetInstance()
	// connect to the database
	c.ConnectToDatabase()
	// fmt.Println(c.GetProfessors())
	// get all departments from the school
	departments := worker.GetDepartments()
	// fmt.Println(departments)
	// get all professors from each department
	for _, department := range departments {
		c.InsertDepartment(department)
		worker.AddProfessorsFromDepartmentToDatabase(c, department.DepartmentBase64Code)
	}
}
