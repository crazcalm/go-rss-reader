package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mmcdole/gofeed"
)

//LoadFeed -- Loads a feed from the database
func LoadFeed(db *sql.DB, id int64) (feed *Feed, err error) {
	url, err := GetFeedURL(db, id)
	if err != nil {
		log.Fatal(err)
	}

	data, err := GetFeedRawData(db, id)
	if err != nil {
		return feed, fmt.Errorf("No data to retrieve: %s", err.Error())
	}

	//Need to convert the data to a gofeed object
	feedParser := gofeed.NewParser()
	feedData, err := feedParser.Parse(strings.NewReader(data))

	var tags []string
	activeTags := AllActiveFeedTags(db, id)

	for _, tag := range activeTags {
		tags = append(tags, tag)
	}

	return &Feed{id, url, tags, feedData}, nil
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
	_, err := db.Exec("UPDATE feeds SET author_id = $2 WHERE id = $1", feedID, authorID)
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

//GetFeedID -- Given a url, it returns the feed id
func GetFeedID(db *sql.DB, url string) (int64, error) {
	var id int64
	row := db.QueryRow("SELECT id FROM feeds WHERE uri = $1", url)
	err := row.Scan(&id)
	if err != nil {
		return id, fmt.Errorf("Error occured while trying to find the feed id for url (%s): %s", url, err.Error())
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
	defer rows.Close()
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
