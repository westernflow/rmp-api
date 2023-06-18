// This file is responsible for transferring data from our postgres database to mongoDB
package main

import (
	"fmt"
	model "rmpParser/models"
	controller "rmpParser/mongoController"

	// postgres

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"

	// mongo
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// load env variables
	err := godotenv.Load("../.env")

	// create mongo connection
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27018"))
	if err != nil {
		fmt.Println(err)
	}
	err = client.Connect(nil)
	if err != nil {
		fmt.Println(err)
	}

	// client.ping to check if connection is successful
	err = client.Ping(nil, nil)
	if err != nil {
		fmt.Println(err)
	}

	// create postgres connection
	// create a controller
	c := controller.GetInstance()

	// connect to the database
	c.ConnectToDatabase()

	// drop existing rmpDB and create new one
	client.Database("rmpDB").Drop(nil)
	mongoDB := client.Database("rmpDB")
	professorCollection := mongoDB.Collection("professors")
	departmentCollection := mongoDB.Collection("departments")
	reviewCollection := mongoDB.Collection("reviews")

	// get all professors from postgres
	professors := c.GetAllProfessors()
	fmt.Println("Got all professors from postgres", len(professors))
	// insert all professors into mongo
	for _, professor := range professors {
		fmt.Println("Inserting professor", professor.Reviews, professor.Departments)
		// build mongodb professor departments
		mongoDepartments := []model.MongoDepartment{}
		for _, department := range professor.Departments {
			// convert departmentID uint to string
			departmentID := fmt.Sprintf("%d", department.ID)
			mongoDepartments = append(mongoDepartments, model.MongoDepartment{Name: department.Name, DepartmentBase64Code: departmentID})
		}

		// build mongodb professor reviews
		mongoReviews := []model.MongoReview{}
		for _, review := range professor.Reviews {
			mongoReviews = append(mongoReviews, model.MongoReview{ProfessorID: professor.RMPId, Professor: professor.Name,
				Quality: review.Quality, Difficulty: review.Difficulty, Date: review.Date, ReviewText: review.ReviewText,
				Helpful: review.Helpful, Clarity: review.Clarity})
		}

		mongoProfessor := model.MongoProfessor{Name: professor.Name, RMPId: professor.RMPId,
			Rating: professor.Rating, Difficulty: professor.Difficulty, Departments: mongoDepartments, Reviews: mongoReviews}

		// fmt.Println("Inserting professor", professor.RMPId)
		professorCollection.InsertOne(nil, mongoProfessor)
	}

	// get all departments from postgres
	departments := c.GetAllDepartments()

	// insert all departments into mongo
	for _, department := range departments {
		fmt.Println("Inserting department", department.Name)
		// convert department to mongoDB department
		departmentID := fmt.Sprintf("%d", department.ID)
		mongoDept := model.MongoDepartment{Name: department.Name, DepartmentBase64Code: departmentID}

		departmentCollection.InsertOne(nil, mongoDept)
	}

	// get all reviews from postgres
	reviews := c.GetAllReviews()

	// insert all reviews into mongo
	for _, review := range reviews {
		fmt.Println("Inserting review", review)
		mongoReview := model.MongoReview{ProfessorID: review.ProfessorID, Professor: review.Professor,
			Quality: review.Quality, Difficulty: review.Difficulty, Date: review.Date, ReviewText: review.ReviewText,
			Helpful: review.Helpful, Clarity: review.Clarity}
		reviewCollection.InsertOne(nil, mongoReview)
	}

	// close the connection
	c.CloseConnection()
}
