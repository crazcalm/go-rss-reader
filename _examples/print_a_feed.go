package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
)

const (
	alertsXML = "https://www.us-cert.gov/ncas/alerts.xml"
	nprXML    = "https://www.npr.org/rss/rss.php?id=1001"
)

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(alertsXML)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("num of items = %d\n", len(feed.Items))
	fmt.Println(feed)

	fmt.Printf("\n\nDescription:\n%s", feed.Items[0].Description)
}
