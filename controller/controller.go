package controller

import (
	// postgres
	"database/sql"
	"fmt"
	model "rmpParser/models"

	_ "github.com/lib/pq"
)

type controller struct { // make ctor private for singleton
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
	db, err := sql.Open("postgres", "postgres://postgres@db:5432/postgres?sslmode=disable")
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

	// Create Department table if it doesn't exist yet
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS  departments (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
		base64Encoding TEXT NOT NULL
);`)
if err != nil {
		fmt.Println(err)
		return
}

// Create Professor table
_, err = db.Exec(`CREATE TABLE IF NOT EXISTS professors (
		rmpId SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		rating FLOAT NOT NULL,
		department TEXT NOT NULL,
		difficulty FLOAT NOT NULL,
		reviews JSONB,
		courseCodes JSONB
);`)
if err != nil {
		fmt.Println(err)
		return
}

// Create Course table
_, err = db.Exec(`CREATE TABLE IF NOT EXISTS  courses (
		department TEXT NOT NULL,
		number TEXT NOT NULL,
		PRIMARY KEY (department, number)
);`)
if err != nil {
		fmt.Println(err)
		return
}

fmt.Println("Tables created successfully")
}

func (c *controller) InsertDepartment(department model.Department) {
	// insert department into database
	_, err := c.db.Exec(`INSERT INTO departments (name, base64Encoding) VALUES ($1, $2)`, department.ID, department.Name)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Inserted department into database")
}

