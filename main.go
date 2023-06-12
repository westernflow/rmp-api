package main

import (
	"github.com/joho/godotenv"

	mongocontroller "rmpParser/mongoController"
	"rmpParser/worker"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	c := mongocontroller.GetInstance()
	// connect to the database
	c.ConnectToDatabase()

	// populate the database (dont do this if you already have data in the database)
	// c.InitializeDatabase()
	// c.PopulateDatabase()
	worker.Scrape()
}
