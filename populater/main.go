package main

import (
	"fmt"
	"rmpParser/controller"
	"rmpParser/handler"
	model "rmpParser/models"
	"rmpParser/worker"

	"github.com/joho/godotenv"
	"os"
)

func main() {
	// load arg variables
	fmt.Println(os.Args)
	dropTables := os.Args[1] == "init"

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	getData(dropTables)

	// this should show all the professors in the database
	fmt.Println(handler.GetProfessors())
}

func getData(dropTables bool) {
	// create a controller
	c := controller.GetInstance()

	// connect to the database
	c.ConnectToDatabase()
	
	// initialize the database
	if dropTables {
		c.InitializeDatabase()
	}

	// get all departments
	departments := worker.GetDepartments()
	
	if !dropTables {
		departments = getDepartmentSplice(departments, "business")
	}

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