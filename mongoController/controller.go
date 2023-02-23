package mongocontroller

import (
	"context"
	"fmt"
	"os"
	model "rmpParser/models"
	worker "rmpParser/worker"

	// mongo
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type controller struct { // make ctor private for singleton
	db          *mongo.Database
	professors  *mongo.Collection
	departments *mongo.Collection
}

var instance *controller

func GetInstance() *controller {
	if instance == nil {
		instance = new(controller)
	}
	return instance
}

func (c *controller) ConnectToDatabase() {
	// Set client options
	fmt.Println("Connecting to database...")
	// load env
	err := godotenv.Load("../.env")
	// get mongo uri from env
	mongoURI := os.Getenv("PROD_MONGODB")
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	fmt.Println("Pinging to check connection...")
	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database!")
	mongoDB := client.Database("rmpDB")
	c.db = mongoDB
}

func (c *controller) InitializeDatabase() {
	fmt.Println("Dropping tables...")
	c.DropTables()
	fmt.Println("Creating tables...")
	c.CreateTables()
	fmt.Println("Database initialized!")
}

func (c *controller) DropTables() {
	// Delete the database
	err := c.db.Drop(context.Background())
	if err != nil {
		panic(err)
	}
}

func (c *controller) CreateTables() {
	// set professors and departments collections
	c.professors = c.db.Collection("professors")
	c.departments = c.db.Collection("departments")
}

func (c *controller) PopulateDatabase(departments []model.MongoDepartment) {
	fmt.Println("Populating database...")
	for i, department := range departments {
		// displays percentages rounded to 2 decimal places
		fmt.Println("Fetching data from department: ", department.Name, "; Percentage done: ", fmt.Sprintf("%.2f", float64(i)/float64(len(departments))*100), "%")
		c.InsertDepartment(department)
		worker.AddProfessorsFromDepartmentToDatabase(c, department)
	}
}

func (c *controller) GetDepartmentByBase64Code(base64Code string) (department model.MongoDepartment, err error) {
	// get department from database given base64Code
	filter := bson.D{
		{Key: "departmentBase64Code", Value: base64Code},
	}
	err = c.departments.FindOne(context.Background(), filter).Decode(&department)
	return
}

func (c *controller) GetAllProfessors() []model.MongoProfessor {
	// get all professors from database

	var professors []model.MongoProfessor

	cursor, err := c.professors.Find(context.Background(), model.MongoProfessor{})
	if err != nil {
		panic(err)
	}
	if err = cursor.All(context.Background(), &professors); err != nil {
		panic(err)
	}

	return professors
}

func (c *controller) CloseConnection() {
	// close connection to database
	c.db.Client().Disconnect(context.Background())
}

func (c *controller) GetAllReviews() []model.MongoReview {
	// query through all professors and get all reviews
	var reviews []model.MongoReview
	professors := c.GetAllProfessors()
	for _, professor := range professors {
		reviews = append(reviews, professor.Reviews...)
	}
	return reviews
}

func (c *controller) GetAllDepartments() []model.MongoDepartment {
	// get all departments from database
	var departments []model.MongoDepartment

	cursor, err := c.departments.Find(context.Background(), model.MongoDepartment{})
	if err != nil {
		panic(err)
	}
	if err = cursor.All(context.Background(), &departments); err != nil {
		panic(err)
	}

	return departments
}

func (c *controller) InsertDepartment(department model.MongoDepartment) {
	// insert department into database
	_, err := c.departments.InsertOne(context.Background(), department)
	if err != nil {
		panic(err)
	}
}

func (c *controller) InsertProfessor(department model.MongoDepartment, professor model.MongoProfessor) {
	// check if the professor already exists in the database
	var existingProfessor model.MongoProfessor
	// query for professor with the same RMPId
	filter := bson.D{
		{Key: "rmpId", Value: professor.RMPId},
	}

	update := bson.M{
		"$addToSet": bson.M{
			"departments": department}}

	result := c.professors.FindOne(context.Background(), filter).Decode(&existingProfessor)

	if result != nil {
		professor.Departments[0].DepartmentBase64Code = department.DepartmentBase64Code
		professor.Departments[0].Name = department.Name
		professor.Departments[0].DepartmentNumber = department.DepartmentNumber
		_, err := c.professors.InsertOne(context.Background(), professor)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Adding department: ", department.Name, " to professor: ", professor.Name)
		departmentExists := false
		for _, existingDepartment := range existingProfessor.Departments {
			if existingDepartment.DepartmentBase64Code == department.DepartmentBase64Code {
				departmentExists = true
				break
			}
		}
		if !departmentExists {
			// append the new department to the professor's departments array
			existingProfessor.Departments = append(existingProfessor.Departments, department)
			// update professor in database
			_, err := c.professors.UpdateOne(context.Background(), filter, update)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (c *controller) GetProfessors() []model.MongoProfessor {
	// get all professors from database
	var professors []model.MongoProfessor

	cursor, err := c.professors.Find(context.Background(), model.MongoProfessor{})
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.Background(), &professors); err != nil {
		panic(err)
	}

	return professors
}

func (c *controller) GetDepartments() []model.MongoDepartment {
	//  get all departments from database
	var departments []model.MongoDepartment

	cursor, err := c.departments.Find(context.Background(), model.MongoDepartment{})
	if err != nil {
		panic(err)
	}

	if err = cursor.All(context.Background(), &departments); err != nil {
		panic(err)
	}

	return departments
}
