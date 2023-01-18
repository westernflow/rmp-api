package main

import (
	"rmpParser/worker"
	"rmpParser/controller"	
)

func main() {
	// create a controller
	c := controller.GetInstance()
	// connect to the database
	c.ConnectToDatabase()

	// get all departments from the school
	departments := worker.GetDepartments()
	// fmt.Println(departments)
	// get all professors from each department
	for _, department := range departments {
		worker.AddProfessorsFromDepartmentToDatabase(department.ID)
	}
}
