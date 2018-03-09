package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

const (
	alertsXML = "https://www.us-cert.gov/ncas/alerts.xml"
	nprXML    = "https://www.npr.org/rss/rss.php?id=1001"
)

func removeSpaces(s string) string {
	return strings.TrimSpace(s)
}

func processTags(s string) string {
	stringBytes := []byte(s)

	if len(s) > 2 {
		//Check for tags
		if bytes.HasPrefix(stringBytes, []byte("<")) && bytes.HasSuffix(stringBytes, []byte(">")) {
			//Tag list
			if bytes.HasPrefix(stringBytes, []byte("<p>")) {
				return "\n"
			}else if bytes.HasPrefix(stringBytes, []byte("</p>")){
				return "\n"
			} else if bytes.HasPrefix(stringBytes, []byte("<h3>")) {
				return "\n"
			} else if bytes.HasPrefix(stringBytes, []byte("</h3>")) {
				return "\n"
			}  else if bytes.HasPrefix(stringBytes, []byte("<ul>")) {
				return "\n"
			} else if bytes.HasPrefix(stringBytes, []byte("</ul>")) {
				return "\n"
			} else if bytes.HasPrefix(stringBytes, []byte("<a")) {
				return " "
			} else if bytes.HasPrefix(stringBytes, []byte("</a>")) {
				return " "
			} else if bytes.HasPrefix(stringBytes, []byte("<li>")) {
				return fmt.Sprintf("    * ")
			} else if bytes.HasPrefix(stringBytes, []byte("</li>")) {
				return "\n"
			}  else if bytes.HasPrefix(stringBytes, []byte("<strong>")) {
				return ""
			} else if bytes.HasPrefix(stringBytes, []byte("</strong>")) {
				return ""
			} else if bytes.HasPrefix(stringBytes, []byte("<br/>")) {
				return "\n"
			}
		}
		
	}
	//fmt.Printf("len of string %d", len(s))
	return fmt.Sprintf("%s", s)
}

func printToScreen(s string){
	if len(s) == 0 {
		return
	}
	fmt.Printf("%s", s)
}

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(alertsXML)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("num of items = %d\n", len(feed.Items))
	//fmt.Println(feed)

	//fmt.Println(feed.Items[0].Description)

	//Create a new reader
	reader := strings.NewReader(feed.Items[0].Description)

	z := html.NewTokenizer(reader)

for {                
	tokenType := z.Next()   
        if tokenType == html.ErrorToken {                
               	//Error case
                log.Fatal("html Parser err token: %d", html.ErrorToken)                        
        }                
        // Process the current token.
        token := z.Token()

		noSpaces := removeSpaces(token.String())

		noTags := processTags(noSpaces)
        printToScreen(noTags)
    }
}
