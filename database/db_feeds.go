package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mmcdole/gofeed"
)

//GetFeedInfo -- Pulls the rss feed from the website and dumps the needed info into the database
func GetFeedInfo(db *sql.DB, feedID int64) (err error) {
	var feedData *gofeed.Feed

	//Get the URl
	url, err := GetFeedURL(db, feedID)
	if err != nil {
		return
	}

	rawData, err := GetFeedDataFromSite(url)
	if err != nil {
		return
	}

	//Add Raw Data to DB
	err = UpdateFeedRawData(db, feedID, rawData)
	if err != nil {
		return
	}

	if !strings.EqualFold(rawData, "") {
		//Need to convert the data to a gofeed object
		feedParser := gofeed.NewParser()
		feedData, err = feedParser.Parse(strings.NewReader(rawData))
		if err != nil {
			return fmt.Errorf("gofeed parser was unable to parse data: %s -- %s", rawData, err.Error())
		}
	}

	//Add Title
	if !strings.EqualFold(feedData.Title, "") {
		err = UpdateFeedTitle(db, feedID, feedData.Title)
	}

	//Add author
	if feedData.Author != nil {
		//If both the name and the email is not blank
		if !strings.EqualFold(feedData.Author.Email, "") || !strings.EqualFold(feedData.Author.Name, "") {
			var authorID int64
			if !AuthorExist(db, feedData.Author.Name, feedData.Author.Email) {
				authorID, err = AddAuthor(db, feedData.Author.Name, feedData.Author.Email)
				if err != nil {
					return err
				}

			} else {
				authorID, err = GetAuthorByNameAndEmail(db, feedData.Author.Name, feedData.Author.Email)
				if err != nil {
					return err
				}
			}
			//Updating the feed author
			err = UpdateFeedAuthor(db, feedID, authorID)
			if err != nil {
				return err
			}
		}
	}

	//Add Episodes
	for _, episode := range feedData.Items {
		var rssHTML string
		if len(episode.Description) > len(episode.Content) {
			rssHTML = episode.Description
		} else {
			rssHTML = episode.Content
		}

		if EpisodeExist(db, episode.Title) {
			//TODO: need to check if this works...
			continue
			//Continue should skipp to the next loop interations
		}

		episodeID, err := AddEpisode(db, feedID, episode.Link, episode.Title, episode.PublishedParsed, rssHTML)
		if err != nil {
			return err
		}

		//Add media content
		media, ok := episode.Extensions["media"]
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

							err = UpdateEpisodeMediaContent(db, episodeID, mediaContent)
							if err != nil {
								return err
							}
						}
					}

				}
			}
		}

		//Add author
		if episode.Author != nil {
			//If both the name and the email is not blank
			if !strings.EqualFold(episode.Author.Email, "") || !strings.EqualFold(episode.Author.Name, "") {
				var authorID int64
				if !AuthorExist(db, episode.Author.Name, episode.Author.Email) {
					authorID, err = AddAuthor(db, episode.Author.Name, episode.Author.Email)
					if err != nil {
						return err
					}

				} else {
					authorID, err = GetAuthorByNameAndEmail(db, episode.Author.Name, episode.Author.Email)
					if err != nil {
						return err
					}
				}

				//Updating the episode author
				err = UpdateEpisodeAuthor(db, episodeID, authorID)
				if err != nil {
					return err
				}
			}
		}
	}
	return
}

//LoadFeed -- Loads a feed from the database
func LoadFeed(db *sql.DB, id int64) (feed *Feed, err error) {
	var feedData *gofeed.Feed

	url, err := GetFeedURL(db, id)
	if err != nil {
		log.Fatal(err)
	}

	title, err := GetFeedTitle(db, id)
	if err != nil {
		log.Fatal(err)
	}
	if strings.EqualFold(title, "") {
		title = url
	}

	data, err := GetFeedRawData(db, id)
	if err != nil {
		return feed, fmt.Errorf("No data to retrieve: %s", err.Error())
	}

	if !strings.EqualFold(data, "") {
		//Need to convert the data to a gofeed object
		feedParser := gofeed.NewParser()
		feedData, err = feedParser.Parse(strings.NewReader(data))
		if err != nil {
			return feed, fmt.Errorf("gofeed parser was unable to parse data: %s -- %s", data, err.Error())
		}
	}

	var tags []string
	activeTags := AllActiveFeedTags(db, id)

	for _, tag := range activeTags {
		tags = append(tags, tag)
	}

	return &Feed{id, url, title, tags, feedData}, nil
}

//GetFeedAuthor -- returns the feed author
func GetFeedAuthor(db *sql.DB, feedID int64) (name, email string, err error) {
	stmt := "SELECT authors.name, authors.email FROM feeds INNER JOIN authors ON authors.id = feeds.author_id WHERE feeds.id = $1"
	row := db.QueryRow(stmt, feedID)
	err = row.Scan(&name, &email)
	return
}

//FeedHasAuthor -- returns true is an author id exists and false otherwise
func FeedHasAuthor(db *sql.DB, feedID int64) (result bool) {
	var count int64
	row := db.QueryRow("SELECT COUNT(author_id) FROM feeds WHERE id = $1", feedID)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count > 0 {
		result = true
	}
	return
}

//GetFeedURL -- returnd the feed's url
func GetFeedURL(db *sql.DB, feedID int64) (url string, err error) {
	row := db.QueryRow("SELECT uri FROM feeds WHERE id = $1", feedID)
	err = row.Scan(&url)
	if err != nil {
		return url, fmt.Errorf("Error occured while trying to find the url for feed id (%d): %s", feedID, err.Error())
	}
	return url, nil
}

