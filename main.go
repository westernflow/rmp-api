package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/valyala/fastjson"

	"rmpParser/controller"
	"rmpParser/handler"
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ApiResponse := events.APIGatewayProxyResponse{}
	// Switch for identifying the HTTP request
	switch request.HTTPMethod {
	case "GET":
		// Obtain the QueryStringParameter
		name := request.QueryStringParameters["name"]
		if name != "" {
			ApiResponse = events.APIGatewayProxyResponse{Body: "Hey " + name + " welcome! ", StatusCode: 200}

		} else {
			ApiResponse = events.APIGatewayProxyResponse{Body: "Error: Query Parameter name missing", StatusCode: 500}
		}

	case "POST":
		// validates json and returns error if not working
		err := fastjson.Validate(request.Body)

		if err != nil {
			body := "Error: Invalid JSON payload ||| " + fmt.Sprint(err) + " Body Obtained" + "||||" + request.Body
			ApiResponse = events.APIGatewayProxyResponse{Body: body, StatusCode: 500}
		} else {
			ApiResponse = events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}
		}

	}
	// Response
	return ApiResponse, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	// lambda.Start(HandleRequest)
	getData()

	// this should show all the professors in the database
	fmt.Println(handler.GetProfessors())
}

func getData() {
	// create a controller
	c := controller.GetInstance()
	// connect to the database
	c.ConnectToDatabase()

	// populate the database (dont do this if you already have data in the database)
	// c.PopulateDatabase()
}
