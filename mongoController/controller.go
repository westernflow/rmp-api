package mongocontroller

import (
	"context"
	"fmt"
	"os"
	model "rmpParser/models"
	"rmpParser/uwomodel"

	// mongo
	"github.com/joho/godotenv"
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
	mongoURI := os.Getenv("LOCAL_MONGODB")
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

}

func (c *controller) PopulateDatabase(given model.TeacherSearchResults) {
	fmt.Println("Populating database...")

	profs := given.Data.Search.Teachers.Edges
	cleanedProfs := []uwomodel.Professor{}
	for _, prof := range profs {
		prof := prof.Node
		fmt.Println(prof)
		uwoReviews := []uwomodel.Review{}
		for _, review := range prof.Ratings.Edges {
			review := review.Node
			// create mongoReview
			mongoReview := uwomodel.Review{
				ProfessorID: prof.ID,
				Quality:     review.QualityRating,
				Clarity:     review.ClarityRating,
				Difficulty:  review.DifficultyRating,
				Helpful:     review.HelpfulRating,
				Date:        review.Date,
				ReviewText:  review.Comment,
			}
			uwoReviews = append(uwoReviews, mongoReview)
		}
		mongoProf := uwomodel.Professor{
			RMPName: prof.FirstName + " " + prof.LastName,
			Reviews: uwoReviews,
		}
		fmt.Println(mongoProf)

		cleanedProfs = append(cleanedProfs, mongoProf)
	}

	// insert professors into database
	c.InsertProfessors(cleanedProfs)
}

// func (c *controller) GetDepartmentByBase64Code(base64Code string) (department model.MongoDepartment, err error) {
// 	// get department from database given base64Code
// 	filter := bson.D{
// 		{Key: "departmentBase64Code", Value: base64Code},
// 	}
// 	err = c.departments.FindOne(context.Background(), filter).Decode(&department)
// 	return
// }

// func (c *controller) GetAllProfessors() []model.MongoProfessor {
// 	// get all professors from database

// 	var professors []model.MongoProfessor

// 	cursor, err := c.professors.Find(context.Background(), model.MongoProfessor{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	if err = cursor.All(context.Background(), &professors); err != nil {
// 		panic(err)
// 	}

// 	return professors
// }

func (c *controller) CloseConnection() {
	// close connection to database
	c.db.Client().Disconnect(context.Background())
}

// func (c *controller) GetAllReviews() []model.MongoReview {
// 	// query through all professors and get all reviews
// 	var reviews []model.MongoReview
// 	// professors := c.GetAllProfessors()
// 	for _, professor := range professors {
// 		reviews = append(reviews, professor.Reviews...)
// 	}
// 	return reviews
// }

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

// array of professors
func (c *controller) InsertProfessors(professors []uwomodel.Professor) {

	for _, prof := range professors {
		_, err := c.professors.InsertOne(context.Background(), prof)
		if err != nil {
			panic(err)
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
