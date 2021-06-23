package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"
)

const layout = "1/02/2006 3:4 PM"

var wantedAuthors map[string]bool = map[string]bool{
	"SaulR80683": true,
}

func main() {
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("boards.fool.com"),
	)

    // Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// #tblMessagesAsp > tbody > tr:nth-child(4) > td.first > span > a
	c.OnHTML("#tblMessagesAsp > tbody > tr", func(e *colly.HTMLElement) {
		temp := e.ChildTexts("td")
		if len(temp) == 0 {
			return
		}
		postLink := e.ChildAttr("td.first > span > a[href]", "href")

		// subject, author, recs, date, numbers := temp[0], temp[1], temp[2], temp[3], temp[4]
		_, author, recs, date, _ := temp[0], temp[1], temp[2], temp[3], temp[4]

		_, err := time.Parse(layout, date)
		if err != nil {
			panic(err)
		}

		var ok bool
		_, ok = wantedAuthors[author]
		if recs == "--" && !ok {
			return
		}
		_recs, err := strconv.Atoi(recs)
		if err != nil {
			panic(err)
		}
		if !ok && _recs < 70 {
			return
		}

		e.Request.Visit(postLink)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://boards.fool.com/sauls-investing-discussions-120980.aspx")

	// TODO send notification to my telegram
	// TODO parse post content
}
