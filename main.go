package main

import (
	// "fmt"
	"database/sql"
	"fmt"
	"rmpParser/worker"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres@db:5432/postgres?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	// ping the database to check if it is alive
	fmt.Println("Pinging the database...")
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to the database!")
	// create the database
	_, err = db.Exec("CREATE DATABASE mydb")
    if err != nil {
        panic(err)
    }

	// get all departments from the school
	departments := worker.GetDepartments()
	// fmt.Println(departments)
	// get all professors from each department
	for _, department := range departments {
		worker.AddProfessorsFromDepartmentToDatabase(department.ID)
	}
}