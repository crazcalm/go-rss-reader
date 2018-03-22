package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/crazcalm/html-to-text"
	"github.com/mmcdole/gofeed"
)

const (
	alertsXML = "https://www.us-cert.gov/ncas/alerts.xml"
	nprXML    = "https://www.npr.org/rss/rss.php?id=1001"
	snXML     = "http://www.leoville.tv/podcasts/sn.xml"
	goXML     = "https://golangweekly.com/rss/1h27nlio"
	linuxXML  = "https://latenightlinux.com/feed/mp3"
	changeXML = "https://changelog.com/gotime/feed"
	xkcdXML   = "https://www.xkcd.com/rss.xml"
	irlXML    = "https://feeds.mozilla-podcasts.org/irl"
)

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(alertsXML)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("num of items = %d\n", len(feed.Items))
	fmt.Println(feed)

	fmt.Printf("\n\nDescription:\n%s", feed.Items[5].Description)
	fmt.Printf("\n\nContent:\n%s", feed.Items[5].Content)

	//fmt.Println(feed.Items[1].Description)

	result, links, err := htmltotext.Translate(strings.NewReader(feed.Items[5].Content))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Results:\n\n%s\n\nLinks:\n\n%s\n\nNum of links: %d\n\n\n\n\n", result, links, len(links))

	media, ok := feed.Items[5].Extensions["media"]
	if !ok {
		log.Fatalf("No media exists")
	}

	content, ok := media["content"]
	if !ok {
		log.Fatalf("Media has not content")
	}

	itemType, ok := content[0].Attrs["type"]
	if !ok {
		log.Fatalf("Content has no attribute type")
	}

	url, ok := content[0].Attrs["url"]
	if !ok {
		log.Fatalf("Content has no attribute url")
	}

	fmt.Println(itemType)
	fmt.Println(url)

}
