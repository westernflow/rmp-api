package controller

import (
	// postgres
	"fmt"
	"os"
	model "rmpParser/models"
	"rmpParser/worker"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type controller struct { // make ctor private for singleton
	db *gorm.DB
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
	uri := os.Getenv("RDS_URI")
	db, err := gorm.Open("postgres", uri)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database!")

	c.db = db
	c.CreateTables()
	fmt.Println("Tables created successfully")
}

func (c *controller) CreateTables() {
	c.db.AutoMigrate(&model.Professor{}, &model.Review{}, &model.Department{}, &model.Course{})
	c.db.Model(&model.Review{}).AddForeignKey("professor_id", "professors(id)", "CASCADE", "CASCADE")
	c.db.Model(&model.Course{}).AddForeignKey("professor_id", "professors(id)", "CASCADE", "CASCADE")
	// create a professor
	// c.db.Create(&model.Professor{Name: "John", Rating: 4.5, Difficulty: 3.5, Department: "CS", RMPId: "1234", Courses: []model.Course{{Number: "CS 123", Department: "CS"}, {Number: "CS 456", Department: "CS"}}})
}

// populate databse with all professors from all departments
func (c *controller) PopulateDatabase() {

	fmt.Println("Populating database...")
	// get all departments from the school
	departments := worker.GetDepartments()

	// get all professors from each department and insert into database
	for _, department := range departments {
		c.InsertDepartment(department)
		worker.AddProfessorsFromDepartmentToDatabase(c, department.DepartmentBase64Code)
	}
}

func (c *controller) InsertDepartment(department model.Department) {
	// insert department into database
	fmt.Println("Inserting department: ", department.Name)
	c.db.Create(&department)
}

func (c *controller) InsertProfessor(professor model.Professor) {
	// insert professor into database
	// fmt.Println("Inserting professor: ", professor)
	// for _, course := range professor.Courses {
	// 	fmt.Println("Inserting course: ", course.Number)
	// }
	c.db.Create(&professor)
}

func (c *controller) InsertReview(review model.Review) {
	// insert review into database
	c.db.Create(&review)
}

func (c *controller) InsertCourse(course model.Course) {
	// insert course into database
	c.db.Create(&course)
}

func (c *controller) GetProfessors() []model.Professor {
	var professors []model.Professor
	c.db.Find(&professors)
	return professors
}

func (c *controller) GetDepartments() []model.Department {
	var departments []model.Department
	c.db.Find(&departments)
	return departments
}
