package worker

import (
	"fmt"
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

// FetchDocument fetches contents of page based on URL
func (scraper *PageScraper) FetchDocument() (document *goquery.Document, err error) {
	// create a new chrome instance
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// navigate to the page
	var body string
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.ratemyprofessors.com/search/teachers?query=*&sid=1491"),
		chromedp.WaitVisible(".SearchResultsPage__SearchResultsWrapper-vhbycj-1"),
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
	fmt.Println("still Scraping for professors")
	// TODO: make these class names consts
	doc.Find(".SearchResultsPage__SearchResultsWrapper-vhbycj-1 gxbBpy").Each(func(i int, s *goquery.Selection) {
		fmt.Println("hi")
		// print the name of each professor
		name := s.Find("div.CardName__StyledCardName-sc-1gyrgim-0 cJdVEK").Text()
		fmt.Println(name)
		professors = append(professors, model.Professor{Name: name})
	})

	return professors
}