package main

import (
	controller "rmpParser/mongoController"
	"rmpParser/worker"

	"github.com/joho/godotenv"
)

func main() {
	// load arg variables
	// fmt.Println(os.Args)
	dropTables := false
	// if len(os.Args) >= 1 {
	// 	dropTables = os.Args[1] == "init"
	// }

	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	getData(dropTables)
}

func getData(dropTables bool) {
	// create a controller
	c := controller.GetInstance()

	// // connect to the database
	c.ConnectToDatabase()

	// // initialize the database
	if dropTables {
		c.InitializeDatabase()
	}

	c.CreateTables()

	// if !dropTables {
	// 	departments = getDepartmentSplice(departments, "business")
	// }

	data := worker.GetKProfessorAtCursor(1, "")
	c.PopulateDatabase(data)
	// fmt.Println(data)
}
