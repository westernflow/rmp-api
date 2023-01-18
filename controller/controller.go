package controller

import (
	// postgres
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

type controller struct {
	db *sql.DB
}

var instance *controller

func GetInstance() *controller {
	if instance == nil {
		instance = new(controller)
	}
	return instance
}

func (c *controller) ConnectToDatabase() {
	fmt.Println("Connecting to database...")
	db, err := sql.Open("postgres", "postgres://postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	// ping the database to make sure it is up and running
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database!")
	c.db = db
}


