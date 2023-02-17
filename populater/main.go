package main

import (
	"fmt"
	"rmpParser/controller"
	"rmpParser/handler"
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
	
	// populate the database (dont do this if you already have data in the database)
	c.PopulateDatabase(departments)
}
