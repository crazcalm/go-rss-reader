package rss

import (
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed"
)

//Episode -- Data structure used to handle each new episode of a feed
type Episode struct {
	Feed   string
	Author *gofeed.Person
	Data   *gofeed.Item
}

//Header -- Returns a string header information important for each episode.
func (e Episode) Header() string {
	result := fmt.Sprintf("Feed: %s\n", e.Feed)

	if !strings.EqualFold(e.Data.Title, "") {
		result += fmt.Sprintf("Title: %s\n", e.Data.Title)
	}

	if !strings.EqualFold(e.Data.Author.Email, "") {
		author := e.Data.Author.Email

		if !strings.EqualFold(e.Data.Author.Name, "") {
			author = fmt.Sprintf("%s (%s)", author, e.Data.Author.Name)
		}
		result += fmt.Sprintf("Author: %s\n", author)
	} else if !strings.EqualFold(e.Author.Email, "") {
		author := e.Author.Email

		if !strings.EqualFold(e.Author.Name, "") {
			author = fmt.Sprintf("%s (%s)", author, e.Author.Name)
		}

		result += fmt.Sprintf("Author: %s\n", author)

	}

	if !strings.EqualFold(e.Data.Link, "") {
		result += fmt.Sprintf("Link: %s\n", e.Data.Link)
	}

	if !strings.EqualFold(e.Data.Published, "") {
		result += fmt.Sprintf("Date: %s\n", e.Data.Published)
	}

	media, ok := e.Data.Extensions["media"]
	if ok {
		content, ok := media["content"]

		if ok {
			for i := 0; i < len(content); i++ {
				var mediaContent string

				url, ok := content[i].Attrs["url"]
				if ok {
					mediaContent += url

					itemType, ok := content[i].Attrs["type"]
					if ok {
						mediaContent = fmt.Sprintf("%s (type: %s)", mediaContent, itemType)
						result += fmt.Sprintf("Podcast Download URL: %s\n", mediaContent)
					}
				}

			}
		}
	}

	/*
		itemType, ok := content[0].Attrs["type"]
		if !ok {
			log.Fatalf("Content has no attribute type")
		}

		url, ok := content[0].Attrs["url"]
		if !ok {
			log.Fatalf("Content has no attribute url")
	*/

	return result
}

//Content -- Formats the body/content of each episode
func (e Episode) Content() string {
	return ""
}

//Links -- A slice of the links presented in the episode
func (e Episode) Links() []string {
	return []string{}
}
