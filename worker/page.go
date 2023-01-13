package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

const profQuery = "query TeacherRatingsPageQuery(\n  $id: ID!\n) {\n  node(id: $id) {\n    __typename\n    ... on Teacher {\n      id\n      legacyId\n      firstName\n      lastName\n      school {\n        legacyId\n        name\n        id\n      }\n      lockStatus\n      ...StickyHeader_teacher\n      ...RatingDistributionWrapper_teacher\n      ...TeacherMetaInfo_teacher\n      ...TeacherInfo_teacher\n      ...SimilarProfessors_teacher\n      ...TeacherRatingTabs_teacher\n    }\n    id\n  }\n}\n\nfragment StickyHeader_teacher on Teacher {\n  ...HeaderDescription_teacher\n  ...HeaderRateButton_teacher\n}\n\nfragment RatingDistributionWrapper_teacher on Teacher {\n  ...NoRatingsArea_teacher\n  ratingsDistribution {\n    total\n    ...RatingDistributionChart_ratingsDistribution\n  }\n}\n\nfragment TeacherMetaInfo_teacher on Teacher {\n  legacyId\n  firstName\n  lastName\n  department\n  school {\n    name\n    city\n    state\n    id\n  }\n}\n\nfragment TeacherInfo_teacher on Teacher {\n  id\n  lastName\n  numRatings\n  ...RatingValue_teacher\n  ...NameTitle_teacher\n  ...TeacherTags_teacher\n  ...NameLink_teacher\n  ...TeacherFeedback_teacher\n  ...RateTeacherLink_teacher\n}\n\nfragment SimilarProfessors_teacher on Teacher {\n  department\n  relatedTeachers {\n    legacyId\n    ...SimilarProfessorListItem_teacher\n    id\n  }\n}\n\nfragment TeacherRatingTabs_teacher on Teacher {\n  numRatings\n  courseCodes {\n    courseName\n    courseCount\n  }\n  ...RatingsList_teacher\n  ...RatingsFilter_teacher\n}\n\nfragment RatingsList_teacher on Teacher {\n  id\n  legacyId\n  lastName\n  numRatings\n  school {\n    id\n    legacyId\n    name\n    city\n    state\n    avgRating\n    numRatings\n  }\n  ...Rating_teacher\n  ...NoRatingsArea_teacher\n  ratings(first: 20) {\n    edges {\n      cursor\n      node {\n        ...Rating_rating\n        id\n        __typename\n      }\n    }\n    pageInfo {\n      hasNextPage\n      endCursor\n    }\n  }\n}\n\nfragment RatingsFilter_teacher on Teacher {\n  courseCodes {\n    courseCount\n    courseName\n  }\n}\n\nfragment Rating_teacher on Teacher {\n  ...RatingFooter_teacher\n  ...RatingSuperHeader_teacher\n  ...ProfessorNoteSection_teacher\n}\n\nfragment NoRatingsArea_teacher on Teacher {\n  lastName\n  ...RateTeacherLink_teacher\n}\n\nfragment Rating_rating on Rating {\n  comment\n  flagStatus\n  createdByUser\n  teacherNote {\n    id\n  }\n  ...RatingHeader_rating\n  ...RatingSuperHeader_rating\n  ...RatingValues_rating\n  ...CourseMeta_rating\n  ...RatingTags_rating\n  ...RatingFooter_rating\n  ...ProfessorNoteSection_rating\n}\n\nfragment RatingHeader_rating on Rating {\n  date\n  class\n  helpfulRating\n  clarityRating\n  isForOnlineClass\n}\n\nfragment RatingSuperHeader_rating on Rating {\n  legacyId\n}\n\nfragment RatingValues_rating on Rating {\n  helpfulRating\n  clarityRating\n  difficultyRating\n}\n\nfragment CourseMeta_rating on Rating {\n  attendanceMandatory\n  wouldTakeAgain\n  grade\n  textbookUse\n  isForOnlineClass\n  isForCredit\n}\n\nfragment RatingTags_rating on Rating {\n  ratingTags\n}\n\nfragment RatingFooter_rating on Rating {\n  id\n  comment\n  adminReviewedAt\n  flagStatus\n  legacyId\n  thumbsUpTotal\n  thumbsDownTotal\n  thumbs {\n    userId\n    thumbsUp\n    thumbsDown\n    id\n  }\n  teacherNote {\n    id\n  }\n}\n\nfragment ProfessorNoteSection_rating on Rating {\n  teacherNote {\n    ...ProfessorNote_note\n    id\n  }\n  ...ProfessorNoteEditor_rating\n}\n\nfragment ProfessorNote_note on TeacherNotes {\n  comment\n  ...ProfessorNoteHeader_note\n  ...ProfessorNoteFooter_note\n}\n\nfragment ProfessorNoteEditor_rating on Rating {\n  id\n  legacyId\n  class\n  teacherNote {\n    id\n    teacherId\n    comment\n  }\n}\n\nfragment ProfessorNoteHeader_note on TeacherNotes {\n  createdAt\n  updatedAt\n}\n\nfragment ProfessorNoteFooter_note on TeacherNotes {\n  legacyId\n  flagStatus\n}\n\nfragment RateTeacherLink_teacher on Teacher {\n  legacyId\n  numRatings\n  lockStatus\n}\n\nfragment RatingFooter_teacher on Teacher {\n  id\n  legacyId\n  lockStatus\n  isProfCurrentUser\n}\n\nfragment RatingSuperHeader_teacher on Teacher {\n  firstName\n  lastName\n  legacyId\n  school {\n    name\n    id\n  }\n}\n\nfragment ProfessorNoteSection_teacher on Teacher {\n  ...ProfessorNote_teacher\n  ...ProfessorNoteEditor_teacher\n}\n\nfragment ProfessorNote_teacher on Teacher {\n  ...ProfessorNoteHeader_teacher\n  ...ProfessorNoteFooter_teacher\n}\n\nfragment ProfessorNoteEditor_teacher on Teacher {\n  id\n}\n\nfragment ProfessorNoteHeader_teacher on Teacher {\n  lastName\n}\n\nfragment ProfessorNoteFooter_teacher on Teacher {\n  legacyId\n  isProfCurrentUser\n}\n\nfragment SimilarProfessorListItem_teacher on RelatedTeacher {\n  legacyId\n  firstName\n  lastName\n  avgRating\n}\n\nfragment RatingValue_teacher on Teacher {\n  avgRating\n  numRatings\n  ...NumRatingsLink_teacher\n}\n\nfragment NameTitle_teacher on Teacher {\n  id\n  firstName\n  lastName\n  department\n  school {\n    legacyId\n    name\n    id\n  }\n  ...TeacherDepartment_teacher\n  ...TeacherBookmark_teacher\n}\n\nfragment TeacherTags_teacher on Teacher {\n  lastName\n  teacherRatingTags {\n    legacyId\n    tagCount\n    tagName\n    id\n  }\n}\n\nfragment NameLink_teacher on Teacher {\n  isProfCurrentUser\n  id\n  legacyId\n  firstName\n  lastName\n  school {\n    name\n    id\n  }\n}\n\nfragment TeacherFeedback_teacher on Teacher {\n  numRatings\n  avgDifficulty\n  wouldTakeAgainPercent\n}\n\nfragment TeacherDepartment_teacher on Teacher {\n  department\n  departmentId\n  school {\n    legacyId\n    name\n    id\n  }\n}\n\nfragment TeacherBookmark_teacher on Teacher {\n  id\n  isSaved\n}\n\nfragment NumRatingsLink_teacher on Teacher {\n  numRatings\n  ...RateTeacherLink_teacher\n}\n\nfragment RatingDistributionChart_ratingsDistribution on ratingsDistribution {\n  r1\n  r2\n  r3\n  r4\n  r5\n}\n\nfragment HeaderDescription_teacher on Teacher {\n  id\n  firstName\n  lastName\n  department\n  school {\n    legacyId\n    name\n    city\n    state\n    id\n  }\n  ...TeacherTitles_teacher\n  ...TeacherBookmark_teacher\n}\n\nfragment HeaderRateButton_teacher on Teacher {\n  ...RateTeacherLink_teacher\n}\n\nfragment TeacherTitles_teacher on Teacher {\n  department\n  school {\n    legacyId\n    name\n    id\n  }\n}\n"

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





// this URL  defines the context for the page to be scraped
const westernProfessorListURL = "https://www.ratemyprofessors.com/search/teachers?query=*&sid=1491"

// Selectors for goquery to select the specific nodes based on CSS classes
const teacherCardSelector = "a.TeacherCard__StyledTeacherCard-syjs0d-0.dLJIlx"
const cardNameSelector = "div.CardName__StyledCardName-sc-1gyrgim-0.cJdVEK"
const cardDifficultySelector = "div.CardFeedback__CardFeedbackNumber-lq6nix-2.hroXqf"
const cardDepartmentSelector = "div.CardSchool__Department-sc-19lmz2k-0.haUIRO"
const cardRatingSectionSelector = "div.CardNumRating__StyledCardNumRating-sc-17t4b9u-0.eWZmyX"
const buttonSelector = "gjQZal"

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
	profId, err := strconv.Atoi(strings.Split(hrefLink,"=")[1]); if err != nil {
		fmt.Println("Error parsing profId")
	}
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