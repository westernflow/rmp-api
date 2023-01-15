package worker

import (
	"fmt"
)

// Create a worker that will scrape the data of each professor on this page `https://www.ratemyprofessors.com/school?sid=1491`
func Scrape () []Professor {
	// create a new page scraper with url set to "https://www.ratemyprofessors.com/school?sid=1491"
	
	scraper := PageScraper{
		URL: "https://www.ratemyprofessors.com/search/teachers?query=*&sid=1491",
	}
	
	// fetch the document from the page scraper
	doc, err := scraper.FetchDocument(); if err != nil {
		fmt.Println("Error fetching document")
	}

	// professors on the page
	fmt.Println("Scraping for professors")
	professors := scraper.scrapeProfessors(doc)

	fmt.Println(professors)
	return professors
}