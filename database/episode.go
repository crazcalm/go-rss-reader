package database

import (
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed"

	"github.com/crazcalm/html-to-text"
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

//Header -- Returns a string header information important for each episode.
func (e Episode) Header() string {
	var author string
	result := fmt.Sprintf("Feed: %s\n", e.Feed)

	if !strings.EqualFold(e.Data.Title, "") {
		result += fmt.Sprintf("Title: %s\n", e.Data.Title)
	}

	if e.Author != nil {

		if !strings.EqualFold(e.Author.Email, "") {
			author = e.Author.Email

			if !strings.EqualFold(e.Author.Name, "") {
				author = fmt.Sprintf("%s (%s)", author, e.Author.Name)
			}

			result += fmt.Sprintf("Author: %s\n", author)
		}

	} else if e.Data.Author != nil {
		if !strings.EqualFold(e.Data.Author.Email, "") {
			author = e.Data.Author.Email

			if !strings.EqualFold(e.Data.Author.Name, "") {
				author = fmt.Sprintf("%s (%s)", author, e.Data.Author.Name)
			}
			result += fmt.Sprintf("Author: %s\n", author)

		} else if !strings.EqualFold(e.Data.Author.Name, "") {
			author = e.Data.Author.Name
			result += fmt.Sprintf("Author: %s\n", author)
		}
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

	return result
}

//Content -- Formats the body/content of each episode
func (e Episode) Content() (string, []string, error) {
	var data string
	if !strings.EqualFold(e.Data.Description, "") {
		data = e.Data.Description
	} else {
		data = e.Data.Content
	}

	content, links, err := htmltotext.Translate(strings.NewReader(data))
	if err != nil {
		return data, links, nil //fmt.Errorf("Error occurred when parsing raw data: (%s). Returning raw data", err.Error())
	}
	linksFormated := e.links(links)

	return fmt.Sprintf("%s\n%s\n\n%s", e.Header(), strings.TrimSpace(content), linksFormated), links, nil
}

//Links -- A slice of the links presented in the episode
func (e Episode) links(links []string) string {
	var result string
	if len(links) != 0 {
		result += fmt.Sprintf("Links:")
	}

	for index, link := range links {
		result += fmt.Sprintf("\n[%d]: %s", index+1, link)
	}
	return result
}
