package worker

import (
	"fmt"
	"strconv"
	// "net/http"

	"rmpParser/models"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/chromedp/chromedp"

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

// this URL  defines the context for the page to be scraped
const westernProfessorListURL = "https://www.ratemyprofessors.com/search/teachers?query=*&sid=1491"

// Selectors for goquery to select the specific nodes based on CSS classes
const teacherCardSelector = "a.TeacherCard__StyledTeacherCard-syjs0d-0.dLJIlx"
const cardNameSelector = "div.CardName__StyledCardName-sc-1gyrgim-0.cJdVEK"
const cardDifficultySelector = "div.CardFeedback__CardFeedbackNumber-lq6nix-2.hroXqf"
const cardDepartmentSelector = "div.CardSchool__Department-sc-19lmz2k-0.haUIRO"
const cardRatingSectionSelector = "div.CardNumRating__StyledCardNumRating-sc-17t4b9u-0.eWZmyX"

// FetchDocument fetches contents of page based on URL
func (scraper *PageScraper) FetchDocument() (document *goquery.Document, err error) {
	// create a new chrome instance
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// navigate to the page
	var body string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(westernProfessorListURL),
		chromedp.WaitVisible(teacherCardSelector),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		fmt.Println(err)
	}

	// Create a Goquery document from the response
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		fmt.Println(err)
	}

	return doc, nil
}

func (scraper *PageScraper) scrapeProfessors(doc *goquery.Document) []model.Professor {
	// create a slice of professors
	var professors []model.Professor

	// use goquery to select the node with this class: "SearchResultsPage__SearchResultsWrapper-vhbycj-1 gxbBpy" and then select the html div node with no class
	doc.Find(teacherCardSelector).Each(func(i int, s *goquery.Selection) {
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

		fmt.Println(name, difficulty, department, rating)
		professors = append(professors, model.Professor{Name: name, Difficulty: difficulty, Department: department, Rating: rating})
	})

	return professors
}