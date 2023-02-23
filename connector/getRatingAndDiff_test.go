package main

import (
	"context"
	"log"
	"os"

	// "os"
	uwomodel "rmpParser/uwomodel"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetRatingAndDiffNoReviews(t *testing.T) {
	// Requires the MongoDB Go Driver
	// https://go.mongodb.org/mongo-driver
	ctx := context.TODO()

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

	professor := uwomodel.Professor{}

	professorsCollection := client.Database("test").Collection("final_professors")
	cursor, err := professorsCollection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	for cursor.Next(ctx) {
		err := cursor.Decode(&professor)
		if err != nil {
			log.Fatal(err)
		}
		// sum all ratings and difficulties and average them
		// check if no reviews
		expectedDiff := 0.0
		expectedRating := 0.0

		if len(professor.Reviews) == 0 {
			return
		}

		for _, review := range professor.Reviews {
			expectedRating += review.Quality
			expectedDiff += review.Difficulty
		}

		expectedRating /= float64(len(professor.Reviews))
		expectedDiff /= float64(len(professor.Reviews))

		updateRatingAndDiff(&professor)

		newRating := professor.Rating
		newDiff := professor.Difficulty

		if newRating != expectedRating {
			t.Errorf("Expected rating to be %f, got %f", expectedRating, newRating)
		}
		if newDiff != expectedDiff {
			t.Errorf("Expected difficulty to be %f, got %f", expectedDiff, newDiff)
		}
	}
}