//GetFeedAuthorID -- returns the feed's author ID
func GetFeedAuthorID(db *sql.DB, feedID int64) (int64, error) {
	var authorID int64
	row := db.QueryRow("SELECT author_id FROM feeds WHERE id = $1", feedID)
	err := row.Scan(&authorID)
	if err != nil {
		return authorID, fmt.Errorf("Error occured while trying to find the author_id for feed id (%d): %s", feedID, err.Error())
	}
	return authorID, nil
}

//UpdateFeedAuthor -- Updates the feed's author
func UpdateFeedAuthor(db *sql.DB, feedID, authorID int64) error {
	_, err := db.Exec("UPDATE feeds SET author_id = $1 WHERE id = $2", authorID, feedID)
	return err
}

//GetFeedRawData -- returns the feed's raw data
func GetFeedRawData(db *sql.DB, feedID int64) (string, error) {
	var rawData string
	row := db.QueryRow("SELECT raw_data FROM feeds WHERE id = $1", feedID)
	err := row.Scan(&rawData)
	if err != nil {
		return rawData, fmt.Errorf("Error occured while trying to find the raw_data for feed id (%d): %s", feedID, err.Error())
	}
	return rawData, nil
}

//UpdateFeedRawData -- Updates the feed's raw data
func UpdateFeedRawData(db *sql.DB, feedID int64, rawData string) error {
	_, err := db.Exec("UPDATE feeds SET raw_data = $1 WHERE id = $2", rawData, feedID)
	return err
}

//GetFeedTitle -- returns the feed title
func GetFeedTitle(db *sql.DB, feedID int64) (string, error) {
	var title string
	row := db.QueryRow("SELECT title FROM feeds WHERE id = $1", feedID)
	err := row.Scan(&title)
	if err != nil {
		return title, fmt.Errorf("Error occured while trying to find the feed title for id (%d): %s", feedID, err.Error())
	}
	return title, nil
}

//UpdateFeedTitle -- Updates the feed title
func UpdateFeedTitle(db *sql.DB, feedID int64, title string) error {
	_, err := db.Exec("UPDATE feeds SET title = $1 WHERE id = $2", title, feedID)
	return err
}

//GetFeedID -- Given a url or title, it returns the feed id
func GetFeedID(db *sql.DB, item string) (int64, error) {
	var id int64
	row := db.QueryRow("SELECT id FROM feeds WHERE uri = $1 OR title = $2", item, item)
	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("Error occured while trying to find the feed id for url/title (%s): %s", item, err.Error())
	}
	return id, nil
}

//AllActiveFeeds -- Returns all active feeds
func AllActiveFeeds(db *sql.DB) map[int64]string {
	var result = make(map[int64]string)

	rows, err := db.Query("SELECT id, uri FROM feeds WHERE deleted = 0")
	if err != nil {
		log.Fatalf("Error happened when trying to get all active feeds: %s", err)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Fatalf("Error happened while trying to close a row: %s", err.Error())
		}
	}()

	for rows.Next() {
		var id int64
		var url string
		err := rows.Scan(&id, &url)
		if err != nil {
			log.Fatalf("Error happened while scanning the rows for the all active feeds function: %s", err.Error())
		}
		result[id] = url
	}
	return result
}

//FilterFeeds -- Takes in a list of feeds and compares them with the feeds listed in the Database.
//Returns all the feeds that are listed as active in the database but where not in the list.
func FilterFeeds(db *sql.DB, feeds map[int64]string) map[int64]string {
	var result = make(map[int64]string)
	allFeeds := AllActiveFeeds(db)

	for dbKey, dbValue := range allFeeds {
		found := false

		for feedKey, feedValue := range feeds {
			if dbKey == feedKey && strings.EqualFold(dbValue, feedValue) {
				found = true
				break
			}
		}

		if !found {
			result[dbKey] = dbValue
		}
	}

	return result
}

//DeleteFeed -- Flips the delete flag on for a feed in the database
func DeleteFeed(db *sql.DB, feedID int64) error {
	_, err := db.Exec("UPDATE feeds SET deleted = 1 WHERE id = $1", feedID)
	return err
}

//UndeleteFeed -- Flips the delete flag off for a feed in the database
func UndeleteFeed(db *sql.DB, feedID int64) error {
	_, err := db.Exec("UPDATE feeds SET deleted = 0 WHERE id = $1", feedID)
	return err
}

//IsFeedDeleted -- Checks to see if the feed is currently marked as deleted
func IsFeedDeleted(db *sql.DB, feedID int64) bool {
	var result bool
	var deleted int64

	row := db.QueryRow("SELECT deleted FROM feeds WHERE id = $1", feedID)
	err := row.Scan(&deleted)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalf("Feed (%d) does not exist: %s", feedID, err.Error())
		} else {
			log.Fatalf("Error happened while trying check the value of the delete flag for feed (%d): %s", feedID, err.Error())
		}
	}

	if deleted == 1 {
		result = true
	} else {
		result = false
	}
	return result
}

//FeedURLExist -- Checks to see if a feed exists
func FeedURLExist(db *sql.DB, url string) bool {
	var id int64
	var result bool

	row := db.QueryRow("SELECT id FROM feeds WHERE uri = $1", url)
	err := row.Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatalf("Error happened when trying to check if the feed (%s) exists: %s", url, err.Error())
		}
	} else {
		result = true
	}
	return result
}

//AddFeedURL -- Adds a feed url to the database
func AddFeedURL(db *sql.DB, url string) (int64, error) {
	var result int64
	feedStmt := "INSERT INTO feeds (uri) VALUES ($1)"

	if FeedURLExist(db, url) {
		return result, fmt.Errorf("Feed already exists")
	}

	dbResult, err := db.Exec(feedStmt, url)
	if err != nil {
		log.Fatal(err)
	}

	result, err = dbResult.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return result, nil
}
