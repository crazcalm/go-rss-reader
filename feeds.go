package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
)

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("http://www.leoville.tv/podcasts/sn.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(feed)
}
