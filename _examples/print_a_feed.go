package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
)

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://www.us-cert.gov/ncas/alerts.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("num of items = %d\n", len(feed.Items))
	fmt.Println(feed)
}
