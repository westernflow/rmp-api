package main

import (
    "encoding/json"
    // "fmt"
    "log"
    "net/http"
    "os"
    "regexp"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

var isbnRegexp = regexp.MustCompile(`[0-9]{3}\-[0-9]{10}`)
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

type book struct {
    ISBN   string `json:"isbn"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Get the `isbn` query string parameter from the request and
    // validate it.
    isbn := req.QueryStringParameters["isbn"]
    if !isbnRegexp.MatchString(isbn) {
        return clientError(http.StatusBadRequest)
    }

    // Fetch the book record from the database based on the isbn value.
    bk, err := getItem(isbn)
    if err != nil {
        return serverError(err)
    }
    if bk == nil {
        return clientError(http.StatusNotFound)
    }

    // The APIGatewayProxyResponse.Body field needs to be a string, so
    // we marshal the book record into JSON.
    js, err := json.Marshal(bk)
    if err != nil {
        return serverError(err)
    }

    // Return a response with a 200 OK status and the JSON book record
    // as the body.
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(js),
    }, nil
}

// Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
    errorLogger.Println(err.Error())

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusInternalServerError,
        Body:       http.StatusText(http.StatusInternalServerError),
    }, nil
}

// Similarly add a helper for send responses relating to client errors.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: status,
        Body:       http.StatusText(status),
    }, nil
}

func main() {
    lambda.Start(show)
}

// package main

// import (
// 	// "fmt"
// 	"rmpParser/controller"
// 	"rmpParser/worker"

// 	"github.com/robfig/cron"

// 	// aws lambda
// 	"github.com/aws/aws-lambda-go/lambda"
// )

// func getData() {
// 	// create a controller
// 	c := controller.GetInstance()
// 	// connect to the database
// 	c.ConnectToDatabase()
// 	// fmt.Println(c.GetProfessors())
// 	// get all departments from the school
// 	departments := worker.GetDepartments()
// 	// fmt.Println(departments)
// 	// get all professors from each department
// 	for _, department := range departments {
// 		c.InsertDepartment(department)
// 		worker.AddProfessorsFromDepartmentToDatabase(c, department.DepartmentBase64Code)
// 	}
// }


// // aws lambda to query for professors and departments in the database
// func main() {
// 	// given http request, return professors and departments
// 	lambda.Start(getData)

// 	// create a cron job that will run getData() every 24 hours
// 	c := cron.New()
// 	c.AddFunc("@every 24h", getData)
// 	c.Start()
// }
