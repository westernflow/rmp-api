package main

import (
	"context"
	_ "fmt"
	"log"
	model "rmpParser/models"
	uwomodel "rmpParser/uwomodel"
	"strconv"

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

	// wipe the final_professors collection
	_, err = client.Database("uwo-tt-api").Collection("final_professors").DeleteMany(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	// get the final_professors collection in the uwo-tt-api database
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

		// the courseString is the course name in the form of "FACULTY NUMBER+SUFFIX" (e.g. "CS 1350A")
		courseString := course.CourseData.Faculty + " " + strconv.Itoa(course.CourseData.Number) + course.CourseData.Suffix

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
			firstInitial := professor[0:1]
			restOfName := professor[3:]

			// query for the professor in the professors collection with the first initial and the rest of the name
			// use this regexp: ^A. X$, where A=firstInitial, X=restOfName
			regexString := "^" + firstInitial + "." + "*" + restOfName + "$"

			// filter for the professor using the regex, this is for checking in rmp professors db
			filter := bson.M{"name": bson.M{"$regex": regexString}}

			// filter for the professor using the regex, this is for checking in final professors db
			rmpFilter := bson.M{"rmpName": bson.M{"$regex": regexString}}

			// this upates the courses array in the final professors db, uses courseString from earlier
			coursesUpdate := bson.D{
				{Key: "$addToSet", Value: bson.D{
					{Key: "currentCourses", Value: courseString}}},
			}

			//  this updates the departments array in the final professors db, uses CourseData.Faculty
			departmentUpdate := bson.D{
				{Key: "$addToSet", Value: bson.D{{Key: "departments", Value: course.CourseData.Faculty}}},
			}

			// check if the professor is already in the final professors collection, and update the courses array
			err := finalProfessorsCollection.FindOneAndUpdate(ctx, rmpFilter, coursesUpdate).Decode(&finalProfessor)
			// if they are in the final professors collection, update the department too
			if err == nil {
				// fmt.Println("found prof in final", finalProfessor.RMPName)

				// update the department, it might be more efficient to do this in 1 find call and 2 update calls or updating
				// with the _id, im not sure, but this works for now

				err = finalProfessorsCollection.FindOneAndUpdate(ctx, rmpFilter, departmentUpdate).Decode(&finalProfessor)
				if err != nil {
					// this should only happen if the professor is in the final professors collection but we couldnt find them or
					// something is wrong with departmentUpdate

					log.Fatal(err)
				}
				break
				// break out of the loop - no need to add the professor to the final professors collection, as it is already there + updated
			}

			// } else {
			// 	fmt.Println("prof not in final", professor)
			// }

			// if not, add them

			professorsCollection.FindOne(ctx, filter).Decode(&rmpProfessor)
			// better than error checking because the error one is a bit weird,
			// this works the same way as the error checking, but is more accurate
			if rmpProfessor.Name != "" {
				// now create the professor object
				finalProfessor.Name = rmpProfessor.Name
				finalProfessor.RMPName = professor
				finalProfessor.Rating = rmpProfessor.Rating
				finalProfessor.Difficulty = rmpProfessor.Difficulty
				finalProfessor.RMPId = rmpProfessor.RMPId
				finalProfessor.Reviews = make([]uwomodel.Review, len(rmpProfessor.Reviews))

				// review conversion
				// this is better than the loop i had earlier because it is more concise
				for i, review := range rmpProfessor.Reviews {
					finalProfessor.Reviews[i] = uwomodel.Review{
						ProfessorID: review.ProfessorID,
						Professor:   review.Professor,
						Quality:     review.Quality,
						Difficulty:  review.Difficulty,
						Date:        review.Date,
						ReviewText:  review.ReviewText,
						Helpful:     review.Helpful,
						Clarity:     review.Clarity,
					}
				}

				// department addition
				finalProfessor.Departments = []string{course.CourseData.Faculty}

				// courses addition
				finalProfessor.CurrentCourses = []string{courseString}

				finalProfessorsCollection.InsertOne(ctx, finalProfessor)
			} else {
				// otherwise we will create the professor object without the rmpprofessor
				// fmt.Println("professor not found", professor)
				finalProfessor := uwomodel.Professor{
					Name:           professor,
					RMPName:        professor,
					Reviews:        []uwomodel.Review{},
					Departments:    []string{course.CourseData.Faculty},
					CurrentCourses: []string{courseString},
					Rating:         0,
					Difficulty:     0,
				}
				finalProfessorsCollection.InsertOne(ctx, finalProfessor)
			}

			// for whatever reason, the insertion doesnt work if its outside the else
			// finalProfessorsCollection.InsertOne(ctx, finalProfessor)
			// create a new professors collection in the uwo-tt-api database
		}
	}
}
