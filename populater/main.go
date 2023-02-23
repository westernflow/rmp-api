package main

import (
	"fmt"
	model "rmpParser/models"
	controller "rmpParser/mongoController"
	"rmpParser/worker"

	"os"

	"github.com/joho/godotenv"
)

func main() {
	// load arg variables
	fmt.Println(os.Args)
	dropTables := false
	if len(os.Args) >= 1 {dropTables = os.Args[1] == "init"}

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	getData(dropTables)
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

	c.PopulateDatabase(departments)
}

func getDepartment(departments []model.MongoDepartment, name string) []model.MongoDepartment {
	for i, department := range departments {
		if department.Name == name {
			return departments[i : i+1]
		}
	}
	return []model.MongoDepartment{}
}

func getDepartmentSplice(departments []model.MongoDepartment, name string) []model.MongoDepartment {
	// finds first occurence of name, returns all departments after that
	for i, department := range departments {
		if department.Name == name {
			return departments[i:]
		}
	}
	return []model.MongoDepartment{}
}
