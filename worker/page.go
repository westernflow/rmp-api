package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	// "net/http"

	"rmpParser/models"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/chromedp/chromedp"
	// "github.com/graphql-go/graphql"

	"context"
	"strings"
)

// PageScraper defines the context for the page to be scraped and the location of scrape resulst
type PageScraper struct {
	Header string
	URL    string
	Status string
	DB     *mongo.Database // switch to dynamo?
	Form   *goquery.Selection
}

// PageResult encompasses data that is passed into channel to be parsed
type PageResult struct {
	Name string
	Doc  *goquery.Document
}

func buildProfessor(node model.ProfessorData) model.Professor {
	var professor model.Professor
	professor.Name = node.FirstName + " " + node.LastName
	professor.RMPId = node.ID
	professor.Rating = node.AvgRating
	professor.Difficulty = node.AvgDifficulty
	professor.Department = node.Department

	// get the reviews from the professor data
	var reviews []model.Review
	for _, edge := range node.Ratings.Edges {
		var review model.Review

		// attempt to parse edge.Node.Class into the course struct
		// if it does not match the following regexp ^[a-zA-z][a-zA-z]+[0-9][0-9][0-9]$ then it is not a course and should be ignored
		
		// first remove spaces
		edge.Node.Class = strings.ReplaceAll(edge.Node.Class, " ", "")
		// then check if it matches the regexp
		if !regexp.MustCompile(`^[a-zA-z][a-zA-z]+[0-9][0-9][0-9]$`).MatchString(edge.Node.Class) {
			continue
		}
		// if it does match, then parse it into the course struct
		// first get course dept -- it is the characters before the first number
		re := regexp.MustCompile(`[0-9]`)
		index := re.FindStringIndex(edge.Node.Class)[0]
		review.Course.Department = edge.Node.Class[:index]
		// then get the course number -- it is the characters after the first number
		review.Course.Number = edge.Node.Class[index:]
		review.Professor = professor.Name
		review.Quality = edge.Node.HelpfulRating
		review.Difficulty = edge.Node.ClarityRating
		review.Date = edge.Node.Date
		review.ReviewText = edge.Node.Comment
		review.Helpful = edge.Node.HelpfulRating
		review.Clarity = edge.Node.ClarityRating
		reviews = append(reviews, review)
	}
	professor.Reviews = reviews
	return professor
}

func GetProfessorData(id string) (professor model.Professor, err error) {
	variables := make(map[string]interface{})
	variables["id"] = id
	request := model.Request{Query: profQuery, Variables: variables}
	// send the graphql request
	// convert request to string
	requestString, err := json.Marshal(request); if err != nil {
		fmt.Println("Error converting request to string:", err)
	}
	
	req, _ := http.NewRequest("POST", "https://www.ratemyprofessors.com/graphql", bytes.NewBuffer(requestString))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic dGVzdDp0ZXN0")
	client := &http.Client{}
	resp, err := client.Do(req); if err != nil {
		fmt.Println("Error sending request:", err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

	// parse the response
	var response model.Response
	err = json.Unmarshal(body, &response); if err != nil {
		fmt.Println("Error parsing response:", err)
	}
	
	// get the professor data from the response
	var newprof = buildProfessor(response.Data.ProfessorData)
	return newprof, err
}

func PopulateDB() {
	variables := make(map[string]interface{})
	variables["id"] = "VGVhY2hlci03OTIy"
	request := model.Request{Query: profQuery, Variables: variables}
	// send the graphql request
	// convert request to string
	requestString, err := json.Marshal(request); if err != nil {
		fmt.Println("Error converting request to string:", err)
	}
	
	req, _ := http.NewRequest("POST", "https://www.ratemyprofessors.com/graphql", bytes.NewBuffer(requestString))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic dGVzdDp0ZXN0")
	client := &http.Client{}
	resp, err := client.Do(req); if err != nil {
		fmt.Println("Error sending request:", err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

// FetchDocument fetches contents of page based on URL
func (scraper *PageScraper) FetchDocument() (document *goquery.Document, err error) {
	// create a new chrome instance
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// navigate to the page
	fmt.Println("Initiating chromedp instance to page")
	var body string
	// use chromedp to navigate to the western professor list page, wait for the teacher card selector to be visible, keep clicking button until it is not visible, and then get the html of the page 
	err = chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(westernProfessorListURL),
		chromedp.WaitVisible(teacherCardSelector),
		// chromedp.Click(buttonSelector, chromedp.NodeVisible),
		// chromedp.WaitNotVisible(buttonSelector),
		chromedp.OuterHTML("html", &body),
	}); if err != nil {
		fmt.Println("Error running chromedp instance:", err)
	}

	fmt.Println("Chromedp instance complete...")

	// Create a Goquery document from the response
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		fmt.Println(err)
	}

	return doc, nil
}

func scrapeProfessorData(s *goquery.Selection) (professor model.Professor) {
	// get the name of the professor
	name := s.Find(cardNameSelector).Text()
		
	// get the difficulty of the professor
	enjoymentAndDifficulty := s.Find(cardDifficultySelector).Text()
	_, difficultyString, err := parseEnjoymentAndDifficulty(enjoymentAndDifficulty); if err != nil {
		fmt.Println("Error parsing enjoyment and difficulty:", err)
	}

	difficulty, err := strconv.ParseFloat(difficultyString, 64); if err != nil {
		fmt.Println("Error parsing difficulty")
	}

	// get the department of the professor
	department := s.Find(cardDepartmentSelector).Text()

	// get the quality card section of the professor
	ratingSection := s.Find(cardRatingSectionSelector)

	// get the quality of the professor from the second div in the selected node
	ratingString := ratingSection.Find("div").Eq(1).Text()

	rating, err := strconv.ParseFloat(ratingString, 64); if err != nil {
		fmt.Println("Error parsing quality")
	}

	// get the ratemyprofessor id from the href of the selected node
	hrefLink, _ := s.Attr("href")

	// parse profId to int
	profId := strings.Split(hrefLink,"=")[1]
	return model.Professor{Name: name, Difficulty: difficulty, Department: department, Rating: rating, RMPId: profId}
}

func (scraper *PageScraper) scrapeProfessors(doc *goquery.Document) []model.Professor {
	// create a slice of professors
	var professors []model.Professor

	// use goquery to select the node with this class: "SearchResultsPage__SearchResultsWrapper-vhbycj-1 gxbBpy" and then select the html div node with no class
	doc.Find(teacherCardSelector).Each(func(i int, s *goquery.Selection) {
		// scrape professor data
		professor := scrapeProfessorData(s)
		professors = append(professors, professor)
	})

	return professors
}