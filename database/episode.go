package database

import (
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed"
	//	"github.com/crazcalm/html-to-text"
)

//Episode -- Data structure used to handle each new episode of a feed
type Episode struct {
	Feed   string
	Author *gofeed.Person
	Data   *gofeed.Item
}

//FormatEpisodeHeader -- formats the episode header
func FormatEpisodeHeader(feedName, episodeTitle, author, episodeLink, episodeDate, mediaContent string) (result string) {
	if !strings.EqualFold(feedName, "") {
		result += fmt.Sprintf("Feed: %s\n", feedName)
	}

	if !strings.EqualFold(episodeTitle, "") {
		result += fmt.Sprintf("Title: %s\n", episodeTitle)
	}

	if !strings.EqualFold(author, "") {
		result += fmt.Sprintf("Author: %s\n", author)
	}

	if !strings.EqualFold(episodeLink, "") {
		result += fmt.Sprintf("Link: %s\n", episodeLink)
	}

	if !strings.EqualFold(episodeDate, "") {
		result += fmt.Sprintf("Date: %s\n", episodeDate)
	}

	if !strings.EqualFold(mediaContent, "") {
		result += fmt.Sprintf("Media Link: %s\n", mediaContent)
	}

	return
}
