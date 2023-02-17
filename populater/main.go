package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"rmpParser/controller"
	"rmpParser/handler"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	// lambda.Start(HandleRequest)
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

	// populate the database (dont do this if you already have data in the database)
	c.PopulateDatabase()
}
