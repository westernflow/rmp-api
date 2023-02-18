// This file is responsible for transferring data from our postgres database to mongoDB

import (
	"fmt"
	"rmpParser/controller"
	model "rmpParser/models"
	"rmpParser/worker"
	// postgres
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// mongo
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// create mongo connection
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println(err)
	}
	err = client.Connect(nil)
	if err != nil {	
		fmt.Println(err)
	}
	
	// create postgres connection
	// create a controller
	c := controller.GetInstance()

	// connect to the database
	c.ConnectToDatabase()

	// create
	
	mongoDB := mongoClient.Database("your_database_name")
	professorCollection := mongoDB.Collection("professors")
	departmentCollection := mongoDB.Collection("departments")
	courseCollection := mongoDB.Collection("courses")
	reviewCollection := mongoDB.Collection("reviews")
}
