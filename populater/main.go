package main

import (
	"fmt"
	"rmpParser/controller"
	"rmpParser/handler"
	model "rmpParser/models"
	"rmpParser/worker"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	getData()

	// this should show all the professors in the database
	fmt.Println(handler.GetProfessors())
}

func getData() {
	// create a controller
	c := controller.GetInstance()

	// connect to the database
	c.ConnectToDatabase()
	
	// initialize the database
	c.InitializeDatabase()

	// get all departments
	departments := worker.GetDepartments()
	
	// historyAndOnwards := getDepartmentSplice(departments, "history")

	// populate the database (dont do this if you already have data in the database)
	c.PopulateDatabase(departments)
}

func getDepartment(departments []model.Department, name string) []model.Department {
	for i, department := range departments {
		if department.Name == name {
			return departments[i:i+1]
		}
	}
	return []model.Department{}
}

func getDepartmentSplice(departments []model.Department, name string) []model.Department {
	// finds first occurence of name, returns all departments after that
	for i, department := range departments {
		if department.Name == name {
			return departments[i:]
		}
	}
	return []model.Department{}
}