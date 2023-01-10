// import (
//   //  "encoding/csv"
//    "fmt"
//   //  "log"
//   //  "os"

//    "github.com/gocolly/colly"
// )
// func ex() {
// 	 c := colly.NewCollector(
//    colly.AllowedDomains("books.toscrape.com"),
// 	)

// 	c.OnRequest(func(r *colly.Request) {
//    fmt.Println("Visiting", r.URL)
// })

// c.OnResponse(func(r *colly.Response) {
// 	fmt.Println(r.StatusCode)
// })

// c.OnHTML("title", func(e *colly.HTMLElement) {
// 	fmt.Println(e.Text)
// })

// c.OnHTML(".product_pod", func(e *colly.HTMLElement) {
// 	title := e.ChildAttr(".image_container img", "alt")
// 	price := e.ChildText(".price_color")
// 	fmt.Println(title, price)
// })

// c.Visit("https://books.toscrape.com/")


//    fmt.Println("Done")
// }