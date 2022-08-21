package handlers

import (
	"github.com/gocolly/colly/v2"
	"log"
)

type Pollen struct {
	Header string
	Date   string
	Count  string
}

// ScrapePollenCount scrapes the canberrapollen website for its daily pollen
// reading.
func ScrapePollenCount() *Pollen {
	var p Pollen
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("response from", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		log.Println("error:", e)
	})
	c.OnHTML("#bar-pollen-div", func(e *colly.HTMLElement) {
		p.Header = e.ChildText("div.pollen-header")
		p.Date = e.ChildText("div.pollen-date")
		p.Count = e.ChildText("#plevel")
	})
	c.Visit("https://canberrapollen.com.au/")
	return &p
}
