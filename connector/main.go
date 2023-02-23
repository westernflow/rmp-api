package main

import (
	"context"
	_ "fmt"
	"log"
	"math"
	"os"

	// "os"
	model "rmpParser/models"
	uwomodel "rmpParser/uwomodel"
	"strconv"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GAMEPLAN:

// 1: WE WILL QUERY THROUGH EACH COURSE IN THE COURSES COLLECTION

// ---
// FOR A GIVEN COURSE,
// -> ASSUME STRUCTURE IS IN THE FORM ([A-Za-z]. ([A-Za-z] )* [A-Za-z])*
// ^ ASSUMING THIS STRUCTURE HOLDS, WE CAN FIND EACH PROFESSOR IN THE LIST BY SPLITTING JUST PRIOR TO THE FIRST INITIAL

// NOW WE WILL HAVE SOME PROFESSORS IN THIS FORM: A. XXX XXX

// ^ MAKE A CHECK TO SEE IF WE HAVE ALREADY SEEN THIS PROFESSOR
// ADD TO THE ARRAY OF COURSES THAT A PROFESSOR TEACHES WITH THE CURRENT GIVEN COURSE
// OTHERWISE, BEGIN THE CREATE PROFESSOR FLOW

// THE CREATE PROFESSOR FLOW WILL BE A SERIES OF STEPS THAT WILL CONSTRUCT THE FINAL PROFESSOR OBJECT WE STORE IN OUR DATABASE
// --- POPULATE AVAILABLE DATA:
// -> WE KNOW THAT THE PROFESSOR TEACHES THE CURRENT GIVEN COURSE, SO WE INITIALIZE COURSESTHEYTEACH := [CURRENTCOURSE]

// --- GET THEIR RMP DATA:
// 1) THE DATA WE WANT FROM RMP IS THEIR FULL NAME, SINCE WE CURRENTLY ONLY HAVE THEIR FIRST INITIAL, AS WELL AS THEIR REVIEWS
// 2) QUERY THE RMPDB WITH THE REGEX FELO MADE WHERE WE ENSURE THAT THEIR FIRST NAME BEGINS WITH THE INITIAL, AND WE ARE ALSO ABLE TO FIND A FULL, CASE INSENSITIVE MATCH WITH THE REST OF THEIR NAME. I.E. NAME.BEGINSWITH(A), NAME.SUBSTRINGIN(XXX XXX)
// 3) IN THE CASE STEP 2) RESULTS IN MULTIPLE MATCHES, WE WILL TAKE THE MATCH WITH THE MOST RATINGS
// -> TRIVIALLY ADD REVIEWS FROM RMP DATA TO OUR NEW DB, AS WELL AS THEIR FULL UPDATED NAME
// --- IF UNABLE TO GET THEIR RMP DATA:
// 1) Fill review with null, leave name as it is

// professor model:
// - courses that they teach ([]strings the course names -> will be href to /course/a[i])
// - full name
// - initialed name
// - reviews []reviews
// - same professor data about helpfulness and shit probably from rmp

// rating = average of quality
// difficulty = average of difficulty

// functiont that takes professor of type uwomodel.proffesor and returns two values,
// rating and diff

func updateRatingAndDiff(professor *uwomodel.Professor) {

	if len(professor.Reviews) == 0 {
		professor.Rating = 0.0
		professor.Difficulty = 0.0
		return
	}

	rating := 0.0
	diff := 0.0
	for _, review := range professor.Reviews {
		rating += review.Quality
		diff += review.Difficulty
	}

	rating = rating / float64(len(professor.Reviews))
	diff = diff / float64(len(professor.Reviews))

	professor.Rating = math.Round(rating*10) / 10
	professor.Difficulty = math.Round(diff*10) / 10
}

func main() {
	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()
	// load .env file from root directory
	err := godotenv.Load("../.env")
	// get the PROD_MONGODB connection string from the .env file
	connectionString := os.Getenv("PROD_MONGODB")
	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// get the courses collection in the uwo-tt-api database
	coursesCollection := client.Database("uwo-tt-api").Collection("courses")

	// get the professors collection in the rmp database
	rmpProfessorsCollection := client.Database("rmpDB").Collection("professors")

	professorsCollection := client.Database("uwo-tt-api").Collection("professors")
	// wipe the final_professors collection
	// _, err = client.Database("uwo-tt-api").Collection("final_professors").DeleteMany(ctx, bson.D{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// get the final_professors collection in the uwo-tt-api database
	_, err = client.Database("uwo-tt-api").Collection("professors_temp").DeleteMany(ctx, bson.D{})
	tempProfessorsCollection := client.Database("uwo-tt-api").Collection("professors_temp")

	// query through each course in the courses collection
	cursor, err := coursesCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	// for each course, we will query through the professors in the professors collection
	professorSet := make(map[string]bool, 10000)
	for cursor.Next(ctx) {
		var section uwomodel.Section
		cursor.Decode(&section)
		sectionData := section.SectionData
		professors := make([]string, 0)
		professorsString := sectionData.Instructor

		// the courseString is the course name in the form of "FACULTY NUMBER+SUFFIX" (e.g. "CS 1350A")
		courseString := section.CourseData.Faculty + " " + strconv.Itoa(section.CourseData.Number) + section.CourseData.Suffix

		currentString := ""
		if (professorsString) == "." {
			continue
		}

		// iterate through each character in the string
		for i := 0; i < len(professorsString); i++ {
			if professorsString[i] == '.' && i != 1 {
				// add the current string not including the last two characters (the space and first initial)
				professor := currentString[:len(currentString)-2]
				if professorSet[professor] == false {
					professorSet[professor] = true
					professors = append(professors, professor)
				}
				currentString = currentString[len(currentString)-1:] + string(professorsString[i])
			} else {
				currentString += string(professorsString[i])
			}
		}

		professors = append(professors, currentString)

		for _, professor := range professors {
			var rmpProfessor model.MongoProfessor
			var finalProfessor uwomodel.Professor
			var existingProfessor uwomodel.Professor
			firstInitial := professor[0:1]
			restOfName := professor[3:]

			// query for the professor in the professors collection with the first initial and the rest of the name
			// use this regexp: ^A. X$, where A=firstInitial, X=restOfName
			regexString := "^" + firstInitial + "." + "*" + restOfName + "$"

			// filter for the professor using the regex, this is for checking in rmp professors db
			filter := bson.M{"name": bson.M{"$regex": regexString}}

			// filter for the professor using the regex, this is for checking in final professors db

			// this upates the courses array in the final professors db, uses courseString from earlier
			coursesUpdate := bson.D{
				{Key: "$addToSet", Value: bson.D{
					{Key: "currentCourses", Value: courseString}}},
			}

			//  this updates the departments array in the final professors db, uses CourseData.Faculty
			departmentUpdate := bson.D{
				{Key: "$addToSet", Value: bson.D{{Key: "departments", Value: section.CourseData.Faculty}}},
			}

			tempErr := tempProfessorsCollection.FindOne(ctx, filter).Decode(&finalProfessor)

			if tempErr == nil {
				tempProfessorsCollection.FindOneAndUpdate(ctx, filter, coursesUpdate).Decode(&finalProfessor)
				tempProfessorsCollection.FindOneAndUpdate(ctx, filter, departmentUpdate).Decode(&finalProfessor)
				// if a professor is found in the professors collection, then add their existing
				// reviews to the final professors collection
				err := professorsCollection.FindOne(ctx, filter).Decode(&existingProfessor)

				// check if the professor is already in the final professors collection, and update the courses array
				if err == nil {
					reviewUpdate := bson.D{
						{Key: "$addToSet", Value: bson.D{{Key: "reviews", Value: existingProfessor.Reviews}}},
					}
					// add all reviews from the existing professor to the final professor
					tempProfessorsCollection.FindOneAndUpdate(ctx, filter, reviewUpdate).Decode(&finalProfessor)
				}
				continue
			}

			// logic to create the professor from scratch
			professorCursor, _ := rmpProfessorsCollection.Find(ctx, filter)

			if professorCursor.Next(ctx) {
				// decode the professor from the rmp professors collection
				professorCursor.Decode(&rmpProfessor)
				finalProfessor.Name = professor            // m millard
				finalProfessor.RMPName = rmpProfessor.Name // max millard
				finalProfessor.RMPId = rmpProfessor.RMPId  // this techncially an array but fuck it
				// finalProfessor.Reviews = make([]uwomodel.Review, len(rmpProfessor.Reviews))

				for _, review := range rmpProfessor.Reviews {
					review := uwomodel.Review{
						ProfessorID: review.ProfessorID,
						Professor:   review.Professor,
						Quality:     review.Quality,
						Difficulty:  review.Difficulty,
						Date:        review.Date,
						ReviewText:  review.ReviewText,
						Helpful:     review.Helpful,
						Clarity:     review.Clarity,
					}
					finalProfessor.Reviews = append(finalProfessor.Reviews, review)
				}

				for professorCursor.Next(ctx) {
					professorCursor.Decode(&rmpProfessor)
					// initialize the final professor
					// merge all reviews to this professor
					for _, review := range rmpProfessor.Reviews {
						review := uwomodel.Review{
							ProfessorID: review.ProfessorID,
							Professor:   review.Professor,
							Quality:     review.Quality,
							Difficulty:  review.Difficulty,
							Date:        review.Date,
							ReviewText:  review.ReviewText,
							Helpful:     review.Helpful,
							Clarity:     review.Clarity,
						}
						finalProfessor.Reviews = append(finalProfessor.Reviews, review)
					}
				}
				// department addition
				finalProfessor.Departments = []string{section.CourseData.Faculty}

				// courses addition
				finalProfessor.CurrentCourses = []string{courseString}
			} else {
				// otherwise we will create the professor object without the rmpprofessor
				// fmt.Println("professor not found", professor)
				finalProfessor = uwomodel.Professor{
					Name:           professor,
					RMPName:        professor,
					Reviews:        []uwomodel.Review{},
					Departments:    []string{section.CourseData.Faculty},
					CurrentCourses: []string{courseString},
				}
			}
			// create a new professors collection in the uwo-tt-api database
			updateRatingAndDiff(&finalProfessor)
			tempProfessorsCollection.InsertOne(ctx, finalProfessor)
		}
	}

	// drop the old professors collection and rename the new one to professors
	professorsCollection.Drop(ctx)
	// copy the temp professors collection to the professors collection
	pipeline := bson.A{
		bson.M{"$match": bson.M{}},
		bson.M{"$out": "professors"},
	}

	tempProfessorsCollection.Aggregate(ctx, pipeline)
	// drop the temp professors collection
	tempProfessorsCollection.Drop(ctx)
}
