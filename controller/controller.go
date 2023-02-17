package controller

import (
	// postgres
	"fmt"
	"os"
	"reflect"
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
}

func (c *controller) InitializeDatabase() {
	fmt.Println("Droping tables...")
	c.DropTables()
	fmt.Println("Creating tables...")
	c.CreateTables()
	fmt.Println("Database initialized!")
}

func (c *controller) DropTables() {
	// drop tables in dependency order
	for _, model := range []interface{}{&model.Review{}, &model.Course{}, &model.Professor{}, &model.Department{}} {
		if err := c.db.DropTableIfExists(model).Error; err != nil {
			fmt.Printf("Error dropping table %v: %v\n", reflect.TypeOf(model).Elem().Name(), err)
			return
		}
		fmt.Printf("Dropped table %v\n", reflect.TypeOf(model).Elem().Name())
	}
}

func (c *controller) CreateTables() {
	c.db.AutoMigrate(&model.Professor{}, &model.Review{}, &model.Department{}, &model.Course{})
	c.db.Model(&model.Professor{}).AddUniqueIndex("idx_rmp_id", "rmp_id")
	c.db.Model(&model.Review{}).AddForeignKey("professor_id", "professors(rmp_id)", "CASCADE", "CASCADE")
}

// populate databse with all professors from all departments
func (c *controller) PopulateDatabase(departments []model.Department) {
	fmt.Println("Populating database...")
	// get all professors from each department and insert into database
	for i, department := range departments {
		// displays percentages rounded to 2 decimal places
		fmt.Println("Fetching data from department: ", department.Name, "; Percentage remaining: ", fmt.Sprintf("%.2f", float64(i)/float64(len(departments))*100), "%")
		c.InsertDepartment(department)
		worker.AddProfessorsFromDepartmentToDatabase(c, department.DepartmentBase64Code)
	}
}

func (c *controller) GetDepartmentByBase64Code(base64Code string) (department model.Department, err error) {
	// get department from database
	err = c.db.Where("department_base64_code = ?", base64Code).First(&department).Error
	return
}

func (c *controller) InsertDepartment(department model.Department) {
	// insert department into database
	fmt.Println("Inserting department: ", department.Name)
	c.db.Create(&department)
}

func (c *controller) InsertProfessor(department model.Department, professor model.Professor) {
	reviews := professor.Reviews
	professor.Reviews = nil
	// insert professor into database
	// first check if a professor exists in the current database with the same RMPID:
	var existingProfessor model.Professor
	c.db.Where("rmp_id = ?", professor.RMPId).First(&existingProfessor)
	if existingProfessor.RMPId != "" {
		// if a professor exists, update the professor with the new data
		fmt.Println("Exist professor", existingProfessor.Name, "with id", existingProfessor.RMPId, existingProfessor.Departments)
		fmt.Println("New professor", professor.Name, "with id", professor.RMPId, professor.Departments)
		fmt.Println("Adding department", department.Name, "to professor", existingProfessor.Name)
		// iterate through each department in existingprofessor.departments and check if it exists in professor.departments
		// if it does not exist, add it to professor.departments
		// check if department.Name exists in professor.Departments
		departmentExist := false
		for _, existingDepartment := range existingProfessor.Departments {
			if existingDepartment == department.Name {
				departmentExist = true
			}
		}
		if !departmentExist {
			existingProfessor.Departments = append(existingProfessor.Departments, department.Name)
		}

		// update the professor in the database
		c.db.Model(&existingProfessor).Update(&professor)
	} else {
		// if a professor does not exist, insert the professor into the database
		c.db.Create(&professor)
	}
	// insert reviews into database
	c.InsertReviews(reviews)
}

func (c *controller) InsertReviews(reviews []model.Review) {
	// insert reviews into database
	for _, review := range reviews {
		c.InsertReview(review)
	}
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
