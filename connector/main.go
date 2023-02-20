package main

import (
	"context"
	"log"
	model "rmpParser/models"
	uwomodel "rmpParser/uwomodel"

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

func main() {
	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27018/")

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
	professorsCollection := client.Database("rmpDB").Collection("professors")
	finalProfessorsCollection := client.Database("uwo-tt-api").Collection("final_professors")

	// query through each course in the courses collection
	cursor, err := coursesCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	// for each course, we will query through the professors in the professors collection
	professorSet := make(map[string]bool, 8000)
	for cursor.Next(ctx) {
		var course uwomodel.Section
		cursor.Decode(&course)
		sectionData := course.SectionData
		professors := make([]string, 0)
		professorsString := sectionData.Instructor
		// fmt.Println(professorsString)

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
		// fmt.Println(professors, len(professors))

		for _, professor := range professors {
			var rmpProfessor model.MongoProfessor
			firstInitial := professor[0:1]
			restOfName := professor[3:]
			// query for the professor in the professors collection with the first initial and the rest of the name
			// use this regexp: ^A. X$, where A=firstInitial, X=restOfName

			regexString := "^" + firstInitial + "." + "*" + restOfName + "$"

			filter := bson.M{"name": bson.M{"$regex": regexString}}
			result := professorsCollection.FindOne(ctx, filter).Decode(&rmpProfessor)
			var finalProfessor uwomodel.Professor
			if result != nil {
				// now create the professor object
				finalProfessor := uwomodel.Professor{}
				finalProfessor.RMPName = rmpProfessor.Name
				finalProfessor.Name = professor

				// review conversion
				finalReviews := make([]uwomodel.Review, 0)
				for _, review := range rmpProfessor.Reviews {
					finalReview := uwomodel.Review{}
					finalReview.ProfessorID = review.ProfessorID
					finalReview.Professor = review.Professor
					finalReview.Quality = review.Quality
					finalReview.Difficulty = review.Difficulty
					finalReview.Date = review.Date
					finalReview.ReviewText = review.ReviewText
					finalReview.Helpful = review.Helpful
					finalReview.Clarity = review.Clarity
					finalReviews = append(finalReviews, finalReview)
				}
				finalProfessor.Reviews = finalReviews

				// departmetn conversion
				finalProfessor.Departments = make([]string, 0)
				finalProfessor.Departments = append(finalProfessor.Departments, course.CourseData.Faculty)

				// courses conversion
				finalCourses := make([]string, 0)
				// COMPSCI 1420A
				finalCourses = append(finalCourses, course.CourseData.Faculty+" "+string(rune(course.CourseData.Number))+course.CourseData.Suffix)
				finalProfessor.CurrentCourses = finalCourses

				// finalProfessor.Departments = []string{course.SectionData.Department}
				// finalProfessor.Courses = []string{course.SectionData.Course}

			} else {
				// otherwise we will create the professor object without the rmpprofessor
				finalProfessor := uwomodel.Professor{}
				finalProfessor.Name = professor
				finalProfessor.RMPName = ""
				finalProfessor.Reviews = make([]uwomodel.Review, 0)
				finalProfessor.Departments = make([]string, 0)
				finalProfessor.Departments = append(finalProfessor.Departments, course.CourseData.Faculty)
				finalProfessor.CurrentCourses = make([]string, 0)
				finalProfessor.CurrentCourses = append(finalProfessor.CurrentCourses, course.CourseData.Faculty+" "+string(rune(course.CourseData.Number))+course.CourseData.Suffix)
				finalProfessor.Rating = 0
				finalProfessor.Difficulty = 0
			}

			// create a new professors collection in the uwo-tt-api database
			finalProfessorsCollection.InsertOne(ctx, finalProfessor)
		}
	}
}

// for each professor,

// for each professor in the course, we will query through the professors in the professors collection
// for _, professor := range course.Professors {
// 	// query through the professors collection to see if we have already seen this professor
// 	var professorExists bool
// 	professorsCollection.FindOne(ctx, bson.D{{"name", professor}}).Decode(&professorExists)
// 	if professorExists {
// 		// if we have already seen this professor, we will update their courses they teach with the current course
// 	}
// }
