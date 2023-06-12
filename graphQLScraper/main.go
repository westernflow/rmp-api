package main

import (
	"context"
	"fmt"
	"log"
	model "rmpParser/models"
	"rmpParser/worker"

	"github.com/machinebox/graphql"
)

func main() {
	AUTH_TOKEN := worker.AuthID
	UNI_ID := worker.WesternID

	graphqlClient := graphql.NewClient("https://www.ratemyprofessors.com/graphql")
	// post request
	graphqlRequest := graphql.NewRequest(worker.ProfQuery)
	graphqlRequest.Var("text", "Allan Gedalof")
	graphqlRequest.Var("schoolID", UNI_ID)

	graphqlRequest.Header.Set("Authorization", "Basic "+AUTH_TOKEN)
	graphqlRequest.Header.Set("Content-Type", "application/json")

	var graphqlResponse model.HomePageData

	if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Query successful")
	}

	fmt.Println(graphqlResponse.Data)
}
